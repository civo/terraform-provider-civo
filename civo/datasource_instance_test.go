package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoInstance_basic(t *testing.T) {
	datasourceName := "data.civo_instance.foobar"
	name := acctest.RandomWithPrefix("instance") + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoInstanceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "hostname", name),
					resource.TestCheckResourceAttrSet(datasourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(datasourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(datasourceName, "pseudo_ip"),
				),
			},
		},
	})
}

func TestAccDataSourceCivoInstanceByID_basic(t *testing.T) {
	datasourceName := "data.civo_instance.foobar"
	name := acctest.RandomWithPrefix("instance") + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoInstanceByIDConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "hostname", name),
					resource.TestCheckResourceAttrSet(datasourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(datasourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(datasourceName, "pseudo_ip"),
				),
			},
		},
	})
}

func testAccDataSourceCivoInstanceConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_instance" "vm" {
	hostname = "%s"
}

data "civo_instance" "foobar" {
	hostname = civo_instance.vm.hostname
}
`, name)
}

func testAccDataSourceCivoInstanceByIDConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_instance" "vm" {
	hostname = "%s"
}

data "civo_instance" "foobar" {
	id = civo_instance.vm.id
}
`, name)
}
