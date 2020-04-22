// Copyright (c) 2016-2020 Hewlett Packard Enterprise Development LP.

package quake

import (
	"time"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform/helper/schema"
	rest "github.com/quattronetworks/quake-client/v1/go-client"
)

const (
	//uBillingMonths = "billing_months"

	hUsage = "host_usage"
	vUsage = "volume_usage"

	//uCost        = "cost"
	//uRateHourly  = "hourly_rate"
	//uRateMonthly = "monthly_rate"
	//uCurrency      = "currency"

	// standard fields for all usage elements
	uProjectID  = "project_id"
	uLocationID = "location_id"
	uAllocated  = "allocated"
	uFreed      = "freed"
	uUsageStart = "start"
	uUsageEnd   = "end"
	uUsageHours = "usage_hours"
	uError      = "error"

	// Host speciific extras
	uMachineSizeName = "machine_size"
	uMachineSizeID   = "machine_size_id"
	uHostName        = "name"
	uHostID          = "id"

	// Volume specific extras
	uVolumeName = "name"
	uVolumeID   = "id"
	uFlavorName = "flavor"
	uFlavorID   = "falvor_id"
	uCapacity   = "capacity"
)

func usageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		uUsageStart: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Start time for usage calculation, format is RFC 3339 e.g. 2018-05-13T07:44:12Z",
		},
		uUsageEnd: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "End time for usage calculation, format is RFC 3339 e.g. 2018-05-13T07:44:12Z",
		},
		hUsage: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					uHostID: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Name of the MachineSize requested when host was created",
					},
					uHostName: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Name of the associated Host",
					},
					uMachineSizeName: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Name of the MachineSize requested when host was created",
					},
					uMachineSizeID: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "ID of the MachineSize requested when host was created",
					},
					// uCost: {
					// 	Type:     schema.TypeString,
					// 	Computed: true,
					// },
					// uRateHourly: {
					// 	Type:     schema.TypeString,
					// 	Computed: true,
					// },
					// uRateMonthly: {
					// 	Type:     schema.TypeString,
					// 	Computed: true,
					// },
					uUsageHours: {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "The difference between the UsageEnd and UsageStart rounded up to the UsageHours",
					},
					uProjectID: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Project ID that created the host",
					},
					uFreed: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Timestamp of when resource machine was freed",
					},
					uAllocated: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Timestamp of when resource machine was allocated",
					},
					uUsageStart: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The start of the usage reporting window or when the host was allocated",
					},
					uUsageEnd: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The end of the usage reporting window or when the host was freed",
					},
					uLocationID: {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		vUsage: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					uVolumeID: {
						Type:     schema.TypeString,
						Computed: true,
					},
					uVolumeName: {
						Type:     schema.TypeString,
						Computed: true,
					},
					uFlavorName: {
						Type:     schema.TypeString,
						Computed: true,
					},
					uFlavorID: {
						Type:     schema.TypeString,
						Computed: true,
					},
					uCapacity: {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					// uCost: {
					// 	Type:     schema.TypeString,
					// 	Computed: true,
					// },
					// uRateHourly: {
					// 	Type:     schema.TypeString,
					// 	Computed: true,
					// },
					// uRateMonthly: {
					// 	Type:     schema.TypeString,
					// 	Computed: true,
					// },
					uUsageHours: {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "The difference between the UsageEnd and UsageStart rounded up to the UsageHours",
					},
					uProjectID: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Project ID that created the volume",
					},
					uFreed: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Timestamp of when the volume was freed",
					},
					uAllocated: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Timestamp of when the volume was allocated",
					},
					uUsageStart: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The start of the usage reporting window or when the volume was allocated",
					},
					uUsageEnd: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The end of the usage reporting window or when the volume was freed",
					},
					uLocationID: {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func dataSourceUsage() *schema.Resource {
	return &schema.Resource{
		Read: resourceQuakeUsageRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: usageSchema(),
	}
}

func resourceQuakeUsageRead(d *schema.ResourceData, meta interface{}) error {
	p := meta.(*Config)
	var gOps *rest.GetOpts
	start, err := time.Parse(time.RFC3339, d.Get(uUsageStart).(string))
	if err != nil {
		return err
	}
	if d.Get(uUsageEnd).(string) != "" {
		end, err := time.Parse(time.RFC3339, d.Get(uUsageEnd).(string))
		if err != nil {
			return err
		}
		gOps = &rest.GetOpts{
			End: optional.NewString(end.String()),
		}
	}
	usage, _, err := p.client.UsageApi.Get(p.context, start.Format(time.RFC3339), gOps)
	if err != nil {
		return err
	}
	var used = make([]map[string]interface{}, 0, len(usage.Hosts))
	for _, use := range usage.Hosts {
		uData := map[string]interface{}{
			uHostID:          use.HostID,
			uHostName:        use.HostName,
			uMachineSizeName: use.MachineSizeName,
			uMachineSizeID:   use.MachineSizeID,
			//uRateHourly:   use.RateHourly,
			//uRateMonthly:  use.RateMonthly,
			//uCost:         use.Cost,
			uUsageHours: use.UsageHours,
			uProjectID:  use.ProjectID,
			uLocationID: use.LocationID,
			uUsageStart: use.UsageStart,
			uUsageEnd:   use.UsageEnd,
			uAllocated:  use.Allocated.String(),
			uFreed:      use.Freed.String(),
			uError:      use.Error,
		}
		used = append(used, uData)
	}
	if err := d.Set(hUsage, used); err != nil {
		return err
	}

	var volUsed = make([]map[string]interface{}, 0, len(usage.Volumes))
	for _, use := range usage.Volumes {
		uData := map[string]interface{}{
			"id":      use.VolumeID,
			vName:     use.VolumeName,
			vSize:     use.Capacity,
			vFlavor:   use.FlavorID,
			vFlavorID: use.FlavorID,
			//uRateHourly:     use.RateHourly,
			//uRateMonthly:    use.RateMonthly,
			//uCost:       use.Cost,
			uUsageHours: use.UsageHours,
			uProjectID:  use.ProjectID,
			uLocationID: use.LocationID,
			uUsageStart: use.UsageStart,
			uUsageEnd:   use.UsageEnd,
			uAllocated:  use.Allocated.String(),
			uFreed:      use.Freed.String(),
			uError:      use.Error,
		}
		volUsed = append(volUsed, uData)
	}
	if err := d.Set(vUsage, volUsed); err != nil {
		return err
	}

	d.SetId("usage")
	return nil
}
