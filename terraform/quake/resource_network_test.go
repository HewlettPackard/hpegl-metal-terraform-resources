// Copyright (c) 2016-2020 Hewlett Packard Enterprise Development LP.

package quake

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccQuakeNetwork(t *testing.T) {
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

func testAccCheckNetworkBasic() string {
	return fmt.Sprintf(`
variable "location" {
	# default = "USA:West Central:FTC DEV 4"  
	default = "USA:Texas:AUSL2"
}
resource "quake_network" "pnet" {
  name               = "pnet-test"              
  location           = var.location
  description        = "tf-net description"
}
`)
}

func testAccCheckNetworkDestroy(t *testing.T, s *terraform.State) error {
	p := testAccProvider.Meta().(*Config)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pnet" {
			continue
		}

		_, _, err := p.client.NetworksApi.GetByID(p.context, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert pnet still exists")
		}
	}

	return nil
}