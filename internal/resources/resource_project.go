// (C) Copyright 2020-2025 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
)

const (
	pName    = "name"
	pProfile = "profile"
	pLimits  = "limits"
	pSites   = "sites"

	pProjectName        = "project_name"
	pProjectDescription = "project_description"
	pCompany            = "company"
	pAddress            = "address"
	pEmail              = "email"
	pEmailVerified      = "email_verified"
	pPhoneNumber        = "phone_number"
	pPhoneVerified      = "phone_number_verified"

	pHosts           = "hosts"
	pVolumes         = "volumes"
	pVolumeCapacity  = "volume_capacity"
	pPrivateNetworks = "private_networks"
	pInstanceTypes   = "instance_types"
	pPermittedImages = "permitted_images"

	pVolumeReplicationEnabled = "volume_replication_enabled"
	pBootFromSANSupport       = "boot_from_san_support"
)

func limitsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		pHosts: {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Maximum number of host allowed in the team.",
		},
		pVolumes: {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Maximum number of volumes allowed in the team.",
		},
		pVolumeCapacity: {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "Total allowable volume capacity (GiB) allowed in the team.",
		},
		pPrivateNetworks: {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Maximum number of private networks allowed in the team.",
		},
		pInstanceTypes: {
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Optional:    true,
			Description: "Map of instance type ID to maximum number of hosts that can be created with that instance type",
		},
	}
}

func profileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		pProjectName: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "A friendly name of the team.",
		},
		pProjectDescription: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A friendly description of the team.",
		},
		pCompany: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The company associated with the team.",
		},
		pAddress: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The company address with the team.",
		},
		pEmail: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Email address.",
		},
		pEmailVerified: {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Email address has been validated.",
		},
		pPhoneNumber: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Phine number.",
		},
		pPhoneVerified: {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Phine number has been validated.",
		},
	}
}

func projectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		pName: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "A friendly name of the project.",
		},

		pProfile: {
			// TODO the V2 SDK doesn't (yet) support TypeMap with Elem *Resource for nested objects
			// This is the currently recommended work-around. See
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/155
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/616
			Type:        schema.TypeList,
			MaxItems:    1,
			Required:    true,
			Description: "Team profile.",
			Elem: &schema.Resource{
				Schema: profileSchema(),
			},
		},
		pLimits: {
			// TODO the V2 SDK doesn't (yet) support TypeMap with Elem *Resource for nested objects
			// This is the currently recommended work-around. See
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/155
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/616
			Type:        schema.TypeList,
			MaxItems:    1,
			Required:    true,
			Description: "Resource limits applied to this team.",
			Elem: &schema.Resource{
				Schema: limitsSchema(),
			},
		},
		pSites: {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "List of Permitted Site IDs",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		pPermittedImages: {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "List of permitted OS service images",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		pVolumeReplicationEnabled: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: "Volume replication is enabled for the project if set.",
		},

		pBootFromSANSupport: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Boot-from-SAN feature is enabled for the project if set.",
		},
	}
}

func ProjectResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalProjectCreate,
		Read:   resourceMetalProjectRead,
		Delete: resourceMetalProjectDelete,
		Update: resourceMetalProjectUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema:      projectSchema(),
		Description: "Provides Project resource. This allows creation, deletion and update of Metal projects.",
	}
}

func resourceMetalProjectCreate(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to create project")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	np := rest.NewProject{
		Name: d.Get(pName).(string),

		VolumeReplicationEnabled: d.Get(pVolumeReplicationEnabled).(bool),
		BootFromSANSupport:       d.Get(pBootFromSANSupport).(bool),
	}

	if list, ok := d.Get(pProfile).([]interface{}); ok && len(list) == 1 {
		if np.Profile, err = getProfile(list[0]); err != nil {
			return fmt.Errorf("failed to create project %s: %v", np.Name, err)
		}
	} else {
		return fmt.Errorf("failed to create project %s: only 1 profile block is allowed", np.Name)
	}

	if list, ok := d.Get(pLimits).([]interface{}); ok && len(list) == 1 {
		if np.Limits, err = getLimits(list[0]); err != nil {
			return fmt.Errorf("failed to create project %s: %v", np.Name, err)
		}
	} else {
		return fmt.Errorf("failed to create project %s: only 1 limit block is allowed", np.Name)
	}

	if f, ok := d.GetOk(pSites); ok {
		s, ok := f.(*schema.Set)
		if !ok {
			err = fmt.Errorf("sites list is not in the expected format")

			return err
		}

		np.PermittedSites = expandStringList(s.List())
	}

	if f, ok := d.GetOk(pPermittedImages); ok {
		s, ok := f.(*schema.Set)
		if !ok {
			err = fmt.Errorf("permitted images list is not in the expected format")

			return err
		}

		np.PermittedOSImages = expandStringList(s.List())
	}

	ctx := p.GetContext()

	// TO DO:
	//  1. Remove 'Space' from the default header list with REST client.
	//  2. Instantiate ProjectsApiAddOpts with space name and pass as 3rd arg.
	//     For now, passing 'nil' as 3rd arg as the 'Space' header is already
	//     configured to be sent as one of the default header with REST client.
	project, _, err := p.Client.ProjectsApi.Add(ctx, np, nil)
	if err != nil {
		return err
	}
	d.SetId(project.ID)

	return resourceMetalProjectRead(d, meta)
}

func getUpdateProfile(profile interface{}) (p rest.UpdateProfile, err error) {
	profileMap, ok := profile.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("wrong profile format")

		return
	}

	return rest.UpdateProfile{
		Address:     safeString(profileMap[pAddress]),
		Company:     safeString(profileMap[pCompany]),
		Email:       safeString(profileMap[pEmail]),
		PhoneNumber: safeString(profileMap[pPhoneNumber]),
		TeamDesc:    safeString(profileMap[pProjectDescription]),
		TeamName:    safeString(profileMap[pProjectName]),
	}, nil
}

func getProfile(profile interface{}) (p rest.Profile, err error) {
	profileMap, ok := profile.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("wrong profile format")
		return
	}

	return rest.Profile{
		Address:     safeString(profileMap[pAddress]),
		Company:     safeString(profileMap[pCompany]),
		Email:       safeString(profileMap[pEmail]),
		PhoneNumber: safeString(profileMap[pPhoneNumber]),
		TeamDesc:    safeString(profileMap[pProjectDescription]),
		TeamName:    safeString(profileMap[pProjectName]),
	}, nil
}

func getUpdateLimits(limits interface{}) (p rest.UpdateLimits, err error) {
	limitsMap, ok := limits.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("wrong limits format")

		return
	}

	return rest.UpdateLimits{
		Hosts:           int32(safeInt(limitsMap[pHosts])),
		Volumes:         int32(safeInt(limitsMap[pVolumes])),
		VolumeCapacity:  int64(safeFloat(limitsMap[pVolumeCapacity])),
		PrivateNetworks: int32(safeInt(limitsMap[pPrivateNetworks])),
		InstanceTypes:   safeMapStrInt32(limitsMap[pInstanceTypes]),
	}, nil
}

func getLimits(limits interface{}) (p rest.Limits, err error) {
	limitsMap, ok := limits.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("wrong limits format")
		return
	}

	return rest.Limits{
		Hosts:           int32(safeInt(limitsMap[pHosts])),
		Volumes:         int32(safeInt(limitsMap[pVolumes])),
		VolumeCapacity:  int64(safeFloat(limitsMap[pVolumeCapacity])),
		PrivateNetworks: int32(safeInt(limitsMap[pPrivateNetworks])),
		InstanceTypes:   safeMapStrInt32(limitsMap[pInstanceTypes]),
	}, nil
}

func resourceMetalProjectRead(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to read project")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	ctx = context.WithValue(ctx, rest.ContextAPIKey, rest.APIKey{Key: d.Id()})
	project, _, err := p.Client.ProjectsApi.GetByID(ctx, d.Id(), nil)
	if err != nil {
		return err
	}
	d.Set(pName, project.Name)

	prof := project.Profile
	pData := map[string]interface{}{
		pAddress:            prof.Address,
		pCompany:            prof.Company,
		pEmail:              prof.Email,
		pEmailVerified:      prof.EmailVerified,
		pPhoneNumber:        prof.PhoneNumber,
		pPhoneVerified:      prof.PhoneVerified,
		pProjectDescription: prof.TeamDesc,
		pProjectName:        prof.TeamName,
	}

	if err = d.Set(pProfile, []interface{}{pData}); err != nil {
		return err
	}

	lim := project.Limits
	lData := map[string]interface{}{
		pHosts:           int(lim.Hosts),
		pVolumes:         int(lim.Volumes),
		pVolumeCapacity:  float64(lim.VolumeCapacity),
		pPrivateNetworks: int(lim.PrivateNetworks),
		pInstanceTypes:   lim.InstanceTypes,
	}

	if err = d.Set(pLimits, []interface{}{lData}); err != nil {
		return err
	}

	if len(project.PermittedSites) > 0 {
		sites := flattenStringList(project.PermittedSites)
		if err = d.Set(pSites, schema.NewSet(schema.HashString, sites)); err != nil {
			return err // nolint:wrapcheck // defer func is wrapping the error.
		}
	}

	if len(project.PermittedOSImages) > 0 {
		images := flattenStringList(project.PermittedOSImages)
		if err = d.Set(pPermittedImages, schema.NewSet(schema.HashString, images)); err != nil {
			return err //nolint:wrapcheck // defer func is wrapping the error.
		}
	}

	if err = d.Set(pVolumeReplicationEnabled, project.VolumeReplicationEnabled); err != nil {
		return err //nolint:wrapcheck // defer func is wrapping the error.
	}

	if err = d.Set(pBootFromSANSupport, project.BootFromSANSupport); err != nil {
		return err //nolint:wrapcheck // defer func is wrapping the error.
	}

	return nil
}

func resourceMetalProjectUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to update project")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return
	}

	ctx := p.GetContext()
	ctx = context.WithValue(ctx, rest.ContextAPIKey, rest.APIKey{Key: d.Id()})
	project, _, err := p.Client.ProjectsApi.GetByID(ctx, d.Id(), nil)
	if err != nil {
		return
	}

	name, ok := d.Get(pName).(string)
	if !ok {
		return fmt.Errorf("name is not in the expected format")
	}

	updateProject := rest.UpdateProject{
		ID:   project.ID,
		ETag: project.ETag,
	}

	updateProject.Name = name

	if list, ok := d.Get(pProfile).([]interface{}); ok && len(list) == 1 {
		if updateProject.Profile, err = getUpdateProfile(list[0]); err != nil {
			return
		}
	} else {
		return fmt.Errorf("only 1 profile block is allowed")
	}

	if list, ok := d.Get(pLimits).([]interface{}); ok && len(list) == 1 {
		if updateProject.Limits, err = getUpdateLimits(list[0]); err != nil {
			return
		}
	} else {
		return fmt.Errorf("only 1 limit block is allowed")
	}

	if f, ok := d.GetOk(pSites); ok {
		s, ok := f.(*schema.Set)
		if !ok {
			err = fmt.Errorf("sites list is not in the expected format")

			return err
		}

		updateProject.PermittedSites = expandStringList(s.List())
	}

	if f, ok := d.GetOk(pPermittedImages); ok {
		s, ok := f.(*schema.Set)
		if !ok {
			err = fmt.Errorf("permitted images list is not in the expected format")

			return err
		}

		updateProject.PermittedOSImages = expandStringList(s.List())
	}

	_, _, err = p.Client.ProjectsApi.Update(ctx, updateProject.ID, updateProject, nil)
	if err != nil {
		return
	}

	return resourceMetalProjectRead(d, meta)
}

func resourceMetalProjectDelete(d *schema.ResourceData, meta interface{}) (err error) {
	defer wrapResourceError(&err, "failed to delete project")

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	ctx = context.WithValue(ctx, rest.ContextAPIKey, rest.APIKey{Key: d.Id()})
	_, err = p.Client.ProjectsApi.Delete(ctx, d.Id(), nil)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
