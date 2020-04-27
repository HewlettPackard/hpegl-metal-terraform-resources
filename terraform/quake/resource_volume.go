// Copyright (c) 2016-2020 Hewlett Packard Enterprise Development LP.

package quake

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	rest "github.com/quattronetworks/quake-client/v1/go-client"
)

const (
	vName        = "name"
	vDescription = "description"
	vLocation    = "location"
	vLocationID  = "location_id"
	vFlavorID    = "flavor_id"
	vFlavor      = "flavor"
	vSize        = "size"
	vState       = "state"
	vStatus      = "status"
)

func volumeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		vName: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "A friendly name of the volume.",
		},

		vFlavorID: {
			Type:        schema.TypeString,
			Required:    false,
			Computed:    true,
			Description: "The flavor of the volume to be created.",
		},

		vFlavor: {
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
			Description: "The flavor of the volume to be created.",
		},

		vDescription: {
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
			Description: "A wordy description of the volume and purpose.",
		},

		vLocation: {
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
			Description: "Location of the volume country:region:data-center.",
		},

		vLocationID: {
			Type:        schema.TypeString,
			Required:    false,
			Computed:    true,
			Description: "LocationID.",
		},

		vSize: {
			Type:        schema.TypeFloat,
			Required:    true,
			Description: "The minimum size of the volume specified in units of GBytes.",
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				sz, ok := val.(float64)
				if !ok {
					errs = append(errs, fmt.Errorf("expected type of %s to be float", key))
					return
				}
				if sz <= 0 {
					errs = append(errs, fmt.Errorf("%q must be greater than 0, got %f", key, sz))
				}
				return
			},
		},

		vState: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The volume provisioning state.",
		},

		vStatus: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The volume provisioning status.",
		},
	}
}

func volumeResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceQuatrroVolumeCreate,
		Read:   resourceQuatrroVolumeRead,
		Update: resourceQuatrroVolumeUpdate,
		Delete: resourceQuatrroVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: volumeSchema(),
	}
}

func resourceQuatrroVolumeCreate(d *schema.ResourceData, meta interface{}) (err error) {
	p := meta.(*Config)
	// Need to create one
	resources := p.availableResources

	var (
		vfID, vfName string
		ok           bool
	)
	if vfID, ok = d.Get(vFlavorID).(string); !ok || vfID == "" {
		// no explicit volume flavor ID is set, so go and try and get one
		// from the volume-flavor-name.
		for _, flavor := range resources.VolumeFlavors {
			if vfName, ok = d.Get(vFlavor).(string); ok {
				if flavor.Name == vfName {
					vfID = flavor.ID
					break
				}
			}
		}
	}
	if vfID == "" {
		return fmt.Errorf("Unable to locate a volume flavor")
	}

	volume := rest.NewVolume{
		Name:        d.Get(vName).(string),
		Capacity:    uint64(d.Get(vSize).(float64)),
		Description: d.Get(vDescription).(string),
		FlavorID:    vfID,
	}

	targetLocation, ok := d.Get(vLocation).(string)
	if !ok || targetLocation == "" {
		return fmt.Errorf("%q must be set", vLocation)
	}

	locations := []string{}
	pieces := strings.Split(targetLocation, ":")
	if len(pieces) != 3 {
		return fmt.Errorf("%q must be of the form country:region:data-center", vLocation)
	}

	found := false
	for _, loc := range resources.Locations {
		if len(pieces) == 3 {
			if string(loc.Country) == pieces[0] && loc.Region == pieces[1] && loc.DataCenter == pieces[2] {
				volume.LocationID = loc.ID
				found = true
				break
			}
		}
		locations = append(locations, fmt.Sprintf("%s:%s:%s", loc.Country, loc.Region, loc.DataCenter))
	}
	if !found {
		return fmt.Errorf("location %q not found in %q", targetLocation, locations)
	}

	v, _, err := p.client.VolumesApi.Add(p.context, volume)
	if err != nil {
		return err
	}
	d.SetId(v.ID)
	for {
		time.Sleep(pollInterval)
		vol, _, err := p.client.VolumesApi.GetByID(p.context, v.ID)
		if err != nil {
			break
		}
		if vol.State != rest.VOLUMESTATE_NEW {
			// The Volume create has been processed by the rack-controller so move on.
			break
		}
	}
	if err = p.refreshAvailableResources(); err != nil {
		return err
	}
	// Now populate additional volume fields.
	return resourceQuatrroVolumeRead(d, meta)
}

func resourceQuatrroVolumeRead(d *schema.ResourceData, meta interface{}) error {
	p := meta.(*Config)

	volume, _, err := p.client.VolumesApi.GetByID(p.context, d.Id())
	if err != nil {
		return err
	}

	d.SetId(volume.ID)
	d.Set(vSize, float64(volume.Capacity/1024/1024))
	d.Set(vName, volume.Name)
	d.Set(vDescription, volume.Description)
	flavorName, _ := p.getVolumeFlavorName(volume.FlavorID)
	d.Set(vFlavor, flavorName)
	d.Set(vFlavorID, volume.FlavorID)
	loc, _ := p.getLocationName(volume.LocationID)
	d.Set(vLocation, loc)
	d.Set(vLocationID, volume.LocationID)
	d.Set(vState, volume.State)
	d.Set(vStatus, volume.Status)
	return nil
}

func resourceQuatrroVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	//@TODO - or not....?
	return resourceQuatrroVolumeRead(d, meta)
}

func resourceQuatrroVolumeDelete(d *schema.ResourceData, meta interface{}) (err error) {
	var volume rest.Volume
	p := meta.(*Config)
	defer func() {
		err = p.refreshAvailableResources()
	}()
	defer func() {
		if err == nil {
			// Volume deletes are async so wait here
			for {
				time.Sleep(pollInterval)
				volume, _, err = p.client.VolumesApi.GetByID(p.context, d.Id())
				if err != nil {
					return
				}
				if volume.State == rest.VOLUMESTATE_DELETED || volume.State == rest.VOLUMESTATE_FAILED {
					d.SetId("")
					return
				}
			}
		}
	}()
	volume, _, err = p.client.VolumesApi.GetByID(p.context, d.Id())
	if err != nil {
		return err
	}
	if volume.State == rest.VOLUMESTATE_DELETED {
		d.SetId("")
		return nil
	}
	_, err = p.client.VolumesApi.Delete(p.context, d.Id())
	return err
}
