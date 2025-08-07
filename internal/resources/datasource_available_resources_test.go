// (C) Copyright 2023, 2025 Hewlett Packard Enterprise Development LP

package resources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/configuration"
)

func Test_addNetworks(t *testing.T) {
	const (
		testVlan  = 200
		testVni   = 12006
		testVlan2 = 200
		testVni2  = 12006
	)

	availableNets := client.AvailableResources{
		Networks: []client.AvailableNetwork{
			{
				ID:   "testID",
				VLAN: int32(testVlan),
				VNI:  int32(testVni),
			},
			{
				ID:       "testID2",
				VLAN:     int32(testVlan2),
				VNI:      int32(testVni2),
				NoIPPool: true,
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
	assert.True(t, ok, "type assertion failed for networks")
	assert.Equal(t, 2, len(networks))

	net, ok := networks[0].(map[string]interface{})
	assert.True(t, ok, "type assertion failed for network 1")

	assert.Equal(t, testVlan, net["vlan"])
	assert.Equal(t, testVni, net["vni"])
	assert.Equal(t, false, net["no_ip_pool"])

	net, ok = networks[1].(map[string]interface{})
	assert.True(t, ok, "type assertion failed for network 2")

	assert.Equal(t, testVlan2, net["vlan"])
	assert.Equal(t, testVni2, net["vni"])
	assert.Equal(t, true, net["no_ip_pool"])
}
