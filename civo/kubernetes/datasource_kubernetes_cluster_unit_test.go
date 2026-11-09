package kubernetes_test

import (
	"testing"

	"github.com/civo/terraform-provider-civo/civo/kubernetes"
)

// TestDataSourceKubernetesCluster_kubeconfigSensitive keeps the data source's kubeconfig,
// a cluster credential, out of plan output, matching the civo_kubernetes_cluster resource.
func TestDataSourceKubernetesCluster_kubeconfigSensitive(t *testing.T) {
	if !kubernetes.DataSourceKubernetesCluster().Schema["kubeconfig"].Sensitive {
		t.Error("kubeconfig on the civo_kubernetes_cluster data source should be Sensitive")
	}
	if !kubernetes.ResourceKubernetesCluster().Schema["kubeconfig"].Sensitive {
		t.Error("kubeconfig on the civo_kubernetes_cluster resource should be Sensitive")
	}
}
