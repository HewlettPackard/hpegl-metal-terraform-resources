// Copyright (c) 2016-2020 Hewlett Packard Enterprise Development LP.

package quake

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rest "github.com/quattronetworks/quake-client/v1/pkg/client"
)

const (
// field names for a Quattro network. These are referenceable from some terraform source
//    resource "quattro_network" "test_net" {
//       name         = "net_test"
//       description  = "Terraform network"
//       location     = "USA:Austin:Demo1"
//    }
)

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
			Description: "Textual representation of the location country:region:enter",
		},
		nKind: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Shared, Private or Custom",
		},
		nHostUse: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Required, Optional or Default",
		},
	}
}

func ProjectNetworkResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceQuattroNetworkCreate,
		Read:   resourceQuattroNetworkRead,
		Delete: resourceQuattroNetworkDelete,
		Update: resourceQuattroNetworkUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: networkSchema(),
	}
}

func resourceQuattroNetworkCreate(d *schema.ResourceData, meta interface{}) (err error) {
	p := meta.(*Config)

	locationID, err := p.getLocationID(d.Get(nLocation).(string))
	if err != nil {
		return err
	}

	newNetwork := rest.NewNetwork{
		Name:        d.Get(nName).(string),
		Description: d.Get(nDescription).(string),
		LocationID:  locationID,
	}

	n, _, err := p.client.NetworksApi.Add(p.context, newNetwork)
	if err != nil {
		return err
	}
	d.SetId(n.ID)
	if err = p.refreshAvailableResources(); err != nil {
		return err
	}
	return resourceQuattroNetworkRead(d, meta)

}

func resourceQuattroNetworkRead(d *schema.ResourceData, meta interface{}) (err error) {
	p := meta.(*Config)
	n, _, err := p.client.NetworksApi.GetByID(p.context, d.Id())
	if err != nil {
		return err
	}
	d.Set(nName, n.Name)
	d.Set(nDescription, n.Description)
	d.Set(nLocationID, n.LocationID)
	// Attempt best-effort to convert the locationID into huma readbale form. Not fatal
	// if we can't
	l, _ := p.getLocationName(n.LocationID)
	d.Set(nLocation, l)
	d.Set(nKind, n.Kind)
	d.Set(nHostUse, n.HostUse)
	return nil
}

func resourceQuattroNetworkUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	p := meta.(*Config)

	n, _, err := p.client.NetworksApi.GetByID(p.context, d.Id())
	if err != nil {
		return err
	}
	n.Name = d.Get(nName).(string)
	n.Description = d.Get(nDescription).(string)

	_, _, err = p.client.NetworksApi.Update(p.context, n.ID, n)
	if err != nil {
		return err
	}

	return resourceQuattroNetworkRead(d, meta)
}

func resourceQuattroNetworkDelete(d *schema.ResourceData, meta interface{}) (err error) {
	p := meta.(*Config)

	_, err = p.client.NetworksApi.Delete(p.context, d.Id())
	if err != nil {
		return err
	}
	d.SetId("")

	return p.refreshAvailableResources()
}
