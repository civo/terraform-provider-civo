package civo

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCivoInstanceSize_basic(t *testing.T) {
	datasourceName := "data.civo_instances_size.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoInstanceSizeConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataSourceCivoInstanceSizeExist(datasourceName),
				),
			},
		},
	})
}

func TestAccDataSourceCivoInstanceSize_WithFilterAndSort(t *testing.T) {
	datasourceName := "data.civo_instances_size.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoInstanceSizeConfigWhitFilterAndSort(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCivoInstanceSizeExist(datasourceName),
					testAccCheckDataSourceCivoInstanceSizeFilteredAndSorted(datasourceName),
				),
			},
		},
	})
}

func testAccCheckDataSourceCivoInstanceSizeExist(n string) resource.TestCheckFunc {
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

func testAccCheckDataSourceCivoInstanceSizeFilteredAndSorted(n string) resource.TestCheckFunc {
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
			if !stringInSlice(name, []string{"g3.large", "g3.xlarge", "g3.2xlarge"}) {
				return fmt.Errorf("Name is not in expected test filter values, got %s", name)
			}

			CPUCore, _ := strconv.ParseFloat(rs.Primary.Attributes[fmt.Sprintf("sizes.%d.cpu_cores", i)], 64)
			if prevCPUCore > 0 {
				return fmt.Errorf("Sizes is not sorted by CPU Core in descending order got %f before %f", CPUCore, prevCPUCore)
			}
			prevCPUCore = CPUCore

		}

		return nil
	}
}

func testAccDataSourceCivoInstanceSizeConfig() string {
	return `
data "civo_instances_size" "foobar" {
}
`
}

func testAccDataSourceCivoInstanceSizeConfigWhitFilterAndSort() string {
	return `
data "civo_instances_size" "foobar" {
	filter {
        key = "name"
        values = ["g3.large"]
	}
	
	sort {
        key = "cpu"
        direction = "desc"
    }
}
`
}
