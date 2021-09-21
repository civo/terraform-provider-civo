/*
Copyright (c) 2020, DigitalOcean
License: MPL (see https://www.mozilla.org/en-US/MPL/)
Source: https://github.com/terraform-providers/terraform-provider-digitalocean
*/

package datalist

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type commonFilter struct {
	key     string
	values  []interface{}
	all     bool
	matchBy string
}

func filterSchema(resultAttributeName string, allowedKeys []string) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(allowedKeys, false),
					Description:  fmt.Sprintf("Filter %s by this key. This may be one of %s.", resultAttributeName, utils.GetCommaSeparatedAllowedKeys(allowedKeys)),
				},
				"values": {
					Type:        schema.TypeList,
					Required:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: fmt.Sprintf("Only retrieves `%s` which keys has value that matches one of the values provided here", resultAttributeName),
				},
				"all": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Set to `true` to require that a field match all of the `values` instead of just one or more of them. This is useful when matching against multi-valued fields such as lists or sets where you want to ensure that all of the `values` are present in the list or set.",
				},
				"match_by": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "exact",
					ValidateFunc: validation.StringInSlice([]string{"exact", "re", "substring"}, false),
					Description:  "One of `exact` (default), `re`, or `substring`. For string-typed fields, specify `re` to match by using the `values` as regular expressions, or specify `substring` to match by treating the `values` as substrings to find within the string field.",
				},
			},
		},
		Optional:    true,
		Description: "One or more key/value pairs on which to filter results",
	}
}

func expandFilters(recordSchema map[string]*schema.Schema, rawFilters []interface{}) ([]commonFilter, error) {
	expandedFilters := make([]commonFilter, len(rawFilters))

	for i, rawFilter := range rawFilters {
		f := rawFilter.(map[string]interface{})

		key := f["key"].(string)
		s, ok := recordSchema[key]
		if !ok {
			return nil, fmt.Errorf("field '%s' does not exist in record schema", key)
		}

		matchBy := "exact"
		if v, ok := f["match_by"].(string); ok {
			matchBy = v
		}

		expandedFilterValues, err := expandFilterValues(f["values"].([]interface{}), s, matchBy)
		if err != nil {
			return nil, err
		}

		all := false
		if v, ok := f["all"]; ok {
			all = v.(bool)
		}

		expandedFilter := commonFilter{
			key:     key,
			values:  expandedFilterValues,
			all:     all,
			matchBy: matchBy,
		}

		expandedFilters[i] = expandedFilter
	}

	return expandedFilters, nil
}

func isPrimitiveType(fieldType schema.ValueType) bool {
	switch fieldType {
	case schema.TypeString,
		schema.TypeBool,
		schema.TypeInt,
		schema.TypeFloat:
		return true
	}

	return false
}

// Expands a single filter value (which is a string) into the Go type that can actually be
// used for comparisons when filtering. This should not be called with container or
// composite types.
func expandPrimitiveFilterValue(
	filterValue string,
	fieldType schema.ValueType,
	matchBy string,
) (interface{}, error) {
	var expandedValue interface{}

	switch fieldType {
	case schema.TypeString:
		switch matchBy {
		case "exact", "substring":
			expandedValue = filterValue
		case "re":
			re, err := regexp.Compile(filterValue)
			if err != nil {
				return nil, fmt.Errorf("unable to parse value as regular expression: %s: %s", filterValue, err)
			}
			expandedValue = re
		default:
			panic("unreachable")
		}

	case schema.TypeBool:
		boolValue, err := strconv.ParseBool(filterValue)
		if err != nil {
			return nil, fmt.Errorf("unable to parse value as bool: %s: %s", filterValue, err)
		}
		expandedValue = boolValue

	case schema.TypeInt:
		intValue, err := strconv.Atoi(filterValue)
		if err != nil {
			return nil, fmt.Errorf("unable to parse value as integer: %s: %s", filterValue, err)
		}
		expandedValue = intValue

	case schema.TypeFloat:
		floatValue, err := strconv.ParseFloat(filterValue, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse value as floating point: %s: %s", filterValue, err)
		}
		expandedValue = floatValue

	default:
		panic("unreachable")
	}

	return expandedValue, nil
}

// Takes the "raw" set of strings provided by the Terraform SDK and converts the values
// into the actual Go types used for comparisons when filtering.
func expandFilterValues(
	rawFilterValues []interface{},
	fieldSchema *schema.Schema,
	matchBy string,
) ([]interface{}, error) {
	expandedFilterValues := make([]interface{}, len(rawFilterValues))

	for i, rawFilterValue := range rawFilterValues {
		filterValue := rawFilterValue.(string)
		var expandedValue interface{}

		if isPrimitiveType(fieldSchema.Type) {
			ev, err := expandPrimitiveFilterValue(filterValue, fieldSchema.Type, matchBy)
			if err != nil {
				return nil, err
			}
			expandedValue = ev
		} else {
			if s, ok := fieldSchema.Elem.(*schema.Schema); ok {
				if isPrimitiveType(s.Type) {
					ev, err := expandPrimitiveFilterValue(filterValue, s.Type, matchBy)
					if err != nil {
						return nil, err
					}
					expandedValue = ev
				} else {
					return nil, fmt.Errorf("cannot filter on a non-primitive type")
				}
			} else {
				return nil, fmt.Errorf("cannot filter on aggregate type with non-Schema element type")
			}
		}

		expandedFilterValues[i] = expandedValue
	}

	return expandedFilterValues, nil
}

func applyFilters(recordSchema map[string]*schema.Schema, records []map[string]interface{}, filters []commonFilter) []map[string]interface{} {
	for _, f := range filters {
		// Handle multiple filters by applying them in order
		var filteredRecords []map[string]interface{}

		filterFunc := func(record map[string]interface{}) bool {
			result := f.all

			for _, filterValue := range f.values {
				thisValueMatches := valueMatches(recordSchema[f.key], record[f.key], filterValue, f.matchBy)
				if f.all {
					result = result && thisValueMatches
				} else {
					result = result || thisValueMatches
				}
			}

			return result
		}

		for _, record := range records {
			if filterFunc(record) {
				filteredRecords = append(filteredRecords, record)
			}
		}

		records = filteredRecords
	}

	return records
}
