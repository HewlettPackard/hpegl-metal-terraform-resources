// (C) Copyright 2020 Hewlett Packard Enterprise Development LP

package registration

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/quattronetworks/quake-client/terraform/quake"
)

const (
	Quake = "hpegl_quake"

	qProject = Quake + "_project"
	qHost    = Quake + "_host"
	qVolume  = Quake + "_volume"
	qSSHKey  = Quake + "_ssh_key"
	qNetwork = Quake + "_network"

	qAvailableResource = Quake + "_available_resources"
	qAvailableImages   = Quake + "_available_images"
	qUsage             = Quake + "_usage"

)

type Registration struct{}

func (r Registration) Name() string {
	return "Quake Service"
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
	}
}