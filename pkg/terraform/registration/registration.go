// (C) Copyright 2020-2021 Hewlett Packard Enterprise Development LP

package registration

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hpe-hcss/quake-client/terraform/quake"
)

const (
	Quake = "hpegl_bmaas"

	qProject = Quake + "_project"
	qHost    = Quake + "_host"
	qVolume  = Quake + "_volume"
	qSSHKey  = Quake + "_ssh_key"
	qNetwork = Quake + "_network"
	qIP      = Quake + "_ip"

	qAvailableResource = Quake + "_available_resources"
	qAvailableImages   = Quake + "_available_images"
	qUsage             = Quake + "_usage"

	// These constants are used to set the optional hpegl provider "bmaas" block field-names
	projectID = "project_id"
	restURL   = "rest_url"
	spaceName = "space_name"
)

type Registration struct{}

func (r Registration) Name() string {
	return "bmaas"
}

func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		qAvailableResource: quake.DataSourceAvailableResources(),
		qAvailableImages:   quake.DataSourceImage(),
		qUsage:             quake.DataSourceUsage(),
	}
}

func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		qHost:   quake.HostResource(),
		qVolume: quake.VolumeResource(),
		//qVolumeAttach: volumeAttachmentResource(),
		qSSHKey:  quake.SshKeyResource(),
		qProject: quake.ProjectResource(),
		qNetwork: quake.ProjectNetworkResource(),
		qIP:      quake.IPResource(),
	}
}

func (r Registration) ProviderSchemaEntry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			projectID: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_BMAAS_PROJECT_ID", ""),
				Description: "The BMaaS project-id to use, can also be set with the HPEGL_BMAAS_PROJECT_ID env-var",
			},
			restURL: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_BMAAS_REST_URL", ""),
				Description: "The BMaaS portal rest-url to use, can also be set with the HPEGL_BMAAS_REST_URL env-var",
			},
			spaceName: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_BMAAS_SPACE_NAME", ""),
				Description: "The space-name to use with BMaaS, only required for project creation operations, can also be set with the HPEGL_BMAAS_SPACE_NAME env-var",
			},
		},
	}
}
