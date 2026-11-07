package dns

import (
	"errors"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceDNSImportNotFound(t *testing.T) {
	tests := []struct {
		name      string
		responses map[string]string
		resource  *schema.Resource
		importer  func(*schema.ResourceData, interface{}) ([]*schema.ResourceData, error)
		id        string
		wantErr   error
	}{
		{
			name:      "domain",
			responses: map[string]string{"/v2/dns": `[]`},
			resource:  ResourceDNSDomainName(),
			importer:  resourceDNSDomainImport,
			id:        "example.com",
			wantErr:   civogo.ErrDNSDomainNotFound,
		},
		{
			name:      "record",
			responses: map[string]string{"/v2/dns/domain-1/records": `[]`},
			resource:  ResourceDNSDomainRecord(),
			importer:  resourceDNSDomainRecordImport,
			id:        "domain-1:record-1",
			wantErr:   civogo.ErrDNSRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server, err := civogo.NewClientForTesting(tt.responses)
			if err != nil {
				t.Fatal(err)
			}
			defer server.Close()

			d := schema.TestResourceDataRaw(t, tt.resource.Schema, map[string]interface{}{})
			d.SetId(tt.id)

			got, err := tt.importer(d, client)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
			if got != nil {
				t.Fatalf("expected no imported state, got %v", got)
			}
		})
	}
}
