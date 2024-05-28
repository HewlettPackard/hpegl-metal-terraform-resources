// (C) Copyright 2023 Hewlett Packard Enterprise Development LP

package resources

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
)

const iServiceImageFile = "os_service_image_file"

func serviceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		iServiceImageFile: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Path to the YAML file containing the service image definition.",
		},
	}
}

func ServiceImageResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalImageCreate,
		Read:   resourceMetalImageRead,
		Delete: resourceMetalImageDelete,
		Update: resourceMetalImageUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema:      serviceSchema(),
		Description: "Provides service image resource. This allows creation, deletion and update of Metal OS service images.",
	}
}

func resourceMetalImageCreate(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "create OS service image")

	filePath := safeString(d.Get(iServiceImageFile))

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file %s, %v", filePath, err)
	}

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return fmt.Errorf("get client, %v", err)
	}

	ctx := p.GetContext()

	svc, _, err := p.Client.ServicesApi.Add(ctx, file, nil)
	if err != nil {
		return err //nolint:wrapcheck // defer func is wrapping the error.
	}

	d.SetId(svc.ID)

	return nil
}

func resourceMetalImageRead(_ *schema.ResourceData, _ interface{}) (err error) {
	return nil
}

func resourceMetalImageDelete(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "delete OS service image")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return fmt.Errorf("get client, %v", err)
	}

	ctx := p.GetContext()

	if _, err = p.Client.ServicesApi.Delete(ctx, d.Id(), nil); err != nil {
		return err //nolint:wrapcheck // defer func is wrapping the error.
	}

	d.SetId("")

	return nil
}

func resourceMetalImageUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "replace OS service image")

	filePath := safeString(d.Get(iServiceImageFile))

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file %s, %v", filePath, err)
	}

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return fmt.Errorf("get client, %v", err)
	}

	ctx := p.GetContext()

	if _, _, err := p.Client.ServicesApi.Update(ctx, d.Id(), file, nil); err != nil {
		return err //nolint:wrapcheck // defer func is wrapping the error.
	}

	return nil
}
