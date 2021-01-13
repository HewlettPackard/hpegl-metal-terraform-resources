// Copyright (c) 2016-2021 Hewlett Packard Enterprise Development LP.

package quake

import (
	"context"
	"fmt"
	"os"
	"strings"

	rest "github.com/quattronetworks/quake-client/v1/pkg/client"
)

// Config holds all the information required to talk to the portal.
type Config struct {
	PortalURL        string
	restURL          string
	token            string
	user             string
	space            string
	client           *rest.APIClient
	terraformVersion string
	context          context.Context
	useGLToken       bool
	// we will cache this here for the life of the provider
	availableResources rest.AvailableResources
}

type CreateOpt func(c *Config)

func WithGLToken(g bool) CreateOpt {
	return func(c *Config) {
		c.useGLToken = g
	}
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
			return makeLocationName(string(loc.Country), loc.Region, loc.DataCenter), nil
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
		locations = append(locations, makeLocationName(string(loc.Country), loc.Region, loc.DataCenter))
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

func makeLocationName(country, region, dataCenter string) string {
	return fmt.Sprintf("%s:%s:%s", country, region, dataCenter)
}

func NewConfig(portalURL string, opts ...CreateOpt) (*Config, error) {
	// create REST client context
	ctx := context.Background()
	config := new(Config)

	config.useGLToken = false

	// run overrides
	for _, opt := range opts {
		if opt != nil {
			opt(config)
		}
	}

	if config.useGLToken {
		gltoken, err := getGLConfig()
		if err != nil {
			return nil, fmt.Errorf("Error reading GL token file:  %w", err)
		}
		config.restURL = gltoken.RestURL
		config.token = gltoken.Token
		config.user = gltoken.ProjectID
		config.space = gltoken.SpaceName
	} else {
		qtoken, err := getQConfig()
		if err != nil {
			return nil, fmt.Errorf("Error reading Q token file:  %w", err)
		}
		if portalURL != "" && portalURL != qtoken.OriginalURL {
			return nil, fmt.Errorf("Provider explicitly states portal is %q yet token is valid for %q", portalURL, qtoken.OriginalURL)
		}
		config.restURL = qtoken.RestURL
		config.token = qtoken.Token
		config.user = qtoken.MemberID
		config.PortalURL = qtoken.OriginalURL
	}

	// add access token for auth to client context as required by the client API
	ctx = context.WithValue(ctx, rest.ContextAccessToken, config.token)

	// Get a new client configuration with basepath set to Quake portal URL and add base version path /rest/v1
	cfg := rest.NewConfiguration()
	cfg.BasePath = config.restURL + "/rest/v1"

	if config.useGLToken {
		// Add required headers if GL authentication method
		if config.user != "" {
			cfg.AddDefaultHeader("Project", config.user)
		}
		if config.space != "" {
			cfg.AddDefaultHeader("Space", config.space)
		}
	} else {
		// Add membership field to header if Q authentication method
		if config.user == "" {
			return config, fmt.Errorf("no valid memberid found for Q access token")
		}
		cfg.AddDefaultHeader("Membership", config.user)
	}

	// get new API client with basepath and auth credentials setup in configuration and context
	config.context = ctx
	config.client = rest.NewAPIClient(cfg)
	return config, nil
}

func getGLConfig() (gljwt *Gljwt, err error) {
	homeDir, _ := os.UserHomeDir()
	workingDir, _ := os.Getwd()
	for _, p := range []string{homeDir, workingDir} {
		gljwt, err = loadGLConfig(p)
		if err == nil {
			break
		}
	}
	return gljwt, err
}

func getQConfig() (qjwt *Qjwt, err error) {
	homeDir, _ := os.UserHomeDir()
	workingDir, _ := os.Getwd()
	for _, p := range []string{homeDir, workingDir} {
		qjwt, err = loadConfig(p)
		if err == nil {
			break
		}
	}
	return qjwt, err
}
