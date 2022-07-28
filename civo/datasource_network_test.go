package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoNetwork_basic(t *testing.T) {
	datasourceName := "data.civo_network.foobar"
	name := acctest.RandomWithPrefix("net-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoNetworkConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "label", name),
				),
			},
		},
	})
}

func testAccDataSourceCivoNetworkConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_network" "foobar" {
	label = "%s"
	region = "LON1"
}

data "civo_network" "foobar" {
	label = civo_network.foobar.name
	region = "LON1"
}
`, name)
}
