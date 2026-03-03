package firewall_test

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
func TestAccCivoVPCFirewall_basic(t *testing.T) {
	var firewall civogo.Firewall

	// generate a random name for each test run
	resName := "civo_vpc_firewall.foobar"
	var firewallName = acctest.RandomWithPrefix("tf-fw")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoVPCFirewallDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: CivoVPCFirewallConfigBasic(firewallName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					CivoVPCFirewallResourceExists(resName, &firewall),
					// verify remote values
					CivoVPCFirewallValues(&firewall, firewallName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", firewallName),
					resource.TestCheckResourceAttrSet(resName, "ingress_rule.#"),
				),
			},
		},
	})
}

func TestAccCivoVPCFirewallWithIngressEgress_basic(t *testing.T) {
	var firewall civogo.Firewall

	// generate a random name for each test run
	resName := "civo_vpc_firewall.foobar"
	var firewallName = acctest.RandomWithPrefix("tf-fw")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoVPCFirewallDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: CivoVPCFirewallConfigWithIngressEgress(firewallName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					CivoVPCFirewallResourceExists(resName, &firewall),
					// verify remote values
					CivoVPCFirewallValues(&firewall, firewallName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", firewallName),
					resource.TestCheckResourceAttrSet(resName, "ingress_rule.#"),
					resource.TestCheckResourceAttrSet(resName, "egress_rule.#"),
					resource.TestCheckResourceAttrSet(resName, "ingress_rule.0.id"),
					resource.TestCheckResourceAttr(resName, "ingress_rule.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resName, "ingress_rule.0.port_range", "443"),
					resource.TestCheckResourceAttrSet(resName, "egress_rule.0.id"),
					resource.TestCheckResourceAttr(resName, "egress_rule.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resName, "egress_rule.0.port_range", "22"),
				),
			},
		},
	})
}

func TestAccCivoVPCFirewall_update(t *testing.T) {
	var firewall civogo.Firewall

	// generate a random name for each test run
	resName := "civo_vpc_firewall.foobar"
	var firewallName = acctest.RandomWithPrefix("tf-fw")
	var firewallNameUpdate = acctest.RandomWithPrefix("rename-fw")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoVPCFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoVPCFirewallConfigBasic(firewallName),
				Check: resource.ComposeTestCheckFunc(
					CivoVPCFirewallResourceExists(resName, &firewall),
					CivoVPCFirewallValues(&firewall, firewallName),
					resource.TestCheckResourceAttr(resName, "name", firewallName),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: CivoVPCFirewallConfigUpdates(firewallNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					CivoVPCFirewallResourceExists(resName, &firewall),
					CivoVPCFirewallUpdated(&firewall, firewallNameUpdate),
					resource.TestCheckResourceAttr(resName, "name", firewallNameUpdate),
				),
			},
		},
	})
}

func CivoVPCFirewallValues(firewall *civogo.Firewall, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if firewall.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, firewall.Name)
		}
		return nil
	}
}

// CheckExampleResourceExists queries the API and retrieves the matching Widget.
func CivoVPCFirewallResourceExists(n string, firewall *civogo.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := acceptance.TestAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindVPCFirewall(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Firewall not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*firewall = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func CivoVPCFirewallUpdated(firewall *civogo.Firewall, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if firewall.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, firewall.Name)
		}
		return nil
	}
}

func CivoVPCFirewallDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_vpc_firewall" {
			continue
		}

		_, err := client.FindVPCFirewall(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Firewall still exists")
		}
	}

	return nil
}

func CivoVPCFirewallConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_vpc_firewall" "foobar" {
	name = "%s"
	create_default_rules = true
	region = "LOCAL"
}`, name)
}

func CivoVPCFirewallConfigWithIngressEgress(name string) string {
	return fmt.Sprintf(`
resource "civo_vpc_firewall" "foobar" {
	name = "%s"
	create_default_rules = false
	region = "LOCAL"

	ingress_rule {
		label = "www https"
		protocol = "tcp"
		port_range = "443"
		cidr = ["192.168.1.1/32", "192.168.10.4/32"]
		action = "allow"
	  }

	  egress_rule {
		label = "ssh"
		protocol = "tcp"
		port_range = "22"
		cidr = ["192.168.1.1/32", "192.168.10.4/32", "192.168.10.10/32"]
		action = "allow"
	  }
}`, name)
}

func CivoVPCFirewallConfigUpdates(name string) string {
	return fmt.Sprintf(`
resource "civo_vpc_firewall" "foobar" {
	name = "%s"
	create_default_rules = true
	region = "LOCAL"
}`, name)
}
