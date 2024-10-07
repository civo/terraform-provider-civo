package volume

import (
	"context"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceVolumeType function returns a schema.Resource that represents a Volume.
// This can be used to query and retrieve details about a specific Volume in the infrastructure using its id or name.
func DataSourceVolumeType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCivoVolumeTypeRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"labels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceCivoVolumeTypeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*civogo.Client)

	// overwrite the region if is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		client.Region = region.(string)
	}

	// Get the name of the volume type from the Terraform configuration
	volumeTypeName := d.Get("name").(string)

	// Fetch all volume types
	volumeTypes, err := client.ListVolumeTypes()
	if err != nil {
		return diag.Errorf("[ERR] failed to retrieve volume type: %s", err)
	}

	// Find the specific volume type by name
	var foundVolumeType *civogo.VolumeType
	for _, vt := range volumeTypes {
		if vt.Name == volumeTypeName {
			foundVolumeType = &vt
			break
		}
	}

	// If the volume type is not found, return an error
	if foundVolumeType == nil {
		return diag.Errorf("volume type with name '%s' not found", volumeTypeName)
	}

	// Set the data source's ID and output fields
	d.SetId(foundVolumeType.Name)
	if err := d.Set("labels", foundVolumeType.Labels); err != nil {
		return diag.Errorf("[ERR] failed to set volume labels: %s", err)
	}

	return nil
}
