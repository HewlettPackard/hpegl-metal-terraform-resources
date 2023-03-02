// (C) Copyright 2023 Hewlett Packard Enterprise Development LP

package resources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/configuration"
)

func Test_addNetworks(t *testing.T) {
	testVlan := 200
	testVni := 12006

	availableNets := client.AvailableResources{
		Networks: []client.AvailableNetwork{
			{
				ID:   "testID",
				VLAN: int32(testVlan),
				VNI:  int32(testVni),
			},
		},
	}

	cfg := &configuration.Config{
		AvailableResources: client.AvailableResources{},
	}

	d := schema.TestResourceDataRaw(t, DataSourceAvailableResources().Schema, map[string]interface{}{})

	// test
	err := addNetworks(cfg, d, availableNets)
	assert.Nil(t, err)

	networks, ok := d.Get(avNetworks).([]interface{})
	assert.True(t, ok, "type assertion failed")
	assert.Equal(t, 1, len(networks))

	net := networks[0].(map[string]interface{})

	assert.Equal(t, testVlan, net["vlan"])
	assert.Equal(t, testVni, net["vni"])
}
