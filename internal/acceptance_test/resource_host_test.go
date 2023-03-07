// (C) Copyright 2020-2023 Hewlett Packard Enterprise Development LP

package acceptance_test

import (
	"fmt"
	"net/http"
	"strconv"
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
	isAsync            = true
	isNotAsync         = false
)

func TestAccResourceHost_Async(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.TestCheckFunc(func(s *terraform.State) error { return testAccCheckHostDestroy(t, s) }),
		Steps: []resource.TestStep{
			// host create step
			{
				Config: testAccCheckHostBasic(isAsync),
				Check:  testWaitUntilHostReady("hpegl_metal_host.test_host"),
			},
			// host update step
			{
				Config: testAccHostUpdateConfig(isAsync),
				Check:  testWaitUntilHostReady("hpegl_metal_host.test_host"),
			},
		},
	})
}

func TestAccResourceHost_Sync(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.TestCheckFunc(func(s *terraform.State) error { return testAccCheckHostDestroy(t, s) }),
		Steps: []resource.TestStep{
			// host create step
			{
				Config: testAccCheckHostBasic(isNotAsync),
				Check:  testVerifyHostReady("hpegl_metal_host.test_host"),
			},
			// host update step
			{
				Config: testAccHostUpdateConfig(isNotAsync),
				Check:  testVerifyHostReady("hpegl_metal_host.test_host"),
			},
		},
	})
}

func testAccCheckHostBasic(async bool) string {
	return hostConfig("create", async)
}

// testAccHostUpdateConfig updates the terraform config by updating
// attributes of the host 'test_host':
//   - description is updated
//   - last network removed from sorted list of networks
//
// Updated config is compared against config specified in testAccCheckHostBasic.
func testAccHostUpdateConfig(async bool) string {
	return hostConfig("update", async)
}

// hostConfig returns the host config to apply for the specified operation.
func hostConfig(op string, async bool) string {
	// common config for create/update
	common := `
provider "hpegl" {
	metal {
	}
	alias = "test"
}

variable "location" {
	default = "USA:Central:AFCDCC1"
}

data "hpegl_metal_available_resources" "compute" {
	provider = hpegl.test
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
	untagged := `"Private"`

	if op == "update" {
		desc = `"hello from Terraform (updated)"`
		nets = `local.updated_networks`
		untagged = `""`
	}

	name := "testAsync"
	if !async {
		name = "testSync"
	}

	// host block
	host := fmt.Sprintf(`
resource "hpegl_metal_host" "test_host" {
	provider           = hpegl.test
	name               = "%s"
	image              = join("@",[local.host_os_flavor, local.host_os_version])
	machine_size       = "${data.hpegl_metal_available_resources.compute.machine_sizes.0.name}"
	ssh                = ["User1 - Linux"]
	networks           = %s
	network_route      = "Public"
	network_untagged   = %s
	location           = var.location
	description        = %s
	host_action_async  = %s
}	
`, name, nets, untagged, desc, strconv.FormatBool(async))

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

			resp.Body.Close()

			hostState = host.State
		}

		if hostState != rest.HOSTSTATE_READY {
			return fmt.Errorf("Host %s state %v != %v", hostID, hostState, rest.HOSTSTATE_READY)
		}

		return nil
	}
}

// testVerifyHostReady checks if the host was created successfully and
// is in the 'Ready' state.
func testVerifyHostReady(resourceStateKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceStateKey]
		if !ok {
			return fmt.Errorf("Host not found: %q", resourceStateKey)
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

		host, resp, err := p.Client.HostsApi.GetByID(ctx, hostID)
		if err != nil {
			return fmt.Errorf("Host: %q not found: %s", hostID, err)
		}

		resp.Body.Close()

		if got, want := host.State, rest.HOSTSTATE_READY; got != want {
			return fmt.Errorf("Host %s state %v != %v", hostID, got, want)
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

		resp.Body.Close()

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
