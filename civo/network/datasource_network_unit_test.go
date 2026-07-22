package network

import (
	"context"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestDataSourceNetworkRead_regionOnlyDoesNotPanic is a regression test for the
// crash reported against provider v1.0.45/v1.2.5, where a config supplying only
// `region` (e.g. `data "civo_network" "x" { region = var.region }`) caused a
// nil pointer dereference: neither the id nor the label branch ran, foundNetwork
// stayed nil, and d.SetId(foundNetwork.ID) segfaulted, taking down the plugin.
//
// The read must now surface a clean diagnostic instead of panicking.
func TestDataSourceNetworkRead_regionOnlyDoesNotPanic(t *testing.T) {
	client, server, err := civogo.NewClientForTesting(map[string]string{})
	if err != nil {
		t.Fatalf("failed to build test client: %s", err)
	}
	defer server.Close()

	d := schema.TestResourceDataRaw(t, DataSourceNetwork().Schema, map[string]interface{}{
		"region": "LON1",
	})

	diags := dataSourceNetworkRead(context.Background(), d, client)

	if !diags.HasError() {
		t.Fatal("expected an error when only region is provided, got none")
	}
}

// TestDataSourceNetworkRead_byID verifies the happy path still populates the
// data source from an id lookup after the nil guard was added.
func TestDataSourceNetworkRead_byID(t *testing.T) {
	client, server, err := civogo.NewClientForTesting(map[string]string{
		"/v2/vpc/networks": `[{"id":"net-1234","name":"my-net","label":"my-net","default":true}]`,
	})
	if err != nil {
		t.Fatalf("failed to build test client: %s", err)
	}
	defer server.Close()

	d := schema.TestResourceDataRaw(t, DataSourceNetwork().Schema, map[string]interface{}{
		"id":     "net-1234",
		"region": "LON1",
	})

	diags := dataSourceNetworkRead(context.Background(), d, client)

	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if got := d.Id(); got != "net-1234" {
		t.Errorf("id = %q, want %q", got, "net-1234")
	}
	if got := d.Get("name").(string); got != "my-net" {
		t.Errorf("name = %q, want %q", got, "my-net")
	}
	if got := d.Get("default").(bool); !got {
		t.Errorf("default = %v, want true", got)
	}
}
