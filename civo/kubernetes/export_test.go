package kubernetes

import (
	"github.com/civo/civogo"
)

// ExportFlattenNodePool exports flattenNodePool for testing
func ExportFlattenNodePool(cluster *civogo.KubernetesCluster) []interface{} {
	return flattenNodePool(cluster)
}
