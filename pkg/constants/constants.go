// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

// Package constants - constants that are used in pkg/client and pkg/resources
package constants

const (
	// ServiceName - the service mnemonic
	ServiceName = "metal"

	// MetalClientMapKey is the key in the map[string]interface{} that is passed down by hpegl used to store *Client
	// This must be unique, hpegl will error-out if it isn't
	MetalClientMapKey = "metalConfig"
)
