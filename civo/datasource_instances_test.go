package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoInstances_basic(t *testing.T) {
	var instanceHostname = acctest.RandomWithPrefix("tf-test")
	var instanceHostname2 = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoInstancesConfig(instanceHostname, instanceHostname2),
			},
			{
				Config: testAccDataSourceCivoInstancesDataConfig(instanceHostname, instanceHostname2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.civo_instances.result", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.civo_instances.result", "instances.0.hostname", instanceHostname),
					resource.TestCheckResourceAttrPair("data.civo_instances.result", "instances.0.id", "civo_instance.foo", "id"),
				),
			},
		},
	})
}

func testAccDataSourceCivoInstancesConfig(name string, name2 string) string {
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

resource "civo_instance" "foo" {
	hostname = "%s"
	size = element(data.civo_instances_size.small.sizes, 0).name
	disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}

resource "civo_instance" "bar" {
	hostname = "%s"
	size = element(data.civo_instances_size.small.sizes, 0).name
	disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}
`, name, name2)
}

func testAccDataSourceCivoInstancesDataConfig(name string, name2 string) string {
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

resource "civo_instance" "foo" {
	hostname = "%s"
	size = element(data.civo_instances_size.small.sizes, 0).name
	disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}

resource "civo_instance" "bar" {
	hostname = "%s"
	size = element(data.civo_instances_size.small.sizes, 0).name
	disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}

data "civo_instances" "result" {
    filter {
        key = "hostname"
        values = ["%s"]
    }
}
`, name, name2, name)
}
