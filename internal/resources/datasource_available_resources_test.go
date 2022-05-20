// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP.

package resources

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAvailableResourcesBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testQuattroAvailableResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testImages("data.hpegl_metal_available_resources.compute"),
				),
			},
		},
	})
}

const testQuattroAvailableResourceBasic = `
data "hpegl_metal_available_resources" "compute" {
	
}
`

func testImages(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource not found: %s", name)
		}
		imagesLen, _ := rs.Primary.Attributes["images.#"]
		if !ok {
			return fmt.Errorf("No attribute 'images' in resource: %s", name)
		}
		if imagesLen == "0" {
			return fmt.Errorf("No 'images.#' found in resource: %s", name)
		}
		return nil
	}
}
