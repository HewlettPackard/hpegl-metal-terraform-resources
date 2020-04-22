// Package main is the executable for the Quake terraform plugin.
package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/quattronetworks/quake-client/terraform/quake"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: quake.Provider,
	})
}
