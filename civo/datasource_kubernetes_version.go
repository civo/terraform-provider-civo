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
		RecordSchema:        KubernetesVersionSchema(),
		ResultAttributeName: "versions",
		FlattenRecord:       flattenKubernetesVersion,
		GetRecords:          getKubernetesVersions,
	}

	return datalist.NewResource(dataListConfig)
}

func getKubernetesVersions(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
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

func flattenKubernetesVersion(version, m interface{}, extra map[string]interface{}) (map[string]interface{}, error) {

	s := version.(civogo.KubernetesVersion)

	flattenedVersion := map[string]interface{}{}
	flattenedVersion["version"] = s.Version
	flattenedVersion["label"] = fmt.Sprintf("v%s", s.Version)
	flattenedVersion["type"] = s.Type
	flattenedVersion["default"] = s.Default
	return flattenedVersion, nil
}

func KubernetesVersionSchema() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"default": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}
