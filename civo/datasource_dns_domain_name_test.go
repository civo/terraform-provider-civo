package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceCivoDnsDomainName(t *testing.T) {
	datasourceName := "data.civo_dns_domain_name.domain"
	domain := acctest.RandomWithPrefix("domian") + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoDNSDomainName(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", domain),
				),
			},
		},
	})
}

func testAccDataSourceCivoDNSDomainName(domain string) string {
	return fmt.Sprintf(`
resource "civo_dns_domain_name" "domain" {
	name = "%[1]s"
}

data "civo_dns_domain_name" "domain" {
	name = civo_dns_domain_name.domain.name
}
`, domain)
}
