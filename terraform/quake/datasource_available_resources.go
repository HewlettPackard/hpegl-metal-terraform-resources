// Copyright (c) 2016-2020 Hewlett Packard Enterprise Development LP.

package quake

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rest "github.com/quattronetworks/quake-client/v1/pkg/client"
)

const (
	// The name are the top level arrays that are available in a terraform block
	// for each time.
	avImages        = "images"
	avSSHKeys       = "ssh_keys"
	avNetworks      = "networks"
	avMachinesSizes = "machine_sizes"
	avVolumes       = "volumes"
	avVolumeFlavors = "volume_flavors"
	avLocations     = "locations"

	// For avImages each terraform block has these attributes.
	iCategory = "category"
	iFlavor   = "flavor"
	iVersion  = "version"

	// For avNetworks each terraform block has these attributes.
	nName        = "name"
	nDescription = "description"
	nKind        = "kind"
	nHostUse     = "host_use"
	nLocation    = "location"
	nLocationID  = "location_id"

	// For avMachineSizes each terraform block has these attributes.
	sName        = "name"
	sQuantity    = "quantity"
	sLocation    = "location"
	sLocationID  = "location_id"
	sDescription = "description"

	// For avVolumeFlavors each terraform block has these attributes.
	fName        = "name"
	fDescription = "description"

	// For avLocations each terraform block has these attributes.
	lCountry = "country"
	lRegion  = "region"
	lCenter  = "data_center"

	// Not avVolumes and avSSHKeys share the schema with the corresponding data sources.
)

func locationResources() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			lCountry: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Country of the location",
			},
			lRegion: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of the location",
			},
			lCenter: {
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

func volumeFlavorResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			fName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the flavor",
			},
			fDescription: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the volume falavor",
			},
		},
	}
}

func imageResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			iCategory: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The image category ",
			},
			iFlavor: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			iVersion: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
		},
	}
}

func machineSizesResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			sName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Machine size name",
			},
			sDescription: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Machine size nadescriptionme",
			},
			sLocationID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The location ID",
			},
			sLocation: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Textual representation of the location country:region:center",
			},
			sQuantity: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of available machines of this size in this location",
			},
		},
	}
}

func existingNetworkResource() *schema.Resource {
	r := &schema.Resource{
		Schema: networkSchema(),
	}
	r.Schema["id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	return r
}

func existingVolumeResource() *schema.Resource {
	r := &schema.Resource{
		Schema: volumeSchema(),
	}
	r.Schema["id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	return r
}

func DataSourceAvailableResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAvailableResourcesRead,
		Schema: map[string]*schema.Schema{
			dsFilter: dataSourceFiltersSchema(),
			avLocations: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     locationResources(),
			},
			avImages: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     imageResource(),
			},
			avSSHKeys: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						sshKeyName: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			avNetworks: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     existingNetworkResource(),
			},
			avMachinesSizes: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     machineSizesResource(),
			},
			avVolumeFlavors: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     volumeFlavorResource(),
			},
			avVolumes: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     existingVolumeResource(),
			},
		},
	}
}

func dataSourceAvailableResourcesRead(d *schema.ResourceData, meta interface{}) (err error) {
	p := meta.(*Config)
	available := p.availableResources

	if err = addLocations(d, available); err != nil {
		return err
	}
	if err = addImages(d, available); err != nil {
		return err
	}
	if err = addSSHKeys(d, available); err != nil {
		return err
	}
	if err = addNetworks(p, d, available); err != nil {
		return err
	}
	if err = addMachineSizes(p, d, available); err != nil {
		return err
	}
	if err = addVolmeFlavors(p, d, available); err != nil {
		return err
	}
	d.SetId("resources")
	return nil
}

func addLocations(d *schema.ResourceData, available rest.AvailableResources) error {
	locations := make([]map[string]interface{}, 0, len(available.Locations))
	for _, loc := range available.Locations {
		iData := map[string]interface{}{
			"id":      loc.ID,
			lCountry:  loc.Country,
			lRegion:   loc.Region,
			lCenter:   loc.DataCenter,
			sLocation: fmt.Sprintf("%s:%s:%s", loc.Country, loc.Region, loc.DataCenter),
		}
		locations = append(locations, iData)
	}
	if err := d.Set(avLocations, locations); err != nil {
		return err
	}
	return nil
}

func addImages(d *schema.ResourceData, available rest.AvailableResources) error {
	images := make([]map[string]interface{}, 0, len(available.Images))
	for _, image := range available.Images {
		iData := map[string]interface{}{
			iFlavor:   image.Flavor,
			iVersion:  image.Version,
			iCategory: image.Category,
			"id":      image.ID,
		}
		images = append(images, iData)
	}
	if err := d.Set(avImages, images); err != nil {
		return err
	}
	return nil
}

func addSSHKeys(d *schema.ResourceData, available rest.AvailableResources) error {
	keys := make([]map[string]interface{}, 0, len(available.SSHKeys))
	for _, key := range available.SSHKeys {
		iData := map[string]interface{}{
			"id":       key.ID,
			sshKeyName: key.Name,
			//sshPublicKey: key.Value,  // Not returned in the available resources API
		}
		keys = append(keys, iData)
	}
	if err := d.Set(avSSHKeys, keys); err != nil {
		return err
	}
	return nil
}

func addNetworks(p *Config, d *schema.ResourceData, available rest.AvailableResources) error {
	networks := make([]map[string]interface{}, 0, len(available.Networks))
	for _, net := range available.Networks {
		iData := map[string]interface{}{
			"id":  net.ID,
			nName: net.Name,
			// Why is this not returned???
			// nDescription: net.Description,
			nKind:       net.Kind,
			nHostUse:    net.HostUse,
			nLocationID: net.LocationID,
		}
		l, _ := p.getLocationName(net.LocationID)
		iData[nLocation] = l
		networks = append(networks, iData)
	}
	if err := d.Set(avNetworks, networks); err != nil {
		return err
	}
	return nil
}

func addMachineSizes(p *Config, d *schema.ResourceData, available rest.AvailableResources) error {
	sizes := make([]map[string]interface{}, 0, len(available.MachineSizes))
	for _, size := range available.MachineSizes {
		var (
			total                int
			locationID, location string
		)
		for _, machines := range available.MachineInventory {
			if machines.SizeID == size.ID {
				total = int(machines.Number)
				locationID = machines.LocationID
				location, _ = p.getLocationName(locationID)
				break
			}
		}
		if total > 0 {
			iData := map[string]interface{}{
				"id":         size.ID,
				sName:        size.Name,
				sDescription: size.Details.Banner1,
				sLocationID:  locationID,
				sLocation:    location,
				sQuantity:    total,
			}
			sizes = append(sizes, iData)
		}
	}
	if err := d.Set(avMachinesSizes, sizes); err != nil {
		return err
	}
	return nil
}

func addVolmeFlavors(p *Config, d *schema.ResourceData, available rest.AvailableResources) error {
	volFalvors := make([]map[string]interface{}, 0, len(available.VolumeFlavors))
	for _, flavor := range available.VolumeFlavors {
		iData := map[string]interface{}{
			"id":         flavor.ID,
			fName:        flavor.Name,
			fDescription: flavor.Details.Banner1,
		}
		volFalvors = append(volFalvors, iData)
	}
	if err := d.Set(avVolumeFlavors, volFalvors); err != nil {
		return err
	}

	existingVols := make([]map[string]interface{}, 0, len(available.Volumes))
	for _, vol := range available.Volumes {
		iData := map[string]interface{}{
			"id":         vol.ID,
			vName:        vol.Name,
			vDescription: vol.Description,
			vSize:        vol.Capacity,
			vLocationID:  vol.LocationID,
			vFlavorID:    vol.FlavorID,
		}
		iData[sLocation], _ = p.getLocationName(vol.LocationID)
		iData[vFlavor], _ = p.getVolumeFlavorName(vol.FlavorID)
		existingVols = append(existingVols, iData)
	}
	if err := d.Set(avVolumes, existingVols); err != nil {
		return err
	}
	return nil
}
