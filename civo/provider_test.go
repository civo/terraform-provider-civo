package civo

import (
	"context"
	"strings"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestProvider tests the provider configuration
func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err.Error())
	}
}

// TestProvider_impl tests the provider implementation
func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

// TestToken tests the token configuration
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
