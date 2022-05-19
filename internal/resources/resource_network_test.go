// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP.

package resources

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/configuration"
)

func TestAccMetalNetwork(t *testing.T) {
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
resource "metal_network" "pnet" {
  name               = "pnet-test"              
  location           = var.location
  description        = "tf-net description"
}
`)
}

func testAccCheckNetworkDestroy(t *testing.T, s *terraform.State) error {
	p := testAccProvider.Meta().(*configuration.Config)
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
