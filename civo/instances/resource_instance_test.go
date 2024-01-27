package instances_test

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
func CivoInstance_basic(t *testing.T) {
	var instance civogo.Instance

	// generate a random name for each test run
	resName := "civo_instance.foobar"
	var instanceHostname = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CivoInstanceDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: CivoInstanceConfigBasic(instanceHostname),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					acceptance.CivoInstanceResourceExists(resName, &instance),
					// verify remote values
					CivoInstanceValues(&instance, instanceHostname),
					// verify local values
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g3.small"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttr(resName, "cpu_cores", "1"),
					resource.TestCheckResourceAttr(resName, "ram_mb", "2048"),
					resource.TestCheckResourceAttr(resName, "disk_gb", "25"),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func CivoInstanceSize_update(t *testing.T) {
	var instance civogo.Instance

	// generate a random name for each test run
	resName := "civo_instance.foobar"
	var instanceHostname = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CivoInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoInstanceConfigBasic(instanceHostname),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CivoInstanceResourceExists(resName, &instance),
					CivoInstanceValues(&instance, instanceHostname),
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g3.small"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttr(resName, "cpu_cores", "1"),
					resource.TestCheckResourceAttr(resName, "ram_mb", "2048"),
					resource.TestCheckResourceAttr(resName, "disk_gb", "25"),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: CivoInstanceConfigUpdates(instanceHostname),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CivoInstanceResourceExists(resName, &instance),
					CivoInstanceUpdated(&instance, instanceHostname),
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g3.medium"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttr(resName, "cpu_cores", "2"),
					resource.TestCheckResourceAttr(resName, "ram_mb", "4096"),
					resource.TestCheckResourceAttr(resName, "disk_gb", "50"),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func CivoInstanceNotes_update(t *testing.T) {
	var instance civogo.Instance

	// generate a random name for each test run
	resName := "civo_instance.foobar"
	var instanceHostname = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CivoInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoInstanceConfigBasic(instanceHostname),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CivoInstanceResourceExists(resName, &instance),
					CivoInstanceValues(&instance, instanceHostname),
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g3.small"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttr(resName, "cpu_cores", "1"),
					resource.TestCheckResourceAttr(resName, "ram_mb", "2048"),
					resource.TestCheckResourceAttr(resName, "disk_gb", "25"),
					resource.TestCheckResourceAttr(resName, "notes", ""),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: CivoInstanceConfigNotes(instanceHostname),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CivoInstanceResourceExists(resName, &instance),
					CivoInstanceUpdated(&instance, instanceHostname),
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g3.small"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttr(resName, "cpu_cores", "1"),
					resource.TestCheckResourceAttr(resName, "ram_mb", "2048"),
					resource.TestCheckResourceAttr(resName, "disk_gb", "25"),
					resource.TestCheckResourceAttr(resName, "notes", "the_test_notes"),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func CivoInstanceFirewall_update(t *testing.T) {
	var instance civogo.Instance

	// generate a random name for each test run
	resName := "civo_instance.foobar"
	var instanceHostname = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CivoInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoInstanceConfigBasic(instanceHostname),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CivoInstanceResourceExists(resName, &instance),
					CivoInstanceValues(&instance, instanceHostname),
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g3.small"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttr(resName, "cpu_cores", "1"),
					resource.TestCheckResourceAttr(resName, "ram_mb", "2048"),
					resource.TestCheckResourceAttr(resName, "disk_gb", "25"),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: CivoInstanceConfigFirewall(instanceHostname),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CivoInstanceResourceExists(resName, &instance),
					CivoInstanceUpdated(&instance, instanceHostname),
					resource.TestCheckResourceAttr(resName, "hostname", instanceHostname),
					resource.TestCheckResourceAttr(resName, "size", "g3.small"),
					resource.TestCheckResourceAttr(resName, "initial_user", "civo"),
					resource.TestCheckResourceAttr(resName, "cpu_cores", "1"),
					resource.TestCheckResourceAttr(resName, "ram_mb", "2048"),
					resource.TestCheckResourceAttr(resName, "disk_gb", "25"),
					resource.TestCheckResourceAttrSet(resName, "firewall_id"),
					resource.TestCheckResourceAttrSet(resName, "initial_password"),
					resource.TestCheckResourceAttrSet(resName, "private_ip"),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func CivoInstanceValues(instance *civogo.Instance, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if instance.Hostname != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, instance.Hostname)
		}
		return nil
	}
}

func CivoInstanceUpdated(instance *civogo.Instance, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if instance.Hostname != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, instance.Hostname)
		}
		return nil
	}
}

func CivoInstanceConfigBasic(hostname string) string {
	return fmt.Sprintf(`
data "civo_size" "small" {
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

resource "civo_instance" "foobar" {
	hostname = "%s"
	region = "FAKE"
	size = element(data.civo_size.small.sizes, 0).name
    disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}`, hostname)
}

func CivoInstanceConfigUpdates(hostname string) string {
	return fmt.Sprintf(`
data "civo_size" "medium" {
	filter {
		key = "name"
		values = ["g3.medium"]
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

resource "civo_instance" "foobar" {
	hostname = "%s"
	size = element(data.civo_instances_size.medium.sizes, 0).name
	disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}`, hostname)
}

func CivoInstanceConfigNotes(hostname string) string {
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
resource "civo_instance" "foobar" {
	hostname = "%s"
	size = element(data.civo_instances_size.small.sizes, 0).name
	disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
	notes = "the_test_notes"
}`, hostname)
}

func CivoInstanceConfigFirewall(hostname string) string {
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

resource "civo_firewall" "foobar" {
	name = "fw-foobar"
}

resource "civo_instance" "foobar" {
	hostname = "%s"
	size = element(data.civo_instances_size.small.sizes, 0).name
	disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
	firewall_id = civo_firewall.foobar.id
}`, hostname)
}
