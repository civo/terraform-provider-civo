package ssh_test

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"strings"
	"testing"

	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"golang.org/x/crypto/ssh"
)

func DataSourceCivoSSHKey_basic(t *testing.T) {
	datasourceName := "data.civo_ssh_key.foobar"
	name := acctest.RandomWithPrefix("sshkey-test")
	pubKey, err := GenerateDataSourceCivoSSHKeyPublic()
	if err != nil {
		t.Fatalf("Unable to generate public key: %v", err)
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoSSHKeyConfig(name, pubKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
				),
			},
		},
	})
}

func GenerateDataSourceCivoSSHKeyPublic() (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", fmt.Errorf("Unable to generate key: %v", err)
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", fmt.Errorf("Unable to generate key: %v", err)
	}

	return strings.TrimSpace(string(ssh.MarshalAuthorizedKey(publicKey))), nil
}

func DataSourceCivoSSHKeyConfig(name string, key string) string {
	return fmt.Sprintf(`
resource "civo_ssh_key" "foobar" {
	name = "%s"
    public_key = "%s"
}

data "civo_ssh_key" "foobar" {
	name = civo_ssh_key.foobar.name
}
`, name, key)
}
