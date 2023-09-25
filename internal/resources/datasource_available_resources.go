// (C) Copyright 2020-2023 Hewlett Packard Enterprise Development LP

package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/configuration"
)

const (
	// The name are the top level arrays that are available in a terraform block
	// for each time.
	avImages            = "images"
	avSSHKeys           = "ssh_keys"
	avNetworks          = "networks"
	avMachinesSizes     = "machine_sizes"
	avVolumes           = "volumes"
	avVolumeFlavors     = "volume_flavors"
	avLocations         = "locations"
	avStoragePools      = "storage_pools"
	avVolumeCollections = "volume_collections"

	// For avImages each terraform block has these attributes.
	iCategory = "category"
	iFlavor   = "flavor"
	iVersion  = "version"

	// For avNetworks each terraform block has these attributes.
	nName        = "name"
	nDescription = "description"
	nHostUse     = "host_use"
	nPurpose     = "purpose"
	nLocation    = "location"
	nLocationID  = "location_id"
	nIPPoolID    = "ip_pool_id"

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

	// For avStoragePools each terraform block has these attributes.
	spName       = "name"
	spLocation   = "location"
	spLocationID = "location_id"
	spCapacity   = "capacity"

	// For avVolumeCollections each terraform block has these attributes.
	vcName        = "name"
	vcLocationID  = "location_id"
	vcDescription = "description"
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

func volumeCollectionResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			vcName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the volume collection",
			},
			vcLocationID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The location ID",
			},
			vcDescription: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the volume collection",
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

func storagePoolResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			spName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the storage pool",
			},
			spLocationID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The location ID",
			},
			spLocation: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Textual representation of the location country:region:center",
			},
			spCapacity: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The available capacity in units of GiB",
			},
		},
	}
}

func DataSourceAvailableResources() *schema.Resource {
	return &schema.Resource{
		Read:        dataSourceAvailableResourcesRead,
		Description: "Provides a list of available resources in a project for creating Hosts and Volumes.",
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
			avStoragePools: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     storagePoolResource(),
			},
			avVolumeCollections: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     volumeCollectionResource(),
			},
		},
	}
}

func dataSourceAvailableResourcesRead(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to read available resources")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}
	available := p.AvailableResources

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

	if err = addStoragePools(p, d, available); err != nil {
		return err
	}

	if err = addVolumeCollections(p, d, available); err != nil {
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

func addNetworks(p *configuration.Config, d *schema.ResourceData, available rest.AvailableResources) error {
	networks := make([]map[string]interface{}, 0, len(available.Networks))
	for _, net := range available.Networks {
		iData := map[string]interface{}{
			"id":         net.ID,
			nName:        net.Name,
			nDescription: net.Description,
			nHostUse:     net.HostUse,
			nPurpose:     net.Purpose,
			nLocationID:  net.LocationID,
			nIPPoolID:    net.IPPoolID,
			nVLAN:        net.VLAN,
			nVNI:         net.VNI,
		}
		l, _ := p.GetLocationName(net.LocationID)
		iData[nLocation] = l
		networks = append(networks, iData)
	}
	if err := d.Set(avNetworks, networks); err != nil {
		return err
	}
	return nil
}

func addMachineSizes(p *configuration.Config, d *schema.ResourceData, available rest.AvailableResources) error {
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
				location, _ = p.GetLocationName(locationID)
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

func addVolmeFlavors(p *configuration.Config, d *schema.ResourceData, available rest.AvailableResources) error {
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
			"id":           vol.ID,
			vName:          vol.Name,
			vDescription:   vol.Description,
			vSize:          vol.Capacity,
			vLocationID:    vol.LocationID,
			vFlavorID:      vol.FlavorID,
			vStoragePoolID: vol.StoragePoolID,
		}
		iData[sLocation], _ = p.GetLocationName(vol.LocationID)
		iData[vFlavor], _ = p.GetVolumeFlavorName(vol.FlavorID)
		iData[vStoragePool], _ = p.GetStoragePoolName(vol.StoragePoolID)
		existingVols = append(existingVols, iData)
	}
	if err := d.Set(avVolumes, existingVols); err != nil {
		return err
	}
	return nil
}

func addStoragePools(p *configuration.Config, d *schema.ResourceData, available rest.AvailableResources) error {
	existingPools := make([]map[string]interface{}, 0, len(available.StoragePools))

	for _, pool := range available.StoragePools {
		iData := map[string]interface{}{
			"id":         pool.ID,
			spName:       pool.Name,
			spLocationID: pool.LocationID,
			spCapacity:   pool.Capacity,
		}

		iData[spLocation], _ = p.GetLocationName(pool.LocationID)
		existingPools = append(existingPools, iData)
	}

	err := d.Set(avStoragePools, existingPools)

	//nolint:wrapcheck // caller defer func is wrapping the error.
	return err
}

func addVolumeCollections(p *configuration.Config, d *schema.ResourceData, available rest.AvailableResources) error {
	existingVCollections := make([]map[string]interface{}, 0, len(available.VolumeCollections))

	for _, vcol := range available.VolumeCollections {
		iData := map[string]interface{}{
			"id":          vcol.ID,
			vcName:        vcol.Name,
			vcLocationID:  vcol.LocationID,
			vcDescription: vcol.Description,
		}

		existingVCollections = append(existingVCollections, iData)
	}

	err := d.Set(avVolumeCollections, existingVCollections)

	//nolint:wrapcheck // caller defer func is wrapping the error.
	return err
}
