package ip_test

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCivoVPCReservedIP_basic(t *testing.T) {
	var ip civogo.IP

	// generate a random name for each test run
	resName := "civo_vpc_reserved_ip.foobar"
	var nameIP = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoVPCReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: CivoVPCReservedIPConfigBasic(nameIP),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					acceptance.CivoReservedIPResourceExists(resName, &ip),
					// verify remote values
					CivoVPCReservedIPValues(&ip, nameIP),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", nameIP),
				),
			},
		},
	})
}

func TestAccCivoVPCReservedIP_update(t *testing.T) {
	var ip civogo.IP

	// generate a random name for each test run
	resName := "civo_vpc_reserved_ip.foobar"
	var nameIP = acctest.RandomWithPrefix("rename-tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoVPCReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoVPCReservedIPConfigUpdates(nameIP),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CivoReservedIPResourceExists(resName, &ip),
					CivoVPCReservedIPValues(&ip, nameIP),
					resource.TestCheckResourceAttr(resName, "name", nameIP),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: CivoVPCReservedIPConfigUpdates(nameIP),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CivoReservedIPResourceExists(resName, &ip),
					CivoVPCReservedIPUpdated(&ip, nameIP),
					resource.TestCheckResourceAttr(resName, "name", nameIP),
				),
			},
		},
	})
}

func CivoVPCReservedIPValues(ip *civogo.IP, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if ip.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, ip.Name)
		}
		return nil
	}
}

func CivoVPCReservedIPUpdated(ip *civogo.IP, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if ip.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, ip.Name)
		}
		return nil
	}
}

func CivoVPCReservedIPDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_vpc_reserved_ip" {
			continue
		}

		_, err := client.FindVPCNetwork(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("IP still exists")
		}
	}

	return nil
}

func CivoVPCReservedIPConfigBasic(label string) string {
	return fmt.Sprintf(`
resource "civo_vpc_reserved_ip" "foobar" {
	name = "%s"
}`, label)
}

func CivoVPCReservedIPConfigUpdates(label string) string {
	return fmt.Sprintf(`
resource "civo_vpc_reserved_ip" "foobar" {
	name = "%s"
}`, label)
}
