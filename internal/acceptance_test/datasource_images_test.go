// (C) Copyright 2020, 2022-2023 Hewlett Packard Enterprise Development LP

package acceptance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testOperatingSystemConfigBasic = `
provider "hpegl" {
	metal {
	}
	alias = "test"
}

data "hpegl_metal_available_images" "example" {
	provider=hpegl.test

	filter {
		name = "category"
		values = ["linux"]
	}
	filter {
		name = "flavor" 
		values = ["ubuntu"]
	}
}`

func TestAccDataSourceAvailableImages_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOperatingSystemConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hpegl_metal_available_images.example", "images.0.flavor", "ubuntu"),
					resource.TestCheckResourceAttr("data.hpegl_metal_available_images.example", "images.0.category", "linux"),
				),
			},
		},
	})
}
