// (C) Copyright 2022, 2024 Hewlett Packard Enterprise Development LP

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name hpegl

package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	testutils "github.com/hewlettpackard/hpegl-metal-terraform-resources/internal/test-utils"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: testutils.ProviderFunc(),
	})
}
