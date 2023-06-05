package civo

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCivoDatabaseVersion_basic(t *testing.T) {
	datasourceName := "data.civo_database_version.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoDatabaseVersionConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataSourceDatabaseVersionExist(datasourceName),
				),
			},
		},
	})
}

func TestAccDataSourceCivoDatabaseVersion_WithFilterAndSort(t *testing.T) {
	datasourceName := "data.civo_database_version.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoDatabaseVersionWhitFilterAndSort(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDatabaseVersionExist(datasourceName),
					testAccCheckDataSourceCivoDatabaseVersionFilteredAndSorted(datasourceName),
				),
			},
		},
	})
}

func testAccCheckDataSourceDatabaseVersionExist(n string) resource.TestCheckFunc {
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

func testAccCheckDataSourceCivoDatabaseVersionFilteredAndSorted(n string) resource.TestCheckFunc {
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

func testAccDataSourceCivoDatabaseVersionConfig() string {
	return `
data "civo_database_version" "foobar" {
}
`
}

func testAccDataSourceCivoDatabaseVersionWhitFilterAndSort() string {
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
