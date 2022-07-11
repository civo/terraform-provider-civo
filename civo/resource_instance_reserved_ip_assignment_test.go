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
func TestAccCivoInstanceReservedIPAssignment_basic(t *testing.T) {
	var ip civogo.IP
	var instance civogo.Instance

	// generate a random name for each test run
	resName := "civo_instance_reserved_ip_assignment.foobar"
	var AttachmentName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCivoInstanceReservedIPAssignmentConfigBasic(AttachmentName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoReservedIPResourceExists("civo_reserved_ip.foo", &ip),
					testAccCheckCivoInstanceResourceExists("civo_instance.vm", &instance),
					// verify local values
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "instance_id"),
					resource.TestCheckResourceAttrSet(resName, "reserved_ip_id"),
				),
			},
		},
	})
}

func testAccCivoInstanceReservedIPAssignmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_instance_reserved_ip_assignment" {
			continue
		}
	}

	return nil
}

func testAccCivoInstanceReservedIPAssignmentConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_instance" "vm" {
	hostname = "instance-%s"
}

resource "civo_reserved_ip" "foo" {
	name = "%s"
}

resource "civo_instance_reserved_ip_assignment" "foobar" {
	instance_id = civo_instance.vm.id
	reserved_ip_id  = civo_instance.foo.id
}
`, name, name)
}
