package network

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// networksAPI answers the network list with status and body.
func networksAPI(t *testing.T, status int, body string) *civogo.Client {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/v2/vpc/networks" {
			t.Errorf("unexpected request %s", req.URL)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(status)
		_, _ = rw.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	client, err := civogo.NewClientForTestingWithServer(server)
	if err != nil {
		t.Fatalf("could not build the test client: %v", err)
	}
	client.Region = "LON1"
	return client
}

func networkState(t *testing.T, id string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceNetwork().Schema, map[string]interface{}{
		"label":  "my-net",
		"region": "LON1",
	})
	d.SetId(id)
	return d
}

func TestResourceNetworkRead_found(t *testing.T) {
	client := networksAPI(t, http.StatusOK, `[{"id":"net-1","name":"cust-net-1","label":"my-net","cidr":"10.0.0.0/24"}]`)
	d := networkState(t, "net-1")

	if diags := resourceNetworkRead(context.Background(), d, client); diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	if d.Id() != "net-1" || d.Get("label").(string) != "my-net" || d.Get("cidr_v4").(string) != "10.0.0.0/24" {
		t.Fatalf("unexpected state: id=%q label=%q cidr_v4=%q", d.Id(), d.Get("label"), d.Get("cidr_v4"))
	}
}

// A network deleted outside Terraform must drop out of state so the plan
// recreates it, instead of keeping the ID with empty attributes.
func TestResourceNetworkRead_deletedOutsideTerraform(t *testing.T) {
	client := networksAPI(t, http.StatusOK, `[{"id":"net-2","label":"other"}]`)
	d := networkState(t, "net-1")

	if diags := resourceNetworkRead(context.Background(), d, client); diags.HasError() || d.Id() != "" {
		t.Fatalf("want network cleared from state, got diags=%v id=%q", diags, d.Id())
	}
}

// A failed list says nothing about the network, so it must stay in state.
func TestResourceNetworkRead_listFails(t *testing.T) {
	client := networksAPI(t, http.StatusInternalServerError, `{"code":"internal_server_error","reason":"boom"}`)
	d := networkState(t, "net-1")

	if diags := resourceNetworkRead(context.Background(), d, client); !diags.HasError() || d.Id() != "net-1" {
		t.Fatalf("want an error and the network kept in state, got diags=%v id=%q", diags, d.Id())
	}
}
