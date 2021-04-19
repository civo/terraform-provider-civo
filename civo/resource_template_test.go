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
func TestAccCivoTemplate_basic(t *testing.T) {
	var template civogo.Template

	// generate a random name for each test run
	resName := "civo_template.foobar"
	var templateName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoTemplateDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoTemplateConfigBasic(templateName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoTemplateResourceExists(resName, &template),
					// verify remote values
					testAccCheckCivoTemplateValues(&template, templateName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", templateName),
					resource.TestCheckResourceAttrSet(resName, "image_id"),
					resource.TestCheckResourceAttrSet(resName, "short_description"),
					resource.TestCheckResourceAttrSet(resName, "description"),
					resource.TestCheckResourceAttrSet(resName, "default_username"),
				),
			},
		},
	})
}

func TestAccCivoTemplate_update(t *testing.T) {
	var template civogo.Template

	// generate a random name for each test run
	resName := "civo_template.foobar"
	var templateName = acctest.RandomWithPrefix("tf-test")
	var templateNameUpdate = acctest.RandomWithPrefix("tf-update")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoTemplateConfigBasic(templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoTemplateResourceExists(resName, &template),
					testAccCheckCivoTemplateValues(&template, templateName),
					resource.TestCheckResourceAttr(resName, "name", templateName),
					resource.TestCheckResourceAttrSet(resName, "image_id"),
					resource.TestCheckResourceAttrSet(resName, "short_description"),
					resource.TestCheckResourceAttrSet(resName, "description"),
					resource.TestCheckResourceAttrSet(resName, "default_username"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoTemplateConfigUpdates(templateName, templateNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoTemplateResourceExists(resName, &template),
					testAccCheckCivoTemplateUpdated(&template, templateNameUpdate),
					resource.TestCheckResourceAttr(resName, "name", templateNameUpdate),
					resource.TestCheckResourceAttrSet(resName, "image_id"),
					resource.TestCheckResourceAttrSet(resName, "short_description"),
					resource.TestCheckResourceAttrSet(resName, "description"),
					resource.TestCheckResourceAttrSet(resName, "default_username"),
				),
			},
		},
	})
}

func testAccCheckCivoTemplateValues(template *civogo.Template, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if template.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, template.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoTemplateResourceExists(n string, template *civogo.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.GetTemplateByCode(rs.Primary.Attributes["code"])
		if err != nil {
			return fmt.Errorf("Template not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*template = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoTemplateUpdated(template *civogo.Template, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if template.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, template.Name)
		}
		return nil
	}
}

func testAccCheckCivoTemplateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_template" {
			continue
		}

		_, err := client.GetTemplateByCode(rs.Primary.Attributes["code"])
		if err == nil {
			return fmt.Errorf("Template still exists")
		}
	}

	return nil
}

func testAccCheckCivoTemplateConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_template" "foobar" {
	code = "%s"
	name = "%s"
	image_id = "38686161-ba25-4899-ac0a-54eaf35239c0"
	short_description = "my-description"
	description = "my-description-long"
	default_username = "admin"
}`, name, name)
}

func testAccCheckCivoTemplateConfigUpdates(name string, updateName string) string {
	return fmt.Sprintf(`
resource "civo_template" "foobar" {
	code = "%s"
	name = "%s"
	image_id = "38686161-ba25-4899-ac0a-54eaf35239c0"
	short_description = "my-description"
	description = "my-description-long"
	default_username = "admin"
}`, name, updateName)
}
