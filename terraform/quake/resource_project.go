// Copyright (c) 2016-2020 Hewlett Packard Enterprise Development LP.

package quake

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rest "github.com/quattronetworks/quake-client/v1/pkg/client"
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

	pHosts          = "hosts"
	pVolumes        = "volumes"
	pVolumeCapacity = "volume_capacity"
)

func limitsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		pHosts: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Maximum number of host allowed in the team.",
		},
		pVolumes: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Maximum number of volumes allowed in the team.",
		},
		pVolumeCapacity: {
			Type:        schema.TypeFloat,
			Computed:    true,
			Description: "Total allowable volume capacity (GiB) allowed in the team.",
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
			// TODO the V2 SDK doesn't (yet) support TypeMap with Elem *Resource
			// This is the currently recommended work-around. See
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/155
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/616
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Team profile.",
			Elem: &schema.Resource{
				Schema: profileSchema(),
			},
		},
		pLimits: {
			// TODO the V2 SDK doesn't (yet) support TypeMap with Elem *Resource
			// This is the currently recommended work-around. See
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/155
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/616
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Resource limits applied to this team.",
			Elem: &schema.Resource{
				Schema: limitsSchema(),
			},
		},
	}
}

func projectResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceQuattroProjectCreate,
		Read:   resourceQuattroProjectRead,
		Delete: resourceQuattroProjectDelete,
		Update: resourceQuattroProjectUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: projectSchema(),
	}
}

func resourceQuattroProjectCreate(d *schema.ResourceData, meta interface{}) (err error) {
	p := meta.(*Config)

	safeString := func(s interface{}) string {
		r, _ := s.(string)
		return r
	}
	safeInt := func(s interface{}) int {
		r, _ := s.(int)
		return r
	}
	safeFloat := func(s interface{}) float64 {
		r, _ := s.(float64)
		return r
	}
	profile := d.Get(pProfile).(map[string]interface{})
	limits := d.Get(pLimits).(map[string]interface{})

	np := rest.NewProject{
		Name: d.Get(pName).(string),
	}
	if profile != nil {
		np.Profile = rest.Profile{
			Address:     safeString(profile[pAddress]),
			Company:     safeString(profile[pCompany]),
			Email:       safeString(profile[pEmail]),
			PhoneNumber: safeString(profile[pPhoneNumber]),
			TeamDesc:    safeString(profile[pProjectDescription]),
			TeamName:    safeString(profile[pProjectName]),
		}
	}

	if limits != nil {
		np.Limits = rest.Limits{
			Hosts:          uint32(safeInt(limits[pHosts])),
			Volumes:        uint32(safeInt(limits[pVolumes])),
			VolumeCapacity: uint64(safeFloat(limits[pVolumeCapacity])),
		}
	}
	project, _, err := p.client.ProjectsApi.Add(p.context, np)
	if err != nil {
		return err
	}
	d.SetId(project.ID)
	return resourceQuattroProjectRead(d, meta)
}

func resourceQuattroProjectRead(d *schema.ResourceData, meta interface{}) (err error) {
	p := meta.(*Config)
	project, _, err := p.client.ProjectsApi.GetByID(p.context, d.Id())
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

	if err = d.Set(pProfile, pData); err != nil {
		return err
	}

	lim := project.Limits
	lData := map[string]interface{}{
		pHosts:          int(lim.Hosts),
		pVolumes:        int(lim.Volumes),
		pVolumeCapacity: float64(lim.VolumeCapacity),
	}

	if err = d.Set(pLimits, lData); err != nil {
		return err
	}
	return nil
}

func resourceQuattroProjectUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	//p := meta.(*Config)
	return resourceQuattroProjectRead(d, meta)
}

func resourceQuattroProjectDelete(d *schema.ResourceData, meta interface{}) (err error) {
	p := meta.(*Config)
	_, err = p.client.ProjectsApi.Delete(p.context, d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
