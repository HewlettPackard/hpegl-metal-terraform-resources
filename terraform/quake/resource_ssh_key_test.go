// (C) Copyright 2016-2021 Hewlett Packard Enterprise Development LP.

package quake

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hpe-hcss/quake-client/pkg/terraform/configuration"
	rest "github.com/hpe-hcss/quake-client/v1/pkg/client"
)

func quatrroSSHKeyConfigBasic(name, publicSSHKey string) string {
	return fmt.Sprintf(`
resource "quake_ssh_key" "test" {
    name       = %q
    public_key = %q
}
`, name, publicSSHKey)
}

func TestAccQuattroSSHKey_Basic(t *testing.T) {
	key := rest.SshKey{}
	keyName := "keySShNameTest"
	keyPublic := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0euPI4c1c5qJcAdHWlV1zI2SGbo136AcL" +
		"0MzkBSxRrm39ve9qrXWYpd50p6uBxG4U4y71MiNC1y5FTmtFyISlIlPR74bESben" +
		"MwUGk++Qliyl0fofjs3DNjiwKAbYEbPrh8taMtZgUDEwbs4EweFmfVqJfnk781vK" +
		"R4A6QVYssv3Q+Wl8XZAEM7keSYZMuPnnaqkU8s2dZQKpPjElMe0yC40U2ZIwQTAg" +
		"Pn2Im1oH4KftTYzhsty2BlFZU3ZTqvb5ocjWzlcLxF2LJxeKin5d8C0jd8w6PZMC" +
		"u7awCSpdcXbicti2mRCe7HPNcdP9FU6hwEEtMsuIBxsGUue6sQCL body@nowhere.com"
	cfg := quatrroSSHKeyConfigBasic(keyName, keyPublic)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckQuattroSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQuattroSSHKeyExists("quake_ssh_key.test", &key),
					resource.TestCheckResourceAttr("quake_ssh_key.test", "public_key", keyPublic),
				),
			},
		},
	})
}

func testAccCheckQuattroSSHKeyExists(n string, out *rest.SshKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		p := testAccProvider.Meta().(*configuration.Config)

		ctx := p.GetContext()
		key, _, err := p.Client.SshkeysApi.GetByID(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if key.ID != rs.Primary.ID {
			return fmt.Errorf("SSHKey not found: %v - %v", rs.Primary.ID, key)
		}
		*out = key
		return nil
	}
}

func testAccCheckQuattroSSHKeyDestroy(s *terraform.State) error {
	p := testAccProvider.Meta().(*configuration.Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "quake_ssh_key" {
			continue
		}

		ctx := p.GetContext()
		if _, _, err := p.Client.SshkeysApi.GetByID(ctx, rs.Primary.ID); err == nil {
			return fmt.Errorf("SSHKey still exists")
		}
	}

	return nil
}
