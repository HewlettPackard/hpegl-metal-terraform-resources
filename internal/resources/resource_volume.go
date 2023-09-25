// (C) Copyright 2020-2023 Hewlett Packard Enterprise Development LP

package resources

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/configuration"
)

const (
	vName          = "name"
	vDescription   = "description"
	vLocation      = "location"
	vLocationID    = "location_id"
	vFlavorID      = "flavor_id"
	vFlavor        = "flavor"
	vSize          = "size"
	vShareable     = "shareable"
	vState         = "state"
	vStatus        = "status"
	vLabels        = "labels"
	vWWN           = "wwn"
	vStoragePool   = "storage_pool"
	vStoragePoolID = "storage_pool_id"
	vCollectionID  = "volume_collection"

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
			ForceNew:    true,
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
			Description: "The minimum size of the volume specified in units of GBytes.",
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				_, ok := val.(float64)
				if !ok {
					errs = append(errs, fmt.Errorf("expected type of %s to be float", key))
					return
				}

				return
			},
		},

		vShareable: {
			Type:        schema.TypeBool,
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
		vLabels: {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The volume labels as (name, value) pairs.",
		},
		vWWN: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The volume serial number.",
		},

		vStoragePool: {
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
			Description: "The storage pool of the volume to be created.",
		},

		vStoragePoolID: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The storage pool of the volume to be created.",
		},
		vCollectionID: {
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
			Description: "The volume collection ID of the volume to be created.",
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

func resourceMetalVolumeCreate(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to create volume")

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

	// handle storage pool inputs
	var (
		vpID, vpName string
	)

	if vpID, ok = d.Get(vStoragePoolID).(string); !ok || vpID == "" {
		// no explicit storage pool ID is set, so try and get one from the storage-pool-name if one was specified.
		if vpName, ok = d.Get(vStoragePool).(string); ok && vpName != "" {
			vpID, _ = p.GetStoragePoolID(vpName)

			if vpID == "" {
				return fmt.Errorf("unable to locate storage pool")
			}
		}
	}

	capacity, ok := d.Get(vSize).(float64)
	if !ok || capacity <= 0 {
		return fmt.Errorf("invalid capacity %v", capacity)
	}

	volume := rest.NewVolume{
		Name:          d.Get(vName).(string),
		Capacity:      int64(capacity),
		Description:   d.Get(vDescription).(string),
		FlavorID:      vfID,
		Shareable:     d.Get(vShareable).(bool),
		StoragePoolID: vpID,
		VolumeCollectionID: func() string {
			if vcid, ok := d.Get(vCollectionID).(string); ok {
				return vcid
			}
			return ""
		}(),
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

	// add tags
	if m, ok := d.Get(vLabels).(map[string]interface{}); ok {
		volume.Labels = convertMap(m)
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

func resourceMetalVolumeRead(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to read volume")

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

	if err = d.Set(vWWN, volume.WWN); err != nil {
		return fmt.Errorf("set WWN: %v", err)
	}

	if volume.Labels != nil {
		tags := make(map[string]string, len(volume.Labels))

		for k, v := range volume.Labels {
			tags[k] = v
		}

		if err := d.Set(vLabels, tags); err != nil {
			return fmt.Errorf("set labels: %v", err)
		}
	}

	if err = d.Set(vStoragePoolID, volume.StoragePoolID); err != nil {
		return fmt.Errorf("set storage pool id: %v", err)
	}

	return nil
}

func resourceMetalVolumeUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to update volume")

	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return
	}

	ctx := c.GetContext()

	vol, _, err := c.Client.VolumesApi.GetByID(ctx, d.Id())
	if err != nil {
		return
	}

	newSize, ok := d.Get(vSize).(float64)
	if !ok {
		return fmt.Errorf("size is not in the expected format")
	}

	updateVol := rest.UpdateVolume{
		ID:   vol.ID,
		ETag: vol.ETag,
	}

	updateVol.Capacity = int64(newSize)

	// add tags
	if m, ok := d.Get(vLabels).(map[string]interface{}); ok {
		updateVol.Labels = convertMap(m)
	}

	_, _, err = c.Client.VolumesApi.Update(ctx, updateVol.ID, updateVol)
	if err != nil {
		return
	}

	pollCount := 0

	for {
		time.Sleep(pollInterval)

		vol, _, err := c.Client.VolumesApi.GetByID(ctx, vol.ID)
		if err != nil {
			return fmt.Errorf("get volume %s: %w", vol.ID, err)
		}

		if vol.SubState != rest.VOLUMESUBSTATE_UPDATE_REQUESTED &&
			vol.SubState != rest.VOLUMESUBSTATE_UPDATING {
			break
		}

		// Fail if volume state hasn't changed after max polls
		if pollCount++; pollCount > pollCountMax {
			return fmt.Errorf("waiting for volume update has timed out")
		}
	}

	if err = c.RefreshAvailableResources(); err != nil {
		return
	}

	return resourceMetalVolumeRead(d, meta)
}

// deleteVAsForVolume deletes all attachments for specified volume.
func deleteVAsForVolume(p *configuration.Config, volID string) error {
	ctx := p.GetContext()

	// Get all attachments
	vas, _, err := p.Client.VolumeAttachmentsApi.List(ctx)
	if err != nil {
		return fmt.Errorf("list volume attachments: %w", err)
	}

	// Initiate attachments deletion for this volume
	for _, va := range vas {
		if va.VolumeID == volID {
			_, err = p.Client.VolumeAttachmentsApi.Delete(ctx, va.ID)
			if err != nil {
				return fmt.Errorf("delete volume attachment %s: %w", va.ID, err)
			}
		}
	}

	// Wait for attachments to be deleted, i.e. volume state to transition out of "visible"
	pollCount := 0

	for {
		time.Sleep(pollInterval)

		volume, _, err := p.Client.VolumesApi.GetByID(ctx, volID)
		if err != nil {
			return fmt.Errorf("get volume %s: %w", volID, err)
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

//nolint:funlen    // Ignoring function length check on existing function
func resourceMetalVolumeDelete(d *schema.ResourceData, meta interface{}) (err error) {
	var volume rest.Volume

	defer wrapResourceError(&err, "failed to delete volume")

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
