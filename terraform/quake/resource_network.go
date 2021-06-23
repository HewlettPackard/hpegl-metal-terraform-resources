// (C) Copyright 2016-2021 Hewlett Packard Enterprise Development LP.

package quake

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rest "github.com/hpe-hcss/quake-client/v1/pkg/client"
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
	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to create network resources %s: %w", strings.Trim(string(nErr.Body()), "\n "), err)

		}
	}()

	p, err := getConfigFromMeta(meta)
	if err != nil {
		return err
	}

	locationID, err := p.GetLocationID(d.Get(nLocation).(string))
	if err != nil {
		return err
	}

	newNetwork := rest.NewNetwork{
		Name:        d.Get(nName).(string),
		Description: d.Get(nDescription).(string),
		LocationID:  locationID,
	}

	ctx := p.GetContext()
	n, _, err := p.Client.NetworksApi.Add(ctx, newNetwork)
	if err != nil {
		return err
	}
	d.SetId(n.ID)
	if err = p.RefreshAvailableResources(); err != nil {
		return err
	}
	return resourceQuattroNetworkRead(d, meta)

}

func resourceQuattroNetworkRead(d *schema.ResourceData, meta interface{}) (err error) {
	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to read network %s: %w", strings.Trim(string(nErr.Body()), "\n "), err)

		}
	}()

	p, err := getConfigFromMeta(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	n, _, err := p.Client.NetworksApi.GetByID(ctx, d.Id())
	if err != nil {
		return err
	}
	d.Set(nName, n.Name)
	d.Set(nDescription, n.Description)
	d.Set(nLocationID, n.LocationID)
	// Attempt best-effort to convert the locationID into huma readbale form. Not fatal
	// if we can't
	l, _ := p.GetLocationName(n.LocationID)
	d.Set(nLocation, l)
	d.Set(nKind, n.Kind)
	d.Set(nHostUse, n.HostUse)
	return nil
}

func resourceQuattroNetworkUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to update network %s: %w", strings.Trim(string(nErr.Body()), "\n "), err)

		}
	}()

	p, err := getConfigFromMeta(meta)
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

	_, _, err = p.Client.NetworksApi.Update(ctx, n.ID, n)
	if err != nil {
		return err
	}

	return resourceQuattroNetworkRead(d, meta)
}

func resourceQuattroNetworkDelete(d *schema.ResourceData, meta interface{}) (err error) {
	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to delete network %s: %w", strings.Trim(string(nErr.Body()), "\n "), err)

		}
	}()

	p, err := getConfigFromMeta(meta)
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
