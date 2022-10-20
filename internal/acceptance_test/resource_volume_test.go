// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

package acceptance_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
)

const (
	testVolCreateSize = 10
	testVolUpdateSize = 12
)

func TestAccResourceVolume_Basic(t *testing.T) {
	var createID, updateID string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.TestCheckFunc(func(s *terraform.State) error { return testAccCheckVolumeDestroy(t, s) }),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVolumeCreateBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("hpegl_metal_volume.test_vol", &createID),
				),
			},
			{
				Config: testAccCheckVolumeUpdateBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("hpegl_metal_volume.test_vol", &updateID),
					resource.TestCheckFunc(func(s *terraform.State) error {
						if createID == updateID {
							return nil
						}

						return fmt.Errorf("Create and Update Volume ID are not same")
					}),
					resource.TestCheckResourceAttr("hpegl_metal_volume.test_vol", "size", strconv.Itoa(testVolUpdateSize)),
				),
			},
		},
	})
}

func testAccCheckVolumeCreateBasic() string {
	return testAccCheckVolumeBasic("create")
}

func testAccCheckVolumeUpdateBasic() string {
	return testAccCheckVolumeBasic("update")
}

func testAccCheckVolumeBasic(op string) string {
	common := `
	provider "hpegl" {
		metal {
		}
	}
	
	variable "location" {
		default = "USA:Matrix:AFCDCC1"
	}
	`
	size := testVolCreateSize

	if op == "update" {
		size = testVolUpdateSize
	}

	res := fmt.Sprintf(`resource "hpegl_metal_volume" "test_vol" {
		name        = "test.volume"
		size        = %d
		flavor      = "Fast"
		description = "hello from Terraform"
		location    = var.location
	  }`, size)

	return common + res
}

func testAccCheckVolumeDestroy(t *testing.T, s *terraform.State) error {
	t.Helper()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hpegl_metal_volume" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume primary ID set")
		}

		volumeID := rs.Primary.ID

		p, err := client.GetClientFromMetaMap(testAccProvider.Meta())
		if err != nil {
			return fmt.Errorf("Error retrieving Metal client: %v", err)
		}

		ctx := p.GetContext()

		volume, _, err := p.Client.VolumesApi.GetByID(ctx, volumeID)
		if err == nil && volume.State != rest.VOLUMESTATE_DELETED {
			return fmt.Errorf("Volume: %v still exists", volume)
		}
	}

	return nil
}

func testAccCheckVolumeExists(resource string, id *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Volume not found: %q", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume primary ID set")
		}

		volumeID := rs.Primary.ID

		p, err := client.GetClientFromMetaMap(testAccProvider.Meta())
		if err != nil {
			return fmt.Errorf("Error retrieving Metal client: %v", err)
		}

		ctx := p.GetContext()

		_, _, err = p.Client.VolumesApi.GetByID(ctx, volumeID)
		if err != nil {
			return fmt.Errorf("Volume: %q not found: %s", volumeID, err)
		}

		*id = volumeID

		return nil
	}
}
