// Copyright (c) 2016-2020 Hewlett Packard Enterprise Development LP.

package quake

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	rest "github.com/quattronetworks/quake-client/v1/go-client"
)

const (
	// Quake defines the name used by terraform to reference this provider.
	Quake = "quake"

	pollInterval = 3 * time.Second
)

const (
	qPortal  = "portal_url"
	qProject = Quake + "_project"
	qHost    = Quake + "_host"
	qVolume  = Quake + "_volume"
	qSSHKey  = Quake + "_ssh_key"
	qNetwork = Quake + "_network"

	// For data sources
	qAvailableResource = Quake + "_available_resources"
	qAvailableImages   = Quake + "_available_images"
	qUsage             = Quake + "_usage"
)

var (
	resourceDefaultTimeouts *schema.ResourceTimeout
)

// Config holds all the information required to talk to the portal.
type Config struct {
	client           *rest.APIClient
	terraformVersion string
	context          context.Context
	// we will cache this here for the life of the provider
	availableResources rest.AvailableResources
}

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
func Provider() terraform.ResourceProvider {
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
		},
		ResourcesMap: map[string]*schema.Resource{
			qHost:   hostResource(),
			qVolume: volumeResource(),
			//qVolumeAttach: volumeAttachmentResource(),
			qSSHKey:  sshKeyResource(),
			qProject: projectResource(),
			qNetwork: projectNetworkResource(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			qAvailableResource: dataSourceAvailableResources(),
			qAvailableImages:   dataSourceImage(),
			qUsage:             dataSourceUsage(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		var (
			err  error
			qjwt *Qjwt
		)
		// Look in the home directory first the .qjwt file
		homeDir, _ := os.UserHomeDir()
		workingDir, _ := os.Getwd()
		for _, p := range []string{homeDir, workingDir} {
			qjwt, err = loadConfig(p)
			if err == nil {
				break
			}
		}
		if err != nil || qjwt == nil {
			return nil, err
		}

		cfg := rest.NewConfiguration()
		// @TODO this 'rest/vi' and 'Membership' is all black-magic and needs something better.
		cfg.BasePath = qjwt.RestURL + "/rest/v1"
		cfg.AddDefaultHeader("Membership", qjwt.MemberID)
		client := rest.NewAPIClient(cfg)
		ctx := context.Background()
		ctx = context.WithValue(ctx, rest.ContextAccessToken, qjwt.Token)
		// sanity check some attributes here
		pURL := d.Get(qPortal).(string)
		if pURL != "" {
			if pURL != qjwt.OriginalURL {
				return nil, fmt.Errorf("Provider explicitly states portal is %q yet token is valid for %q", pURL, qjwt.OriginalURL)
			}
		}
		d.Set(qPortal, qjwt.OriginalURL)
		//d.Set(qUser, loginInfo.User)
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		c := &Config{
			client:           client,
			terraformVersion: terraformVersion,
			context:          ctx,
		}
		if err = c.refreshAvailableResources(); err != nil {
			return nil, err
		}
		return c, nil
	}
	return provider
}

func (c *Config) refreshAvailableResources() error {
	resources, _, err := c.client.AvailableResourcesApi.List(c.context)
	if err != nil {
		return err
	}
	c.availableResources = resources
	return nil
}

func (c *Config) getLocationName(locationID string) (string, error) {
	for _, loc := range c.availableResources.Locations {
		if loc.ID == locationID {
			return fmt.Sprintf("%s:%s:%s", loc.Country, loc.Region, loc.DataCenter), nil
		}
	}
	return "", fmt.Errorf("LocationID %s not found", locationID)
}

func (c *Config) getLocationID(locationName string) (locationID string, err error) {
	locations := []string{}
	pieces := strings.Split(locationName, ":")

	for _, loc := range c.availableResources.Locations {
		if len(pieces) == 3 {
			if string(loc.Country) == pieces[0] && loc.Region == pieces[1] && loc.DataCenter == pieces[2] {
				return loc.ID, nil
			}
		}
		locations = append(locations, fmt.Sprintf("%s:%s:%s", loc.Country, loc.Region, loc.DataCenter))
	}
	return "", fmt.Errorf("location %q not found in %q", locationName, locations)
}

func (c *Config) getVolumeFlavorName(flavorID string) (string, error) {
	for _, vf := range c.availableResources.VolumeFlavors {
		if flavorID == vf.ID {
			return vf.Name, nil
		}
	}
	return "", fmt.Errorf("VolumeFalvorID %s not found", flavorID)
}
