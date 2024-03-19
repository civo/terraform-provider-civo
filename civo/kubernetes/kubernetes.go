package kubernetes

import (
	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	corev1 "k8s.io/api/core/v1"
)

// nodePoolSchema function to define the node pool schema
func nodePoolSchema(isResource bool) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"label": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: utils.ValidateNameOnlyContainsAlphanumericCharacters,
			Description:      "Node pool label, if you don't provide one, we will generate one for you",
		},
		"node_count": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Number of nodes in the nodepool",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"size": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Size of the nodes in the nodepool",
		},
		"instance_names": {
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Instance names in the nodepool",
		},
		"public_ip_node_pool": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Node pool belongs to the public ip node pool",
		},
		"labels": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"taint": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
					"effect": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"NoSchedule",
							"PreferNoSchedule",
							"NoExecute",
						}, false),
					},
				},
			},
		},
	}

	if isResource {
		// add the cluster id to the schema
		s["cluster_id"] = &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			Description:  "The ID of your cluster",
			ValidateFunc: validation.StringIsNotEmpty,
		}
	}

	return s
}

// function to flatten all instances inside the cluster
func flattenNodePool(cluster *civogo.KubernetesCluster) []interface{} {
	if cluster.Pools == nil {
		return nil
	}

	flattenedPool := make([]interface{}, 0)

	poolInstanceNames := make([]string, 0)
	poolInstanceNames = append(poolInstanceNames, cluster.Pools[0].InstanceNames...)

	rawPool := map[string]interface{}{
		"label":               cluster.Pools[0].ID,
		"node_count":          cluster.Pools[0].Count,
		"size":                cluster.Pools[0].Size,
		"instance_names":      poolInstanceNames,
		"public_ip_node_pool": cluster.Pools[0].PublicIPNodePool,
	}

	flattenedPool = append(flattenedPool, rawPool)

	return flattenedPool
}

// function to flatten all applications inside the cluster
func flattenInstalledApplication(apps []civogo.KubernetesInstalledApplication) []interface{} {
	if apps == nil {
		return nil
	}

	flattenedInstalledApplication := make([]interface{}, 0)
	for _, app := range apps {
		rawInstalledApplication := map[string]interface{}{
			"application": app.Name,
			"version":     app.Version,
			"installed":   app.Installed,
			"category":    app.Category,
		}

		flattenedInstalledApplication = append(flattenedInstalledApplication, rawInstalledApplication)
	}

	return flattenedInstalledApplication
}

// expandNodePools function to expand the node pools
func expandNodePools(nodePools []interface{}) []civogo.KubernetesClusterPoolConfig {
	expandedNodePools := make([]civogo.KubernetesClusterPoolConfig, 0, len(nodePools))
	for _, rawPool := range nodePools {
		pool := rawPool.(map[string]interface{})

		poolID := uuid.NewString()
		if pool["label"].(string) != "" {
			poolID = pool["label"].(string)
		}

		// Initialize labels map only if they are provided and valid
		var labels map[string]string
		if rawLabels, ok := pool["labels"].(map[string]interface{}); ok {
			labels = make(map[string]string, len(rawLabels))
			for k, v := range rawLabels {
				if strVal, ok := v.(string); ok {
					labels[k] = strVal
				}
			}
		}

		// Initialize taints slice only if they are provided and valid
		var taints []corev1.Taint
		if taintSet, ok := pool["taint"].(*schema.Set); ok {
			for _, taintInterface := range taintSet.List() {
				taintMap := taintInterface.(map[string]interface{})
				taints = append(taints, corev1.Taint{
					Key:    taintMap["key"].(string),
					Value:  taintMap["value"].(string),
					Effect: corev1.TaintEffect(taintMap["effect"].(string)),
				})
			}
		}

		cr := civogo.KubernetesClusterPoolConfig{
			ID:     poolID,
			Size:   pool["size"].(string),
			Count:  pool["node_count"].(int),
			Labels: labels,
			Taints: taints,
		}

		if pool["public_ip_node_pool"].(bool) {
			cr.PublicIPNodePool = pool["public_ip_node_pool"].(bool)
		}

		expandedNodePools = append(expandedNodePools, cr)
	}

	return expandedNodePools
}
