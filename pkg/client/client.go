// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/configuration"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/constants"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/client"
)

// Assert that InitialiseClient satisfies the client.Initialisation interface
var _ client.Initialisation = (*InitialiseClient)(nil)

// InitialiseClient is imported by hpegl from each service repo
type InitialiseClient struct{}

// NewClient takes an argument of all of the provider.ConfigData, and returns an interface{} and error
// If there is no error interface{} will contain *Client.
// The hpegl provider will put *Client at the value of keyForGLClientMap (returned by ServiceName) in
// the map of clients that it creates and passes down to provider code.  hpegl executes NewClient for each service.
func (i InitialiseClient) NewClient(r *schema.ResourceData) (interface{}, error) {
	var err error

	defer func() {
		var nErr = rest.GenericOpenAPIError{}
		if errors.As(err, &nErr) {
			err = fmt.Errorf("failed to configure provider %s: %w", strings.Trim(nErr.Message(), "\n "), err)
		}
	}()

	// Get metal settings from the service block
	// If an error is returned it means that the service block doesn't exist
	// In this case, return nil since we can't tell if the lack of a service block is a real
	// error or is intentional - i.e. if the user has forgotten to add the block or if the
	// user doesn't intend to use terraform to run against this service.  GetClientFromMetaMap
	// below will be used to handle the situation where Client is nil.
	metalMap, err := client.GetServiceSettingsMap(constants.ServiceName, r)
	if err != nil {
		return nil, nil
	}

	// Initialize the metal client
	metalConfig, err := configuration.NewConfig("",
		configuration.WithGLToken(metalMap["gl_token"].(bool)),
		configuration.WithRole(metalMap["glp_role"].(string)),
		configuration.WithWorkspace(metalMap["glp_workspace"].(string)),
	)
	if err != nil {
		return nil, fmt.Errorf("error in creating metal client: %s", err)
	}

	// Cache Available Resources in a project when the scope is Project level
	if !metalConfig.IsHosterContext() {
		if err := metalConfig.RefreshAvailableResources(); err != nil {
			return nil, fmt.Errorf("error in refreshing available resources for metal: %s", err)
		}
	}

	return metalConfig, nil
}

// ServiceName is used to return the value of MetalClientMapKey, for use by hpegl
func (i InitialiseClient) ServiceName() string {
	return constants.MetalClientMapKey
}

// GetClientFromMetaMap is a convenience function used by provider code to extract *Client from the
// meta argument passed-in by terraform
func GetClientFromMetaMap(meta interface{}) (*configuration.Config, error) {
	cli := meta.(map[string]interface{})[constants.MetalClientMapKey]
	if cli == nil {
		return nil, fmt.Errorf(
			"client is not initialised, make sure that %s block is defined in hpegl provider stanza", constants.ServiceName)
	}

	return cli.(*configuration.Config), nil
}
