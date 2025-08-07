// (C) Copyright 2022,2024-2025 Hewlett Packard Enterprise Development LP

package acceptance_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
)

func TestAccResourceProject_Basic(t *testing.T) {
	// as-of Project creation is only supported when using GL IAM token.
	// so, skipping test if it is explicitly disabled.
	if os.Getenv("HPEGL_METAL_GL_TOKEN") == "false" {
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.TestCheckFunc(func(s *terraform.State) error { return testAccCheckProjectDestroy(t, s) }),
		Steps: []resource.TestStep{
			{
				// create step
				Config: testAccCheckProjectCreateBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectExists("hpegl_metal_project.project1"),
					resource.TestCheckResourceAttr("hpegl_metal_project.project1", "boot_from_san_support", "true"),
				),
			},
			{
				// update step
				Config: testAccCheckProjectUpdateBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectExists("hpegl_metal_project.project1"),
					resource.TestCheckResourceAttr("hpegl_metal_project.project1", "boot_from_san_support", "false"),
				),
			},
		},
	})
}

func testAccCheckProjectCreateBasic() string {
	return projectBasic("create")
}

func testAccCheckProjectUpdateBasic() string {
	return projectBasic("update")
}

func projectBasic(op string) string {
	common := `
	provider "hpegl" {
		metal {
		}
		alias = "test"
	}`

	name := `"TestHoster1-SimProject1"`
	company := `"HPE"`
	hosts := 10
	sites := `["1ad98170-993e-4bfc-8b84-e689ea9a429b"]`
	bfsSupport := true

	if op == "update" {
		name = `"TestHoster1-SimProject1-Update"`
		company = `"HPE-Update"`
		hosts = 5
		sites = `["22473578-e18b-4753-a2e6-ba405b8abc32", "1ad98170-993e-4bfc-8b84-e689ea9a429b"]`
		bfsSupport = false
	}

	res := fmt.Sprintf(`
	resource "hpegl_metal_project" "project1" {
		provider            = hpegl.test
		name = %s
		profile {
		company             = %s
		address             = "Area51"
		email               = "acme@intergalactic.universe"
		phone_number        = "+112 234 1245 3245"
		project_description = "Primitive Life"
		project_name        = "Umbrella Corporation"
		}
		limits {
		hosts            = %d
		volumes          = 10
		volume_capacity  = 300
		private_networks = 30
		}
		sites            = %s
		volume_replication_enabled = true
		boot_from_san_support      = %v
	}`, name, company, hosts, sites, bfsSupport)

	return common + res
}

func testAccCheckProjectDestroy(t *testing.T, s *terraform.State) error {
	t.Helper()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hpegl_metal_Project" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Project primary ID set")
		}

		ProjectID := rs.Primary.ID

		p, err := client.GetClientFromMetaMap(testAccProvider.Meta())
		if err != nil {
			return fmt.Errorf("Error retrieving Metal client: %v", err)
		}

		ctx := p.GetContext()

		_, res, err := p.Client.ProjectsApi.GetByID(ctx, ProjectID, nil)
		if err == nil {
			return fmt.Errorf("Project %v still exists", ProjectID)
		}

		res.Body.Close()
	}

	return nil
}

func testAccCheckProjectExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Project not found: %q", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Project primary ID set")
		}

		ProjectID := rs.Primary.ID

		p, err := client.GetClientFromMetaMap(testAccProvider.Meta())
		if err != nil {
			return fmt.Errorf("Error retrieving Metal client: %v", err)
		}

		ctx := p.GetContext()

		ret, res, err := p.Client.ProjectsApi.GetByID(ctx, ProjectID, nil)
		if err != nil {
			return fmt.Errorf("Project %v not found: %s", ProjectID, err)
		}

		res.Body.Close()

		if ret.ID != ProjectID {
			return fmt.Errorf("Project not found: %v - %v", rs.Primary.ID, ret.ID)
		}

		return nil
	}
}
