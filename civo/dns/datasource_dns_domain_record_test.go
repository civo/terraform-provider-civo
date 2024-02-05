package dns_test

import (
	"fmt"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourceCivoDNSDomainRecord_basic is a basic test case for a DNS domain record data source.
func TestAccDataSourceCivoDNSDomainRecord_basic(t *testing.T) {
	datasourceName := "data.civo_dns_domain_record.record"
	domain := acctest.RandomWithPrefix("recordtest") + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoDNSDomainRecordConfigBasic(domain),
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

func DataSourceCivoDNSDomainRecordConfigBasic(domain string) string {
	return fmt.Sprintf(`
resource "civo_dns_domain_name" "domain" {
	name = "%[1]s"
}

resource "civo_dns_domain_record" "record" {
	domain_id = civo_dns_domain_name.domain.id
    type = "A"
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
