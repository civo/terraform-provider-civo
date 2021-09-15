package civo

import (
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source to get from the api a specific template
// using the code of the image
func dataSourceDiskImage() *schema.Resource {

	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: templateSchema(),
		ExtraQuerySchema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResultAttributeName: "diskimages",

		// use functions from "datasource_template.go" file
		FlattenRecord: flattenTemplate,
		GetRecords:    getTemplates,
	}

	return datalist.NewResource(dataListConfig)
}
