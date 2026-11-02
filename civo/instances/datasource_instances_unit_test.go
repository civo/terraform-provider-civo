package instances

import (
	"context"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestDataSourceInstances_filterByTags is a regression test for #429, where filtering
// civo_instances on `tags` panicked with "interface {} is []string, not *schema.Set"
// because the flatten passes the instance's []string for a TypeSet field, taking down
// the plugin. The filter must now return the matching instances.
func TestDataSourceInstances_filterByTags(t *testing.T) {
	client, server, err := civogo.NewClientForTesting(map[string]string{
		"/v2/instances": `{"page": 1, "per_page": 100, "pages": 1, "items": [
			{"id": "inst-1", "hostname": "web-1", "tags": ["prod", "web"]},
			{"id": "inst-2", "hostname": "db-1", "tags": ["staging"]},
			{"id": "inst-3", "hostname": "untagged"}
		]}`,
	})
	if err != nil {
		t.Fatalf("failed to build test client: %s", err)
	}
	defer server.Close()

	ds := DataSourceInstances()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"filter": []interface{}{
			map[string]interface{}{"key": "tags", "values": []interface{}{"prod"}},
		},
	})

	if diags := ds.ReadContext(context.Background(), d, client); diags.HasError() {
		t.Fatalf("read failed: %v", diags)
	}

	if got := d.Get("instances.#").(int); got != 1 {
		t.Fatalf("expected 1 instance tagged prod, got %d", got)
	}
	if got := d.Get("instances.0.hostname"); got != "web-1" {
		t.Errorf("expected web-1, got %v", got)
	}
}
