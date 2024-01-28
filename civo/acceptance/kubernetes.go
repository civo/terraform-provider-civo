package acceptance

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// CivoKubernetesClusterDestroy is used to destroy the kubernetes cluster created during the test
func CivoKubernetesClusterDestroy(s *terraform.State) error {
	client := TestAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_kubernetes_cluster" {
			continue
		}

		_, err := client.GetKubernetesCluster(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("kubernetes Cluster still exists")
		}
	}

	return nil
}
