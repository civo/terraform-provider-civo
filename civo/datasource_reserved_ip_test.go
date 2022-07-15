package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceReservedIP_basic(t *testing.T) {
	datasourceName := "data.civo_reserved_ip.foobar"
	name := acctest.RandomWithPrefix("ip-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceReservedIPConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
				),
			},
		},
	})
}

func testAccDataSourceReservedIPConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_reserved_ip" "foobar" {
	name = "%s"
}

data "civo_network" "foobar" {
	name = civo_network.foobar.name
}
`, name)
}
