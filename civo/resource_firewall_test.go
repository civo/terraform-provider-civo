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
func TestAccCivoFirewall_basic(t *testing.T) {
	var firewall civogo.Firewall

	// generate a random name for each test run
	resName := "civo_firewall.foobar"
	var firewallName = acctest.RandomWithPrefix("tf-fw")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoFirewallDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoFirewallConfigBasic(firewallName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoFirewallResourceExists(resName, &firewall),
					// verify remote values
					testAccCheckCivoFirewallValues(&firewall, firewallName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", firewallName),
				),
			},
		},
	})
}

func TestAccCivoFirewall_update(t *testing.T) {
	var firewall civogo.Firewall

	// generate a random name for each test run
	resName := "civo_firewall.foobar"
	var firewallName = acctest.RandomWithPrefix("tf-fw")
	var firewallNameUpdate = acctest.RandomWithPrefix("rename-fw")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoFirewallConfigBasic(firewallName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoFirewallResourceExists(resName, &firewall),
					testAccCheckCivoFirewallValues(&firewall, firewallName),
					resource.TestCheckResourceAttr(resName, "name", firewallName),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoFirewallConfigUpdates(firewallNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoFirewallResourceExists(resName, &firewall),
					testAccCheckCivoFirewallUpdated(&firewall, firewallNameUpdate),
					resource.TestCheckResourceAttr(resName, "name", firewallNameUpdate),
				),
			},
		},
	})
}

func testAccCheckCivoFirewallValues(firewall *civogo.Firewall, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if firewall.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, firewall.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoFirewallResourceExists(n string, firewall *civogo.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindFirewall(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Firewall not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*firewall = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoFirewallUpdated(firewall *civogo.Firewall, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if firewall.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, firewall.Name)
		}
		return nil
	}
}

func testAccCheckCivoFirewallDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_firewall" {
			continue
		}

		_, err := client.FindFirewall(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Firewall still exists")
		}
	}

	return nil
}

func testAccCheckCivoFirewallConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_firewall" "foobar" {
	name = "%s"
}`, name)
}

func testAccCheckCivoFirewallConfigUpdates(name string) string {
	return fmt.Sprintf(`
resource "civo_firewall" "foobar" {
	name = "%s"
}`, name)
}
