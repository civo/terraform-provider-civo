package kubernetes

import (
	"github.com/civo/civogo"
	corev1 "k8s.io/api/core/v1"
)

// ExportFlattenNodePool exports flattenNodePool for testing
func ExportFlattenNodePool(cluster *civogo.KubernetesCluster) []interface{} {
	return flattenNodePool(cluster)
}

// ExportUpdateNodePool exports updateNodePool for testing
func ExportUpdateNodePool(s []civogo.KubernetesClusterPoolConfig, id string, count int, labels map[string]string, taints []corev1.Taint) []civogo.KubernetesClusterPoolConfig {
	return updateNodePool(s, id, count, labels, taints)
}
