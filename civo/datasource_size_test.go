package civo

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCivoSize_basic(t *testing.T) {
	datasourceName := "data.civo_size.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoSizeConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataSourceCivoSizeExist(datasourceName),
				),
			},
		},
	})
}

func TestAccDataSourceCivoSize_WithFilterAndSort(t *testing.T) {
	datasourceName := "data.civo_size.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoSizeConfigWhitFilterAndSort(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCivoSizeExist(datasourceName),
					testAccCheckDataSourceCivoSizeFilteredAndSorted(datasourceName),
				),
			},
		},
	})
}

func testAccCheckDataSourceCivoSizeExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		rawTotal := rs.Primary.Attributes["sizes.#"]
		total, err := strconv.Atoi(rawTotal)
		if err != nil {
			return err
		}

		if total < 1 {
			return fmt.Errorf("No Civo sizes retrieved")
		}

		return nil
	}
}

func testAccCheckDataSourceCivoSizeFilteredAndSorted(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		rawTotal := rs.Primary.Attributes["sizes.#"]
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

		var prevCPUCore float64
		for i := 0; i < total; i++ {
			name := rs.Primary.Attributes[fmt.Sprintf("sizes.%d.name", i)]
			if !stringInSlice(name, []string{"g2.large", "g2.xlarge", "g2.2xlarge"}) {
				return fmt.Errorf("Name is not in expected test filter values")
			}

			CPUCore, _ := strconv.ParseFloat(rs.Primary.Attributes[fmt.Sprintf("sizes.%d.cpu_cores", i)], 64)
			if prevCPUCore > 0 {
				return fmt.Errorf("Sizes is not sorted by CPU Core in descending order")
			}
			prevCPUCore = CPUCore

		}

		return nil
	}
}

func testAccDataSourceCivoSizeConfig() string {
	return `
data "civo_size" "foobar" {
}
`
}

func testAccDataSourceCivoSizeConfigWhitFilterAndSort() string {
	return `
data "civo_size" "foobar" {
	filter {
        key = "name"
        values = ["large"]
	}
	
	sort {
        key = "cpu_cores"
        direction = "desc"
    }
}
`
}
