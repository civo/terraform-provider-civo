package kubernetes_test

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/civo/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// example.Widget represents a concrete Go type that represents an API resource
func CivoKubernetesCluster_basic(t *testing.T) {
	var kubernetes civogo.KubernetesCluster

	// generate a random name for each test run
	resName := "civo_kubernetes_cluster.foobar"
	var kubernetesClusterName = acctest.RandomWithPrefix("tf-test") + ".example"

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
					// verify remote values
					CivoKubernetesClusterValues(&kubernetes, kubernetesClusterName),
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
					resource.TestCheckResourceAttrSet(resName, "cluster_type"),
				),
			},
		},
	})
}

func CivoKubernetesClusterCNI(t *testing.T) {
	var kubernetes civogo.KubernetesCluster

	// generate a random name for each test run
	resName := "civo_kubernetes_cluster.foobar"
	var kubernetesClusterName = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CivoInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: CivoKubernetesClusterConfigCNI(kubernetesClusterName),
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					CivoKubernetesClusterResourceExists(resName, &kubernetes),
					// verify remote values
					CivoKubernetesClusterValues(&kubernetes, kubernetesClusterName),
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
					resource.TestCheckResourceAttrSet(resName, "cluster_type"),
				),
			},
		},
	})
}

func CivoKubernetesClusterValues(kubernetes *civogo.KubernetesCluster, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if kubernetes.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, kubernetes.Name)
		}
		return nil
	}
}

// CheckExampleResourceExists queries the API and retrieves the matching Widget.
func CivoKubernetesClusterResourceExists(n string, kubernetes *civogo.KubernetesCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := acceptance.TestAccProvider.Meta().(*civogo.Client)
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

func CivoKubernetesClusterConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_firewall" "default" {
	name = "%s"
	create_default_rules = true
	region = "FAKE"
}

resource "civo_kubernetes_cluster" "foobar" {
	name = "%s"
	firewall_id = civo_firewall.default.id
	pools {
		node_count = 2
		size = "g4s.kube.small"
	}
}`, name, name)
}

func CivoKubernetesClusterConfigCNI(name string) string {
	return fmt.Sprintf(`
resource "civo_firewall" "default" {
	name = "%s"
	create_default_rules = true
	region = "FAKE"
}

resource "civo_kubernetes_cluster" "foobar" {
	name = "%s"
	firewall_id = civo_firewall.default.id
	region = "FAKE"
	pools {
		node_count = 2
		size = "g4s.kube.small"
	}
	cni = "cilium"
}`, name, name)
}
