package dns_test

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
func TestAccCivoDNSDomainName_basic(t *testing.T) {
	var domain civogo.DNSDomain

	// generate a random name for each test run
	resName := "civo_dns_domain_name.foobar"
	var domainName = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoDNSDomainNameDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: CivoDNSDomainNameConfigBasic(domainName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					CivoDNSDomainNameResourceExists(resName, &domain),
					// verify remote values
					CivoDNSDomainNameValues(&domain, domainName),
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
	var domainName = acctest.RandomWithPrefix("tf-test") + ".example"
	var domainNameUpdate = acctest.RandomWithPrefix("renamed-tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoDNSDomainNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoDNSDomainNameConfigBasic(domainName),
				Check: resource.ComposeTestCheckFunc(
					CivoDNSDomainNameResourceExists(resName, &domain),
					CivoDNSDomainNameValues(&domain, domainName),
					resource.TestCheckResourceAttr(resName, "name", domainName),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: CivoDNSDomainNameConfigUpdates(domainNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					CivoDNSDomainNameResourceExists(resName, &domain),
					CivoDNSDomainNameUpdated(&domain, domainNameUpdate),
					resource.TestCheckResourceAttr(resName, "name", domainNameUpdate),
				),
			},
		},
	})
}

func CivoDNSDomainNameValues(domain *civogo.DNSDomain, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if domain.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, domain.Name)
		}
		return nil
	}
}

// CheckExampleResourceExists queries the API and retrieves the matching Widget.
func CivoDNSDomainNameResourceExists(n string, domain *civogo.DNSDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := acceptance.TestAccProvider.Meta().(*civogo.Client)
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

func CivoDNSDomainNameUpdated(domain *civogo.DNSDomain, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if domain.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, domain.Name)
		}
		return nil
	}
}

func CivoDNSDomainNameDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*civogo.Client)

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

func CivoDNSDomainNameConfigBasic(domain string) string {
	return fmt.Sprintf(`
resource "civo_dns_domain_name" "foobar" {
	name = "%s"
}`, domain)
}

func CivoDNSDomainNameConfigUpdates(domain string) string {
	return fmt.Sprintf(`
resource "civo_dns_domain_name" "foobar" {
	name = "%s"
}`, domain)
}
