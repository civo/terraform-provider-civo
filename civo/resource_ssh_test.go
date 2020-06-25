package civo

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// example.Widget represents a concrete Go type that represents an API resource
func TestAccCivoSSHKey_basic(t *testing.T) {
	var SSHKey civogo.SSHKey

	// generate a random name for each test run
	resName := "civo_ssh_key.foobar"
	var SSHKeyName = acctest.RandomWithPrefix("tf-test")
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("civo@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoSSHKeyConfigBasic(SSHKeyName, publicKeyMaterial),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoSSHKeyResourceExists(resName, &SSHKey),
					// verify remote values
					testAccCheckCivoSSHKeyValues(&SSHKey, SSHKeyName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", SSHKeyName),
					resource.TestCheckResourceAttr(resName, "public_key", publicKeyMaterial),
				),
			},
		},
	})
}

func TestAccCivoSSHKey_update(t *testing.T) {
	var SSHKey civogo.SSHKey

	// generate a random name for each test run
	resName := "civo_ssh_key.foobar"
	var SSHKeyName = acctest.RandomWithPrefix("tf-test")
	var SSHKeyNameUpdate = acctest.RandomWithPrefix("rename-tf-test")
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("civo@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoSSHKeyConfigBasic(SSHKeyName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoSSHKeyResourceExists(resName, &SSHKey),
					testAccCheckCivoSSHKeyValues(&SSHKey, SSHKeyName),
					resource.TestCheckResourceAttr(resName, "name", SSHKeyName),
					resource.TestCheckResourceAttr(resName, "public_key", publicKeyMaterial),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoSSHKeyConfigUpdates(SSHKeyNameUpdate, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCivoSSHKeyResourceExists(resName, &SSHKey),
					testAccCheckCivoSSHKeyUpdated(&SSHKey, SSHKeyNameUpdate),
					resource.TestCheckResourceAttr(resName, "name", SSHKeyNameUpdate),
					resource.TestCheckResourceAttr(resName, "public_key", publicKeyMaterial),
				),
			},
		},
	})
}

func testAccCheckCivoSSHKeyValues(SSHKey *civogo.SSHKey, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if SSHKey.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, SSHKey.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoSSHKeyResourceExists(n string, SSHKey *civogo.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindSSHKey(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Ssh key not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*SSHKey = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoSSHKeyUpdated(SSHKey *civogo.SSHKey, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if SSHKey.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, SSHKey.Name)
		}
		return nil
	}
}

func testAccCheckCivoSSHKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_ssh_key" {
			continue
		}

		_, err := client.FindSSHKey(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Ssh key still exists")
		}
	}

	return nil
}

func testAccCheckCivoSSHKeyConfigBasic(name string, key string) string {
	return fmt.Sprintf(`
resource "civo_ssh_key" "foobar" {
	name = "%s"
	public_key = "%s"
}`, name, key)
}

func testAccCheckCivoSSHKeyConfigUpdates(name string, key string) string {
	return fmt.Sprintf(`
resource "civo_ssh_key" "foobar" {
	name = "%s"
	public_key = "%s"
}`, name, key)
}
