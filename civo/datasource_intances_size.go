package civo

import (
	"fmt"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SizeList is a temporal struct to save all size
type SizeList struct {
	Name        string
	Description string
	Type        string
	CPU         int
	RAM         int
	DisK        int
	Selectable  bool
}

// Data source to get and filter all instances size
// use to define the size in resourceInstance
func dataSourceInstancesSize() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		Description:         "Retrieves information about the instance sizes that Civo supports, with the ability to filter the results.",
		RecordSchema:        instancesSizeSchema(),
		ResultAttributeName: "sizes",
		FlattenRecord:       flattenInstancesSize,
		GetRecords:          getInstancesSizes,
	}

	instanceSizeResource := datalist.NewResource(dataListConfig)
	instanceSizeResource.DeprecationMessage = "Use the civo_size datasource instead"

	return instanceSizeResource

}

func getInstancesSizes(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	apiClient := m.(*civogo.Client)

	sizes := []interface{}{}
	partialSizes, err := apiClient.ListInstanceSizes()
	if err != nil {
		return nil, fmt.Errorf("[ERR] error retrieving sizes: %s", err)
	}

	sizeList := []SizeList{}

	for _, v := range partialSizes {
		if !v.Selectable {
			continue
		}

		typeName := ""

		switch {
		case strings.Contains(v.Name, "db"):
			typeName = "database"
		case strings.Contains(v.Name, "kube") || strings.Contains(v.Name, "k3s"):
			typeName = "kubernetes"
		default:
			typeName = "instance"
		}

		sizeList = append(sizeList, SizeList{
			Name:        v.Name,
			Description: v.Description,
			Type:        typeName,
			CPU:         v.CPUCores,
			RAM:         v.RAMMegabytes,
			DisK:        v.DiskGigabytes,
			Selectable:  v.Selectable,
		})
	}

	for _, partialSize := range sizeList {
		sizes = append(sizes, partialSize)
	}

	return sizes, nil
}

func flattenInstancesSize(size, m interface{}, extra map[string]interface{}) (map[string]interface{}, error) {

	s := size.(SizeList)

	flattenedSize := map[string]interface{}{}
	flattenedSize["name"] = s.Name
	flattenedSize["type"] = s.Type
	flattenedSize["cpu"] = s.CPU
	flattenedSize["ram"] = s.RAM
	flattenedSize["disk"] = s.DisK
	flattenedSize["description"] = s.Description
	flattenedSize["selectable"] = s.Selectable

	return flattenedSize, nil
}

func instancesSizeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The name of the instance size",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "A human name of the instance size",
		},
		"cpu": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total of CPU in the instance",
		},
		"ram": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total of RAM of the instance",
		},
		"disk": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The instance size of SSD",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "A description of the instance size",
		},
		"selectable": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "If can use the instance size",
		},
	}
}
