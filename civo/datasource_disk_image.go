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
		Description:  "Get information on an disk image for use in other resources (e.g. creating a instance) with the ability to filter the results.",
		ExtraQuerySchema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If is used, all disk image will be from this region. Required if no region is set in provider.",
			},
		},
		ResultAttributeName: "diskimages",

		// use functions from "datasource_template.go" file
		FlattenRecord: flattenTemplate,
		GetRecords:    getTemplates,
	}

	return datalist.NewResource(dataListConfig)
}
