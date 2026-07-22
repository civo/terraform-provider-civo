package kubernetes_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/civo/terraform-provider-civo/civo/kubernetes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	corev1 "k8s.io/api/core/v1"
)

// example.Widget represents a concrete Go type that represents an API resource
func TestAccCivoKubernetesClusterNodePool_basic(t *testing.T) {
	var kubernetes civogo.KubernetesCluster
	var kubernetesNodePool civogo.KubernetesPool

	// generate a random name for each test run
	resName := "civo_kubernetes_cluster.foobar"
	resPoolName := "civo_kubernetes_node_pool.foobar"
	var kubernetesClusterName = acctest.RandomWithPrefix("tf-test") + "-example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CivoKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: CivoKubernetesClusterConfigBasic(kubernetesClusterName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					CivoKubernetesClusterResourceExists(resName, &kubernetes),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: CivoKubernetesClusterConfigBasic(kubernetesClusterName) + CivoKubernetesClusterNodePoolConfigBasic(),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					CivoKubernetesClusterNodePoolResourceExists(resPoolName, &kubernetes, &kubernetesNodePool),
					// verify remote values
					CivoKubernetesClusterNodePoolValues(&kubernetesNodePool, "g4s.kube.small"),
					// verify local values
					// resource.TestCheckResourceAttr(resPoolName, "cluster_id", kubernetes.ID),
					resource.TestCheckResourceAttr(resPoolName, "node_count", "3"),
					resource.TestCheckResourceAttr(resPoolName, "size", "g4s.kube.small"),
				),
			},
		},
	})
}

func CivoKubernetesClusterNodePoolValues(kubernetes *civogo.KubernetesPool, value string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if kubernetes.Size != value {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", value, kubernetes.Size)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func CivoKubernetesClusterNodePoolResourceExists(n string, kubernetes *civogo.KubernetesCluster, pool *civogo.KubernetesPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := acceptance.TestAccProvider.Meta().(*civogo.Client)
		resp, err := client.GetKubernetesCluster(kubernetes.ID)
		if err != nil {
			return fmt.Errorf("Kuberenetes Cluster not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		var id int
		for k, v := range resp.Pools {
			if v.ID == rs.Primary.ID {
				id = k
				break
			}
		}

		*pool = resp.Pools[id]

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func CivoKubernetesClusterNodePoolConfigBasic() string {
	return `
resource "civo_kubernetes_node_pool" "foobar" {
	cluster_id = civo_kubernetes_cluster.foobar.id
	node_count = 3
	size = "g4s.kube.small"
	region = "LON1"
	depends_on = [civo_kubernetes_cluster.foobar]
}`
}

func TestFlattenNodePool(t *testing.T) {
	cases := []struct {
		name     string
		cluster  *civogo.KubernetesCluster
		expected []interface{}
	}{
		{
			name: "Creating a cluster with labels+taints on the default pool",
			cluster: &civogo.KubernetesCluster{
				Pools: []civogo.KubernetesPool{
					{
						ID:               "pool-1",
						Count:            3,
						Size:             "g4s.kube.medium",
						PublicIPNodePool: true,
						InstanceNames:    []string{"node-1", "node-2", "node-3"},
						Labels: map[string]string{
							"env": "production",
						},
						Taints: []corev1.Taint{
							{
								Key:    "gpu",
								Value:  "true",
								Effect: "NoSchedule",
							},
						},
					},
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"label":               "pool-1",
					"node_count":          3,
					"size":                "g4s.kube.medium",
					"instance_names":      []string{"node-1", "node-2", "node-3"},
					"public_ip_node_pool": true,
					"labels": map[string]string{
						"env": "production",
					},
					"taint": []map[string]interface{}{
						{
							"key":    "gpu",
							"value":  "true",
							"effect": "NoSchedule",
						},
					},
				},
			},
		},
		{
			name: "Verifying no drift on a second plan/Removing labels+taints (empty in API)",
			cluster: &civogo.KubernetesCluster{
				Pools: []civogo.KubernetesPool{
					{
						ID:               "pool-1",
						Count:            3,
						Size:             "g4s.kube.medium",
						PublicIPNodePool: false,
						InstanceNames:    []string{"node-1"},
						Labels:           map[string]string{},
						Taints:           []corev1.Taint{},
					},
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"label":               "pool-1",
					"node_count":          3,
					"size":                "g4s.kube.medium",
					"instance_names":      []string{"node-1"},
					"public_ip_node_pool": false,
					// When labels/taints are empty in API, flattenNodePool typically omits them
					// which correctly clears them out of Terraform state avoiding drift if they were removed.
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := kubernetes.ExportFlattenNodePool(tc.cluster)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("expected: %#v, got: %#v", tc.expected, actual)
			}
		})
	}
}
