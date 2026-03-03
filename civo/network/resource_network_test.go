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

// example.Widget represents a concrete Go type that represents an API resource
func TestAccCivoVPCNetwork_basic(t *testing.T) {
	var network civogo.Network

	// generate a random name for each test run
	resName := "civo_vpc_network.foobar"
	var networkLabel = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoVPCNetworkDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: CivoVPCNetworkConfigBasic(networkLabel),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					CivoVPCNetworkResourceExists(resName, &network),
					// verify remote values
					CivoVPCNetworkValues(&network, networkLabel),
					// verify local values
					resource.TestCheckResourceAttr(resName, "label", networkLabel),
					resource.TestCheckResourceAttr(resName, "default", "false"),
				),
			},
		},
	})
}

func TestAccCivoVPCNetwork_update(t *testing.T) {
	var network civogo.Network

	// generate a random name for each test run
	resName := "civo_vpc_network.foobar"
	var networkLabel = acctest.RandomWithPrefix("rename-tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoVPCNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoVPCNetworkConfigUpdates(networkLabel),
				Check: resource.ComposeTestCheckFunc(
					CivoVPCNetworkResourceExists(resName, &network),
					CivoVPCNetworkValues(&network, networkLabel),
					resource.TestCheckResourceAttr(resName, "label", networkLabel),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: CivoVPCNetworkConfigUpdates(networkLabel),
				Check: resource.ComposeTestCheckFunc(
					CivoVPCNetworkResourceExists(resName, &network),
					CivoVPCNetworkUpdated(&network, networkLabel),
					resource.TestCheckResourceAttr(resName, "label", networkLabel),
				),
			},
		},
	})
}

func CivoVPCNetworkValues(network *civogo.Network, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if network.Label != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, network.Label)
		}
		return nil
	}
}

// CheckExampleResourceExists queries the API and retrieves the matching Widget.
func CivoVPCNetworkResourceExists(n string, network *civogo.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := acceptance.TestAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindVPCNetwork(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Network not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*network = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func CivoVPCNetworkUpdated(network *civogo.Network, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if network.Label != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, network.Label)
		}
		return nil
	}
}

func CivoVPCNetworkDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_vpc_network" {
			continue
		}

		_, err := client.FindVPCNetwork(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Network still exists")
		}
	}

	return nil
}

func CivoVPCNetworkConfigBasic(label string) string {
	return fmt.Sprintf(`
resource "civo_vpc_network" "foobar" {
	label = "%s"
}`, label)
}

func CivoVPCNetworkConfigUpdates(label string) string {
	return fmt.Sprintf(`
resource "civo_vpc_network" "foobar" {
	label = "%s"
}`, label)
}
