package volume_test

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// example.Widget represents a concrete Go type that represents an API resource
func CivoVolumeAttachment_basic(t *testing.T) {
	var volume civogo.Volume
	var instance civogo.Instance

	// generate a random name for each test run
	resName := "civo_volume_attachment.foobar"
	var VolumeAttachmentName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: CivoVolumeAttachmentConfigBasic(VolumeAttachmentName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					CivoVolumeResourceExists("civo_volume.foo", &volume),
					acceptance.CivoInstanceResourceExists("civo_instance.vm", &instance),
					// verify local values
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "instance_id"),
					resource.TestCheckResourceAttrSet(resName, "volume_id"),
				),
			},
		},
	})
}

func CivoVolumeAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_volume_attachment" {
			continue
		}
	}

	return nil
}

func CivoVolumeAttachmentConfigBasic(name string) string {
	return fmt.Sprintf(`
data "civo_instances_size" "small" {
	filter {
		key = "name"
		values = ["g3.small"]
		match_by = "re"
	}

	filter {
		key = "type"
		values = ["instance"]
	}

}

# Query instance disk image
data "civo_disk_image" "debian" {
	filter {
		key = "name"
		values = ["debian-10"]
	}
}

data "civo_network" "default" {
	label = "default"
	region = "LON1"
}

resource "civo_instance" "vm" {
	hostname = "instance-%s"
	size = element(data.civo_instances_size.small.sizes, 0).name
	disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}

resource "civo_volume" "foo" {
	name = "%s"
	size_gb = 10
	network_id = data.civo_network.default.id
}

resource "civo_volume_attachment" "foobar" {
	instance_id = civo_instance.vm.id
	volume_id  = civo_volume.foo.id
	region = "LON1"
}
`, name, name)
}
