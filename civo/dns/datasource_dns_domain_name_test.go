package dns_test

import (
	"fmt"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccCivoDNSDomainNameDataSource_basic is a basic test case for a DNS Domain Name data source.
func TestAccDataSourceCivoDnsDomainName(t *testing.T) {
	datasourceName := "data.civo_dns_domain_name.domain"
	domain := acctest.RandomWithPrefix("domian") + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoDNSDomainName(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", domain),
				),
			},
		},
	})
}

func DataSourceCivoDNSDomainName(domain string) string {
	return fmt.Sprintf(`
resource "civo_dns_domain_name" "domain" {
	name = "%[1]s"
}

data "civo_dns_domain_name" "domain" {
	name = civo_dns_domain_name.domain.name
}
`, domain)
}
