package network_test

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCivoVPCSubnet_basic(t *testing.T) {
	var subnet civogo.Subnet

	resName := "civo_vpc_subnet.foobar"
	var subnetName = acctest.RandomWithPrefix("tf-subnet")
	var networkLabel = acctest.RandomWithPrefix("tf-net")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoVPCSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoVPCSubnetConfigBasic(networkLabel, subnetName),
				Check: resource.ComposeTestCheckFunc(
					CivoVPCSubnetResourceExists(resName, &subnet),
					resource.TestCheckResourceAttr(resName, "name", subnetName),
					resource.TestCheckResourceAttrSet(resName, "network_id"),
				),
			},
		},
	})
}

func CivoVPCSubnetResourceExists(n string, subnet *civogo.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		client := acceptance.TestAccProvider.Meta().(*civogo.Client)
		networkID := rs.Primary.Attributes["network_id"]
		resp, err := client.GetVPCSubnet(networkID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("VPC Subnet not found: (%s) %s", rs.Primary.ID, err)
		}

		*subnet = *resp
		return nil
	}
}

func CivoVPCSubnetDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_vpc_subnet" {
			continue
		}

		networkID := rs.Primary.Attributes["network_id"]
		_, err := client.GetVPCSubnet(networkID, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPC Subnet still exists")
		}
	}

	return nil
}

func CivoVPCSubnetConfigBasic(networkLabel, subnetName string) string {
	return fmt.Sprintf(`
resource "civo_vpc_network" "foobar" {
	label = "%s"
}

resource "civo_vpc_subnet" "foobar" {
	name       = "%s"
	network_id = civo_vpc_network.foobar.id
}
`, networkLabel, subnetName)
}
