package instances

import (
	"context"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestResourceInstanceReservedIPRead_drift covers a reserved IP assignment that
// changed outside Terraform. The read used to return an error in both cases, so
// plan and refresh failed instead of planning to recreate the assignment.
func TestResourceInstanceReservedIPRead_drift(t *testing.T) {
	tests := []struct {
		name       string
		ips        string
		wantGoneID bool
	}{
		{
			name:       "still assigned",
			ips:        `{"page": 1, "per_page": 100, "pages": 1, "items": [{"id": "ip-1", "name": "web", "ip": "192.0.2.10", "assigned_to": {"id": "inst-1", "type": "instance"}}]}`,
			wantGoneID: false,
		},
		{
			name:       "reserved ip deleted",
			ips:        `{"page": 1, "per_page": 100, "pages": 1, "items": []}`,
			wantGoneID: true,
		},
		{
			name:       "reserved ip unassigned",
			ips:        `{"page": 1, "per_page": 100, "pages": 1, "items": [{"id": "ip-1", "name": "web", "ip": "192.0.2.10"}]}`,
			wantGoneID: true,
		},
		{
			name:       "reserved ip moved to another instance",
			ips:        `{"page": 1, "per_page": 100, "pages": 1, "items": [{"id": "ip-1", "name": "web", "ip": "192.0.2.10", "assigned_to": {"id": "inst-2", "type": "instance"}}]}`,
			wantGoneID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server, err := civogo.NewClientForTesting(map[string]string{
				"/v2/vpc/ips": tt.ips,
			})
			if err != nil {
				t.Fatalf("failed to build test client: %s", err)
			}
			defer server.Close()

			d := schema.TestResourceDataRaw(t, ResourceInstanceReservedIPAssignment().Schema, map[string]interface{}{
				"reserved_ip_id": "ip-1",
				"instance_id":    "inst-1",
				"region":         "LON1",
			})
			d.SetId("assignment-1")

			if diags := resourceInstanceReservedIPRead(context.Background(), d, client); diags.HasError() {
				t.Fatalf("read returned an error: %v", diags)
			}

			if gone := d.Id() == ""; gone != tt.wantGoneID {
				t.Errorf("removed from state = %t, want %t", gone, tt.wantGoneID)
			}
		})
	}
}
