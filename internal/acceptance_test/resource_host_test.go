// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

package acceptance_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hewlettpackard/hpegl-metal-terraform-resources/pkg/client"
)

const (
	hostStateReadyWait = 30 * time.Second
	hostStatePollCount = 4
)

func TestAccResourceHost_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.TestCheckFunc(func(s *terraform.State) error { return testAccCheckHostDestroy(t, s) }),
		Steps: []resource.TestStep{
			// host create step
			{
				Config: testAccCheckHostBasic(),
				Check:  testWaitUntilHostReady("hpegl_metal_host.test_host"),
			},
			// host update step
			{
				Config: testAccHostUpdateConfig(),
				Check:  testWaitUntilHostReady("hpegl_metal_host.test_host"),
			},
		},
	})
}

func testAccCheckHostBasic() string {
	return hostConfig("create")
}

// testAccHostUpdateConfig updates the terraform config by updating
// attributes of the host 'test_host':
//   - description is updated
//   - last network removed from sorted list of networks
// Updated config is compared against config specified in testAccCheckHostBasic.
func testAccHostUpdateConfig() string {
	return hostConfig("update")
}

// hostConfig returns the host config to apply for the specified operation.
func hostConfig(op string) string {
	// common config for create/update
	common := `
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
	sorted_networks = sort(local.networks)
	updated_networks_length = length(local.sorted_networks) - 1
	updated_networks = ([for i, net in local.sorted_networks : net if i < local.updated_networks_length])
}
`
	// description, networks
	desc := `"hello from Terraform"`
	nets := `local.sorted_networks`

	if op == "update" {
		desc = `"hello from Terraform (updated)"`
		nets = `local.updated_networks`
	}

	// host block
	host := fmt.Sprintf(`
resource "hpegl_metal_host" "test_host" {
	name               = "test"
	image              = join("@",[local.host_os_flavor, local.host_os_version])
	machine_size       = "${data.hpegl_metal_available_resources.compute.machine_sizes.0.name}"
	ssh                = ["User1 - Linux"]
	networks           = %s
	network_route      = "Public"
	location           = var.location
	description        = %s
}	
`, nets, desc)

	return common + host
}

// testWaitUntilHostReady checks if the host was created successfully and
// is in the 'Ready' state.
func testWaitUntilHostReady(rsrc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rsrc]
		if !ok {
			return fmt.Errorf("Host not found: %q", rsrc)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No host primary ID set")
		}

		hostID := rs.Primary.ID

		p, err := client.GetClientFromMetaMap(testAccProvider.Meta())
		if err != nil {
			return fmt.Errorf("Error retrieving Metal client: %v", err)
		}

		ctx := p.GetContext()

		hostState := rest.HOSTSTATE_NEW
		for i := 0; i < hostStatePollCount && hostState != rest.HOSTSTATE_READY; i++ {
			time.Sleep(hostStateReadyWait)

			host, resp, err := p.Client.HostsApi.GetByID(ctx, hostID)
			if err != nil {
				return fmt.Errorf("Host: %q not found: %s", hostID, err)
			}
			defer resp.Body.Close()

			hostState = host.State
		}

		if hostState != rest.HOSTSTATE_READY {
			return fmt.Errorf("Host %s state %v != %v", hostID, hostState, rest.HOSTSTATE_READY)
		}

		return nil
	}
}

func testAccCheckHostDestroy(t *testing.T, s *terraform.State) error {
	t.Helper()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "test_host" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No host primary ID set")
		}

		hostID := rs.Primary.ID

		p, err := client.GetClientFromMetaMap(testAccProvider.Meta())
		if err != nil {
			return fmt.Errorf("Error retrieving Metal client: %v", err)
		}

		ctx := p.GetContext()

		_, resp, err := p.Client.HostsApi.GetByID(ctx, hostID)
		if err != nil {
			return fmt.Errorf("Error retrieving host %s: %v", hostID, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Host %s exists. Response code: %d", hostID, resp.StatusCode)
		}
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
