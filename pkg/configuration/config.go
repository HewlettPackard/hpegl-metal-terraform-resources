// (C) Copyright 2020-2023 Hewlett Packard Enterprise Development LP

package configuration

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/constants"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/gltform"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/retrieve"
)

// KeyForGLClientMap is used by the GL terraform provider to set the key in the
// map of clients that it creates.  The Metal client will be found using the key
// returned here.
func KeyForGLClientMap() string {
	return constants.MetalClientMapKey
}

// Config holds all the information required to talk to the portal.
type Config struct {
	restURL    string
	token      string
	user       string
	space      string
	trf        retrieve.TokenRetrieveFuncCtx
	useGLToken bool
	context    context.Context
	// Exported fields
	PortalURL string
	Client    *rest.APIClient
	// we will cache this here for the life of the provider
	AvailableResources rest.AvailableResources
}

type CreateOpt func(c *Config)

func WithGLToken(g bool) CreateOpt {
	return func(c *Config) {
		c.useGLToken = g
	}
}

// WithTRF this create option is for use by the hpegl terraform provider
// It is used to pass-in a token retrieve function which is used to get
// a GL IAM token.  Behind the scenes tokens are generated and refreshed
// if they are due to expire
func WithTRF(trf retrieve.TokenRetrieveFuncCtx) CreateOpt {
	return func(c *Config) {
		c.trf = trf
	}
}

func (c *Config) RefreshAvailableResources() error {
	ctx := c.GetContext()

	resources, _, err := c.Client.AvailableResourcesApi.List(ctx)
	if err != nil {
		return err
	}

	c.AvailableResources = resources
	return nil
}

func (c *Config) GetLocationName(locationID string) (string, error) {
	for _, loc := range c.AvailableResources.Locations {
		if loc.ID == locationID {
			return makeLocationName(string(loc.Country), loc.Region, loc.DataCenter), nil
		}
	}

	return "", fmt.Errorf("LocationID %s not found", locationID)
}

func (c *Config) GetLocationID(locationName string) (locationID string, err error) {
	locations := []string{}
	pieces := strings.Split(locationName, ":")

	for _, loc := range c.AvailableResources.Locations {
		if len(pieces) == 3 {
			if string(loc.Country) == pieces[0] && loc.Region == pieces[1] && loc.DataCenter == pieces[2] {
				return loc.ID, nil
			}
		}
		locations = append(locations, makeLocationName(string(loc.Country), loc.Region, loc.DataCenter))
	}
	return "", fmt.Errorf("location %q not found in %q", locationName, locations)
}

func (c *Config) GetVolumeFlavorName(flavorID string) (string, error) {
	for _, vf := range c.AvailableResources.VolumeFlavors {
		if flavorID == vf.ID {
			return vf.Name, nil
		}
	}
	return "", fmt.Errorf("VolumeFalvorID %s not found", flavorID)
}

func (c *Config) GetStoragePoolName(storagePoolID string) (string, error) {
	for _, sp := range c.AvailableResources.StoragePools {
		if storagePoolID == sp.ID {
			return sp.Name, nil
		}
	}

	return "", fmt.Errorf("StoragePoolID %s not found", storagePoolID)
}

func (c *Config) GetStoragePoolID(storagePoolName string) (string, error) {
	for _, sp := range c.AvailableResources.StoragePools {
		if storagePoolName == sp.Name {
			return sp.ID, nil
		}
	}

	return "", fmt.Errorf("StoragePoolName %s not found", storagePoolName)
}

// GetContext is used to retrieve the context.
// If the token retrieve function is nil the context in Config is returned
// If there is a token retrieve function it is executed to retrieve a GL IAM token, which is
// placed in the context before it is returned.
// If we get an error we log it and return the context with the new token.
func (c *Config) GetContext() context.Context {
	if c.trf == nil {
		return c.context
	}

	token, err := c.trf(c.context)
	if err != nil {
		log.Printf("error in retrieving token %s", err)
	}

	return context.WithValue(c.context, rest.ContextAccessToken, token)
}

func makeLocationName(country, region, dataCenter string) string {
	return fmt.Sprintf("%s:%s:%s", country, region, dataCenter)
}

func NewConfig(portalURL string, opts ...CreateOpt) (*Config, error) {
	// create REST Client Context
	ctx := context.Background()
	config := new(Config)

	config.useGLToken = false
	config.trf = nil

	// run overrides
	for _, opt := range opts {
		if opt != nil {
			opt(config)
		}
	}

	if config.useGLToken || config.trf != nil {
		// Use GetGLConfig from gltform
		glconfig, err := gltform.GetGLConfig()
		if err != nil {
			return nil, fmt.Errorf("error reading GL token file:  %w", err)
		}

		config.restURL = glconfig.RestURL
		config.token = glconfig.Token
		config.user = glconfig.ProjectID
		config.space = glconfig.SpaceName
	} else {
		qtoken, err := getQConfig()
		if err != nil {
			return nil, fmt.Errorf("error reading Metal token file:  %w", err)
		}

		if portalURL != "" && portalURL != qtoken.OriginalURL {
			return nil, fmt.Errorf("Provider explicitly states portal is %q yet token is valid for %q", portalURL, qtoken.OriginalURL)
		}
		config.restURL = qtoken.RestURL
		config.token = qtoken.Token
		config.user = qtoken.MemberID
		config.PortalURL = qtoken.OriginalURL
	}

	// add access token for auth to Client Context as required by the Client API
	ctx = context.WithValue(ctx, rest.ContextAccessToken, config.token)

	// Get a new Client configuration with basepath set to Metal portal URL and add base version path /rest/v1
	cfg := rest.NewConfiguration()

	basePath, err := url.JoinPath(config.restURL, "/rest/v1")
	if err != nil {
		return nil, fmt.Errorf("configuration error: %v", err)
	}

	cfg.BasePath = basePath

	if config.useGLToken || config.trf != nil {
		if err := validateGLConfig(*config); err != nil {
			return config, fmt.Errorf("configuration error: %v", err)
		}

		// Add required headers if GL authentication method
		if config.user != "" {
			cfg.AddDefaultHeader("Project", config.user)
		}

		if config.space != "" {
			cfg.AddDefaultHeader("Space", config.space)
		}
	} else {
		if err := validateMetalConfig(*config); err != nil {
			return config, fmt.Errorf("configuration error: %v", err)
		}

		// Add membership field to header if Q authentication method
		cfg.AddDefaultHeader("Membership", config.user)
	}

	// get new API Client with basepath and auth credentials setup in configuration and Context
	config.context = ctx
	config.Client = rest.NewAPIClient(cfg)
	return config, nil
}

func validateMetalConfig(config Config) error {
	if config.restURL == "" {
		return fmt.Errorf("rest_url is not set")
	}

	if config.user == "" {
		return fmt.Errorf("member_id is not set")
	}

	if config.token == "" {
		return fmt.Errorf(("jwt is not set"))
	}

	return nil
}

func validateGLConfig(config Config) error {
	if config.restURL == "" {
		return fmt.Errorf("rest_url is not set")
	}

	// when token retrieval function is registered, then token isn't required as config param, which
	// is the case when initialized from hpegl.
	if config.trf == nil && config.token == "" {
		return fmt.Errorf("access_token is not set")
	}

	return nil
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

// IsHosterContext determines whether the provider configuration
// is project scope or hoster scope when GL IAM token is used.
// Project operations with Metal token not supported yet via Terraform, so
// 'false' is returned in this case.
// The scope determines what Metal APIs are allowed. Project resource create/delete/update
// requires Hoster scope configuration and, the remaining resources require Project scope
// configuration Project.
func (c *Config) IsHosterContext() bool {
	if c.useGLToken || c.trf != nil {
		// if project_id is not set, then it is Hoster scope
		if c.user == "" {
			return true
		}
	}

	return false
}

// GetVolumeCollectionID returns volume collection ID from volume collection name.
func (c *Config) GetVolumeCollectionID(vcolName string) (string, error) {
	for _, vc := range c.AvailableResources.VolumeCollections {
		if vcolName == vc.Name {
			return vc.ID, nil
		}
	}

	return "", fmt.Errorf("volume collection  %s not found", vcolName)
}

// GetVolumeCollectionName return the ID of the volume collection
func (c *Config) GetVolumeCollectionName(vcolID string) (string, error) {
	if vcolID == "" {
		return "", nil
	}
	for _, vc := range c.AvailableResources.VolumeCollections {
		if vcolID == vc.ID {
			return vc.Name, nil
		}
	}

	return "", fmt.Errorf("volume collection %s not found", vcolID)
}
