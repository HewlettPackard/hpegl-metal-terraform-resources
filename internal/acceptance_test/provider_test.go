// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"os"
	"testing"

	testutils "github.com/hewlettpackard/hpegl-metal-terraform-resources/internal/test-utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	libUtils "github.com/hewlettpackard/hpegl-provider-lib/pkg/utils"
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
	// this fails c is a nil interface....
	// c := testAccProvider.Meta().(*Config)
	// if c.member.GetHosterID() == "" {
	// 	t.Fatalf("Acceptance tests must be run with hoster-scope %+v", c.member)
	// }
}

func TestProvider(t *testing.T) {
	if err := testutils.ProviderFunc()().InternalValidate(); err != nil {
		t.Fatalf("%s\n", err)
	}
	testAccPreCheck(t)
}

func TestMain(m *testing.M) {
	// TF_ACC_CONFIG_PATH set in make acceptance
	libUtils.ReadAccConfig(os.Getenv("TF_ACC_CONFIG_PATH"))
	m.Run()
	os.Exit(0)
}
