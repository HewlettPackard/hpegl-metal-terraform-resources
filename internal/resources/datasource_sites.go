// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
)

const (
	// The name are the top level arrays that are available in a terraform block
	// for each time.

	lSites = "sites"

	// For lSites each terraform block has these attributes.
	sCountry = "country"
	sRegion  = "region"
	sCenter  = "data_center"
)

func sitesResources() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			sCountry: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Country of the location",
			},
			sRegion: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of the location",
			},
			sCenter: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data center of the location",
			},
			sLocation: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location combination center:region:country",
			},
		},
	}
}

func DataSourceSites() *schema.Resource {
	return &schema.Resource{
		Read:        dataSourceSitesRead,
		Description: "Provides a list of sites.",
		Schema: map[string]*schema.Schema{
			lSites: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sitesResources(),
			},
		},
	}
}

func dataSourceSitesRead(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to read projects-info")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	info, _, err := p.Client.ProjectsInfoApi.List(ctx, nil)
	if err != nil {
		return err
	}

	if err = addSites(d, info.Summary.Locations); err != nil {
		return err
	}
	d.SetId("sites")
	return nil
}

func addSites(d *schema.ResourceData, locs []rest.LocationInfo) error {
	locations := make([]map[string]interface{}, 0, len(locs))
	for _, loc := range locs {
		iData := map[string]interface{}{
			"id":      loc.ID,
			sCountry:  loc.Country,
			sRegion:   loc.Region,
			sCenter:   loc.DataCenter,
			sLocation: fmt.Sprintf("%s:%s:%s", loc.Country, loc.Region, loc.DataCenter),
		}
		locations = append(locations, iData)
	}
	if err := d.Set(lSites, locations); err != nil {
		return err
	}
	return nil
}
