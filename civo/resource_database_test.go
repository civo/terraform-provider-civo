package civo

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCivoDatabase_basic(t *testing.T) {
	var database civogo.Database

	// generate a random name for each test run
	resName := "civo_database.foobar"
	var databaseName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoDatabaseConfigBasic(databaseName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoDatabaseResourceExists(resName, &database),
					// verify remote values
					testAccCheckCivoDatabaseValues(&database, databaseName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", databaseName),
					resource.TestCheckResourceAttrSet(resName, "size"),
					resource.TestCheckResourceAttrSet(resName, "nodes"),
					resource.TestCheckResourceAttr(resName, "status", "ready"),
				),
			},
		},
	})
}

func TestAccCivoDatabase_update(t *testing.T) {
	var database civogo.Database

	// generate a random name for each test run
	resName := "civo_database.foobar"
	var databaseName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoDatabaseConfigBasic(databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoDatabaseResourceExists(resName, &database),
					testAccCheckCivoDatabaseValues(&database, databaseName),
					resource.TestCheckResourceAttr(resName, "name", databaseName),
					resource.TestCheckResourceAttr(resName, "status", "ready"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoDatabaseConfigUpdates(databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoDatabaseResourceExists(resName, &database),
					testAccCheckCivoDatabaseUpdated(&database, databaseName),
					resource.TestCheckResourceAttr(resName, "name", databaseName),
					resource.TestCheckResourceAttr(resName, "status", "ready"),
				),
			},
		},
	})
}

func testAccCheckCivoDatabaseValues(database *civogo.Database, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if database.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, database.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoDatabaseResourceExists(n string, database *civogo.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindDatabase(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Database not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*database = *resp

		return nil
	}
}

func testAccCheckCivoDatabaseUpdated(database *civogo.Database, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if database.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, database.Name)
		}
		return nil
	}
}

func testAccCheckCivoDatabaseDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_database" {
			continue
		}

		_, err := client.FindDatabase(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Database still exists")
		}
	}

	return nil
}

func testAccCheckCivoDatabaseConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_database" "foobar" {
	name = "%s"
	size = "g4s.kube.small"
	nodes = 2
}`, name)
}

func testAccCheckCivoDatabaseConfigUpdates(name string) string {
	return fmt.Sprintf(`
resource "civo_database" "foobar" {
	name = "%s"
	size = "g4s.kube.small"
	nodes = 2
}`, name)
}
