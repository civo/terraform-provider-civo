package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoApplication_basic(t *testing.T) {
	datasourceName := "data.civo_application.foobar"
	name := acctest.RandomWithPrefix("app")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoAppConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
					resource.TestCheckResourceAttr(datasourceName, "process_info.0.process_type", "web"),
					resource.TestCheckResourceAttr(datasourceName, "process_info.0.process_count", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceCivoAppByID_basic(t *testing.T) {
	datasourceName := "data.civo_application.foobar"
	name := acctest.RandomWithPrefix("app")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoAppByIDConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
					resource.TestCheckResourceAttr(datasourceName, "process_info.0.process_type", "web"),
					resource.TestCheckResourceAttr(datasourceName, "process_info.0.process_count", "1"),
				),
			},
		},
	})
}

func testAccDataSourceCivoAppConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_application" "my-app" {
	name = "%s"
	process_info {
		process_type = web
		process_count = 1
	}
}

data "civo_application" "foobar" {
	name = civo_application.my-app.name
}
`, name)
}

func testAccDataSourceCivoAppByIDConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_application" "my-app" {
	name = "%s"
	process_info {
		process_type = web
		process_count = 1
	}
}

data "civo_application" "foobar" {
	id = civo_application.my-app.id
}
`, name)
}
