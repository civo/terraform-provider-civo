package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceCivoDNSDomainRecord_basic(t *testing.T) {
	datasourceName := "data.civo_dns_domain_record.record"
	domain := acctest.RandomWithPrefix("recordtest") + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoDNSDomainRecordConfigBasic(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", "www"),
					resource.TestCheckResourceAttr(datasourceName, "type", "a"),
					resource.TestCheckResourceAttr(datasourceName, "ttl", "600"),
					resource.TestCheckResourceAttr(datasourceName, "value", "192.168.1.1"),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "domain_id"),
				),
			},
		},
	})
}

func testAccDataSourceCivoDNSDomainRecordConfigBasic(domain string) string {
	return fmt.Sprintf(`
resource "civo_dns_domain_name" "domain" {
	name = "%[1]s"
}

resource "civo_dns_domain_record" "record" {
	domain_id = civo_dns_domain_name.domain.id
    type = "a"
    name = "www"
    value = "192.168.1.1"
    ttl = 600
}

data "civo_dns_domain_record" "record" {
	domain_id = civo_dns_domain_name.domain.id
	name = civo_dns_domain_record.record.name
}
`, domain)
}
