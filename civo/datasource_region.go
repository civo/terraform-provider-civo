package civo

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source to get and filter all regions
// use to define the region in the rest of the resource or datasource
func dataSourceRegion() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		Description: "Retrieves information about the region that Civo supports, with the ability to filter the results.",
		RecordSchema: map[string]*schema.Schema{
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The code of the region",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human name of the region",
			},
			"country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The country of the region",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If the region is the default region, this will return `true`",
			},
		},
		ResultAttributeName: "regions",
		FlattenRecord:       flattenRegions,
		GetRecords:          getRegios,
	}

	return datalist.NewResource(dataListConfig)

}

func getRegios(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	apiClient := m.(*civogo.Client)

	regions := []interface{}{}
	partialRegions, err := apiClient.ListRegions()
	if err != nil {
		return nil, fmt.Errorf("[ERR] error retrieving regions: %s", err)
	}

	for _, partialRegion := range partialRegions {
		regions = append(regions, partialRegion)
	}

	return regions, nil
}

func flattenRegions(region, m interface{}, extra map[string]interface{}) (map[string]interface{}, error) {

	s := region.(civogo.Region)

	flattenedRegion := map[string]interface{}{}
	flattenedRegion["code"] = s.Code
	flattenedRegion["name"] = s.Name
	flattenedRegion["country"] = s.Country
	flattenedRegion["default"] = s.Default

	return flattenedRegion, nil
}
