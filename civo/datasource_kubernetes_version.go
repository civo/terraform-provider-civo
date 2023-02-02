package civo

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source to get and filter all kubernetes version
// available in the server, use to define the version at the
// moment of the cluster creation in resourceKubernetesCluster
func dataSourceKubernetesVersion() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		Description:         "Provides access to the available Civo Kubernetes versions, with the ability to filter the results.",
		RecordSchema:        kubernetesVersionSchema(),
		ResultAttributeName: "versions",
		FlattenRecord:       flattenKubernetesVersion,
		GetRecords:          getKubernetesVersions,
	}

	return datalist.NewResource(dataListConfig)
}

func getKubernetesVersions(m interface{}, _ map[string]interface{}) ([]interface{}, error) {
	apiClient := m.(*civogo.Client)

	versions := []interface{}{}
	partialVersions, err := apiClient.ListAvailableKubernetesVersions()
	if err != nil {
		return nil, fmt.Errorf("[ERR] error retrieving all versions: %s", err)
	}

	for _, partialSize := range partialVersions {
		versions = append(versions, partialSize)
	}

	return versions, nil
}

func flattenKubernetesVersion(version, _ interface{}, _ map[string]interface{}) (map[string]interface{}, error) {

	s := version.(civogo.KubernetesVersion)

	flattenedVersion := map[string]interface{}{}
	flattenedVersion["version"] = s.Version
	flattenedVersion["label"] = fmt.Sprintf("v%s", s.Version)
	flattenedVersion["type"] = s.ClusterType
	flattenedVersion["default"] = s.Default
	return flattenedVersion, nil
}

func kubernetesVersionSchema() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "A version of the Kubernetes",
		},
		"label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The label of this version",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The type of the cluster, can be `talos` or `k3s`",
		},
		"default": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "If is the default version used in all cluster, this will return `true`",
		},
	}
}
