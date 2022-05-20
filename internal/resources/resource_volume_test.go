// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP.

package resources

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/configuration"
)

func TestAccQuattroVolume(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.TestCheckFunc(func(s *terraform.State) error { return testAccCheckVolumeDestroy(t, s) }),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVolumeBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("hpegl_metal_volume.test_vol"),
				),
			},
		},
	})
}

func testAccCheckVolumeBasic() string {
	return fmt.Sprintf(`
resource "hpegl_metal_volume" "test_vol" {
  name        = "test.volume"
  size        = 10
  flavor      = "Fast"
  description = "hello from Terraform"
  location    = "USA:Texas:AUSL2"
}
`)
}

func testAccCheckVolumeDestroy(t *testing.T, s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hpegl_metal_volume" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume primary ID set")
		}
		volumeID := rs.Primary.ID
		p := testAccProvider.Meta().(*configuration.Config)
		ctx := p.GetContext()
		volume, _, err := p.Client.VolumesApi.GetByID(ctx, volumeID)
		if err == nil && volume.State != rest.VOLUMESTATE_DELETED {
			return fmt.Errorf("Volume: %v still exists", volume)
		}
	}

	return nil
}

func testAccCheckVolumeExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Volume not found: %q", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume primary ID set")
		}
		volumeID := rs.Primary.ID
		p := testAccProvider.Meta().(*configuration.Config)
		ctx := p.GetContext()
		_, _, err := p.Client.VolumesApi.GetByID(ctx, volumeID)
		if err != nil {
			return fmt.Errorf("Volume: %q not found: %s", volumeID, err)
		}
		return nil
	}
}
