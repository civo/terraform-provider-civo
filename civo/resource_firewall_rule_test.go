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
func TestAccCivoFirewallRule_basic(t *testing.T) {
	var firewallRule civogo.FirewallRule

	// generate a random name for each test run
	resName := "civo_firewall_rule.testrule"
	var firewalName = acctest.RandomWithPrefix("tf-fw-rule")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoFirewallRuleConfigBasic(firewalName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoFirewallRuleResourceExists(resName, &firewallRule),
					// verify remote values
					testAccCheckCivoFirewallRuleValues(&firewallRule),
					// verify local values
					resource.TestCheckResourceAttr(resName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resName, "start_port", "80"),
				),
			},
		},
	})
}

func TestAccCivoFirewallRule_update(t *testing.T) {
	var firewallRule civogo.FirewallRule

	// generate a random name for each test run
	resName := "civo_firewall_rule.testrule"
	var firewalName = acctest.RandomWithPrefix("tf-fw-rule")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoFirewallRuleConfigBasic(firewalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoFirewallRuleResourceExists(resName, &firewallRule),
					resource.TestCheckResourceAttr(resName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resName, "start_port", "80"),
					resource.TestCheckResourceAttr(resName, "label", "web"),
					resource.TestCheckResourceAttr(resName, "action", "allow"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoFirewallRuleConfigUpdates(firewalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoFirewallRuleResourceExists(resName, &firewallRule),
					testAccCheckCivoFirewallRuleUpdated(&firewallRule),
					resource.TestCheckResourceAttr(resName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resName, "start_port", "443"),
					resource.TestCheckResourceAttr(resName, "label", "web_server"),
					resource.TestCheckResourceAttr(resName, "action", "allow"),
				),
			},
		},
	})
}

func testAccCheckCivoFirewallRuleValues(firewall *civogo.FirewallRule) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if firewall.Protocol != "tcp" {
			return fmt.Errorf("bad protocol, expected \"%s\", got: %#v", "tcp", firewall.Protocol)
		}
		if firewall.StartPort != "80" {
			return fmt.Errorf("bad port, expected \"%s\", got: %#v", "80", firewall.StartPort)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoFirewallRuleResourceExists(n string, firewall *civogo.FirewallRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindFirewallRule(rs.Primary.Attributes["firewall_id"], rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Firewall rule not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*firewall = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoFirewallRuleUpdated(firewall *civogo.FirewallRule) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if firewall.Protocol != "tcp" {
			return fmt.Errorf("bad protocol, expected \"%s\", got: %#v", "tcp", firewall.Protocol)
		}
		if firewall.StartPort != "443" {
			return fmt.Errorf("bad port, expected \"%s\", got: %#v", "443", firewall.StartPort)
		}
		return nil
	}
}

func testAccCheckCivoFirewallRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_firewall_rule" {
			continue
		}

		_, err := client.FindFirewall(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Firewall rule still exists")
		}
	}

	return nil
}

func testAccCheckCivoFirewallRuleConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_firewall" "foobar" {
	name = "%s"
}

resource "civo_firewall_rule" "testrule" {
	firewall_id = civo_firewall.foobar.id
	protocol = "tcp"
	start_port = "80"
	end_port = "80"
	cidr = ["192.168.1.2/32"]
	direction = "ingress"
	action = "allow"
	label = "web"
}

`, name)
}

func testAccCheckCivoFirewallRuleConfigUpdates(name string) string {
	return fmt.Sprintf(`
resource "civo_firewall" "foobar" {
	name = "%s"
}

resource "civo_firewall_rule" "testrule" {
	firewall_id = civo_firewall.foobar.id
	protocol = "tcp"
	start_port = "443"
	end_port = "443"
	cidr = ["192.168.1.2/32"]
	direction = "ingress"
	action = "allow"
	label = "web_server"
}
`, name)
}
