// (C) Copyright 2020-2023 Hewlett Packard Enterprise Development LP

package resources

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
)

const (
	nIPPool = "ip_pool"
	nVLAN   = "vlan"
	nVNI    = "vni"

	poolName         = "name"
	poolDescription  = "description"
	poolVer          = "ip_ver"
	poolBaseIP       = "base_ip"
	poolNetmask      = "netmask"
	poolDefaultRoute = "default_route"
	poolSources      = "sources"
	poolDNS          = "dns"
	poolProxy        = "proxy"
	poolNoProxy      = "no_proxy"
	poolNTP          = "ntp"

	sBaseIP = "base_ip"
	sCount  = "count"
)

func ipPoolSourcesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		sBaseIP: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Base IP for the source.",
		},
		sCount: {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Number of IPs to include starting from the base.",
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				c, ok := val.(int)
				if !ok {
					errs = append(errs, fmt.Errorf("expected type of %s to be int got %v", key, reflect.TypeOf(c)))
					return
				}

				if c <= 0 {
					errs = append(errs, fmt.Errorf("%q must be greater than 0, got %v", key, c))
				}
				return
			},
		},
	}
}

func ipPoolSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		poolName: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "A friendly name of the IP pool.",
		},
		poolDescription: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A description of the IP pool.",
		},
		poolVer: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "IP version of the pool (IPv4 or IPv6).",
		},
		poolBaseIP: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Base IP of the pool.",
		},
		poolNetmask: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Netmask of the IP pool.",
		},
		poolDefaultRoute: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Default route of the IP pool.",
		},
		poolSources: {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: ipPoolSourcesSchema(),
			},
			Description: "IP ranges that are to be included in the pool within the base IP and netmask",
		},
		poolDNS: {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "DNS servers to be specified in each allocation from the pool",
		},
		poolProxy: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Web-proxy for external internet access should this pool actually be behind a firewall.",
		},
		poolNoProxy: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IPs or CIDRs for which proxy requests are not made.",
		},
		poolNTP: {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "NTP servers of the IP pool",
		},
	}
}

func networkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		nName: {
			Type:     schema.TypeString,
			Required: true,
		},
		nDescription: {
			Type:     schema.TypeString,
			Optional: true,
		},
		nLocationID: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The PodID of the network",
		},
		nLocation: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Textual representation of the location country:region:enter",
		},
		nHostUse: {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     rest.NETWORKHOSTUSE_OPTIONAL,
			Description: "Required, Optional or Default",
		},
		nPurpose: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Purpose of the network. Allowed values: " + getSupportedNetworkPurpose(),
		},
		nIPPoolID: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP pool ID associated with the network",
		},
		nIPPool: {
			// TODO the V2 SDK doesn't (yet) support TypeMap with Elem *Resource for nested objects
			// This is the currently recommended work-around. See
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/155
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/616
			Type:     schema.TypeSet,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: ipPoolSchema(),
			},
			Description: "Create the specified IP Pool to be used for the network",
		},
		nVLAN: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The VLAN ID of the network",
		},
		nVNI: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The VNI ID of the network",
		},
	}
}

func ProjectNetworkResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalNetworkCreate,
		Read:   resourceMetalNetworkRead,
		Delete: resourceMetalNetworkDelete,
		Update: resourceMetalNetworkUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema:      networkSchema(),
		Description: "Provides Network resource. This allows creation, deletion and update of Metal networks.",
	}
}

func resourceMetalNetworkCreate(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to create network resources")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	locationID, err := p.GetLocationID(d.Get(nLocation).(string))
	if err != nil {
		return err
	}

	var ippool *rest.NewIpPool
	if set, ok := d.Get(nIPPool).(*schema.Set); ok && len(set.List()) != 0 {
		ippool = getIPPool(set)
	}

	newNetwork := rest.NewNetwork{
		Name:        d.Get(nName).(string),
		Description: d.Get(nDescription).(string),
		LocationID:  locationID,
		NewIPPool:   ippool,
	}

	if hostUse, ok := d.Get(nHostUse).(string); ok {
		newNetwork.HostUse = rest.NetworkHostUse(hostUse)
	}

	if purpose, ok := d.Get(nPurpose).(string); ok {
		newNetwork.Purpose = rest.NetworkPurpose(purpose)
	}

	ctx := p.GetContext()
	n, _, err := p.Client.NetworksApi.Add(ctx, newNetwork)
	if err != nil {
		return err
	}

	d.SetId(n.ID)

	if err = d.Set(nIPPoolID, n.IPPoolID); err != nil {
		return err
	}

	if err = p.RefreshAvailableResources(); err != nil {
		return err
	}

	return resourceMetalNetworkRead(d, meta)
}

func getIPPool(set *schema.Set) (ipPool *rest.NewIpPool) {
	ipPool = &rest.NewIpPool{}

	for _, elem := range set.List() {
		pool := elem.(map[string]interface{})

		ipPool.Name = safeString(pool[poolName])
		ipPool.Description = safeString(pool[poolDescription])
		ipPool.IPVersion = rest.IpVer(safeString(pool[poolVer]))
		ipPool.BaseIP = safeString(pool[poolBaseIP])
		ipPool.Netmask = rest.Netmask(safeString(pool[poolNetmask]))
		ipPool.DefaultRoute = safeString(pool[poolDefaultRoute])
		ipPool.Proxy = safeString(pool[poolProxy])
		ipPool.NoProxy = safeString(pool[poolNoProxy])

		var ipSources []rest.IpSource

		if sources, ok := pool[poolSources].([]interface{}); ok {
			for _, source := range sources {
				if s, ok := source.(map[string]interface{}); ok {
					ipSources = append(ipSources, rest.IpSource{
						Base:  safeString(s[sBaseIP]),
						Count: int32(s[sCount].(int)),
					})
				}
			}
		}

		ipPool.Sources = ipSources

		var pDNS []string

		if dns, ok := pool[poolDNS].([]interface{}); ok {
			for _, d := range dns {
				pDNS = append(pDNS, safeString(d))
			}
		}

		ipPool.DNS = pDNS

		var pNTP []string

		if ntp, ok := pool[poolDNS].([]interface{}); ok {
			for _, n := range ntp {
				pNTP = append(pNTP, safeString(n))
			}
		}

		ipPool.NTP = pNTP
	}

	return
}

func resourceMetalNetworkRead(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to read network")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	n, _, err := p.Client.NetworksApi.GetByID(ctx, d.Id())
	if err != nil {
		return err
	}

	if err = d.Set(nName, n.Name); err != nil {
		return err
	}

	if err = d.Set(nDescription, n.Description); err != nil {
		return err
	}

	if err = d.Set(nLocationID, n.LocationID); err != nil {
		return err
	}
	// Attempt best-effort to convert the locationID into huma readbale form. Not fatal
	// if we can't
	l, _ := p.GetLocationName(n.LocationID)

	if err = d.Set(nLocation, l); err != nil {
		return err
	}

	if err = d.Set(nHostUse, n.HostUse); err != nil {
		return err
	}

	if err = d.Set(nPurpose, n.Purpose); err != nil {
		// nolint:wrapcheck // defer func is wrapping the error.
		return err
	}

	if err = d.Set(nIPPoolID, n.IPPoolID); err != nil {
		return err
	}

	if err = d.Set(nVLAN, n.VLAN); err != nil {
		return err
	}

	if err = d.Set(nVNI, n.VNI); err != nil {
		return err
	}

	return nil
}

func resourceMetalNetworkUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to update network")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()

	n, _, err := p.Client.NetworksApi.GetByID(ctx, d.Id())
	if err != nil {
		return err
	}
	n.Name = d.Get(nName).(string)
	n.Description = d.Get(nDescription).(string)

	if hostUse, ok := d.Get(nHostUse).(string); ok {
		n.HostUse = rest.NetworkHostUse(hostUse)
	}

	if purpose, ok := d.Get(nPurpose).(string); ok {
		n.Purpose = rest.NetworkPurpose(purpose)
	}

	_, _, err = p.Client.NetworksApi.Update(ctx, n.ID, n)
	if err != nil {
		return err
	}

	return resourceMetalNetworkRead(d, meta)
}

//nolint: dupl   // Ignoring issues in the existing code
func resourceMetalNetworkDelete(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to delete network")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	_, err = p.Client.NetworksApi.Delete(ctx, d.Id())
	if err != nil {
		return err
	}
	d.SetId("")

	return p.RefreshAvailableResources()
}

// getSupportedNetworkPurpose returns a string containing supported network purpose values.
func getSupportedNetworkPurpose() string {
	return fmt.Sprintf("%v, %v, %v, %v, %v, %v, %v, %v, %v or %v",
		rest.NETWORKPURPOSE_BACKUP,
		rest.NETWORKPURPOSE_STORAGE,
		rest.NETWORKPURPOSE_VM_KERNEL,
		rest.NETWORKPURPOSE_VM_NSX_T,
		rest.NETWORKPURPOSE_V_MOTION,
		rest.NETWORKPURPOSE_V_CHA,
		rest.NETWORKPURPOSE_VM_FT,
		rest.NETWORKPURPOSE_I_SCSI_A,
		rest.NETWORKPURPOSE_I_SCSI_B,
		rest.NETWORKPURPOSE_TELEMETRY,
		rest.NETWORKPURPOSE_EXTERNAL,
	)
}
