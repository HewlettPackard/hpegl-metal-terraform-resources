// (C) Copyright 2016-2022 Hewlett Packard Enterprise Development LP

package quake

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hpe-hcss/quake-client/pkg/terraform/configuration"
	rest "github.com/hpe-hcss/quake-client/v1/pkg/client"
)

const (
	// Quake defines the name used by terraform to reference this provider.
	Quake = "quake"

	pollInterval = 3 * time.Second
)

const (
	qPortal     = "portal_url"
	qUseGLToken = "gl_token"
	qProject    = Quake + "_project"
	qHost       = Quake + "_host"
	qVolume     = Quake + "_volume"
	qSSHKey     = Quake + "_ssh_key"
	qNetwork    = Quake + "_network"
	qIP         = Quake + "_ip"

	// For data sources
	qAvailableResource = Quake + "_available_resources"
	qAvailableImages   = Quake + "_available_images"
	qUsage             = Quake + "_usage"
)

var (
	resourceDefaultTimeouts *schema.ResourceTimeout
)

func init() {
	d := time.Minute * 60
	resourceDefaultTimeouts = &schema.ResourceTimeout{
		Create:  schema.DefaultTimeout(d),
		Update:  schema.DefaultTimeout(d),
		Delete:  schema.DefaultTimeout(d),
		Default: schema.DefaultTimeout(d),
	}
}

// Provider returns the QuattroLabs terrform rovider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			qProject: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ProjectID",
			},
			qPortal: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Fully qualified URL to the portal",
			},
			qUseGLToken: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Toggle use of GL tokens (bool), default is false",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			qHost:   HostResource(),
			qVolume: VolumeResource(),
			//qVolumeAttach: volumeAttachmentResource(),
			qSSHKey:  SshKeyResource(),
			qProject: ProjectResource(),
			qNetwork: ProjectNetworkResource(),
			qIP:      IPResource(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			qAvailableResource: DataSourceAvailableResources(),
			qAvailableImages:   DataSourceImage(),
			qUsage:             DataSourceUsage(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		var err error

		defer func() {
			var nErr = rest.GenericOpenAPIError{}
			if errors.As(err, &nErr) {
				err = fmt.Errorf("failed to configure provider %s: %w", strings.Trim(nErr.Message(), "\n "), err)
			}
		}()

		config, err := configuration.NewConfig(d.Get(qPortal).(string), configuration.WithGLToken(d.Get(qUseGLToken).(bool)))
		if err != nil {
			return nil, err
		}

		d.Set(qPortal, config.PortalURL)
		//d.Set(qUser, loginInfo.User)
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		// Cache Available Resources in a project when the scope is Project level
		if !config.IsHosterContext() {
			if err = config.RefreshAvailableResources(); err != nil {
				return nil, err
			}
		}

		return config, nil
	}
	return provider
}
