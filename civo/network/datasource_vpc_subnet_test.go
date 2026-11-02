package network_test

import (
	"fmt"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoVPCSubnet_basic(t *testing.T) {
	datasourceName := "data.civo_vpc_subnet.foobar"
	subnetName := acctest.RandomWithPrefix("subnet-test")
	networkLabel := acctest.RandomWithPrefix("net-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoVPCSubnetConfig(networkLabel, subnetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", subnetName),
					resource.TestCheckResourceAttrSet(datasourceName, "network_id"),
				),
			},
		},
	})
}

func DataSourceCivoVPCSubnetConfig(networkLabel, subnetName string) string {
	return fmt.Sprintf(`
resource "civo_vpc_network" "foobar" {
	label = "%s"
}

resource "civo_vpc_subnet" "foobar" {
	name       = "%s"
	network_id = civo_vpc_network.foobar.id
}

data "civo_vpc_subnet" "foobar" {
	name       = civo_vpc_subnet.foobar.name
	network_id = civo_vpc_network.foobar.id
}
`, networkLabel, subnetName)
}
