package kubernetes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// twoRegionAPI fakes NYC1 and LON1, with cluster "c1" holding pool "p1" only in
// clusterRegion (none when empty). Any other cluster lookup is a 404.
func twoRegionAPI(t *testing.T, clusterRegion string) *civogo.Client {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		region := req.URL.Query().Get("region")
		switch {
		case req.URL.Path == "/v2/regions":
			rw.Write([]byte(`[{"code":"NYC1"},{"code":"LON1"}]`))
		case req.URL.Path == "/v2/kubernetes/clusters/c1" && region == clusterRegion:
			rw.Write([]byte(`{"id":"c1"}`))
		case req.URL.Path == "/v2/kubernetes/clusters/c1/pools/p1" && region == clusterRegion:
			rw.Write([]byte(`{"id":"p1","count":2,"size":"g4s.kube.small"}`))
		case strings.HasPrefix(req.URL.Path, "/v2/kubernetes/clusters/"):
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(`{"code":"database_kubernetes_cluster_not_found","reason":"not found"}`))
		default:
			t.Errorf("unexpected request %s", req.URL)
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}))
	t.Cleanup(server.Close)

	client, err := civogo.NewClientForTestingWithServer(server)
	if err != nil {
		t.Fatalf("could not build the test client: %v", err)
	}
	client.Region = "NYC1"
	return client
}

// A pool in state from before `region` existed, whose cluster was deleted
// outside Terraform, must drop out of state so the plan recreates it.
func TestNodePoolReadLegacyStateClusterGone(t *testing.T) {
	client := twoRegionAPI(t, "")
	d := schema.TestResourceDataRaw(t, nodePoolSchema(true), map[string]interface{}{"cluster_id": "c1"})
	d.SetId("p1")

	if diags := resourceKubernetesClusterNodePoolRead(context.Background(), d, client); len(diags) != 0 || d.Id() != "" {
		t.Fatalf("want pool cleared from state, got diags=%v id=%q", diags, d.Id())
	}
}

// With no `region`, the pool must follow its cluster into a region other than
// the provider's (#415).
func TestNodePoolReadInfersRegionFromCluster(t *testing.T) {
	client := twoRegionAPI(t, "LON1")
	d := schema.TestResourceDataRaw(t, nodePoolSchema(true), map[string]interface{}{"cluster_id": "c1"})
	d.SetId("p1")

	if diags := resourceKubernetesClusterNodePoolRead(context.Background(), d, client); len(diags) != 0 {
		t.Fatalf("unexpected diags: %v", diags)
	}
	if d.Id() != "p1" {
		t.Fatalf("pool dropped from state; id = %q", d.Id())
	}
	if got := d.Get("region").(string); got != "LON1" {
		t.Errorf("region = %q, want LON1 (the cluster's)", got)
	}
	if got := d.Get("node_count").(int); got != 2 {
		t.Errorf("node_count = %d, want 2 (read from LON1)", got)
	}
}

// poolLookupAPI fakes cluster "c1" listing the given pools, while every pool
// lookup answers with poolStatus.
func poolLookupAPI(t *testing.T, clusterPools string, poolStatus int) *civogo.Client {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch {
		case req.URL.Path == "/v2/kubernetes/clusters/c1":
			_, _ = rw.Write([]byte(`{"id":"c1","pools":` + clusterPools + `}`))
		case strings.HasPrefix(req.URL.Path, "/v2/kubernetes/clusters/c1/pools/"):
			rw.WriteHeader(poolStatus)
			_, _ = rw.Write([]byte(`{"code":"database_cluster_pool_not_found","reason":"not found"}`))
		default:
			t.Errorf("unexpected request %s", req.URL)
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}))
	t.Cleanup(server.Close)

	client, err := civogo.NewClientForTestingWithServer(server)
	if err != nil {
		t.Fatalf("could not build the test client: %v", err)
	}
	client.Region = "NYC1"
	return client
}

// A pool deleted outside Terraform must drop out of state so the plan
// recreates it, instead of failing the refresh.
func TestNodePoolReadPoolDeletedOutsideTerraform(t *testing.T) {
	client := poolLookupAPI(t, `[{"id":"p2"}]`, http.StatusNotFound)
	d := schema.TestResourceDataRaw(t, nodePoolSchema(true), map[string]interface{}{"cluster_id": "c1", "region": "NYC1"})
	d.SetId("p1")

	if diags := resourceKubernetesClusterNodePoolRead(context.Background(), d, client); len(diags) != 0 || d.Id() != "" {
		t.Fatalf("want pool cleared from state, got diags=%v id=%q", diags, d.Id())
	}
}

// A failed lookup of a pool the cluster still lists is an API error, and the
// pool must stay in state.
func TestNodePoolReadPoolLookupFails(t *testing.T) {
	client := poolLookupAPI(t, `[{"id":"p1"}]`, http.StatusInternalServerError)
	d := schema.TestResourceDataRaw(t, nodePoolSchema(true), map[string]interface{}{"cluster_id": "c1", "region": "NYC1"})
	d.SetId("p1")

	if diags := resourceKubernetesClusterNodePoolRead(context.Background(), d, client); !diags.HasError() || d.Id() != "p1" {
		t.Fatalf("want an error and the pool kept in state, got diags=%v id=%q", diags, d.Id())
	}
}
