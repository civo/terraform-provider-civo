package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoSubnet_basic(t *testing.T) {
	datasourceName := "data.civo_subnet.foobar"
	name := acctest.RandomWithPrefix("subnet-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoSubnetConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
				),
			},
		},
	})
}

func testAccDataSourceCivoSubnetConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_subnet" "foobar" {
	name = "%s"
	subnetSize = "50"
}

data "civo_subnet" "foobar" {
	name = civo_subnet.foobar.name
	region = "50"
}
`, name)
}
