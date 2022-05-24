// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

package acceptance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	testutils "github.com/hewlettpackard/hpegl-metal-terraform-resources/internal/test-utils"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = testutils.ProviderFunc()()
	testAccProviders = map[string]*schema.Provider{
		"hpegl": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	t.Helper()
}

func TestProvider(t *testing.T) {
	if err := testutils.ProviderFunc()().InternalValidate(); err != nil {
		t.Fatalf("%s\n", err)
	}
	testAccPreCheck(t)
}
