package civo

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// example.Widget represents a concrete Go type that represents an API resource
func TestAccCivoVolumeAttachment_basic(t *testing.T) {
	var volume civogo.Volume
	var instance civogo.Instance

	// generate a random name for each test run
	resName := "civo_volume_attachment.foobar"
	var VolumeAttachmentName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoVolumeAttachmentConfigBasic(VolumeAttachmentName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoVolumeResourceExists("civo_volume.foo", &volume),
					testAccCheckCivoInstanceResourceExists("civo_instance.vm", &instance),
					// verify local values
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "instance_id"),
					resource.TestCheckResourceAttrSet(resName, "volume_id"),
				),
			},
		},
	})
}

func testAccCheckCivoVolumeAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_volume_attachment" {
			continue
		}
	}

	return nil
}

func testAccCheckCivoVolumeAttachmentConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_instance" "vm" {
	hostname = "instance-%s"
}

resource "civo_volume" "foo" {
	name = "%s"
	size_gb = 60
	bootable = false
}

resource "civo_volume_attachment" "foobar" {
	instance_id = civo_instance.vm.id
	volume_id  = civo_volume.foo.id
}
`, name, name)
}
