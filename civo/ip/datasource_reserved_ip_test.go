package ip_test

import (
	"fmt"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func DataSourceReservedIP_basic(t *testing.T) {
	datasourceName := "data.civo_reserved_ip.foobar"
	name := acctest.RandomWithPrefix("ip-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceReservedIPConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
					resource.TestCheckResourceAttrSet(datasourceName, "ip"),
					resource.TestCheckResourceAttrSet(datasourceName, "region"),
				),
			},
		},
	})
}

func DataSourceReservedIPConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_reserved_ip" "newip" {
	name = "%s"
	region = "LON1"
}

data "civo_reserved_ip" "foobar" {
	name = civo_reserved_ip.newip.name
}
`, name)
}
