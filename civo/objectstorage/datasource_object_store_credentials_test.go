package objectstorage_test

import (
	"fmt"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoObjectStoreCredential_basic(t *testing.T) {
	datasourceName := "data.civo_object_store_credential.foobar"
	name := acctest.RandomWithPrefix("objectstorecredential")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoObjectStoreCredentialConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
					resource.TestCheckResourceAttrSet(datasourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "secret_access_key"),
					resource.TestCheckResourceAttr(datasourceName, "status", "ready"),
				),
			},
		},
	})
}

func DataSourceCivoObjectStoreCredentialConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_object_store_credential" "foobar" {
	name = "%s"
}

data "civo_object_store_credential" "foobar" {
	name = civo_object_store_credential.foobar.name
}
`, name)
}
