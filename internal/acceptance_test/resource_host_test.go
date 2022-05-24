// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

package acceptance_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const hostCreateWait = 1 * time.Minute

func TestAccResourceHost_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.TestCheckFunc(func(s *terraform.State) error { return testAccCheckHostDestroy(t, s) }),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckHostBasic(),
				Check:  testWaitUntilHostCreated(),
			},
		},
	})
}

func testAccCheckHostBasic() string {
	return `
provider "hpegl" {
	metal {
	}
}

variable "location" {
	default = "USA:Central:V2DCC01"
}

data "hpegl_metal_available_resources" "compute" {
}

locals  {
	host_os_flavor = "${data.hpegl_metal_available_resources.compute.images.0.flavor}"
	host_os_version = "${data.hpegl_metal_available_resources.compute.images.0.version}"
	networks = ([for net in "${data.hpegl_metal_available_resources.compute.networks}": 
		net.name if net.location == var.location] )
}

resource "hpegl_metal_host" "test_host" {
  name               = "test"
  image              = join("@",[local.host_os_flavor, local.host_os_version])
  machine_size       = "${data.hpegl_metal_available_resources.compute.machine_sizes.0.name}"
  ssh                = ["User1 - Linux"]
  networks           = sort(local.networks)
  location           = var.location
  description        = "hello from Terraform"
}
`
}

func testWaitUntilHostCreated() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// TODO
		// 1. Poll the state instead of fixed wait time
		time.Sleep(hostCreateWait)

		return nil
	}
}

func testAccCheckHostDestroy(t *testing.T, s *terraform.State) error {
	t.Helper()

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
