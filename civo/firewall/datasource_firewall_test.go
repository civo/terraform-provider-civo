package firewall_test

import (
	"fmt"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoFirewall_basic(t *testing.T) {
	datasourceName := "data.civo_firewall.foobar"
	name := acctest.RandomWithPrefix("net-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoFirewallConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
				),
			},
		},
	})
}

func DataSourceCivoFirewallConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_firewall" "foobar" {
	name = "%s"
}

data "civo_firewall" "foobar" {
	name = civo_firewall.foobar.name
	region = "LON1"
}
`, name)
}
