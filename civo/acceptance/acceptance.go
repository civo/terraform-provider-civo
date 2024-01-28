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

// TestAccProvider is a global instance of the provider under test. 
// It is used in acceptance tests to configure resources.
var TestAccProvider *schema.Provider

// TestAccProviders is a map of provider instances keyed by their name. 
// It is used in acceptance tests where multiple providers are in play.
var TestAccProviders map[string]*schema.Provider

// TestAccProviderFactories is a map of functions that return a provider instance and an error. 
// It is used in acceptance tests where the provider needs to be configured in a certain way.
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

// TestProvider - Test the provider itself
func TestProvider(t *testing.T) {
	if err := civo.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}


// TestProviderImpl is a test function to ensure that the Provider function of the civo package 
// returns an instance of the *schema.Provider type. It doesn't test any behavior of the provider itself.
func TestProviderImpl(t *testing.T) {
	var _ *schema.Provider = civo.Provider()
}

// TestToken - Test the provider token
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

// DiagnosticsToString - Convert diag.Diagnostics to string
func DiagnosticsToString(diags diag.Diagnostics) string {
	diagsAsStrings := make([]string, len(diags))
	for i, diag := range diags {
		diagsAsStrings[i] = diag.Summary
	}

	return strings.Join(diagsAsStrings, "; ")
}

// TestAccPreCheck - Check if the environment variables are set
func TestAccPreCheck(t *testing.T) {
	if v := os.Getenv("CIVO_TOKEN"); v == "" {
		t.Fatal("CIVO_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("CIVO_REGION"); v == "" {
		t.Fatal("CIVO_REGION must be set for acceptance tests")
	}
}
