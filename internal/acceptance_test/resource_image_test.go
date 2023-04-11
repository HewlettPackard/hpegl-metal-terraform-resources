// (C) Copyright 2023 Hewlett Packard Enterprise Development LP

package acceptance_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
)

func TestAccResourceImage_Basic(t *testing.T) {
	// as-of Project creation is only supported when using GL IAM token.
	// so, skipping test if it is explicitly disabled.
	if os.Getenv("HPEGL_METAL_GL_TOKEN") == "false" {
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.TestCheckFunc(func(s *terraform.State) error { return testAccCheckImageDestroy(t, s) }),
		Steps: []resource.TestStep{
			{
				// create step
				Config: testAccCheckImageCreateBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExists("hpegl_metal_image.image1"),
				),
			},
			{
				// update step
				Config: testAccCheckImageUpdateBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExists("hpegl_metal_image.image1"),
				),
			},
		},
	})
}

func testAccCheckImageCreateBasic() string {
	return imageBasic("create")
}

func testAccCheckImageUpdateBasic() string {
	return imageBasic("update")
}

func imageBasic(op string) string {
	common := `
	provider "hpegl" {
		metal {
		}
		alias = "test"
	}`

	file := "./service.yml"

	if op == "update" {
		file = "./servicev2.yml"
	}

	res := fmt.Sprintf(`
	resource "hpegl_metal_image" "image1" {
		os_service_image_file = %s
	}`, file)

	return common + res
}

func testAccCheckImageDestroy(t *testing.T, s *terraform.State) error {
	t.Helper()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hpegl_metal_image" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no image primary ID set")
		}

		imageID := rs.Primary.ID

		p, err := client.GetClientFromMetaMap(testAccProvider.Meta())
		if err != nil {
			return fmt.Errorf("Error retrieving Metal client: %v", err)
		}

		ctx := p.GetContext()

		_, res, err := p.Client.ServicesApi.GetByID(ctx, imageID)
		if err == nil {
			return fmt.Errorf("image %v still exists", imageID)
		}

		res.Body.Close()
	}

	return nil
}

func testAccCheckImageExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("image not found: %q", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no image primary ID set")
		}

		imageID := rs.Primary.ID

		p, err := client.GetClientFromMetaMap(testAccProvider.Meta())
		if err != nil {
			return fmt.Errorf("Error retrieving Metal client: %v", err)
		}

		ctx := p.GetContext()

		ret, res, err := p.Client.ServicesApi.GetByID(ctx, imageID)
		if err != nil {
			return fmt.Errorf("image %v not found: %s", imageID, err)
		}

		res.Body.Close()

		if ret.ID != imageID {
			return fmt.Errorf("image not found: %v - %v", rs.Primary.ID, ret.ID)
		}

		return nil
	}
}
