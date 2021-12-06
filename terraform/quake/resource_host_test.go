// (C) Copyright 2016-2020, 2021 Hewlett Packard Enterprise Development LP

package quake

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccQuakeHost(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.TestCheckFunc(func(s *terraform.State) error { return testAccCheckHostDestroy(t, s) }),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckHostBasic(),
			},
		},
	})
}

func testAccCheckHostBasic() string {
	return fmt.Sprintf(`
variable "location" {
	# default = "USA:West Central:FTC DEV 4"  
	default = "USA:Texas:AUSL2"
}
data "quake_available_resources" "compute" {
	
}
resource "quake_host" "test_host" {
  name               = "test"
  image              = "${data.quake_available_resources.compute.images.0.image}"
  machine_size       = "Any"
  ssh                = ["User1 - Linux"]
  networks           = [for net in "${data.quake_available_resources.compute.networks}": net.name if net.location == var.location]               
  location           = var.location
  description        = "hello from Terraform"
}
`)
}

func testAccCheckHostDestroy(t *testing.T, s *terraform.State) error {
	//apiClient := testAccProvider.Meta().(*QuakeProvider)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "test_host" {
			continue
		}

		// _, err := apiClient.GetItem(rs.Primary.ID)
		// if err == nil {
		// 	return fmt.Errorf("Alert still exists")
		// }
		// notFoundErr := "not found"
		// expectedErr := regexp.MustCompile(notFoundErr)
		// if !expectedErr.Match([]byte(err.Error())) {
		// 	return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		// }
	}

	return nil
}

// func TestAccItem_Basic(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckItemDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccCheckItemBasic(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckExampleItemExists("example_item.test_item"),
// 					resource.TestCheckResourceAttr(
// 						"example_item.test_item", "name", "test"),
// 					resource.TestCheckResourceAttr(
// 						"example_item.test_item", "description", "hello"),
// 					resource.TestCheckResourceAttr(
// 						"example_item.test_item", "tags.#", "2"),
// 					resource.TestCheckResourceAttr("example_item.test_item", "tags.1931743815", "tag1"),
// 					resource.TestCheckResourceAttr("example_item.test_item", "tags.1477001604", "tag2"),
// 				),
// 			},
// 		},
// 	})
// }
