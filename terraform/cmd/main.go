// Package main is the executable for the Quake terraform plugin.
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/hpe-hcss/quake-client/terraform/quake"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: quake.Provider,
	})
}
