package network

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestResourceVPCSubnetImport covers the importer added for
// civo/terraform-provider-civo#419: the docs have always documented
// `NETWORK_ID:SUBNET_ID`, but the resource used a passthrough importer, which
// left `network_id` empty and made the first read after an import fail.
func TestResourceVPCSubnetImport(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		wantErr   bool
		wantID    string
		wantNetID string
	}{
		{
			name:      "network and subnet id",
			id:        "net-1234:subnet-5678",
			wantID:    "subnet-5678",
			wantNetID: "net-1234",
		},
		{"subnet id only", "subnet-5678", true, "", ""},
		{"empty network id", ":subnet-5678", true, "", ""},
		{"empty subnet id", "net-1234:", true, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, ResourceVPCSubnet().Schema, map[string]interface{}{})
			d.SetId(tt.id)

			got, err := resourceVPCSubnetImport(context.Background(), d, nil)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error for id %q, got none", tt.id)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(got) != 1 {
				t.Fatalf("expected 1 resource, got %d", len(got))
			}
			if got[0].Id() != tt.wantID {
				t.Errorf("id = %q, want %q", got[0].Id(), tt.wantID)
			}
			if netID := got[0].Get("network_id").(string); netID != tt.wantNetID {
				t.Errorf("network_id = %q, want %q", netID, tt.wantNetID)
			}
		})
	}
}
