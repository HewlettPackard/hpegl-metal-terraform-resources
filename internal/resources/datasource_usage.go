// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP.

package resources

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"

	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
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
	uReady      = "ready"
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
			Optional:    true,
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
					uReady: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Timestamp of when resource machine was ready",
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
					uError: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Any error message",
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
					uReady: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Timestamp of when resource machine was ready",
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
					uError: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Any error message",
					},
				},
			},
		},
	}
}

func DataSourceUsage() *schema.Resource {
	return &schema.Resource{
		Read: resourceMetalUsageRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema:      usageSchema(),
		Description: "Get a usage report",
	}
}

func resourceMetalUsageRead(d *schema.ResourceData, meta interface{}) (err error) {
	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to get usage %s: %w", strings.Trim(nErr.Message(), "\n "), err)

		}
	}()

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	var gOps *rest.UsageReportsApiGetOpts

	start, err := time.Parse(time.RFC3339, d.Get(uUsageStart).(string))
	if err != nil {
		return err
	}
	if d.Get(uUsageEnd).(string) != "" {
		end, err := time.Parse(time.RFC3339, d.Get(uUsageEnd).(string))
		if err != nil {
			return err
		}

		gOps = &rest.UsageReportsApiGetOpts{
			End: optional.NewString(end.String()),
		}
	}

	ctx := p.GetContext()
	usage, _, err := p.Client.UsageReportsApi.Get(ctx, start.Format(time.RFC3339), gOps)
	if err != nil {
		return err
	}
	used := make([]map[string]interface{}, 0, len(usage.Hosts))
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
			uUsageStart: use.UsageStart.Format(time.RFC3339),
			uUsageEnd:   use.UsageEnd.Format(time.RFC3339),
			uAllocated:  use.Allocated.Format(time.RFC3339),
			uFreed:      use.Freed.Format(time.RFC3339),
			uReady:      use.Ready.Format(time.RFC3339),
			uError:      use.Error,
		}
		// Patch up zero times to be emptry strings
		if use.Freed.IsZero() {
			uData[uFreed] = ""
		}
		if use.Ready.IsZero() {
			uData[uReady] = uData[uAllocated]
		}
		used = append(used, uData)
	}
	if err := d.Set(hUsage, used); err != nil {
		return err
	}

	volUsed := make([]map[string]interface{}, 0, len(usage.Volumes))
	for _, use := range usage.Volumes {
		uData := map[string]interface{}{
			"id":        use.VolumeID,
			vName:       use.VolumeName,
			uCapacity:   use.Capacity,
			uFlavorID:   use.FlavorID,
			uFlavorName: use.FlavorName,
			//uRateHourly:     use.RateHourly,
			//uRateMonthly:    use.RateMonthly,
			//uCost:       use.Cost,
			uUsageHours: use.UsageHours,
			uProjectID:  use.ProjectID,
			uLocationID: use.LocationID,
			uUsageStart: use.UsageStart.Format(time.RFC3339),
			uUsageEnd:   use.UsageEnd.Format(time.RFC3339),
			uAllocated:  use.Allocated.Format(time.RFC3339),
			uFreed:      use.Freed.Format(time.RFC3339),
			uReady:      use.Ready.Format(time.RFC3339),
			uError:      use.Error,
		}
		// Patch up zero times to be emptry strings
		if use.Freed.IsZero() {
			uData[uFreed] = ""
		}
		if use.Ready.IsZero() {
			uData[uReady] = uData[uAllocated]
		}
		volUsed = append(volUsed, uData)
	}
	if err := d.Set(vUsage, volUsed); err != nil {
		return err
	}

	d.SetId("usage")
	return nil
}
