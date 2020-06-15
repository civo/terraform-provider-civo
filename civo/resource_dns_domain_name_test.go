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
func TestAccCivoDNSDomainName_basic(t *testing.T) {
	var domain civogo.DNSDomain

	// generate a random name for each test run
	resName := "civo_dns_domain_name.foobar"
	var domainName = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoDNSDomainNameDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoDNSDomainNameConfigBasic(domainName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoDNSDomainNameResourceExists(resName, &domain),
					// verify remote values
					testAccCheckCivoDNSDomainNameValues(&domain, domainName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", domainName),
				),
			},
		},
	})
}

func TestAccCivoDNSDomainName_update(t *testing.T) {
	var domain civogo.DNSDomain

	// generate a random name for each test run
	resName := "civo_dns_domain_name.foobar"
	var domainName = acctest.RandomWithPrefix("renamed-tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoDNSDomainNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoDNSDomainNameConfigUpdates(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoDNSDomainNameResourceExists(resName, &domain),
					testAccCheckCivoDNSDomainNameValues(&domain, domainName),
					resource.TestCheckResourceAttr(resName, "name", domainName),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoDNSDomainNameConfigUpdates(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoDNSDomainNameResourceExists(resName, &domain),
					testAccCheckCivoDNSDomainNameUpdated(&domain, domainName),
					resource.TestCheckResourceAttr(resName, "name", domainName),
				),
			},
		},
	})
}

func testAccCheckCivoDNSDomainNameValues(domain *civogo.DNSDomain, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if domain.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, domain.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoDNSDomainNameResourceExists(n string, domain *civogo.DNSDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindDNSDomain(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Domain not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*domain = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoDNSDomainNameUpdated(domain *civogo.DNSDomain, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if domain.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, domain.Name)
		}
		return nil
	}
}

func testAccCheckCivoDNSDomainNameDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_dns_domain_name" {
			continue
		}

		_, err := client.FindDNSDomain(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Domain still exists")
		}
	}

	return nil
}

func testAccCheckCivoDNSDomainNameConfigBasic(domain string) string {
	return fmt.Sprintf(`
resource "civo_dns_domain_name" "foobar" {
	name = "%s"
}`, domain)
}

func testAccCheckCivoDNSDomainNameConfigUpdates(domain string) string {
	return fmt.Sprintf(`
resource "civo_dns_domain_name" "foobar" {
	name = "%s"
}`, domain)
}
