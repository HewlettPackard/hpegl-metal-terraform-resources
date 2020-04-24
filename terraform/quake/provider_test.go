// Copyright (c) 2016-2020 Hewlett Packard Enterprise Development LP.

package quake

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var (
	testAccProviders map[string]terraform.ResourceProvider
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		Quake: testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	// this fails c is a nil interface....
	// c := testAccProvider.Meta().(*Config)
	// if c.member.GetHosterID() == "" {
	// 	t.Fatalf("Acceptance tests must be run with hoster-scope %+v", c.member)
	// }
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("%s\n", err)
	}
	testAccPreCheck(t)
}

func TestProviderInterface(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}
