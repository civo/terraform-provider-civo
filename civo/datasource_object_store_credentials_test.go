package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoObjectStoreCredential_basic(t *testing.T) {
	datasourceName := "data.civo_object_store_credential.foobar"
	name := acctest.RandomWithPrefix("objectstorecredential")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoObjectStoreCredentialConfig(name),
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

func testAccDataSourceCivoObjectStoreCredentialConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_object_store_credential" "foobar" {
	name = "%s"
}

data "civo_object_store_credential" "foobar" {
	name = civo_object_store_credential.foobar.name
}
`, name)
}
