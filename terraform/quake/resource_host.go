// Copyright (c) 2016-2020 Hewlett Packard Enterprise Development LP.

package quake

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	rest "github.com/quattronetworks/quake-client/v1/go-client"
)

const (
	// field names for a Quattro host. These are referenceable from some terraform source
	//    resource "quattro_host" "test_host" {
	//       name              = "test"
	//       description      = "hello from Terraform"
	//       image_flavor      = "coreos"
	//       image_version     = "1800.6.0"
	//       machie_size       = "Very Small"
	//       ssh               = ["Chuck's Mac as Team One member"]
	//       networks          = ["Private", "Public"]
	//       location          = "Demo Pod" //"USA:Austin:Demo1"
	//    }
	hName          = "name"
	hDescription   = "description"
	hFlavor        = "image_flavor"
	hImageVersion  = "image_version"
	hLocation      = "location"
	hLocationID    = "location_id"
	hNetworks      = "networks"
	hNetworkIDs    = "network_ids"
	hSSHKeys       = "ssh"
	hSSHKeyIDs     = "ssh_ids"
	hSize          = "machine_size"
	hSizeID        = "machine_size_id"
	hConnections   = "connections"
	hUserData      = "user_data"
	hCHAPUser      = "chap_user"
	hCHAPSecret    = "chap_secret"
	hInitiatorName = "initiator_name"
	hVolumes       = "volumes"
	hState         = "state"
)

func hostSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		hName: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Any friendly name to identify the host that will become the OS hostname in lower case.",
		},
		hFlavor: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "A reference to the image that will be provisioned, eg 'ubuntu'.",
		},
		hImageVersion: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "A specific flavor version, eg '18.04'.",
		},
		hSSHKeys: {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "A list of SSH keys that will be pushed to the host.",
		},
		hSSHKeyIDs: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		hSize: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Some generic sizing information for the machine like 'Small', 'Very Large'.",
		},
		hSizeID: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Machine size ID",
		},
		hLocation: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The location of where the machine will be provisioned, of the form 'country:region:centre', eg 'USA:Texas:AUSL2'.",
		},
		hUserData: {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Any yaml compliant string that will be merged into cloud-init for this host.",
		},
		hLocationID: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "UUID of the location",
		},
		hNetworks: {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "List of network name e.g. ['Public', 'Private'].",
		},
		hNetworkIDs: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "List of network UUIDs.",
		},
		hDescription: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A wordy description of the machine and purpose.",
		},
		hConnections: {
			Type:        schema.TypeMap,
			Computed:    true,
			Description: "A map of network connections and assigned IP addreses, eg {'Private':'10.83.0.17'}.",
		},
		hCHAPUser: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The iSCSI CHAP name for this host.",
		},
		hCHAPSecret: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The iSCSI CHAP secret for this host.",
		},
		hInitiatorName: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The iSCSI initiator name for this host.",
		},
		hVolumes: {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Create volumes and connect them to this host.",
			Elem: &schema.Resource{
				Schema: volumeSchema(),
			},
		},
		hState: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The current state of the host",
		},
	}
}

func hostResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceQuattroHostCreate,
		Read:   resourceQuattroHostRead,
		Delete: resourceQuattroHostDelete,
		Update: resourceQuattroHostUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: hostSchema(),
	}
}

func resourceQuattroHostCreate(d *schema.ResourceData, meta interface{}) (err error) {
	p := meta.(*Config)
	// get available resources
	resources := p.availableResources

	host := rest.NewHost{
		Name:        d.Get(hName).(string),
		Description: d.Get(hDescription).(string),
		UserData:    d.Get(hUserData).(string),
	}

	// 1) verify that flavor and version are sane
	flavorFound := false
	versionFound := false
	targetImageFlavor := d.Get(hFlavor).(string)
	targetImageVersion := d.Get(hImageVersion).(string)
	flavors := []string{}
	for _, image := range resources.Images {
		if image.Flavor == targetImageFlavor {
			flavorFound = true
			if image.Version == targetImageVersion {
				versionFound = true
				host.ServiceID = image.ID
			}
		}
		flavors = append(flavors, fmt.Sprintf("%s@%s", image.Flavor, image.Version))
	}
	if !flavorFound {
		return fmt.Errorf("image flavor %q not found in %q", targetImageFlavor, flavors)
	}
	if !versionFound {
		return fmt.Errorf("image version %q of flavor %q not found in %q", targetImageVersion, targetImageFlavor, flavors)
	}

	// 2) verify that machine size exists and get id
	sizes := []string{}
	targetSize := d.Get(hSize).(string)
	for _, mSize := range resources.MachineSizes {
		if mSize.Name == targetSize {
			host.MachineSizeID = mSize.ID
			break
		}
		sizes = append(sizes, mSize.Name)
	}
	if host.MachineSizeID == "" {
		return fmt.Errorf("machine size %q not found in %q", targetSize, sizes)
	}

	// 3) verify that all of the ssh keys exist and get ids
	var keyIDs, allKeys []string
	for _, name := range convertStringArr(d.Get(hSSHKeys).([]interface{})) {
		found := false
		for _, sshKey := range resources.SSHKeys {
			if sshKey.Name == name {
				found = true
				keyIDs = append(keyIDs, sshKey.ID)
				break
			}
			allKeys = append(allKeys, sshKey.Name)
		}
		if !found {
			return fmt.Errorf("SSH key %q not found in %q", name, allKeys)
		}
	}
	host.SSHKeyIDs = keyIDs

	// 4) parse location, verify it exists, and get id
	locationID, err := p.getLocationID(d.Get(hLocation).(string))
	if err != nil {
		return err
	}
	host.LocationID = locationID

	// Add networks
	processedNetworks := []string{}
	availableNetworks := []string{}
	podNetMap := make(map[string]string)
	podNetMapCount := make(map[string]int)
	podNetIDMap := make(map[string]string)
	for _, podNet := range resources.Networks {
		if podNet.LocationID == host.LocationID {
			podNetMap[podNet.Name] = podNet.ID
			podNetMapCount[podNet.Name] = podNetMapCount[podNet.Name] + 1
			podNetIDMap[podNet.ID] = podNet.Name
			availableNetworks = append(availableNetworks, podNet.Name)
		}
	}
	for _, network := range convertStringArr(d.Get(hNetworks).([]interface{})) {
		if _, ok := podNetIDMap[network]; ok {
			processedNetworks = append(processedNetworks, network)
			continue
		}
		if id, ok := podNetMap[network]; ok {
			if podNetMapCount[network] > 1 {
				return fmt.Errorf("network %q is ambiguous in location %q %s", network, host.LocationID, availableNetworks)
			}
			processedNetworks = append(processedNetworks, id)
			continue
		}
		return fmt.Errorf("network %q not found in location %q %s", network, host.LocationID, availableNetworks)
	}
	if len(processedNetworks) == 0 {
		return fmt.Errorf("no networks in %q found in %q", d.Get(hNetworks), availableNetworks)
	}
	host.NetworkIDs = processedNetworks

	// 6 add any new volumes
	if vols, vOK := d.Get(hVolumes).(*schema.Set); vOK {
		for _, element := range vols.List() {
			var vID, vfID string
			vID = ""
			vol := element.(map[string]interface{})
			// Determine if this volume is an existing volume
			if vName, ok := vol[vName].(string); ok && vName != "" {
				for _, volume := range resources.Volumes {
					if vName == volume.Name {
						vID = volume.ID
					}
				}
			}
			vfID, ok := vol[vFlavorID].(string)
			if !ok || vfID == "" {
				flavorName, ok := vol[vFlavor]
				if ok && flavorName != "" {
					for _, flavor := range resources.VolumeFlavors {
						if flavor.Name == flavorName {
							vfID = flavor.ID
							break
						}
					}
				} else {
					return fmt.Errorf("volume %q needs %q or a %q to be set", vol[vName], vFlavorID, vFlavor)
				}
			}
			if vID != "" {
				host.VolumeIDs = append(host.VolumeIDs, vID)
			} else {
				host.NewVolumes = append(host.NewVolumes, rest.AddVolume{
					Name:        vol[vName].(string),
					Description: vol[vDescription].(string),
					Capacity:    uint64(vol[vSize].(float64)),
					FlavorID:    vfID,
				})
			}
		}
	}

	// Create it
	h, _, err := p.client.HostsApi.Add(p.context, host)
	if err != nil {
		return fmt.Errorf("failed to create host %+v: %v", host, err)
	}

	d.SetId(h.ID)
	return resourceQuattroHostRead(d, meta)
}

func resourceQuattroHostRead(d *schema.ResourceData, meta interface{}) error {
	p := meta.(*Config)
	host, _, err := p.client.HostsApi.GetByID(p.context, d.Id())
	if err != nil {
		return err
	}
	d.Set(hName, host.Name)
	d.Set(hState, host.State)
	d.Set(hFlavor, host.ServiceFlavor)
	d.Set(hImageVersion, host.ServiceVersion)
	d.Set(hSSHKeyIDs, host.SSHAuthorizedKeys)
	d.Set(hSizeID, host.MachineSizeID)
	d.Set(hSize, host.MachineSizeName)
	d.Set(hUserData, host.UserData)
	loc, _ := p.getLocationName(host.LocationID)
	d.Set(hLocation, loc)
	d.Set(hLocationID, host.LocationID)
	d.Set(hNetworkIDs, host.NetworkIDs)
	d.Set(hDescription, host.Description)
	hCons := make(map[string]interface{})
	for _, con := range host.Connections {
		for _, hNet := range con.Networks {
			hCons[hNet.Name] = hNet.IP
		}
	}
	d.Set(hConnections, hCons)
	d.Set(hCHAPUser, host.ISCSIConfig.CHAPUser)
	d.Set(hCHAPSecret, host.ISCSIConfig.CHAPSecret)
	d.Set(hInitiatorName, host.ISCSIConfig.InitiatorName)
	return nil
}

func resourceQuattroHostUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceQuattroHostRead(d, meta)
}

func resourceQuattroHostDelete(d *schema.ResourceData, meta interface{}) (err error) {
	p := meta.(*Config)

	defer func() {
		if err == nil {
			d.SetId("")
		}
	}()
	host, _, err := p.client.HostsApi.GetByID(p.context, d.Id())
	if err != nil {
		return err
	}
	if host.State == rest.HOSTSTATE_DELETED {
		return nil
	}
	if host.State == rest.HOSTSTATE_NEW {
		_, err = p.client.HostsApi.Delete(p.context, d.Id())
		return err
	}
	if host.PowerStatus == rest.HOSTPOWERSTATE_ON {
		_, _, err = p.client.HostsApi.PowerOff(p.context, d.Id())
		if err != nil {
			return err
		}
		for host.PowerStatus != rest.HOSTPOWERSTATE_OFF {
			time.Sleep(pollInterval)
			host, _, err = p.client.HostsApi.GetByID(p.context, d.Id())
			if err != nil {
				return err
			}
		}
	}
	_, err = p.client.HostsApi.Delete(p.context, d.Id())
	if err == nil {
		for host.State != rest.HOSTSTATE_DELETED {
			time.Sleep(pollInterval)
			if host, _, err = p.client.HostsApi.GetByID(p.context, d.Id()); err != nil {
				return err
			}
		}
	}
	return err
}
