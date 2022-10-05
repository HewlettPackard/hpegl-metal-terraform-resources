// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
)

func Test_setConnectionsValues(t *testing.T) {
	someName := "someName"
	someVlan := int32(22)
	someIP := "someip"
	someSubnet := "somesubnet"
	someGateway := "somegateway"

	conns := []client.HostConnection{
		{
			Networks: []client.HostNetworkConnection{
				{
					Name:    someName,
					IP:      someIP,
					Subnet:  someSubnet,
					Gateway: someGateway,
					VLAN:    someVlan,
				},
			},
		},
	}

	d := schema.TestResourceDataRaw(t, hostSchema(), map[string]interface{}{})

	// test
	err := setConnectionsValues(d, conns)
	assert.Nil(t, err)

	connIPs, ok := d.Get(hConnections).(map[string]interface{})
	assert.True(t, ok, "type assertion failed")
	assert.Equal(t, 1, len(connIPs))
	assert.Equal(t, someIP, connIPs[someName])

	connSubNets, ok := d.Get(hConnectionsSubnet).(map[string]interface{})
	assert.True(t, ok, "type assertion failed")
	assert.Equal(t, 1, len(connSubNets))
	assert.Equal(t, someSubnet, connSubNets[someName])

	connGateways, ok := d.Get(hConnectionsGateway).(map[string]interface{})
	assert.True(t, ok, "type assertion failed")
	assert.Equal(t, 1, len(connGateways))
	assert.Equal(t, someGateway, connGateways[someName])
}
