package civo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
	t.Run("reading token from token attribute", func(t *testing.T) {
		const testToken = "123456789"

		raw := map[string]interface{}{
			"token": testToken,
		}
		configureProvider(t, raw)
	})

	t.Run("reading token from environment variable", func(t *testing.T) {
		const testToken = "env12345"
		oldToken := os.Getenv("CIVO_TOKEN")
		os.Setenv("CIVO_TOKEN", testToken)
		defer os.Setenv("CIVO_TOKEN", oldToken) // Restore the original value

		raw := map[string]interface{}{}
		configureProvider(t, raw)
	})

	t.Run("reading token from credential file", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "civo-provider-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		credentialFile := filepath.Join(tempDir, "credential.json")
		testToken := "file3409"
		credContent := fmt.Sprintf(`{"apikeys":{"CIVO_TOKEN":"%s"}, "meta":{"current_apikey":"CIVO_TOKEN"}}`, testToken)
		err = os.WriteFile(credentialFile, []byte(credContent), 0600)
		if err != nil {
			t.Fatalf("Failed to write credentials file: %v", err)
		}

		raw := map[string]interface{}{
			"credential_file": credentialFile,
		}
		configureProvider(t, raw)
	})

	t.Run("reading token from CLI config", func(t *testing.T) {
		// Create a mock CLI config file
		tempDir, err := os.MkdirTemp("", "civo-cli-config-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		testToken := "cliconfig12345"
		cliConfigContent := fmt.Sprintf(`
			{
				"apikeys": {
					"CIVO_TOKEN": "%s"
				},
				"meta": {
					"admin": false,
					"current_apikey": "CIVO_TOKEN",
					"default_region": "NYC1",
					"latest_release_check": "2024-07-07T14:07:51.996201195+05:30",
					"url": "https://api.civo.com",
					"last_command_executed": "2024-06-20T15:09:00.212548723+05:30"
				},
				"region_to_features": null
			}`, testToken)

		cliConfigFile := filepath.Join(tempDir, ".civo.json")
		err = os.WriteFile(cliConfigFile, []byte(cliConfigContent), 0600)
		if err != nil {
			t.Fatalf("Failed to write CLI config file: %v", err)
		}

		// Temporarily set HOME to our temp directory
		oldHome := os.Getenv("HOME")
		os.Setenv("HOME", tempDir)
		defer os.Setenv("HOME", oldHome)

		raw := map[string]interface{}{}
		configureProvider(t, raw)
	})

}

func diagnosticsToString(diags diag.Diagnostics) string {
	diagsAsStrings := make([]string, len(diags))
	for i, diag := range diags {
		diagsAsStrings[i] = diag.Summary
	}

	return strings.Join(diagsAsStrings, "; ")
}

func configureProvider(t testing.TB, raw map[string]interface{}) {
	t.Helper()

	rawProvider := Provider()
	diags := rawProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configure failed: %s", diagnosticsToString(diags))
	}
}
