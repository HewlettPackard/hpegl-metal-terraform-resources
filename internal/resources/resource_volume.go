// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/configuration"
)

const (
	vName        = "name"
	vDescription = "description"
	vLocation    = "location"
	vLocationID  = "location_id"
	vFlavorID    = "flavor_id"
	vFlavor      = "flavor"
	vSize        = "size"
	vShareable   = "shareable"
	vState       = "state"
	vStatus      = "status"

	// volume Info constants.
	vID          = "id"
	vDiscoveryIP = "discovery_ip"
	vTargetIQN   = "target_iqn"
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
			Computed:    true,
			Description: "The flavor of the volume to be created.",
		},

		vFlavor: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
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
			Required:    true,
			ForceNew:    true,
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
			ForceNew:    true,
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

		vShareable: {
			Type:        schema.TypeBool,
			Required:    false,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: "The volume can be shared by multiple hosts if set.",
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

func VolumeResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalVolumeCreate,
		Read:   resourceMetalVolumeRead,
		Update: resourceMetalVolumeUpdate,
		Delete: resourceMetalVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema:      volumeSchema(),
		Description: "Provides Volume resource. This allows creation, deletion and update of Metal volumes.",
	}
}

//nolint: funlen    // Ignoring function length check on existing function
func resourceMetalVolumeCreate(d *schema.ResourceData, meta interface{}) (err error) {
	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to create volume %s: %w", strings.Trim(nErr.Message(), "\n "), err)

		}
	}()

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}
	// Need to create one
	resources := p.AvailableResources

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
		return fmt.Errorf("unable to locate a volume flavor")
	}

	volume := rest.NewVolume{
		Name:        d.Get(vName).(string),
		Capacity:    int64(d.Get(vSize).(float64)),
		Description: d.Get(vDescription).(string),
		FlavorID:    vfID,
		Shareable:   d.Get(vShareable).(bool),
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

	ctx := p.GetContext()
	v, _, err := p.Client.VolumesApi.Add(ctx, volume)
	if err != nil {
		return err
	}
	d.SetId(v.ID)
	for {
		time.Sleep(pollInterval)

		ctx = p.GetContext()
		vol, _, err := p.Client.VolumesApi.GetByID(ctx, v.ID)
		if err != nil {
			break
		}
		if vol.State != rest.VOLUMESTATE_NEW && vol.State != rest.VOLUMESTATE_ALLOCATING {
			// The Volume create has been processed by the rack-controller so move on.
			break
		}
	}
	if err = p.RefreshAvailableResources(); err != nil {
		return err
	}

	// Now populate additional volume fields.
	return resourceMetalVolumeRead(d, meta)
}

//nolint: funlen    // Ignoring function length check on existing function
func resourceMetalVolumeRead(d *schema.ResourceData, meta interface{}) (err error) {
	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to read volume %s: %w", strings.Trim(nErr.Message(), "\n "), err)

		}
	}()

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	volume, _, err := p.Client.VolumesApi.GetByID(ctx, d.Id())
	if err != nil {
		return err
	}

	d.SetId(volume.ID)
	d.Set(vSize, float64(volume.Capacity/1024/1024))
	d.Set(vName, volume.Name)
	d.Set(vDescription, volume.Description)
	flavorName, _ := p.GetVolumeFlavorName(volume.FlavorID)
	d.Set(vFlavor, flavorName)
	d.Set(vFlavorID, volume.FlavorID)
	loc, _ := p.GetLocationName(volume.LocationID)
	d.Set(vLocation, loc)
	d.Set(vLocationID, volume.LocationID)
	if err = d.Set(vShareable, volume.Shareable); err != nil {
		return err
	}
	d.Set(vState, volume.State)
	d.Set(vStatus, volume.Status)

	return nil
}

//nolint: funlen    // Ignoring function length check on existing function
func resourceMetalVolumeUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to update volume %s: %w", strings.Trim(nErr.Message(), "\n "), err)

		}
	}()

	//@TODO - or not....?
	return resourceMetalVolumeRead(d, meta)
}

// deleteVAsForVolume deletes all attachments for specified volume.
func deleteVAsForVolume(p *configuration.Config, volID string) error {

	ctx := p.GetContext()

	// Get all attachments
	vas, _, err := p.Client.VolumeAttachmentsApi.List(ctx)
	if err != nil {
		return err
	}

	// Initiate attachments deletion for this volume
	for _, va := range vas {
		if va.VolumeID == volID {
			_, err = p.Client.VolumeAttachmentsApi.Delete(ctx, va.ID)
			if err != nil {
				return err
			}
		}
	}

	// Wait for attachments to be deleted, i.e. volume state to transition out of "visible"
	pollCount := 0

	for {
		time.Sleep(pollInterval)

		volume, _, err := p.Client.VolumesApi.GetByID(ctx, volID)
		if err != nil {
			return err
		}

		if volume.State != rest.VOLUMESTATE_VISIBLE {
			break
		}

		// Fail if volume state hasn't changed after max polls
		if pollCount++; pollCount > pollCountMax {
			return fmt.Errorf("waiting for volume state change has timed out")
		}
	}

	return nil
}

//nolint: funlen    // Ignoring function length check on existing function
func resourceMetalVolumeDelete(d *schema.ResourceData, meta interface{}) (err error) {
	var volume rest.Volume

	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to delete volume %s: %w", strings.Trim(nErr.Message(), "\n "), err)

		}
	}()

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	defer func() {
		// This is the last in the deferred chain to fire. If there has been no
		// preceding error we will refresh the available resources and return
		// any possible error that may have caused.
		if err == nil {
			err = p.RefreshAvailableResources()
		}
	}()

	defer func() {
		// Volume deletes are asynchronous in Metal svc and we can not delete terraform's
		// reference to the volume until it has really gone from Metal svc. If we delete the
		// reference too early, or in the presence of errors, we will never be able to retry
		// the delete operation from Terraform (since it has no reference to the resource).
		if err == nil {
			// Volume deletes are async so wait here until Metal svc reports that the volume has really gone.
			for {
				time.Sleep(pollInterval)

				ctx := p.GetContext()
				volume, _, err = p.Client.VolumesApi.GetByID(ctx, d.Id())
				if err != nil {
					return
				}
				switch volume.State {
				case rest.VOLUMESTATE_DELETED:
					// Success; delete terraform reference.
					d.SetId("")
					return

				case rest.VOLUMESTATE_FAILED:
					// Metal svc has finished a delete attempts but failed. Retain the reference to
					// the volume since it technically still exists so that terraform can attempt
					// another delete at a later time.
					err = fmt.Errorf("unable to delete volume")
					return
				}
			}
		}
	}()

	ctx := p.GetContext()
	volume, _, err = p.Client.VolumesApi.GetByID(ctx, d.Id())
	if err != nil {
		return err
	}

	// Nothing to do if volume is already deleted
	if volume.State == rest.VOLUMESTATE_DELETED {
		d.SetId("")
		return nil
	}

	// Delete attachments if volume is visible
	if volume.State == rest.VOLUMESTATE_VISIBLE {
		err = deleteVAsForVolume(p, d.Id())
		if err != nil {
			return err
		}
	}

	_, err = p.Client.VolumesApi.Delete(ctx, d.Id())
	return err
}

func volumeInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		vName: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "A friendly name of the volume attached.",
		},

		vID: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The ID the volume attached.",
		},

		vDiscoveryIP: {
			Type:        schema.TypeString,
			Required:    false,
			Computed:    true,
			Description: "iSCSI Discovery IP.",
		},

		vTargetIQN: {
			Type:        schema.TypeString,
			Required:    false,
			Computed:    true,
			Description: "iSCSI Target IQN.",
		},
	}
}
