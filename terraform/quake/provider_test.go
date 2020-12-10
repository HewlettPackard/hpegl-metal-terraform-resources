// Copyright (c) 2016-2020 Hewlett Packard Enterprise Development LP.

package quake

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
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
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("%s\n", err)
	}
	testAccPreCheck(t)
}

func TestProviderInterface(t *testing.T) {
	var _ *schema.Provider = Provider()
}
