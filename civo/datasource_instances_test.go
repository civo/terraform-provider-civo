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
resource "civo_instance" "foo" {
	hostname = "%s"
}

resource "civo_instance" "bar" {
	hostname = "%s"
}
`, name, name2)
}

func testAccDataSourceCivoInstancesDataConfig(name string, name2 string) string {
	return fmt.Sprintf(`
resource "civo_instance" "foo" {
	hostname = "%s"
}

resource "civo_instance" "bar" {
	hostname = "%s"
}

data "civo_instances" "result" {
    filter {
        key = "hostname"
        values = ["%s"]
    }
}
`, name, name2, name)
}
