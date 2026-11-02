package volume_test

import (
	"fmt"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoVolume_basic(t *testing.T) {
	datasourceName := "data.civo_volume.foobar"
	name := acctest.RandomWithPrefix("ds-test")
	volumeType := "encrypted-standard"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoVolumeConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
					resource.TestCheckResourceAttrSet(datasourceName, "size_gb"),
					resource.TestCheckResourceAttr(datasourceName, "volume_type", volumeType),
				),
			},
		},
	})
}

func DataSourceCivoVolumeConfig(name string) string {
	return fmt.Sprintf(`
data "civo_network" "default" {
	label = "default"
	region = "LON1"
}

resource "civo_volume" "newvolume" {
	name = "%s"
	size_gb = 10
	network_id = data.civo_network.default.id
	volume_type = "encrypted-standard"
}

data "civo_volume" "foobar" {
	name = civo_volume.newvolume.name
}
`, name)
}
