// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package testutils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/HewlettPackard/hpegl-metal-terraform-resources/pkg/client"
	"github.com/HewlettPackard/hpegl-metal-terraform-resources/pkg/registration"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/provider"
)

func ProviderFunc() plugin.ProviderFunc {
	return provider.NewProviderFunc(provider.ServiceRegistrationSlice(registration.Registration{}), providerConfigure)
}

func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc { // nolint staticcheck
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		cli, err := client.InitialiseClient{}.NewClient(d)
		if err != nil {
			return nil, diag.Errorf("error in creating client: %s", err)
		}

		return map[string]interface{}{
			client.InitialiseClient{}.ServiceName(): cli,
		}, nil
	}
}
