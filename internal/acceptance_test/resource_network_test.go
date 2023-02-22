// (C) Copyright 2020-2023 Hewlett Packard Enterprise Development LP

package acceptance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
)

func TestAccResourceNetwork_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.TestCheckFunc(func(s *terraform.State) error { return testAccCheckNetworkDestroy(t, s) }),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetworkBasic(),
			},
		},
	})
}

func TestAccResourceNetwork_OptFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.TestCheckFunc(func(s *terraform.State) error { return testAccCheckNetworkDestroy(t, s) }),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetworkOptFields(),
			},
		},
	})
}

func testAccCheckNetworkBasic() string {
	return `
provider "hpegl" {
	metal {
	}
	alias = "test"
}

variable "location" {
	default = "USA:Central:AFCDCC1"
}

resource "hpegl_metal_network" "pnet" {
  provider           = hpegl.test
  name               = "pnet-test"              
  location           = var.location
  description        = "tf-net description"
}`
}

func testAccCheckNetworkOptFields() string {
	return `
provider "hpegl" {
	metal {
	}
	alias = "test"
}

variable "location" {
	default = "USA:Central:AFCDCC1"
}

resource "hpegl_metal_network" "pnet" {
  provider           = hpegl.test
  name               = "pnet-test"              
  location           = var.location
  description        = "tf-net description"
  host_use           = "Default"
  purpose            = "Storage"
}`
}

func testAccCheckNetworkDestroy(t *testing.T, s *terraform.State) error {
	t.Helper()

	p, err := client.GetClientFromMetaMap(testAccProvider.Meta())
	if err != nil {
		return fmt.Errorf("Error retrieving Metal client: %v", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pnet" {
			continue
		}

		ctx := p.GetContext()

		_, _, err := p.Client.NetworksApi.GetByID(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert pnet still exists")
		}
	}

	return nil
}
