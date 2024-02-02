package database_test

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// CivoDatabase_basic is used to test the database resource
func TestAccCivoDatabase_basic(t *testing.T) {
	var database civogo.Database

	// generate a random name for each test run
	resName := "civo_database.foobar"
	var databaseName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: CivoDatabaseConfigBasic(databaseName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					CivoDatabaseResourceExists(resName, &database),
					// verify remote values
					CivoDatabaseValues(&database, databaseName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", databaseName),
					resource.TestCheckResourceAttrSet(resName, "size"),
					resource.TestCheckResourceAttrSet(resName, "nodes"),
					resource.TestCheckResourceAttrSet(resName, "engine"),
					resource.TestCheckResourceAttrSet(resName, "version"),
					resource.TestCheckResourceAttr(resName, "status", "Ready"),
				),
			},
		},
	})
}

// CivoDatabase_update is used to test the database resource
func TestAccCivoDatabase_update(t *testing.T) {
	var database civogo.Database

	// generate a random name for each test run
	resName := "civo_database.foobar"
	var databaseName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoDatabaseConfigBasic(databaseName),
				Check: resource.ComposeTestCheckFunc(
					CivoDatabaseResourceExists(resName, &database),
					CivoDatabaseValues(&database, databaseName),
					resource.TestCheckResourceAttr(resName, "name", databaseName),
					resource.TestCheckResourceAttr(resName, "status", "Ready"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: CivoDatabaseConfigUpdates(databaseName),
				Check: resource.ComposeTestCheckFunc(
					CivoDatabaseResourceExists(resName, &database),
					CivoDatabaseUpdated(&database, databaseName),
					resource.TestCheckResourceAttr(resName, "name", databaseName),
					resource.TestCheckResourceAttr(resName, "status", "Ready"),
				),
			},
		},
	})
}

// CivoDatabaseConfig is used to configure the database resource
func CivoDatabaseValues(database *civogo.Database, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if database.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, database.Name)
		}
		return nil
	}
}

// CivoDatabaseResourceExists - Check if the database resource exist
func CivoDatabaseResourceExists(n string, database *civogo.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := acceptance.TestAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindDatabase(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Database not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*database = *resp

		return nil
	}
}

// CivoDatabaseUpdated - Check if the database resource is updated
func CivoDatabaseUpdated(database *civogo.Database, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if database.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, database.Name)
		}
		return nil
	}
}

// CivoDatabaseDestroy is used to destroy the database created during the test
func CivoDatabaseDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*civogo.Client)

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

// CivoDatabaseConfigBasic is used to configure the database resource
func CivoDatabaseConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_database" "foobar" {
	name = "%s"
	size = "g3.db.xsmall"
	engine = "Postgres"
	version = "13"
	nodes = 2
}`, name)
}

// CivoDatabaseConfigUpdates is used to configure the database resource
func CivoDatabaseConfigUpdates(name string) string {
	return fmt.Sprintf(`
resource "civo_database" "foobar" {
	name = "%s"
	size = "g3.db.xsmall"
	engine = "Postgres"
	version = "13"
	nodes = 2
}`, name)
}
