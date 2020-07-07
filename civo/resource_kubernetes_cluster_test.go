package civo

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
					resource.TestCheckResourceAttr(resName, "num_target_nodes", "2"),
					resource.TestCheckResourceAttr(resName, "target_nodes_size", "g2.small"),
					resource.TestCheckResourceAttr(resName, "instances.0.cpu_cores", "1"),
					resource.TestCheckResourceAttr(resName, "instances.0.ram_mb", "2048"),
					resource.TestCheckResourceAttr(resName, "instances.0.disk_gb", "25"),
					// resource.TestCheckResourceAttrSet(resName, "instances"),
					resource.TestCheckResourceAttrSet(resName, "kubeconfig"),
					resource.TestCheckResourceAttrSet(resName, "api_endpoint"),
					resource.TestCheckResourceAttrSet(resName, "master_ip"),
					resource.TestCheckResourceAttrSet(resName, "dns_entry"),
					resource.TestCheckResourceAttrSet(resName, "built_at"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func TestAccCivoKubernetesClusterSize_update(t *testing.T) {
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
				Config: testAccCheckCivoKubernetesClusterConfigBasic(kubernetesClusterName),
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoKubernetesClusterResourceExists(resName, &kubernetes),
					// verify remote values
					testAccCheckCivoKubernetesClusterValues(&kubernetes, kubernetesClusterName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", kubernetesClusterName),
					resource.TestCheckResourceAttr(resName, "num_target_nodes", "2"),
					resource.TestCheckResourceAttr(resName, "target_nodes_size", "g2.small"),
					resource.TestCheckResourceAttrSet(resName, "kubeconfig"),
					resource.TestCheckResourceAttrSet(resName, "api_endpoint"),
					resource.TestCheckResourceAttrSet(resName, "master_ip"),
					resource.TestCheckResourceAttrSet(resName, "dns_entry"),
					resource.TestCheckResourceAttrSet(resName, "built_at"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoKubernetesClusterConfigSize(kubernetesClusterName),
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoKubernetesClusterResourceExists(resName, &kubernetes),
					// verify remote values
					testAccCheckCivoKubernetesClusterValues(&kubernetes, kubernetesClusterName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", kubernetesClusterName),
					resource.TestCheckResourceAttr(resName, "num_target_nodes", "4"),
					resource.TestCheckResourceAttr(resName, "target_nodes_size", "g2.small"),
					resource.TestCheckResourceAttrSet(resName, "kubeconfig"),
					resource.TestCheckResourceAttrSet(resName, "api_endpoint"),
					resource.TestCheckResourceAttrSet(resName, "master_ip"),
					resource.TestCheckResourceAttrSet(resName, "dns_entry"),
					resource.TestCheckResourceAttrSet(resName, "built_at"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func TestAccCivoKubernetesClusterTags_update(t *testing.T) {
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
				Config: testAccCheckCivoKubernetesClusterConfigBasic(kubernetesClusterName),
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoKubernetesClusterResourceExists(resName, &kubernetes),
					// verify remote values
					testAccCheckCivoKubernetesClusterValues(&kubernetes, kubernetesClusterName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", kubernetesClusterName),
					resource.TestCheckResourceAttr(resName, "num_target_nodes", "2"),
					resource.TestCheckResourceAttr(resName, "target_nodes_size", "g2.small"),
					resource.TestCheckResourceAttrSet(resName, "kubeconfig"),
					resource.TestCheckResourceAttrSet(resName, "api_endpoint"),
					resource.TestCheckResourceAttrSet(resName, "master_ip"),
					resource.TestCheckResourceAttrSet(resName, "dns_entry"),
					resource.TestCheckResourceAttrSet(resName, "built_at"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoKubernetesClusterConfigTags(kubernetesClusterName),
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoKubernetesClusterResourceExists(resName, &kubernetes),
					// verify remote values
					testAccCheckCivoKubernetesClusterValues(&kubernetes, kubernetesClusterName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "name", kubernetesClusterName),
					resource.TestCheckResourceAttr(resName, "num_target_nodes", "2"),
					resource.TestCheckResourceAttr(resName, "target_nodes_size", "g2.small"),
					resource.TestCheckResourceAttr(resName, "tags", "foo"),
					resource.TestCheckResourceAttrSet(resName, "kubeconfig"),
					resource.TestCheckResourceAttrSet(resName, "api_endpoint"),
					resource.TestCheckResourceAttrSet(resName, "master_ip"),
					resource.TestCheckResourceAttrSet(resName, "dns_entry"),
					resource.TestCheckResourceAttrSet(resName, "built_at"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
				),
			},
		},
	})
}

func testAccCheckCivoKubernetesClusterValues(kubernetes *civogo.KubernetesCluster, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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
		resp, err := client.GetKubernetesClusters(rs.Primary.ID)
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

		_, err := client.GetKubernetesClusters(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Kubernetes Cluster still exists")
		}
	}

	return nil
}

func testAccCheckCivoKubernetesClusterConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_kubernetes_cluster" "foobar" {
	name = "%s"
	num_target_nodes = 2
}`, name)
}

func testAccCheckCivoKubernetesClusterConfigSize(name string) string {
	return fmt.Sprintf(`
resource "civo_kubernetes_cluster" "foobar" {
	name = "%s"
	num_target_nodes = 4
}`, name)
}

func testAccCheckCivoKubernetesClusterConfigTags(name string) string {
	return fmt.Sprintf(`
resource "civo_kubernetes_cluster" "foobar" {
	name = "%s"
	num_target_nodes = 2
	tags = "foo"
}`, name)
}
