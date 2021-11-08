package civo

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TemplateDisk is a temporal struct to get all template in one place
type TemplateDisk struct {
	ID      string
	Name    string
	Version string
	Label   string
}

// Data source to get from the api a specific template
// using the code of the image
func dataSourceDiskImage() *schema.Resource {

	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: diskimageSchema(),
		Description:  "Get information on an disk image for use in other resources (e.g. creating a instance) with the ability to filter the results.",
		ExtraQuerySchema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If is used, all disk image will be from this region. Required if no region is set in provider.",
			},
		},
		ResultAttributeName: "diskimages",
		FlattenRecord:       flattenDiskimage,
		GetRecords:          getDiskimages,
	}

	return datalist.NewResource(dataListConfig)
}

func getDiskimages(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	region, ok := extra["region"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to find `region` key from query data")
	}

	if region != "" {
		apiClient.Region = region
	}

	templateDiskList := []TemplateDisk{}

	diskImage, err := apiClient.ListDiskImages()
	if err != nil {
		return nil, fmt.Errorf("[ERR] error retrieving all Disk Images: %s", err)
	}

	for _, v := range diskImage {
		templateDiskList = append(templateDiskList, TemplateDisk{ID: v.ID, Name: v.Name, Version: v.Version, Label: v.Label})
	}

	templates := []interface{}{}
	for _, partialSize := range templateDiskList {
		templates = append(templates, partialSize)
	}

	return templates, nil
}

func flattenDiskimage(template, m interface{}, extra map[string]interface{}) (map[string]interface{}, error) {

	s := template.(TemplateDisk)

	flattenedTemplate := map[string]interface{}{}
	flattenedTemplate["id"] = s.ID
	flattenedTemplate["name"] = s.Name
	flattenedTemplate["version"] = s.Version
	flattenedTemplate["label"] = s.Label

	return flattenedTemplate, nil
}

func diskimageSchema() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of disk image",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of disk image",
		},
		"version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Version of disk image",
		},
		"label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Label of disk image",
		},
	}
}
