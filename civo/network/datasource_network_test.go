package network_test

import (
	"fmt"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func DataSourceCivoNetwork_basic(t *testing.T) {
	datasourceName := "data.civo_network.foobar"
	name := acctest.RandomWithPrefix("net-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoNetworkConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "label", name),
				),
			},
		},
	})
}

func DataSourceCivoNetworkConfig(name string) string {
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
