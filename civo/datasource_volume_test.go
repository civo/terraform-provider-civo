package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceCivoVolume_basic(t *testing.T) {
	datasourceName := "data.civo_volume.foobar"
	name := acctest.RandomWithPrefix("ds-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoVolumeConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
				),
			},
		},
	})
}

func testAccDataSourceCivoVolumeConfig(name string) string {
	return fmt.Sprintf(`
resource "civo_volume" "foobar" {
	name = "%s"
	size_gb = 60
	bootable = false
}

data "civo_volume" "foobar" {
	name = civo_volume.foobar.name
}
`, name)
}
