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
func TestAccCivoDNSDomainNameRecord_basic(t *testing.T) {
	var domainRecord civogo.DNSRecord

	// generate a random name for each test run
	resName := "civo_dns_domain_record.www"
	var domainName = acctest.RandomWithPrefix("tf-test-record") + ".example"
	var recordName = acctest.RandomWithPrefix("record")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoDNSDomainNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoDNSDomainNameRecordConfigBasic(domainName, recordName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoDNSDomainNameRecordResourceExists(resName, &domainRecord),
					// verify remote values
					testAccCheckCivoDNSDomainNameRecordValues(&domainRecord, recordName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", recordName),
				),
			},
		},
	})
}

func TestAccCivoDNSDomainNameRecord_update(t *testing.T) {
	var domainRecord civogo.DNSRecord

	// generate a random name for each test run
	resName := "civo_dns_domain_record.www"
	var domainName = acctest.RandomWithPrefix("renamed-tf-test-record") + ".example"
	var recordName = acctest.RandomWithPrefix("renamed-record")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoDNSDomainNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoDNSDomainNameRecordConfigUpdates(domainName, recordName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoDNSDomainNameRecordResourceExists(resName, &domainRecord),
					testAccCheckCivoDNSDomainNameRecordValues(&domainRecord, recordName),
					resource.TestCheckResourceAttr(resName, "name", recordName),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoDNSDomainNameRecordConfigUpdates(domainName, recordName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoDNSDomainNameRecordResourceExists(resName, &domainRecord),
					testAccCheckCivoDNSDomainNameRecordUpdated(&domainRecord, recordName),
					resource.TestCheckResourceAttr(resName, "name", recordName),
				),
			},
		},
	})
}

func testAccCheckCivoDNSDomainNameRecordValues(domainRecord *civogo.DNSRecord, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if domainRecord.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, domainRecord.Name)
		}
		return nil
	}
}

func testAccCheckCivoDNSDomainNameRecordResourceExists(n string, domainRecord *civogo.DNSRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.GetDNSRecord(rs.Primary.Attributes["domain_id"], rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Domain record not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*domainRecord = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoDNSDomainNameRecordUpdated(domainRecord *civogo.DNSRecord, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if domainRecord.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, domainRecord.Name)
		}
		return nil
	}
}

func testAccCheckCivoDNSDomainNameRecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_dns_domain_record" {
			continue
		}

		_, err := client.GetDNSRecord(rs.Primary.Attributes["domain_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Domain record still exists")
		}
	}

	return nil
}

func testAccCheckCivoDNSDomainNameRecordConfigBasic(domain string, record string) string {
	return fmt.Sprintf(`
resource "civo_dns_domain_name" "foobar" {
	name = "%s"
}

resource "civo_dns_domain_record" "www" {
    domain_id = civo_dns_domain_name.foobar.id
    type = "a"
    name = "%s"
    value = "10.10.10.1"
    ttl = 600
}
`, domain, record)
}

func testAccCheckCivoDNSDomainNameRecordConfigUpdates(domain string, record string) string {
	return fmt.Sprintf(`
resource "civo_dns_domain_name" "foobar" {
	name = "%s"
}

resource "civo_dns_domain_record" "www" {
    domain_id = civo_dns_domain_name.foobar.id
    type = "a"
    name = "%s"
    value = "10.10.10.1"
    ttl = 600
}
`, domain, record)
}
