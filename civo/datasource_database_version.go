package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DatabaseVersion is a temporal struct to save all versions
type DatabaseVersion struct {
	Engine  string
	Version string
	Default bool
}

// Data source to get and filter all database version
// use to define the engine and version in resourceDatabase
func dataDatabaseVersion() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		Description:         "Retrieves information about the database versions that Civo supports, with the ability to filter the results.",
		RecordSchema:        versionSchema(),
		ResultAttributeName: "versions",
		FlattenRecord:       flattenVersion,
		GetRecords:          getVersion,
	}

	return datalist.NewResource(dataListConfig)
}

func getVersion(m interface{}, _ map[string]interface{}) ([]interface{}, error) {
	apiClient := m.(*civogo.Client)

	versions := []interface{}{}
	partialVersions, err := apiClient.ListDBVersions()
	if err != nil {
		return nil, fmt.Errorf("[ERR] error retrieving version: %s", err)
	}

	versionList := []DatabaseVersion{}
	for k, v := range partialVersions {
		for _, version := range v {
			versionList = append(versionList, DatabaseVersion{
				Engine:  k,
				Version: version.SoftwareVersion,
				Default: version.Default,
			})
		}
	}

	for _, version := range versionList {
		versions = append(versions, version)
	}

	return versions, nil
}

func flattenVersion(versions, _ interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
	s := versions.(DatabaseVersion)

	flattenedSize := map[string]interface{}{}
	flattenedSize["engine"] = s.Engine
	flattenedSize["version"] = s.Version
	flattenedSize["default"] = s.Default

	return flattenedSize, nil
}

func versionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"engine": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The engine of the database",
		},
		"version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The version of the database",
		},
		"default": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "If the version is the default",
		},
	}
}
