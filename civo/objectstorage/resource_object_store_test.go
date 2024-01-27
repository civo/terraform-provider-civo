package objectstorage_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// example.Widget represents a concrete Go type that represents an API resource
func CivoObjectStore_basic(t *testing.T) {
	var store civogo.ObjectStore

	// generate a random name for each test run
	resName := "civo_object_store.foobar"
	var storeName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: CivoObjectStoreConfigBasic(storeName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					CivoObjectStoreResourceExists(resName, &store),
					// verify remote values
					CivoObjectStoreValues(&store, storeName),
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

func CivoObjectStore_update(t *testing.T) {
	var store civogo.ObjectStore

	// generate a random name for each test run
	resName := "civo_object_store.foobar"
	var storeName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoObjectStoreConfigBasic(storeName),
				Check: resource.ComposeTestCheckFunc(
					CivoObjectStoreResourceExists(resName, &store),
					CivoObjectStoreValues(&store, storeName),
					resource.TestCheckResourceAttr(resName, "name", storeName),
					resource.TestCheckResourceAttrSet(resName, "access_key_id"),
					resource.TestCheckResourceAttr(resName, "max_size_gb", strconv.Itoa(500)),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: CivoObjectStoreConfigUpdates(storeName),
				Check: resource.ComposeTestCheckFunc(
					CivoObjectStoreResourceExists(resName, &store),
					CivoObjectStoreUpdated(&store, storeName),
					resource.TestCheckResourceAttr(resName, "name", storeName),
					resource.TestCheckResourceAttrSet(resName, "access_key_id"),
					resource.TestCheckResourceAttr(resName, "max_size_gb", strconv.Itoa(1000)),
				),
			},
		},
	})
}

func CivoObjectStoreValues(store *civogo.ObjectStore, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if store.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, store.Name)
		}
		return nil
	}
}

// CheckExampleResourceExists queries the API and retrieves the matching Widget.
func CivoObjectStoreResourceExists(n string, store *civogo.ObjectStore) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := acceptance.TestAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindObjectStore(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Object Store not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*store = *resp

		return nil
	}
}

func CivoObjectStoreUpdated(store *civogo.ObjectStore, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if store.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, store.Name)
		}
		return nil
	}
}

func CivoObjectStoreDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*civogo.Client)

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

func CivoObjectStoreConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_object_store" "foobar" {
	name = "%s"
	max_size_gb = 500
	region = "FAKE"
}`, name)
}

func CivoObjectStoreConfigUpdates(name string) string {
	return fmt.Sprintf(`
resource "civo_object_store" "foobar" {
	name = "%s"
	max_size_gb = 1000
	region = "FAKE"
}`, name)
}
