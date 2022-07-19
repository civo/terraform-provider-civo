package civo

import (
	"context"
	"os"
	"strings"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"civo": testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"civo": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestToken(t *testing.T) {
	rawProvider := Provider()
	raw := map[string]interface{}{
		"token": "123456789",
	}

	diags := rawProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configure failed: %s", diagnosticsToString(diags))
	}
}

func diagnosticsToString(diags diag.Diagnostics) string {
	diagsAsStrings := make([]string, len(diags))
	for i, diag := range diags {
		diagsAsStrings[i] = diag.Summary
	}

	return strings.Join(diagsAsStrings, "; ")
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CIVO_TOKEN"); v == "" {
		t.Fatal("CIVO_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("CIVO_REGION"); v == "" {
		t.Fatal("CIVO_REGION must be set for acceptance tests")
	}
}
