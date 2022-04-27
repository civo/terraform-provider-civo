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
func TestAccCivoApplication_basic(t *testing.T) {
	var app civogo.Application

	// generate a random name for each test run
	resName := "civo_application.foobar"
	var appName = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoApplicationDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoApplicationConfigBasic(appName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoApplicationResourceExists(resName, &app),
					// verify remote values
					testAccCheckCivoApplicationValues(&app, appName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", appName),
					resource.TestCheckResourceAttr(resName, "size", "small"),
				),
			},
		},
	})
}

func TestAccCivoApplicationSize_update(t *testing.T) {
	var app civogo.Application

	// generate a random name for each test run
	resName := "civo_application.foobar"
	var appName = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoApplicationConfigBasic(appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoApplicationResourceExists(resName, &app),
					testAccCheckCivoApplicationValues(&app, appName),
					resource.TestCheckResourceAttr(resName, "name", appName),
					resource.TestCheckResourceAttr(resName, "size", "small"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoApplicationConfigUpdates(appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoApplicationResourceExists(resName, &app),
					testAccCheckCivoApplicationUpdated(&app, appName),
					resource.TestCheckResourceAttr(resName, "name", appName),
					resource.TestCheckResourceAttr(resName, "size", "medium"),
				),
			},
		},
	})
}

func testAccCheckCivoApplicationValues(app *civogo.Application, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if app.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, app.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoApplicationResourceExists(n string, app *civogo.Application) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.GetApplication(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("App not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*app = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoApplicationUpdated(app *civogo.Application, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if app.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, app.Name)
		}
		return nil
	}
}

func testAccCheckCivoApplicationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_application" {
			continue
		}

		_, err := client.GetApplication(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Application still exists")
		}
	}

	return nil
}

func testAccCheckCivoApplicationConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_application" "foobar" {
	name = "%s"
	size = "small"
}`, name)
}

func testAccCheckCivoApplicationConfigUpdates(name string) string {
	return fmt.Sprintf(`
resource "civo_application" "foobar" {
	name = "%s"
	size = "medium"
}`, name)
}
