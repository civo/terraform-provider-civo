package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoDatabase_basic(t *testing.T) {
	datasourceName := "data.civo_database.foobar"
	name := acctest.RandomWithPrefix("database")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoDatabaseConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
					resource.TestCheckResourceAttrSet(datasourceName, "size"),
					resource.TestCheckResourceAttrSet(datasourceName, "nodes"),
					resource.TestCheckResourceAttrSet(datasourceName, "engine"),
					resource.TestCheckResourceAttrSet(datasourceName, "version"),
					resource.TestCheckResourceAttr(datasourceName, "status", "Ready"),
				),
			},
		},
	})
}

func testAccDataSourceCivoDatabaseConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_database" "foobar" {
	name = "%s"
	size = "g3.db.xsmall"
	engine = "Postgres"
	version = "13"
	nodes = 2
}

data "civo_database" "foobar" {
	name = civo_database.foobar.name
}`, name)
}
