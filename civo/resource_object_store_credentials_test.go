package civo

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCivoObjectStoreCredential_basic(t *testing.T) {
	var storeCredential civogo.ObjectStoreCredential

	// generate a random name for each test run
	resName := "civo_object_store_credential.foobar"
	var storeCredentialName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoObjectStoreCredentialDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoObjectStoreCredentialConfigBasic(storeCredentialName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoObjectStoreCredentialResourceExists(resName, &storeCredential),
					// verify remote values
					testAccCheckCivoObjectStoreCredentialValues(&storeCredential, storeCredentialName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", storeCredentialName),
					resource.TestCheckResourceAttrSet(resName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resName, "secret_access_key"),
					resource.TestCheckResourceAttr(resName, "status", "ready"),
				),
			},
		},
	})
}

func TestAccCivoObjectStoreCredentialWhitCustomCredential_basic(t *testing.T) {
	var storeCredential civogo.ObjectStoreCredential

	// generate a random name for each test run
	resName := "civo_object_store_credential.foobar"
	var storeCredentialName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoObjectStoreCredentialDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoObjectStoreCredentialWhitCustomCredentialBasic(storeCredentialName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoObjectStoreCredentialResourceExists(resName, &storeCredential),
					// verify remote values
					testAccCheckCivoObjectStoreCredentialValues(&storeCredential, storeCredentialName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", storeCredentialName),
					resource.TestCheckResourceAttr(resName, "access_key_id", "1234567890"),
					resource.TestCheckResourceAttr(resName, "secret_access_key", "1234567890"),
					resource.TestCheckResourceAttr(resName, "status", "ready"),
				),
			},
		},
	})
}

func TestAccCivoObjectStoreCredential_update(t *testing.T) {
	var storeCredential civogo.ObjectStoreCredential

	// generate a random name for each test run
	resName := "civo_object_store_credential.foobar"
	var storeCredentialName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoObjectStoreCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoObjectStoreCredentialConfigBasic(storeCredentialName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoObjectStoreCredentialResourceExists(resName, &storeCredential),
					testAccCheckCivoObjectStoreCredentialValues(&storeCredential, storeCredentialName),
					resource.TestCheckResourceAttr(resName, "name", storeCredentialName),
					resource.TestCheckResourceAttrSet(resName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resName, "secret_access_key"),
					resource.TestCheckResourceAttr(resName, "status", "ready"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoObjectStoreCredentialConfigUpdates(storeCredentialName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoObjectStoreCredentialResourceExists(resName, &storeCredential),
					testAccCheckCivoObjectStoreCredentialUpdated(&storeCredential, storeCredentialName),
					resource.TestCheckResourceAttr(resName, "name", storeCredentialName),
					resource.TestCheckResourceAttrSet(resName, "access_key_id"),
					resource.TestCheckResourceAttr(resName, "secret_access_key", "1234567890"),
					resource.TestCheckResourceAttr(resName, "status", "ready"),
				),
			},
		},
	})
}

func testAccCheckCivoObjectStoreCredentialValues(storeCredential *civogo.ObjectStoreCredential, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if storeCredential.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, storeCredential.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoObjectStoreCredentialResourceExists(n string, storeCredential *civogo.ObjectStoreCredential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindObjectStoreCredential(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Object Store Credential not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*storeCredential = *resp

		return nil
	}
}

func testAccCheckCivoObjectStoreCredentialUpdated(storeCredential *civogo.ObjectStoreCredential, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if storeCredential.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, storeCredential.Name)
		}
		return nil
	}
}

func testAccCheckCivoObjectStoreCredentialDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_object_store_credential" {
			continue
		}

		_, err := client.FindObjectStoreCredential(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Object Store Credential still exists")
		}
	}

	return nil
}

func testAccCheckCivoObjectStoreCredentialConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_object_store_credential" "foobar" {
	name = "%s"
}`, name)
}

func testAccCheckCivoObjectStoreCredentialWhitCustomCredentialBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_object_store_credential" "foobar" {
	name = "%s"
	access_key_id = "1234567890"
	secret_access_key = "1234567890"
}`, name)
}

func testAccCheckCivoObjectStoreCredentialConfigUpdates(name string) string {
	return fmt.Sprintf(`
resource "civo_object_store_credential" "foobar" {
	name = "%s"
	secret_access_key = "1234567890"
}`, name)
}
