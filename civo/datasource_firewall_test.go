package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoFirewall_basic(t *testing.T) {
	datasourceName := "data.civo_firewall.foobar"
	name := acctest.RandomWithPrefix("net-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoFirewallConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
				),
			},
		},
	})
}

func testAccDataSourceCivoFirewallConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_firewall" "foobar" {
	label = "%s"
}

data "civo_firewall" "foobar" {
	name = civo_firewall.foobar.name
	region = "LON1"
}
`, name)
}
