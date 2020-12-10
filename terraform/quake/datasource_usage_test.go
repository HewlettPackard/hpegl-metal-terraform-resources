package quake

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUsages_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUsageConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.quake_usage.used", "id"),
				),
			},
		},
	})
}

var testUsageConfigBasic string

func init() {

	testUsageConfigBasic = fmt.Sprintf(`
data "quake_usage" "used" {
	start = %q
}
`, time.Now().Format(time.RFC3339))
}
