package size

import (
	"fmt"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Size is a temporal struct to save all size
type Size struct {
	Name        string
	Description string
	Type        string
	CPU         int
	RAM         int
	GPU         int
	GPUType     string
	DisK        int
	Selectable  bool
}

// DataSourceSize function returns a schema.Resource that represents an Instance Size.
// This can be used to query and retrieve details about a specific Instance Size in the infrastructure.
// The retrieved Instance Size can then be used to define the size for other resources or data sources.
func DataSourceSize() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		Description:         "Retrieves information about the sizes that Civo supports, with the ability to filter the results.",
		RecordSchema:        sizeSchema(),
		ResultAttributeName: "sizes",
		FlattenRecord:       flattenSize,
		GetRecords:          getSizes,
	}

	return datalist.NewResource(dataListConfig)

}

func getSizes(m interface{}, _ map[string]interface{}) ([]interface{}, error) {
	apiClient := m.(*civogo.Client)

	sizes := []interface{}{}
	partialSizes, err := apiClient.ListInstanceSizes()
	if err != nil {
		return nil, fmt.Errorf("[ERR] error retrieving sizes: %s", err)
	}

	sizeList := []Size{}

	for _, v := range partialSizes {
		if !v.Selectable {
			continue
		}

		sizeList = append(sizeList, Size{
			Name:        v.Name,
			Description: v.Description,
			Type:        strings.ToLower(v.Type),
			CPU:         v.CPUCores,
			RAM:         v.RAMMegabytes,
			DisK:        v.DiskGigabytes,
			GPU:         v.GPUCount,
			GPUType:     v.GPUType,
			Selectable:  v.Selectable,
		})
	}

	for _, partialSize := range sizeList {
		sizes = append(sizes, partialSize)
	}

	return sizes, nil
}

func flattenSize(size, _ interface{}, _ map[string]interface{}) (map[string]interface{}, error) {

	s := size.(Size)

	flattenedSize := map[string]interface{}{}
	flattenedSize["name"] = s.Name
	flattenedSize["type"] = s.Type
	flattenedSize["cpu"] = s.CPU
	flattenedSize["ram"] = s.RAM
	flattenedSize["disk"] = s.DisK
	flattenedSize["gpu"] = s.GPU
	flattenedSize["gpu_type"] = s.GPUType
	flattenedSize["description"] = s.Description
	flattenedSize["selectable"] = s.Selectable

	return flattenedSize, nil
}

func sizeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The name of the size",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "A human name of the size",
		},
		"cpu": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total of CPU",
		},
		"ram": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total of RAM",
		},
		"disk": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The size of SSD",
		},
		"gpu": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total of GPU",
		},
		"gpu_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "GPU type",
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
