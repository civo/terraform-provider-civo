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
func TestAccCivoSnapshot_basic(t *testing.T) {
	var snapshot civogo.Snapshot

	// generate a random name for each test run
	resName := "civo_snapshot.foobar"
	var snapshotName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoSnapshotConfigBasic(snapshotName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoSnapshotResourceExists(resName, &snapshot),
					// verify remote values
					testAccCheckCivoSnapshotValues(&snapshot, snapshotName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", snapshotName),
				),
			},
		},
	})
}

func TestAccCivoSnapshot_update(t *testing.T) {
	var snapshot civogo.Snapshot

	// generate a random name for each test run
	resName := "civo_snapshot.foobar"
	var snapshotName = acctest.RandomWithPrefix("tf-test")
	var snapshotNameUpdate = acctest.RandomWithPrefix("rename-tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoSnapshotConfigBasic(snapshotName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoSnapshotResourceExists(resName, &snapshot),
					testAccCheckCivoSnapshotValues(&snapshot, snapshotName),
					resource.TestCheckResourceAttr(resName, "name", snapshotName),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoSnapshotConfigUpdates(snapshotNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoSnapshotResourceExists(resName, &snapshot),
					testAccCheckCivoSnapshotUpdated(&snapshot, snapshotNameUpdate),
					resource.TestCheckResourceAttr(resName, "name", snapshotNameUpdate),
				),
			},
		},
	})
}

func testAccCheckCivoSnapshotValues(snapshot *civogo.Snapshot, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if snapshot.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, snapshot.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoSnapshotResourceExists(n string, snapshot *civogo.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindSnapshot(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Snapshot not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*snapshot = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoSnapshotUpdated(snapshot *civogo.Snapshot, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if snapshot.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, snapshot.Name)
		}
		return nil
	}
}

func testAccCheckCivoSnapshotDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_snapshot" {
			continue
		}

		_, err := client.FindSnapshot(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Snapshot still exists")
		}
	}

	return nil
}

func testAccCheckCivoSnapshotConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_instance" "vm" {
	hostname = "instance-%s"
}

resource "civo_snapshot" "foobar" {
	name = "%s"
	instance_id = civo_instance.vm.id
}`, name, name)
}

func testAccCheckCivoSnapshotConfigUpdates(name string) string {
	return fmt.Sprintf(`
resource "civo_instance" "vm" {
	hostname = "instance-%s"
}

resource "civo_snapshot" "foobar" {
	name = "%s"
	instance_id = civo_instance.vm.id
}`, name, name)
}
