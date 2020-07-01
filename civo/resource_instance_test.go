package civo

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// example.Widget represents a concrete Go type that represents an API resource
func TestAccCivoInstance_basic(t *testing.T) {
	var instance civogo.Instance

	// generate a random name for each test run
	resName := "civo_instance.foobar"
	var instanceHostname = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoInstanceDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoInstanceConfigBasic(instanceHostname),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoInstanceResourceExists(resName, &instance),
					// verify remote values
					testAccCheckCivoInstanceValues(&instance, instanceHostname),
					// verify local values
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g2.xsmall"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "pseudo_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func TestAccCivoInstanceSize_update(t *testing.T) {
	var instance civogo.Instance

	// generate a random name for each test run
	resName := "civo_instance.foobar"
	var instanceHostname = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoInstanceConfigBasic(instanceHostname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoInstanceResourceExists(resName, &instance),
					testAccCheckCivoInstanceValues(&instance, instanceHostname),
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g2.xsmall"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "pseudo_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoInstanceConfigUpdates(instanceHostname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoInstanceResourceExists(resName, &instance),
					testAccCheckCivoInstanceUpdated(&instance, instanceHostname),
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g2.large"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "pseudo_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func TestAccCivoInstanceNotes_update(t *testing.T) {
	var instance civogo.Instance

	// generate a random name for each test run
	resName := "civo_instance.foobar"
	var instanceHostname = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoInstanceConfigBasic(instanceHostname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoInstanceResourceExists(resName, &instance),
					testAccCheckCivoInstanceValues(&instance, instanceHostname),
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g2.xsmall"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttr(resName, "notes", ""),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "pseudo_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoInstanceConfigNotes(instanceHostname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoInstanceResourceExists(resName, &instance),
					testAccCheckCivoInstanceUpdated(&instance, instanceHostname),
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g2.xsmall"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttr(resName, "notes", "the_test_notes"),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "pseudo_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func TestAccCivoInstanceFirewall_update(t *testing.T) {
	var instance civogo.Instance
	var firewall civogo.Firewall

	// generate a random name for each test run
	resName := "civo_instance.foobar"
	var instanceHostname = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoInstanceConfigBasic(instanceHostname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoInstanceResourceExists(resName, &instance),
					testAccCheckCivoInstanceValues(&instance, instanceHostname),
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g2.xsmall"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "pseudo_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoInstanceConfigFirewall(instanceHostname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoInstanceResourceExists(resName, &instance),
					testAccCheckCivoFirewallResourceExists("civo_firewall.foobar", &firewall),
					testAccCheckCivoInstanceUpdated(&instance, instanceHostname),
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g2.xsmall"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttrSet(resName, "firewall_id"),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "pseudo_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func testAccCheckCivoInstanceValues(instance *civogo.Instance, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Hostname != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, instance.Hostname)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoInstanceResourceExists(n string, instance *civogo.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.GetInstance(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Instance not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*instance = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoInstanceUpdated(instance *civogo.Instance, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Hostname != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, instance.Hostname)
		}
		return nil
	}
}

func testAccCheckCivoInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_instance" {
			continue
		}

		_, err := client.GetInstance(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Instance still exists")
		}
	}

	return nil
}

func testAccCheckCivoInstanceConfigBasic(hostname string) string {
	return fmt.Sprintf(`
resource "civo_instance" "foobar" {
	hostname = "%s"
	size = "g2.xsmall"
}`, hostname)
}

func testAccCheckCivoInstanceConfigUpdates(hostname string) string {
	return fmt.Sprintf(`
resource "civo_instance" "foobar" {
	hostname = "%s"
	size = "g2.large"
}`, hostname)
}

func testAccCheckCivoInstanceConfigNotes(hostname string) string {
	return fmt.Sprintf(`
resource "civo_instance" "foobar" {
	hostname = "%s"
	size = "g2.xsmall"
	notes = "the_test_notes"
}`, hostname)
}

func testAccCheckCivoInstanceConfigFirewall(hostname string) string {
	return fmt.Sprintf(`
resource "civo_firewall" "foobar" {
	name = "fw-foobar"
}

resource "civo_instance" "foobar" {
	hostname = "%s"
	size = "g2.xsmall"
	firewall_id = civo_firewall.foobar.id
}`, hostname)
}
