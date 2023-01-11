package civo

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// example.Widget represents a concrete Go type that represents an API resource
func TestAccCivoSubnet_basic(t *testing.T) {
	var subnet civogo.Subnet

	// generate a random name for each test run
	resName := "civo_subnet.foobar"
	var subnetName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoSubnetDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoSubnetConfigBasic(subnetName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoSubnetResourceExists(resName, &subnet),
					// verify remote values
					testAccCheckCivoSubnetValues(&subnet, subnetName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", subnetName),
					resource.TestCheckResourceAttr(resName, "default", "false"),
				),
			},
		},
	})
}

func TestAccCivoSubnet_update(t *testing.T) {
	var subnet civogo.Subnet

	// generate a random name for each test run
	resName := "civo_subnet.foobar"
	var subnetName = acctest.RandomWithPrefix("rename-tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoSubnetConfigUpdates(subnetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoSubnetResourceExists(resName, &subnet),
					testAccCheckCivoSubnetValues(&subnet, subnetName),
					resource.TestCheckResourceAttr(resName, "name", subnetName),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoSubnetConfigUpdates(subnetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoSubnetResourceExists(resName, &subnet),
					testAccCheckCivoSubnetUpdated(&subnet, subnetName),
					resource.TestCheckResourceAttr(resName, "name", subnetName),
				),
			},
		},
	})
}

func testAccCheckCivoSubnetValues(subnet *civogo.Subnet, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if subnet.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, subnet.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoSubnetResourceExists(n string, subnet *civogo.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindSubnet(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Subnet not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*subnet = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoSubnetUpdated(subnet *civogo.Subnet, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if subnet.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, subnet.Name)
		}
		return nil
	}
}

func testAccCheckCivoSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_subnet" {
			continue
		}

		_, err := client.FindSubnet(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}

func testAccCheckCivoSubnetConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_subnet" "foobar" {
	name = "%s"
	subnetSize = "50"
}`, name)
}

func testAccCheckCivoSubnetConfigUpdates(name string) string {
	return fmt.Sprintf(`
resource "civo_subnet" "foobar" {
	name = "%s"
	subnetSize = "50"
}`, name)
}
