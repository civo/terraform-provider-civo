package size_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCivoSize_basic(t *testing.T) {
	datasourceName := "data.civo_size.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoSizeConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					DataSourceCivoSizeExist(datasourceName),
				),
			},
		},
	})
}

func TestAccDataSourceCivoSize_WithFilterAndSort(t *testing.T) {
	datasourceName := "data.civo_size.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoSizeConfigWhitFilterAndSort(),
				Check: resource.ComposeTestCheckFunc(
					DataSourceCivoSizeExist(datasourceName),
					DataSourceCivoSizeFilteredAndSorted(datasourceName),
				),
			},
		},
	})
}

func DataSourceCivoSizeExist(n string) resource.TestCheckFunc {
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

func DataSourceCivoSizeFilteredAndSorted(n string) resource.TestCheckFunc {
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
				return fmt.Errorf("Name is not in expected test filter values")
			}

			CPUCore, _ := strconv.ParseFloat(rs.Primary.Attributes[fmt.Sprintf("sizes.%d.cpu", i)], 64)
			if prevCPUCore > 0 {
				return fmt.Errorf("Sizes is not sorted by CPU Core in descending order")
			}
			prevCPUCore = CPUCore

		}

		return nil
	}
}

func DataSourceCivoSizeConfig() string {
	return `
data "civo_size" "foobar" {
}
`
}

func DataSourceCivoSizeConfigWhitFilterAndSort() string {
	return `
data "civo_size" "foobar" {
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
