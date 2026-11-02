package kubernetes_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/civo/terraform-provider-civo/civo/kubernetes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDataSourceCivoKubernetesClusters_basic(t *testing.T) {
	datasourceName := "data.civo_kubernetes_clusters.foobar"
	name := acctest.RandomWithPrefix("k8s")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceCivoKubernetesClustersConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "clusters.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "clusters.0.name", name),
					resource.TestCheckResourceAttrPair(datasourceName, "clusters.0.id", "civo_kubernetes_cluster.my-cluster", "id"),
					resource.TestCheckResourceAttr(datasourceName, "clusters.0.region", "LON1"),
					resource.TestCheckResourceAttr(datasourceName, "clusters.0.pools.0.node_count", "2"),
					resource.TestCheckResourceAttrSet(datasourceName, "clusters.0.api_endpoint"),
					resource.TestCheckResourceAttrSet(datasourceName, "clusters.0.master_ip"),
				),
			},
		},
	})
}

func DataSourceCivoKubernetesClustersConfig(name string) string {
	return fmt.Sprintf(`
data "civo_firewall" "default" {
	name = "default-default"
	region = "LON1"
}

resource "civo_kubernetes_cluster" "my-cluster" {
	name = "%s"
	region = "LON1"
	firewall_id = data.civo_firewall.default.id
	pools {
		node_count = 2
		size = "g4s.kube.small"
	}
}

data "civo_kubernetes_clusters" "foobar" {
	region = "LON1"
	filter {
		key = "name"
		values = [civo_kubernetes_cluster.my-cluster.name]
	}
}
`, name)
}

const fakeClustersResponse = `{
	"page": 1, "per_page": 100, "pages": 1,
	"items": [
		{
			"id": "11111111-0000-0000-0000-000000000001", "name": "alpha", "cluster_type": "k3s",
			"status": "ACTIVE", "ready": true, "kubernetes_version": "1.30.5-k3s1", "cni_plugin": "flannel",
			"api_endpoint": "https://10.0.0.1:6443", "master_ip": "10.0.0.1", "dns_entry": "alpha.k8s.civo.com",
			"kubeconfig": "apiVersion: v1 # secret", "tags": ["prod"], "network_id": "net-1", "firewall_id": "fw-1",
			"pools": [{"id": "pool-a", "count": 2, "size": "g4s.kube.small", "instance_names": ["alpha-node-1", "alpha-node-2"],
				"labels": {"team": "stale", "kubernetes.civo.com/civo-node-pool": "pool-a"},
				"taints": [{"key": "stale", "value": "x", "effect": "NoExecute"}]}],
			"required_pools": [
				{"id": "pool-a", "count": 2, "size": "g4s.kube.small",
					"labels": {"team": "payments", "kubernetes.civo.com/civo-node-pool": "pool-a"},
					"taints": [{"key": "dedicated", "value": "gpu", "effect": "NoSchedule"}]},
				{"id": "pool-b", "count": 1, "size": "g4s.kube.large", "labels": {"tier": "gpu"}}
			],
			"installed_applications": [{"application": "traefik2-nodeport", "name": "traefik2-nodeport", "version": "2.9.4", "installed": true, "category": "architecture"}]
		},
		{
			"id": "22222222-0000-0000-0000-000000000002", "name": "beta", "cluster_type": "talos",
			"status": "ACTIVE", "ready": false, "api_endpoint": "https://10.0.0.2:6443"
		}
	]
}`

const fakeRegionsResponse = `[
	{"code": "LON1", "name": "London", "default": false},
	{"code": "NYC1", "name": "New York", "default": true}
]`

// fakeClustersAPI serves the cluster and region lists and records the region query
// parameter of each cluster request. clientRegion is the region the provider was configured with.
func fakeClustersAPI(t *testing.T, clientRegion string) (*civogo.Client, func() []string) {
	t.Helper()

	var mu sync.Mutex
	var regions []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v2/regions":
			_, _ = w.Write([]byte(fakeRegionsResponse))
		case "/v2/kubernetes/clusters":
			mu.Lock()
			regions = append(regions, r.URL.Query().Get("region"))
			mu.Unlock()
			_, _ = w.Write([]byte(fakeClustersResponse))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := civogo.NewClientWithURL("fake-key", srv.URL, clientRegion)
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	return client, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return append([]string(nil), regions...)
	}
}

func readClusters(t *testing.T, client *civogo.Client, config map[string]interface{}) *schema.ResourceData {
	t.Helper()

	ds := kubernetes.DataSourceKubernetesClusters()
	d := schema.TestResourceDataRaw(t, ds.Schema, config)
	if diags := ds.ReadContext(context.Background(), d, client); diags.HasError() {
		t.Fatalf("read failed: %v", diags)
	}
	return d
}

func clusterNames(d *schema.ResourceData) []string {
	names := []string{}
	for i := 0; i < d.Get("clusters.#").(int); i++ {
		names = append(names, d.Get(fmt.Sprintf("clusters.%d.name", i)).(string))
	}
	return names
}

func TestDataSourceKubernetesClusters_listsAllClusters(t *testing.T) {
	client, regions := fakeClustersAPI(t, "LON1")

	d := readClusters(t, client, map[string]interface{}{})

	if got := d.Get("clusters.#").(int); got != 2 {
		t.Fatalf("expected 2 clusters, got %d", got)
	}
	want := map[string]interface{}{
		"clusters.0.name":                                 "alpha",
		"clusters.0.id":                                   "11111111-0000-0000-0000-000000000001",
		"clusters.0.region":                               "LON1",
		"clusters.0.cluster_type":                         "k3s",
		"clusters.0.api_endpoint":                         "https://10.0.0.1:6443",
		"clusters.0.ready":                                true,
		"clusters.0.tags.#":                               1,
		"clusters.0.pools.0.label":                        "pool-a",
		"clusters.0.pools.0.node_count":                   2,
		"clusters.0.pools.0.instance_names.#":             2,
		"clusters.0.pools.0.labels.%":                     1,
		"clusters.0.pools.0.labels.team":                  "payments",
		"clusters.0.pools.0.taint.#":                      1,
		"clusters.0.pools.#":                              2,
		"clusters.0.pools.1.label":                        "pool-b",
		"clusters.0.pools.1.node_count":                   1,
		"clusters.0.pools.1.size":                         "g4s.kube.large",
		"clusters.0.pools.1.instance_names.#":             0,
		"clusters.0.pools.1.labels.tier":                  "gpu",
		"clusters.0.pools.1.taint.#":                      0,
		"clusters.0.installed_applications.0.application": "traefik2-nodeport",
		"clusters.1.name":                                 "beta",
		"clusters.1.cluster_type":                         "talos",
	}
	for key, value := range want {
		// compare printed values: the SDK returns counts such as "labels.%" as a different type than int
		if got := d.Get(key); fmt.Sprint(got) != fmt.Sprint(value) {
			t.Errorf("%s = %v, want %v", key, got, value)
		}
	}
	if got := regions(); len(got) != 1 || got[0] != "LON1" {
		t.Errorf("expected one request scoped to the provider's region LON1, got %v", got)
	}

	// the labels and taint set by the user come from the matching required pool, not the built pool
	taints := d.Get("clusters.0.pools.0.taint").(*schema.Set).List()
	if len(taints) != 1 {
		t.Fatalf("expected 1 taint, got %v", taints)
	}
	want = map[string]interface{}{"key": "dedicated", "value": "gpu", "effect": "NoSchedule"}
	if got := taints[0].(map[string]interface{}); !reflect.DeepEqual(got, want) {
		t.Errorf("taint = %v, want %v", got, want)
	}
}

func TestDataSourceKubernetesClusters_filters(t *testing.T) {
	tests := []struct {
		name   string
		filter map[string]interface{}
		want   []string
	}{
		{name: "string field", filter: map[string]interface{}{"key": "name", "values": []interface{}{"beta"}}, want: []string{"beta"}},
		{name: "bool field", filter: map[string]interface{}{"key": "ready", "values": []interface{}{"true"}}, want: []string{"alpha"}},
		{name: "set field", filter: map[string]interface{}{"key": "tags", "values": []interface{}{"prod"}}, want: []string{"alpha"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := fakeClustersAPI(t, "LON1")

			d := readClusters(t, client, map[string]interface{}{"filter": []interface{}{tt.filter}})

			if got := clusterNames(d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got clusters %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataSourceKubernetesClusters_sort(t *testing.T) {
	client, _ := fakeClustersAPI(t, "LON1")

	d := readClusters(t, client, map[string]interface{}{
		"sort": []interface{}{map[string]interface{}{"key": "name", "direction": "desc"}},
	})

	if got, want := clusterNames(d), []string{"beta", "alpha"}; !reflect.DeepEqual(got, want) {
		t.Errorf("got clusters %v, want %v", got, want)
	}
}

func TestDataSourceKubernetesClusters_region(t *testing.T) {
	client, regions := fakeClustersAPI(t, "LON1")

	d := readClusters(t, client, map[string]interface{}{"region": "FRA1"})

	if got := regions(); len(got) != 1 || got[0] != "FRA1" {
		t.Errorf("expected one request scoped to FRA1, got %v", got)
	}
	if got := d.Get("clusters.0.region"); got != "FRA1" {
		t.Errorf("expected clusters to report region FRA1, got %v", got)
	}
}

func TestDataSourceKubernetesClusters_defaultRegion(t *testing.T) {
	// No region on the provider or the data source: the API answers from the
	// account's default region, which should be looked up and reported
	client, regions := fakeClustersAPI(t, "")

	d := readClusters(t, client, map[string]interface{}{})

	if got := regions(); len(got) != 1 || got[0] != "" {
		t.Errorf("expected one unscoped cluster request, got %v", got)
	}
	if got := d.Get("clusters.0.region"); got != "NYC1" {
		t.Errorf("expected clusters to report the account default region NYC1, got %v", got)
	}
}

func TestDataSourceKubernetesClusters_defaultRegionLookupFails(t *testing.T) {
	// The region is informational, so a failed default region lookup must not fail the read
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/kubernetes/clusters" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(fakeClustersResponse))
			return
		}
		http.Error(w, `{"code":"internal_error"}`, http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)
	client, err := civogo.NewClientWithURL("fake-key", srv.URL, "")
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	d := readClusters(t, client, map[string]interface{}{})

	if got := d.Get("clusters.#").(int); got != 2 {
		t.Fatalf("expected 2 clusters, got %d", got)
	}
	if got := d.Get("clusters.0.region"); got != "" {
		t.Errorf("expected an empty region when the lookup fails, got %v", got)
	}
}

func TestDataSourceKubernetesClusters_schema(t *testing.T) {
	ds := kubernetes.DataSourceKubernetesClusters()
	record := ds.Schema["clusters"].Elem.(*schema.Resource).Schema
	if _, ok := record["kubeconfig"]; ok {
		t.Error("the clusters list should not include kubeconfig, which is a credential")
	}

	// nested object fields can't be filtered on, so they should not be offered as filter keys
	filterKey := ds.Schema["filter"].Elem.(*schema.Resource).Schema["key"]
	for _, key := range []string{"pools", "installed_applications"} {
		if _, errs := filterKey.ValidateFunc(key, "filter.0.key"); len(errs) == 0 {
			t.Errorf("expected %q to be rejected as a filter key", key)
		}
	}
	if _, errs := filterKey.ValidateFunc("tags", "filter.0.key"); len(errs) != 0 {
		t.Errorf("expected tags to be a valid filter key, got %v", errs)
	}
}

func TestDataSourceKubernetesClusters_emptyAccountSkipsRegionLookup(t *testing.T) {
	var mu sync.Mutex
	regionLookups := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v2/regions":
			mu.Lock()
			regionLookups++
			mu.Unlock()
			_, _ = w.Write([]byte(fakeRegionsResponse))
		case "/v2/kubernetes/clusters":
			_, _ = w.Write([]byte(`{"page": 1, "per_page": 100, "pages": 1, "items": []}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	client, err := civogo.NewClientWithURL("fake-key", srv.URL, "")
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	d := readClusters(t, client, map[string]interface{}{})

	if got := d.Get("clusters.#").(int); got != 0 {
		t.Errorf("expected no clusters, got %d", got)
	}
	mu.Lock()
	defer mu.Unlock()
	if regionLookups != 0 {
		t.Errorf("expected no default region lookup when there are no clusters, got %d", regionLookups)
	}
}

func TestDataSourceKubernetesCluster_poolLabelsAndTaints(t *testing.T) {
	// civo_kubernetes_cluster shares flattenDataSourceNodePool, so it gets the same pools
	client, _ := fakeClustersAPI(t, "LON1")

	ds := kubernetes.DataSourceKubernetesCluster()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"name": "alpha"})
	if diags := ds.ReadContext(context.Background(), d, client); diags.HasError() {
		t.Fatalf("read failed: %v", diags)
	}

	want := map[string]interface{}{
		"pools.#":             2,
		"pools.0.label":       "pool-a",
		"pools.0.labels.%":    1,
		"pools.0.labels.team": "payments",
		"pools.0.taint.#":     1,
		"pools.1.label":       "pool-b",
		"pools.1.labels.tier": "gpu",
		"pools.1.node_count":  1,
	}
	for key, value := range want {
		if got := d.Get(key); fmt.Sprint(got) != fmt.Sprint(value) {
			t.Errorf("%s = %v, want %v", key, got, value)
		}
	}
}
