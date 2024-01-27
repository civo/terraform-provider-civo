package ssh_test

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// example.Widget represents a concrete Go type that represents an API resource
func CivoSSHKey_basic(t *testing.T) {
	var SSHKey civogo.SSHKey

	// generate a random name for each test run
	resName := "civo_ssh_key.foobar"
	var SSHKeyName = acctest.RandomWithPrefix("tf-test")
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("civo@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: CivoSSHKeyConfigBasic(SSHKeyName, publicKeyMaterial),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					CivoSSHKeyResourceExists(resName, &SSHKey),
					// verify remote values
					CivoSSHKeyValues(&SSHKey, SSHKeyName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", SSHKeyName),
					resource.TestCheckResourceAttr(resName, "public_key", publicKeyMaterial),
				),
			},
		},
	})
}

func CivoSSHKey_update(t *testing.T) {
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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CivoSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoSSHKeyConfigBasic(SSHKeyName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					CivoSSHKeyResourceExists(resName, &SSHKey),
					CivoSSHKeyValues(&SSHKey, SSHKeyName),
					resource.TestCheckResourceAttr(resName, "name", SSHKeyName),
					resource.TestCheckResourceAttr(resName, "public_key", publicKeyMaterial),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: CivoSSHKeyConfigUpdates(SSHKeyNameUpdate, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					CivoSSHKeyResourceExists(resName, &SSHKey),
					CivoSSHKeyUpdated(&SSHKey, SSHKeyNameUpdate),
					resource.TestCheckResourceAttr(resName, "name", SSHKeyNameUpdate),
					resource.TestCheckResourceAttr(resName, "public_key", publicKeyMaterial),
				),
			},
		},
	})
}

func CivoSSHKeyValues(SSHKey *civogo.SSHKey, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if SSHKey.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, SSHKey.Name)
		}
		return nil
	}
}

// CheckExampleResourceExists queries the API and retrieves the matching Widget.
func CivoSSHKeyResourceExists(n string, SSHKey *civogo.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := acceptance.TestAccProvider.Meta().(*civogo.Client)
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

func CivoSSHKeyUpdated(SSHKey *civogo.SSHKey, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if SSHKey.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, SSHKey.Name)
		}
		return nil
	}
}

func CivoSSHKeyDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*civogo.Client)

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

func CivoSSHKeyConfigBasic(name string, key string) string {
	return fmt.Sprintf(`
resource "civo_ssh_key" "foobar" {
	name = "%s"
	public_key = "%s"
}`, name, key)
}

func CivoSSHKeyConfigUpdates(name string, key string) string {
	return fmt.Sprintf(`
resource "civo_ssh_key" "foobar" {
	name = "%s"
	public_key = "%s"
}`, name, key)
}
