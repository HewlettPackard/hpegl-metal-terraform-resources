// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"

	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
)

const (
	pName    = "name"
	pProfile = "profile"
	pLimits  = "limits"

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
	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to create project %s: %w", strings.Trim(nErr.Message(), "\n "), err)

		}
	}()

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	np := rest.NewProject{
		Name: d.Get(pName).(string),
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
	}, nil
}

func resourceMetalProjectRead(d *schema.ResourceData, meta interface{}) (err error) {
	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to read project %s: %w", strings.Trim(nErr.Message(), "\n "), err)

		}
	}()

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	project, _, err := p.Client.ProjectsApi.GetByID(ctx, d.Id())
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
	}

	if err = d.Set(pLimits, []interface{}{lData}); err != nil {
		return err
	}
	return nil
}

func resourceMetalProjectUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to update project %s: %w", strings.Trim(nErr.Message(), "\n "), err)

		}
	}()

	//p := meta.(*Config)
	return resourceMetalProjectRead(d, meta)
}

func resourceMetalProjectDelete(d *schema.ResourceData, meta interface{}) (err error) {
	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to delete project %s: %w", strings.Trim(nErr.Message(), "\n "), err)

		}
	}()

	p, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return err
	}

	ctx := p.GetContext()
	_, err = p.Client.ProjectsApi.Delete(ctx, d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
