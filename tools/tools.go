//go:build tools
// +build tools

// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package tools

import (
	// document generation
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
