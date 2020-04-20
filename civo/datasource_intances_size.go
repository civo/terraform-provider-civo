package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"regexp"
	"strconv"
)

// Data source to get and filter all instances size
// use to define the size in resourceInstance
func dataSourceInstancesSize() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInstancesSizeRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			// computed attributes
			"sizes": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "a list of backend instances, each containing an instance_id, protocol (http or https) and port",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
		},
	}
}

func dataSourceInstancesSizeRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return fmt.Errorf("one of filters must be assigned")
	}

	if filtersOk {
		log.Printf("[INFO] Getting the instances size")
		resp, err := apiClient.ListInstanceSizes()
		if err != nil {
			return fmt.Errorf("no instances size was found in the server")
		}

		log.Printf("[INFO] Finding the size of the instances")
		size, err := findInstancesSizeByFilter(resp, filters.(*schema.Set))
		if err != nil {
			return fmt.Errorf("no instances size was found in the server, %s", err)
		}

		d.SetId(resource.UniqueId())

		if err := d.Set("sizes", size); err != nil {
			return fmt.Errorf("unable to set `%s` attribute: %s", "sizes", err)
		}
	}

	return nil
}

func findInstancesSizeByFilter(sizes []civogo.InstanceSize, set *schema.Set) ([]civogo.InstanceSize, error) {
	results := make([]civogo.InstanceSize, 0)

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
		for _, valueSize := range sizes {

			// filter for the name
			if valueFilters.Name == "name" {
				if valueFilters.Regex {
					r, _ := regexp.Compile(valueFilters.Values[0])
					if r.MatchString(valueSize.Name) {
						results = append(results, valueSize)
					}
				} else {
					if valueSize.Name == valueFilters.Values[0] {
						results = append(results, valueSize)
					}
				}
			}

			// filter for the CPU
			if valueFilters.Name == "cpu" {
				if valueFilters.Regex {
					r, _ := regexp.Compile(valueFilters.Values[0])
					if r.MatchString(strconv.Itoa(valueSize.CPUCores)) {
						results = append(results, valueSize)
					}
				} else {
					if strconv.Itoa(valueSize.CPUCores) == valueFilters.Values[0] {
						results = append(results, valueSize)
					}
				}
			}

			// filter for the RAM
			if valueFilters.Name == "ram" {
				if valueFilters.Regex {
					r, _ := regexp.Compile(valueFilters.Values[0])
					if r.MatchString(strconv.Itoa(valueSize.RAMMegabytes)) {
						results = append(results, valueSize)
					}
				} else {
					if strconv.Itoa(valueSize.RAMMegabytes) == valueFilters.Values[0] {
						results = append(results, valueSize)
					}
				}
			}

			// filter for the Disk
			if valueFilters.Name == "disk" {
				if valueFilters.Regex {
					r, _ := regexp.Compile(valueFilters.Values[0])
					if r.MatchString(strconv.Itoa(valueSize.DiskGigabytes)) {
						results = append(results, valueSize)
					}
				} else {
					if strconv.Itoa(valueSize.DiskGigabytes) == valueFilters.Values[0] {
						results = append(results, valueSize)
					}
				}
			}

		}
	}

	if len(results) > 1 {
		return results, nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no instances sizes found for your search")
	}
	return nil, fmt.Errorf("too many instances sizes found (found %d, expected 1)", len(results))
}
