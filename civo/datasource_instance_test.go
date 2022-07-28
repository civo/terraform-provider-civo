package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoInstance_basic(t *testing.T) {
	datasourceName := "data.civo_instance.foobar"
	name := acctest.RandomWithPrefix("instance") + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoInstanceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "hostname", name),
					resource.TestCheckResourceAttrSet(datasourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(datasourceName, "public_ip"),
				),
			},
		},
	})
}

func TestAccDataSourceCivoInstanceByID_basic(t *testing.T) {
	datasourceName := "data.civo_instance.foobar"
	name := acctest.RandomWithPrefix("instance") + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoInstanceByIDConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "hostname", name),
					resource.TestCheckResourceAttrSet(datasourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(datasourceName, "public_ip"),
				),
			},
		},
	})
}

func testAccDataSourceCivoInstanceConfig(name string) string {
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

resource "civo_instance" "vm" {
	hostname = "%s"
	size = element(data.civo_instances_size.small.sizes, 0).name
    disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}

data "civo_instance" "foobar" {
	hostname = civo_instance.vm.hostname
}
`, name)
}

func testAccDataSourceCivoInstanceByIDConfig(name string) string {
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

resource "civo_instance" "vm" {
	hostname = "%s"
	size = element(data.civo_instances_size.small.sizes, 0).name
	disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}

data "civo_instance" "foobar" {
	id = civo_instance.vm.id
}
`, name)
}
