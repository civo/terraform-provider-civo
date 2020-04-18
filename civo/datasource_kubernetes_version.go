package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

// Data source to get and filter all kubernetes version
// available in the server, use to define the version at the
// moment of the cluster creation in resourceKubernetesCluster
func dataSourceKubernetesVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKubernetesVersionRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			// computed attributes
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
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKubernetesVersionRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("one of filters must be assigned")
	}

	if filtersOk {
		log.Printf("[INFO] Getting all versions of kubernetes")
		resp, err := apiClient.ListAvailableKubernetesVersions()
		if err != nil {
			return fmt.Errorf("no version was found in the server")
		}

		log.Printf("[INFO] Finding the version of kubernetes")
		version, err := findKubernetesVersionByFilter(resp, filters.(*schema.Set))
		if err != nil {
			return fmt.Errorf("no version was found in the server, %s", err)
		}

		d.SetId(version.Version)
		d.Set("version", version.Version)
		d.Set("label", fmt.Sprintf("v%s", version.Version))
		d.Set("type", version.Type)
		d.Set("default", version.Default)
	}

	return nil
}

func findKubernetesVersionByFilter(version []civogo.KubernetesVersion, set *schema.Set) (*civogo.KubernetesVersion, error) {
	results := make([]civogo.KubernetesVersion, 0)

	var filters []Filter

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}
		filters = append(filters, Filter{Name: m["name"].(string), Values: filterValues, Regex: m["regex"].(bool)})
	}

	for _, valueFilters := range filters {
		for _, valueVersion := range version {

			// Filter for version
			if valueFilters.Name == "version" {
				if valueVersion.Version == valueFilters.Values[0] {
					results = append(results, valueVersion)
				}
			}

			// Filter for type
			if valueFilters.Name == "type" {
				if valueVersion.Type == valueFilters.Values[0] {
					results = append(results, valueVersion)
				}
			}
		}
	}

	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no kubernetes version found for your search")
	}
	return nil, fmt.Errorf("too many kubernetes version found (found %d, expected 1)", len(results))
}
