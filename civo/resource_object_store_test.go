package civo

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// example.Widget represents a concrete Go type that represents an API resource
func TestAccCivoObjectStore_basic(t *testing.T) {
	var store civogo.ObjectStore

	// generate a random name for each test run
	resName := "civo_object_store.foobar"
	var storeName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoObjectStoreConfigBasic(storeName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoObjectStoreResourceExists(resName, &store),
					// verify remote values
					testAccCheckCivoObjectStoreValues(&store, storeName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", storeName),
					resource.TestCheckResourceAttrSet(resName, "max_size_gb"),
					resource.TestCheckResourceAttrSet(resName, "bucket_url"),
					resource.TestCheckResourceAttrSet(resName, "access_key_id"),
					resource.TestCheckResourceAttr(resName, "status", "ready"),
				),
			},
		},
	})
}

func TestAccCivoObjectStore_update(t *testing.T) {
	var store civogo.ObjectStore

	// generate a random name for each test run
	resName := "civo_object_store.foobar"
	var storeName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoObjectStoreConfigBasic(storeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoObjectStoreResourceExists(resName, &store),
					testAccCheckCivoObjectStoreValues(&store, storeName),
					resource.TestCheckResourceAttr(resName, "name", storeName),
					resource.TestCheckResourceAttrSet(resName, "access_key_id"),
					resource.TestCheckResourceAttr(resName, "max_size_gb", strconv.Itoa(500)),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoObjectStoreConfigUpdates(storeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoObjectStoreResourceExists(resName, &store),
					testAccCheckCivoObjectStoreUpdated(&store, storeName),
					resource.TestCheckResourceAttr(resName, "name", storeName),
					resource.TestCheckResourceAttrSet(resName, "access_key_id"),
					resource.TestCheckResourceAttr(resName, "max_size_gb", strconv.Itoa(1000)),
				),
			},
		},
	})
}

func testAccCheckCivoObjectStoreValues(store *civogo.ObjectStore, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if store.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, store.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoObjectStoreResourceExists(n string, store *civogo.ObjectStore) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindObjectStore(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Object Store not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*store = *resp

		return nil
	}
}

func testAccCheckCivoObjectStoreUpdated(store *civogo.ObjectStore, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if store.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, store.Name)
		}
		return nil
	}
}

func testAccCheckCivoObjectStoreDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_object_store" {
			continue
		}

		_, err := client.FindObjectStore(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Object Store still exists")
		}
	}

	return nil
}

func testAccCheckCivoObjectStoreConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_object_store" "foobar" {
	name = "%s"
	max_size_gb = 500
	region = "FAKE"
}`, name)
}

func testAccCheckCivoObjectStoreConfigUpdates(name string) string {
	return fmt.Sprintf(`
resource "civo_object_store" "foobar" {
	name = "%s"
	max_size_gb = 1000
	region = "FAKE"
}`, name)
}
