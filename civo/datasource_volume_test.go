package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
					resource.TestCheckResourceAttrSet(datasourceName, "size_gb"),
				),
			},
		},
	})
}

func testAccDataSourceCivoVolumeConfig(name string) string {
	return fmt.Sprintf(`
data "civo_network" "default" {
	label = "default"
	region = "LON1"
}

resource "civo_volume" "newvolume" {
	name = "%s"
	size_gb = 10
	network_id = data.civo_network.default.id
}

data "civo_volume" "foobar" {
	name = civo_volume.newvolume.name
}
`, name)
}
