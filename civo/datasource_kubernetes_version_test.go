package civo

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCivoKubernetesVersion_basic(t *testing.T) {
	datasourceName := "data.civo_kubernetes_version.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoKubernetesVersionConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataSourceCivoKubernetesVersionExist(datasourceName),
				),
			},
		},
	})
}

func TestAccDataSourceCivoKubernetesVersion_WithFilter(t *testing.T) {
	datasourceName := "data.civo_kubernetes_version.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoKubernetesVersionConfigWhitFilter(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCivoKubernetesVersionExist(datasourceName),
					testAccCheckDataSourceCivoKubernetesVersionFiltered(datasourceName),
				),
			},
		},
	})
}

func testAccCheckDataSourceCivoKubernetesVersionExist(n string) resource.TestCheckFunc {
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
			return fmt.Errorf("No Civo kubernetes versions retrieved")
		}

		return nil
	}
}

func testAccCheckDataSourceCivoKubernetesVersionFiltered(n string) resource.TestCheckFunc {
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
			name := rs.Primary.Attributes[fmt.Sprintf("versions.%d.type", i)]
			if !stringInSlice(name, []string{"talos"}) {
				return fmt.Errorf("Type is not in expected test filter values")
			}

		}

		return nil
	}
}

func testAccDataSourceCivoKubernetesVersionConfig() string {
	return `
data "civo_kubernetes_version" "foobar" {
}
`
}

func testAccDataSourceCivoKubernetesVersionConfigWhitFilter() string {
	return `
data "civo_kubernetes_version" "foobar" {
	filter {
        key = "type"
        values = ["talos"]
	}
}
`
}
