package dns

import (
	"errors"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const testDNSDomainsListResponse = `[{"id": "12345", "account_id": "1", "name": "example.com"}, {"id": "12346", "account_id": "1", "name": "example.net"}]`

func TestResourceDNSDomainImport(t *testing.T) {
	tests := []struct {
		name       string
		responses  map[string]string
		id         string
		wantErr    error
		wantID     string
		wantName   string
		wantAcctID string
	}{
		{
			name:       "by domain name",
			responses:  map[string]string{"/v2/dns": testDNSDomainsListResponse},
			id:         "example.com",
			wantID:     "12345",
			wantName:   "example.com",
			wantAcctID: "1",
		},
		{
			name:       "by domain ID",
			responses:  map[string]string{"/v2/dns": testDNSDomainsListResponse},
			id:         "12346",
			wantID:     "12346",
			wantName:   "example.net",
			wantAcctID: "1",
		},
		{
			name:      "missing domain name",
			responses: map[string]string{"/v2/dns": `[]`},
			id:        "missing.example",
			wantErr:   civogo.ErrDNSDomainNotFound,
		},
		{
			name:      "missing domain ID",
			responses: map[string]string{"/v2/dns": testDNSDomainsListResponse},
			id:        "99999",
			wantErr:   civogo.ErrDNSDomainNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server, err := civogo.NewClientForTesting(tt.responses)
			if err != nil {
				t.Fatal(err)
			}
			defer server.Close()

			d := schema.TestResourceDataRaw(t, ResourceDNSDomainName().Schema, map[string]interface{}{})
			d.SetId(tt.id)

			got, err := resourceDNSDomainImport(d, client)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				if got != nil {
					t.Fatalf("expected no imported state, got %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != 1 {
				t.Fatalf("expected 1 resource, got %d", len(got))
			}
			if got[0].Id() != tt.wantID {
				t.Errorf("id = %q, want %q", got[0].Id(), tt.wantID)
			}
			if name := got[0].Get("name").(string); name != tt.wantName {
				t.Errorf("name = %q, want %q", name, tt.wantName)
			}
			if acct := got[0].Get("account_id").(string); acct != tt.wantAcctID {
				t.Errorf("account_id = %q, want %q", acct, tt.wantAcctID)
			}
		})
	}
}

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
