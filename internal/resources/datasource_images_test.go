// (C) Copyright 2020 Hewlett Packard Enterprise Development LP.

package resources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImages_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOperatingSystemConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.metal_available_images.example", "id"),
				),
			},
		},
	})
}

const testOperatingSystemConfigBasic = `
data "metal_available_images" "example" {
	filter {
		name = "category"
		values = ["linux"]
	}
	filter {
		name = "flavor" 
		values = ["coreos"]
	}
}`
