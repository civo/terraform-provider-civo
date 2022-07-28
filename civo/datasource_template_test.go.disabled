package civo

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCivoTemplate_basic(t *testing.T) {
	datasourceName := "data.civo_template.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoTemplatesConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataSourceCivoTemplateExist(datasourceName),
				),
			},
		},
	})
}

func TestAccDataSourceCivoTemplate_WithFilter(t *testing.T) {
	datasourceName := "data.civo_template.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoTemplatesConfigWhitFilter(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCivoTemplateExist(datasourceName),
					testAccCheckDataSourceCivoTemplatesFiltered(datasourceName),
				),
			},
		},
	})
}

func testAccCheckDataSourceCivoTemplateExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		rawTotal := rs.Primary.Attributes["templates.#"]
		total, err := strconv.Atoi(rawTotal)
		if err != nil {
			return err
		}

		if total < 1 {
			return fmt.Errorf("No Civo templates retrieved")
		}

		return nil
	}
}

func testAccCheckDataSourceCivoTemplatesFiltered(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		rawTotal := rs.Primary.Attributes["templates.#"]
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

		var prevCode string
		for i := 0; i < total; i++ {
			code := rs.Primary.Attributes[fmt.Sprintf("templates.%d.code", i)]
			if !stringInSlice(code, []string{"debian-stretch", "debian-buster"}) {
				return fmt.Errorf("Code is not in expected test filter values")
			}
			if prevCode != "" && prevCode < code {
				return fmt.Errorf("Template is not sorted by code")
			}
			prevCode = code
		}

		return nil
	}
}

func testAccDataSourceCivoTemplatesConfig() string {
	return `
data "civo_template" "foobar" {
}
`
}

func testAccDataSourceCivoTemplatesConfigWhitFilter() string {
	return `
data "civo_template" "foobar" {
	filter {
        key = "code"
        values = ["debian"]
   }
}
`
}
