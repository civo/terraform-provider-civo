package objectstorage_test

import (
	"fmt"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func DataSourceCivoObjectStore_basic(t *testing.T) {
	datasourceName := "data.civo_object_store.foobar"
	name := acctest.RandomWithPrefix("objectstore")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoObjectStoreConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
					resource.TestCheckResourceAttrSet(datasourceName, "max_size_gb"),
					resource.TestCheckResourceAttrSet(datasourceName, "region"),
					resource.TestCheckResourceAttrSet(datasourceName, "bucket_url"),
					resource.TestCheckResourceAttrSet(datasourceName, "access_key_id"),
					resource.TestCheckResourceAttr(datasourceName, "status", "ready"),
				),
			},
		},
	})
}

func DataSourceCivoObjectStoreConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_object_store" "foobar" {
	name = "%s"
	max_size_gb = 500
	region = "FAKE"
}

data "civo_object_store" "foobar" {
	name = civo_object_store.foobar.name
}
`, name)
}
