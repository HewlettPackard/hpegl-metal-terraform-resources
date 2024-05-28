// (C) Copyright 2020-2023 Hewlett Packard Enterprise Development LP

package registration

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hewlettpackard/hpegl-metal-terraform-resources/internal/resources"
)

const (
	mPrefix = "hpegl_metal"

	qProject = mPrefix + "_project"
	qHost    = mPrefix + "_host"
	qVolume  = mPrefix + "_volume"
	qSSHKey  = mPrefix + "_ssh_key"
	qNetwork = mPrefix + "_network"
	qIP      = mPrefix + "_ip"
	qImage   = mPrefix + "_image"

	qAvailableResource = mPrefix + "_available_resources"
	qAvailableImages   = mPrefix + "_available_images"

	// These constants are used to set the optional hpegl provider "metal" block field-names
	projectID    = "project_id"
	restURL      = "rest_url"
	spaceName    = "space_name"
	glToken      = "gl_token"
	glpRole      = "glp_role"
	glpWorkspace = "glp_workspace"
)

type Registration struct{}

func (r Registration) Name() string {
	return "metal"
}

func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		qAvailableResource: resources.DataSourceAvailableResources(),
		qAvailableImages:   resources.DataSourceImage(),
	}
}

func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		qHost:   resources.HostResource(),
		qVolume: resources.VolumeResource(),
		//qVolumeAttach: volumeAttachmentResource(),
		qSSHKey:  resources.SshKeyResource(),
		qProject: resources.ProjectResource(),
		qNetwork: resources.ProjectNetworkResource(),
		qIP:      resources.IPResource(),
		qImage:   resources.ServiceImageResource(),
	}
}

func (r Registration) ProviderSchemaEntry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			projectID: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_METAL_PROJECT_ID", ""),
				Description: "The Metal project-id to use, can also be set with the HPEGL_METAL_PROJECT_ID env-var",
			},
			restURL: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_METAL_REST_URL", ""),
				Description: "The Metal portal rest-url to use, can also be set with the HPEGL_METAL_REST_URL env-var",
			},
			spaceName: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_METAL_SPACE_NAME", ""),
				Description: `The space-name to use with Metal, only required for project creation operations,
				can also be set with the HPEGL_METAL_SPACE_NAME env-var`,
			},
			glToken: {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_METAL_GL_TOKEN", true),
				Description: `Field indicating whether the token is GreenLake (GLCS or GLP) IAM issued token or Metal Service issued one,
				can also be set with the HPEGL_METAL_GL_TOKEN env-var`,
			},
			glpRole: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_METAL_GLP_ROLE", true),
				Description: `Field indicating the GLP role to be used, can also be set with the HPEGL_METAL_GLP_ROLE env-var`,
			},
			glpWorkspace: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_METAL_GLP_WORKSPACE", true),
				Description: `Field indicating the GLP workspace to be used, can also be set with the HPEGL_METAL_GLP_WORKSPACE env-var`,
			},
		},
	}
}
