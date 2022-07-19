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
func TestAccCivoKubernetesCluster_basic(t *testing.T) {
	var kubernetes civogo.KubernetesCluster

	// generate a random name for each test run
	resName := "civo_kubernetes_cluster.foobar"
	var kubernetesClusterName = acctest.RandomWithPrefix("tf-test") + ".example"

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
					// verify remote values
					testAccCheckCivoKubernetesClusterValues(&kubernetes, kubernetesClusterName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", kubernetesClusterName),

					// resource.TestCheckResourceAttrSet(resName, "instances"),
					resource.TestCheckResourceAttr(resName, "pools.0.node_count", "2"),
					resource.TestCheckResourceAttr(resName, "pools.0.size", "g4s.kube.small"),
					resource.TestCheckResourceAttrSet(resName, "kubeconfig"),
					resource.TestCheckResourceAttrSet(resName, "api_endpoint"),
					resource.TestCheckResourceAttrSet(resName, "master_ip"),
					resource.TestCheckResourceAttrSet(resName, "dns_entry"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func TestAccCivoKubernetesClusterCNI(t *testing.T) {
	var kubernetes civogo.KubernetesCluster

	// generate a random name for each test run
	resName := "civo_kubernetes_cluster.foobar"
	var kubernetesClusterName = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCivoKubernetesClusterConfigCNI(kubernetesClusterName),
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoKubernetesClusterResourceExists(resName, &kubernetes),
					// verify remote values
					testAccCheckCivoKubernetesClusterValues(&kubernetes, kubernetesClusterName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", kubernetesClusterName),
					resource.TestCheckResourceAttr(resName, "cni", "cilium"),
					resource.TestCheckResourceAttr(resName, "pools.0.node_count", "2"),
					resource.TestCheckResourceAttr(resName, "pools.0.size", "g4s.kube.small"),
					resource.TestCheckResourceAttrSet(resName, "kubeconfig"),
					resource.TestCheckResourceAttrSet(resName, "api_endpoint"),
					resource.TestCheckResourceAttrSet(resName, "master_ip"),
					resource.TestCheckResourceAttrSet(resName, "dns_entry"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func testAccCheckCivoKubernetesClusterValues(kubernetes *civogo.KubernetesCluster, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if kubernetes.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, kubernetes.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoKubernetesClusterResourceExists(n string, kubernetes *civogo.KubernetesCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.GetKubernetesCluster(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Kuberenetes Cluster not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*kubernetes = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

func testAccCheckCivoKubernetesClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_kubernetes_cluster" {
			continue
		}

		_, err := client.GetKubernetesCluster(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Kubernetes Cluster still exists")
		}
	}

	return nil
}

func testAccCheckCivoKubernetesClusterConfigBasic(name string) string {
	return fmt.Sprintf(`
data "civo_firewall" "default" {
	name = "default-default"
	region = "LON1"
}

resource "civo_kubernetes_cluster" "foobar" {
	name = "%s"
	firewall_id = data.civo_firewall.default.id
	pools {
		node_count = 2
		size = "g4s.kube.small"
	}
}`, name)
}

func testAccCheckCivoKubernetesClusterConfigCNI(name string) string {
	return fmt.Sprintf(`
data "civo_firewall" "default" {
	name = "default-default"
	region = "LON1"
}

resource "civo_kubernetes_cluster" "foobar" {
	name = "%s"
	firewall_id = data.civo_firewall.default.id
	pools {
		node_count = 2
		size = "g4s.kube.small"
	}
	cni = "cilium"
}`, name)
}
