package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source to get from the api a specific template
// using the code of the image
func dataSourceTemplate() *schema.Resource {

	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"short_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		FilterKeys: []string{
			"code",
			"name",
		},
		SortKeys: []string{
			"code",
			"name",
		},
		ResultAttributeName: "templates",
		FlattenRecord:       flattenTemplate,
		GetRecords:          getTemplates,
	}

	return datalist.NewResource(dataListConfig)

}

func getTemplates(m interface{}) ([]interface{}, error) {
	apiClient := m.(*civogo.Client)

	templates := []interface{}{}
	partialTemplates, err := apiClient.ListTemplates()
	if err != nil {
		return nil, fmt.Errorf("[ERR] error retrieving all templates: %s", err)
	}

	for _, partialSize := range partialTemplates {
		templates = append(templates, partialSize)
	}

	return templates, nil
}

func flattenTemplate(template, m interface{}) (map[string]interface{}, error) {

	s := template.(civogo.Template)

	flattenedTemplate := map[string]interface{}{}
	flattenedTemplate["id"] = s.ID
	flattenedTemplate["code"] = s.Code
	flattenedTemplate["name"] = s.Name
	flattenedTemplate["volume_id"] = s.VolumeID
	flattenedTemplate["image_id"] = s.ImageID
	flattenedTemplate["short_description"] = s.ShortDescription
	flattenedTemplate["description"] = s.Description
	flattenedTemplate["default_username"] = s.DefaultUsername
	flattenedTemplate["cloud_config"] = s.CloudConfig

	return flattenedTemplate, nil
}
