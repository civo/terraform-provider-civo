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
func TestAccCivoNetwork_basic(t *testing.T) {
	var network civogo.Network

	// generate a random name for each test run
	resName := "civo_network.foobar"
	var networkLabel = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoNetworkDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoNetworkConfigBasic(networkLabel),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoNetworkResourceExists(resName, &network),
					// verify remote values
					testAccCheckCivoNetworkValues(&network, networkLabel),
					// verify local values
					resource.TestCheckResourceAttr(resName, "label", networkLabel),
					resource.TestCheckResourceAttr(resName, "default", "false"),
				),
			},
		},
	})
}

func TestAccCivoNetwork_update(t *testing.T) {
	var network civogo.Network

	// generate a random name for each test run
	resName := "civo_network.foobar"
	var networkLabel = acctest.RandomWithPrefix("rename-tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoNetworkConfigUpdates(networkLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoNetworkResourceExists(resName, &network),
					testAccCheckCivoNetworkValues(&network, networkLabel),
					resource.TestCheckResourceAttr(resName, "label", networkLabel),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoNetworkConfigUpdates(networkLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoNetworkResourceExists(resName, &network),
					testAccCheckCivoNetworkUpdated(&network, networkLabel),
					resource.TestCheckResourceAttr(resName, "label", networkLabel),
				),
			},
		},
	})
}

func testAccCheckCivoNetworkValues(network *civogo.Network, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if network.Label != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, network.Label)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoNetworkResourceExists(n string, network *civogo.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindNetwork(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Network not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*network = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoNetworkUpdated(network *civogo.Network, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if network.Label != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, network.Label)
		}
		return nil
	}
}

func testAccCheckCivoNetworkDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_network" {
			continue
		}

		_, err := client.FindNetwork(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Network still exists")
		}
	}

	return nil
}

func testAccCheckCivoNetworkConfigBasic(label string) string {
	return fmt.Sprintf(`
resource "civo_network" "foobar" {
	label = "%s"
}`, label)
}

func testAccCheckCivoNetworkConfigUpdates(label string) string {
	return fmt.Sprintf(`
resource "civo_network" "foobar" {
	label = "%s"
}`, label)
}
