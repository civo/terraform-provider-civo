package acceptance

import (
	"context"
	"os"
	"strings"

	"testing"

	"github.com/civo/terraform-provider-civo/civo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TestAccProvider *schema.Provider
var TestAccProviders map[string]*schema.Provider
var TestAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	TestAccProvider = civo.Provider()
	TestAccProviders = map[string]*schema.Provider{
		"civo": TestAccProvider,
	}
	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"civo": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := civo.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = civo.Provider()
}

func TestToken(t *testing.T) {
	rawProvider := civo.Provider()
	raw := map[string]interface{}{
		"token": "123456789",
	}

	diags := rawProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configure failed: %s", DiagnosticsToString(diags))
	}
}

func DiagnosticsToString(diags diag.Diagnostics) string {
	diagsAsStrings := make([]string, len(diags))
	for i, diag := range diags {
		diagsAsStrings[i] = diag.Summary
	}

	return strings.Join(diagsAsStrings, "; ")
}

func TestAccPreCheck(t *testing.T) {
	if v := os.Getenv("CIVO_TOKEN"); v == "" {
		t.Fatal("CIVO_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("CIVO_REGION"); v == "" {
		t.Fatal("CIVO_REGION must be set for acceptance tests")
	}
}
