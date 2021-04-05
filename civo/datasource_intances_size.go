package civo

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source to get and filter all instances size
// use to define the size in resourceInstance
func dataSourceInstancesSize() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        instancesSizeSchema(),
		ResultAttributeName: "sizes",
		FlattenRecord:       flattenInstancesSize,
		GetRecords:          getInstancesSizes,
	}

	return datalist.NewResource(dataListConfig)

}

func getInstancesSizes(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	apiClient := m.(*civogo.Client)

	sizes := []interface{}{}
	partialSizes, err := apiClient.ListInstanceSizes()
	if err != nil {
		return nil, fmt.Errorf("[ERR] error retrieving sizes: %s", err)
	}

	for _, partialSize := range partialSizes {
		sizes = append(sizes, partialSize)
	}

	return sizes, nil
}

func flattenInstancesSize(size, m interface{}, extra map[string]interface{}) (map[string]interface{}, error) {

	s := size.(civogo.InstanceSize)

	flattenedSize := map[string]interface{}{}
	flattenedSize["name"] = s.Name
	flattenedSize["nice_name"] = s.NiceName
	flattenedSize["cpu_cores"] = s.CPUCores
	flattenedSize["ram_mb"] = s.RAMMegabytes
	flattenedSize["disk_gb"] = s.DiskGigabytes
	flattenedSize["description"] = s.Description
	flattenedSize["selectable"] = s.Selectable

	return flattenedSize, nil
}

func instancesSizeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"nice_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cpu_cores": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ram_mb": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"disk_gb": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"selectable": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}
