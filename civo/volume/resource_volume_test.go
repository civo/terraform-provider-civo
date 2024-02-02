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
func TestAccCivoVolume_basic(t *testing.T) {
	var volume civogo.Volume

	// generate a random name for each test run
	resName := "civo_volume.foobar"
	var VolumeName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoVolumeDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: CivoVolumeConfigBasic(VolumeName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					CivoVolumeResourceExists(resName, &volume),
					// verify remote values
					CivoVolumeValues(&volume, VolumeName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", VolumeName),
					resource.TestCheckResourceAttr(resName, "size_gb", "10"),
				),
			},
		},
	})
}

func CivoVolumeValues(volume *civogo.Volume, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if volume.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, volume.Name)
		}
		return nil
	}
}

// CheckExampleResourceExists queries the API and retrieves the matching Widget.
func CivoVolumeResourceExists(n string, volume *civogo.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := acceptance.TestAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindVolume(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Volume not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*volume = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func CivoVolumeDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_volume" {
			continue
		}

		_, err := client.FindVolume(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}

func CivoVolumeConfigBasic(name string) string {
	return fmt.Sprintf(`
data "civo_network" "default" {
	label = "default"
	region = "LON1"
}

resource "civo_volume" "foobar" {
	name = "%s"
	size_gb = 10
	network_id = data.civo_network.default.id
	region = "LON1"
}`, name)
}
