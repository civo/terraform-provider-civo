package civo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCivoKubernetesCluster_basic(t *testing.T) {
	datasourceName := "data.civo_kubernetes_cluster.foobar"
	name := acctest.RandomWithPrefix("k8s")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoKubernetesClusterConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
					resource.TestCheckResourceAttr(datasourceName, "pools.0.node_count", "2"),
					resource.TestCheckResourceAttr(datasourceName, "pools.0.size", "g4s.kube.small"),
					resource.TestCheckResourceAttrSet(datasourceName, "pools.0.instance_names.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "kubeconfig"),
					resource.TestCheckResourceAttrSet(datasourceName, "api_endpoint"),
					resource.TestCheckResourceAttrSet(datasourceName, "master_ip"),
				),
			},
		},
	})
}

func TestAccDataSourceCivoKubernetesClusterByID_basic(t *testing.T) {
	datasourceName := "data.civo_kubernetes_cluster.foobar"
	name := acctest.RandomWithPrefix("k8s")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCivoKubernetesClusterByIDConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", name),
					resource.TestCheckResourceAttr(datasourceName, "pools.0.node_count", "2"),
					resource.TestCheckResourceAttr(datasourceName, "pools.0.size", "g4s.kube.small"),
					resource.TestCheckResourceAttrSet(datasourceName, "pools.0.instance_names.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "kubeconfig"),
					resource.TestCheckResourceAttrSet(datasourceName, "api_endpoint"),
					resource.TestCheckResourceAttrSet(datasourceName, "master_ip"),
				),
			},
		},
	})
}

func testAccDataSourceCivoKubernetesClusterConfig(name string) string {
	return fmt.Sprintf(`
data "civo_firewall" "default" {
	name = "default-default"
	region = "LON1"
}

resource "civo_kubernetes_cluster" "my-cluster" {
	name = "%s"
	firewall_id = data.civo_firewall.default.id
	pools {
		node_count = 2
		size = "g4s.kube.small"
	}
}

data "civo_kubernetes_cluster" "foobar" {
	name = civo_kubernetes_cluster.my-cluster.name
}
`, name)
}

func testAccDataSourceCivoKubernetesClusterByIDConfig(name string) string {
	return fmt.Sprintf(`
data "civo_firewall" "default" {
	name = "default-default"
	region = "LON1"
}

resource "civo_kubernetes_cluster" "my-cluster" {
	name = "%s"
	firewall_id = data.civo_firewall.default.id
	pools {
		node_count = 2
		size = "g4s.kube.small"
	}
}

data "civo_kubernetes_cluster" "foobar" {
	id = civo_kubernetes_cluster.my-cluster.id
}
`, name)
}
