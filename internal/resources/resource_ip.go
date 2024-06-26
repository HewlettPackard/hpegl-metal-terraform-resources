// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
)

const (
	ipPoolID = "ip_pool_id"
	address  = "ip"
	ipUsage  = "usage"
)

func ipSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ipPoolID: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "IP pool ID from which the address will be allocated",
		},
		address: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "IP address to be allocated",
		},
		ipUsage: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Describe usage for the allocated IP",
		},
	}
}

func IPResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPCreate,
		Read:   resourceIPRead,
		Delete: resourceIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema:      ipSchema(),
		Description: "Provides an IP resource. This allows allocation, deallocation of IP addresses from IP Pool. ",
	}
}

func resourceIPCreate(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to create IP resources")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	poolID := d.Get(ipPoolID).(string)
	ip := safeString(d.Get(address))
	allocation := rest.IpAllocation{
		Base:  ip,
		Count: 1,
		Usage: safeString(d.Get(ipUsage)),
	}

	ctx := p.GetContext()

	ipPools, _, err := p.Client.IppoolsApi.List(ctx, nil)
	if err != nil {
		return err
	}

	for _, ipPool := range ipPools {
		if poolID == ipPool.ID || poolID == ipPool.Name {
			poolID = ipPool.ID
			break
		}
	}

	if _, _, err := p.Client.IppoolsApi.AllocateIPs(ctx, poolID, []rest.IpAllocation{allocation}, nil); err != nil {
		return err
	}

	d.SetId(createIPResourceID(poolID, ip))

	return resourceIPRead(d, meta)
}

func resourceIPRead(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to read IP resources")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	poolID := extractIPPoolID(d.Id())
	allocIP := extractIP(d.Id())

	ippool, _, err := p.Client.IppoolsApi.GetByID(ctx, poolID, nil)
	if err != nil {
		return err
	}

	var usage, ip string

	for _, record := range ippool.UseRecords {
		if record.Base == allocIP {
			usage = record.Usage
			ip = record.Base
		}
	}

	if err = d.Set(ipPoolID, ippool.ID); err != nil {
		return err
	}

	if err = d.Set(address, ip); err != nil {
		return err
	}

	if err = d.Set(ipUsage, usage); err != nil {
		return err
	}

	return nil
}

func resourceIPDelete(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to delete IP resources")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	poolID := extractIPPoolID(d.Id())
	ip := extractIP(d.Id())

	if _, _, err = p.Client.IppoolsApi.ReturnIPs(ctx, poolID, []string{ip}, nil); err != nil {
		return err
	}

	return nil
}

func createIPResourceID(ipPoolID, ip string) string {
	return ipPoolID + ":" + ip
}

func extractIPPoolID(resourceID string) string {
	return strings.Split(resourceID, ":")[0]
}

func extractIP(resourceID string) string {
	return strings.Split(resourceID, ":")[1]
}
