package civo

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoKubernetesClusterConfigBasic(kubernetesClusterName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoKubernetesClusterResourceExists(resName, &kubernetes),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoKubernetesClusterConfigBasic(kubernetesClusterName) + testAccCheckCivoKubernetesClusterNodePoolConfigBasic(),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoKubernetesClusterNodePoolResourceExists(resPoolName, &kubernetes, &kubernetesNodePool),
					// verify remote values
					testAccCheckCivoKubernetesClusterNodePoolValues(&kubernetesNodePool, "g4s.kube.small"),
					// verify local values
					// resource.TestCheckResourceAttr(resPoolName, "cluster_id", kubernetes.ID),
					resource.TestCheckResourceAttr(resPoolName, "node_count", "3"),
					resource.TestCheckResourceAttr(resPoolName, "size", "g4s.kube.small"),
				),
			},
		},
	})
}

func testAccCheckCivoKubernetesClusterNodePoolValues(kubernetes *civogo.KubernetesPool, value string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if kubernetes.Size != value {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", value, kubernetes.Size)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoKubernetesClusterNodePoolResourceExists(n string, kubernetes *civogo.KubernetesCluster, pool *civogo.KubernetesPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
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

func testAccCheckCivoKubernetesClusterNodePoolConfigBasic() string {
	return `
resource "civo_kubernetes_node_pool" "foobar" {
	cluster_id = civo_kubernetes_cluster.foobar.id
	node_count = 3
	size = "g4s.kube.small"
	region = "LON1"
	depends_on = [civo_kubernetes_cluster.foobar]
}`
}
