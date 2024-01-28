package database_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// DataSourceCivoDatabaseVersion_basic - Test the data source for database version
func DataSourceCivoDatabaseVersion_basic(t *testing.T) {
	datasourceName := "data.civo_database_version.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoDatabaseVersionConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					DataSourceDatabaseVersionExist(datasourceName),
				),
			},
		},
	})
}

// DataSourceCivoDatabaseVersion_WithFilterAndSort - Test the data source for database version with filter and sort
func DataSourceCivoDatabaseVersion_WithFilterAndSort(t *testing.T) {
	datasourceName := "data.civo_database_version.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoDatabaseVersionWhitFilterAndSort(),
				Check: resource.ComposeTestCheckFunc(
					DataSourceDatabaseVersionExist(datasourceName),
					DataSourceCivoDatabaseVersionFilteredAndSorted(datasourceName),
				),
			},
		},
	})
}

// DataSourceDatabaseVersionExist - Check if the data source exist
func DataSourceDatabaseVersionExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		rawTotal := rs.Primary.Attributes["versions.#"]
		total, err := strconv.Atoi(rawTotal)
		if err != nil {
			return err
		}

		if total < 1 {
			return fmt.Errorf("No Civo database version retrieved")
		}

		return nil
	}
}

// DataSourceCivoDatabaseVersionFilteredAndSorted - Check if the data source is filtered and sorted
func DataSourceCivoDatabaseVersionFilteredAndSorted(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		rawTotal := rs.Primary.Attributes["versions.#"]
		total, err := strconv.Atoi(rawTotal)
		if err != nil {
			return err
		}

		stringInSlice := func(value string, slice []string) bool {
			for _, item := range slice {
				if item == value {
					return true
				}
			}
			return false
		}

		for i := 0; i < total; i++ {
			name := rs.Primary.Attributes[fmt.Sprintf("versions.%d.engine", i)]
			if !stringInSlice(name, []string{"Mysql", "PostgreSQL"}) {
				return fmt.Errorf("engine is not in expected test filter values")
			}
		}

		return nil
	}
}

// DataSourceCivoDatabaseVersionConfig - Config for the data source
func DataSourceCivoDatabaseVersionConfig() string {
	return `
data "civo_database_version" "foobar" {
}
`
}

// DataSourceCivoDatabaseVersionWhitFilterAndSort - Config for the data source with filter and sort	
func DataSourceCivoDatabaseVersionWhitFilterAndSort() string {
	return `
data "civo_database_version" "foobar" {
	filter {
        key = "engine"
        values = ["mysql"]
	}
	
	sort {
        key = "engine"
        direction = "desc"
    }
}
`
}
