package network

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// networkAPI is a minimal fake of the VPC network endpoints used by
// resourceNetworkCreate/Read. It records the body of the network create
// request so the tests can check what was sent to the API.
func networkAPI(t *testing.T, network string) (*civogo.Client, *map[string]interface{}) {
	t.Helper()
	created := map[string]interface{}{}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch {
		case req.Method == http.MethodPost && req.URL.Path == "/v2/vpc/networks":
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Errorf("failed to read request body: %s", err)
			}
			if err := json.Unmarshal(body, &created); err != nil {
				t.Errorf("failed to decode request body %q: %s", body, err)
			}
			rw.Write([]byte(`{"id":"net-1234","label":"my-net","result":"success"}`))
		case req.Method == http.MethodPost && req.URL.Path == "/v2/vpc/firewalls":
			rw.Write([]byte(`{"id":"fw-1234","name":"my-net-default","result":"success"}`))
		case req.Method == http.MethodGet && req.URL.Path == "/v2/vpc/networks":
			rw.Write([]byte("[" + network + "]"))
		default:
			t.Errorf("unexpected request %s %s", req.Method, req.URL.Path)
			rw.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	client, err := civogo.NewClientForTestingWithServer(server)
	if err != nil {
		t.Fatalf("failed to build test client: %s", err)
	}
	return client, &created
}

// TestResourceNetworkCreate_ipv6 covers civo/terraform-provider-civo#416:
// civo_vpc_subnet can only be created in a network with IPv6 enabled, but the
// resource had no way to ask the API for one.
func TestResourceNetworkCreate_ipv6(t *testing.T) {
	client, created := networkAPI(t, `{"id":"net-1234","name":"my-net","label":"my-net",
		"cidr":"10.0.0.0/24","ipv6_enabled":true,"cidr_v6":"2a00:1234:5678:100::/56",
		"nameservers_v6":["2001:4860:4860::8888"]}`)

	d := schema.TestResourceDataRaw(t, ResourceNetwork().Schema, map[string]interface{}{
		"label":          "my-net",
		"region":         "LON1",
		"ipv6_enabled":   true,
		"nameservers_v6": []interface{}{"2001:4860:4860::8888"},
	})

	if diags := resourceNetworkCreate(context.Background(), d, client); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}

	if got := (*created)["ipv6_enabled"]; got != true {
		t.Errorf("request ipv6_enabled = %v, want true", got)
	}
	if got, want := (*created)["nameservers_v6"], []interface{}{"2001:4860:4860::8888"}; !reflect.DeepEqual(got, want) {
		t.Errorf("request nameservers_v6 = %v, want %v", got, want)
	}

	if got := d.Get("ipv6_enabled").(bool); !got {
		t.Errorf("ipv6_enabled = %v, want true", got)
	}
	if got, want := d.Get("cidr_v6").(string), "2a00:1234:5678:100::/56"; got != want {
		t.Errorf("cidr_v6 = %q, want %q", got, want)
	}
	if got, want := d.Get("nameservers_v6").([]interface{}), []interface{}{"2001:4860:4860::8888"}; !reflect.DeepEqual(got, want) {
		t.Errorf("nameservers_v6 = %v, want %v", got, want)
	}
}

// TestResourceNetworkCreate_ipv4Only checks that a network without
// ipv6_enabled sends the same request as before (no IPv6 settings).
func TestResourceNetworkCreate_ipv4Only(t *testing.T) {
	client, created := networkAPI(t, `{"id":"net-1234","name":"my-net","label":"my-net","cidr":"10.0.0.0/24"}`)

	d := schema.TestResourceDataRaw(t, ResourceNetwork().Schema, map[string]interface{}{
		"label":  "my-net",
		"region": "LON1",
	})

	if diags := resourceNetworkCreate(context.Background(), d, client); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}

	if got := (*created)["ipv6_enabled"]; got != nil {
		t.Errorf("request ipv6_enabled = %v, want it unset", got)
	}
	if got := (*created)["nameservers_v6"]; got != nil {
		t.Errorf("request nameservers_v6 = %v, want it unset", got)
	}
	if got := d.Get("ipv6_enabled").(bool); got {
		t.Errorf("ipv6_enabled = %v, want false", got)
	}
	if got := d.Get("cidr_v6").(string); got != "" {
		t.Errorf("cidr_v6 = %q, want empty", got)
	}
}
