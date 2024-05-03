// (C) Copyright 2020-2023 Hewlett Packard Enterprise Development LP

package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
)

const (
	sshKeyName   = "name"
	sshPublicKey = "public_key"
)

func sshKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		sshKeyName: {
			Type:     schema.TypeString,
			Required: true,
		},

		sshPublicKey: {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func SshKeyResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalSSHKeyCreate,
		Read:   resourceMetalSSHKeyRead,
		Update: resourceMetalSSHKeyUpdate,
		Delete: resourceMetalSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema:      sshKeySchema(),
		Description: "Provides SSH resource. This allows creation, deletion and update of Metal SSHKeys",
	}
}

func resourceMetalSSHKeyCreate(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to create ssh_key")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}
	r := rest.NewSshKey{
		Name: d.Get(sshKeyName).(string),
		Key:  d.Get(sshPublicKey).(string),
	}
	ctx := p.GetContext()
	key, _, err := p.Client.SshkeysApi.Add(ctx, r, nil)
	if err != nil {
		return err
	}
	d.SetId(key.ID)

	if err = p.RefreshAvailableResources(); err != nil {
		return err
	}

	return resourceMetalSSHKeyRead(d, meta)
}

func resourceMetalSSHKeyRead(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to read ssh_key")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	ssh, _, err := p.Client.SshkeysApi.GetByID(ctx, d.Id(), nil)
	if err != nil {
		return err
	}
	d.Set(sshKeyName, ssh.Name)
	d.Set(sshPublicKey, ssh.Key)
	return nil
}

func resourceMetalSSHKeyUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to update ssh_key")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	// Read existing
	ctx := p.GetContext()
	ssh, _, err := p.Client.SshkeysApi.GetByID(ctx, d.Id(), nil)
	if err != nil {
		return err
	}

	updateSSH := rest.UpdateSshKey{
		ID:   ssh.ID,
		ETag: ssh.ETag,
	}

	// Modify
	if name, ok := d.Get(sshKeyName).(string); ok && name != "" {
		updateSSH.Name = name
	}

	if public, ok := d.Get(sshPublicKey).(string); ok && public != "" {
		updateSSH.Key = public
	}

	// Update
	ctx = p.GetContext()
	if _, _, err = p.Client.SshkeysApi.Update(ctx, updateSSH.ID, updateSSH, nil); err != nil {
		return err
	}

	return resourceMetalSSHKeyRead(d, meta)
}

// nolint: dupl   // Ignoring issues in the existing code
func resourceMetalSSHKeyDelete(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to delete ssh_key")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	_, err = p.Client.SshkeysApi.Delete(ctx, d.Id(), nil)
	if err != nil {
		return err
	}
	d.SetId("")
	return p.RefreshAvailableResources()
}
