// (C) Copyright 2020, 2022 Hewlett Packard Enterprise Development LP

package acceptance_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUsages_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUsageConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hpegl_metal_usage.used", "id"),
				),
			},
		},
	})
}

func testUsageConfigBasic() string {
	return fmt.Sprintf(`
	provider "hpegl" {
		metal {
		}
	}
	
	data "hpegl_metal_usage" "used" {
		start = %q
	}
	`, time.Now().Format(time.RFC3339))
}
