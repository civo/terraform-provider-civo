package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoObjectStore_basic(t *testing.T) {
	datasourceName := "data.civo_object_store.foobar"
	name := acctest.RandomWithPrefix("objectstore")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoObjectStoreConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
					resource.TestCheckResourceAttrSet(datasourceName, "max_size_gb"),
					resource.TestCheckResourceAttrSet(datasourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "secret_access_key"),
					resource.TestCheckResourceAttrSet(datasourceName, "endpoint"),
					resource.TestCheckResourceAttr(datasourceName, "status", "ready"),
				),
			},
		},
	})
}

func testAccDataSourceCivoObjectStoreConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_object_store" "foobar" {
	name = "%s"
	max_size_gb = 500
}

data "civo_object_store" "foobar" {
	name = civo_object_store.foobar.name
}
`, name)
}
