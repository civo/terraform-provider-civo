package civo

import (
	"fmt"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Data source to get from the api a specific domain
// using the id or the name of the domain
func dataSourceVolume() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Get information on a volume for use in other resources. This data source provides all of the volumes properties as configured on your Civo account.",
			"An error will be raised if the provided volume name does not exist in your Civo account.",
		}, "\n\n"),
		Read: dataSourceVolumeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The name of the volume",
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The region where volume is running",
			},
			// Computed resource
			"size_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the volume (in GB)",
			},
			"mount_point": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The mount point of the volume",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date of the creation of the volume",
			},
		},
	}
}

func dataSourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	var foundVolume *civogo.Volume

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the volume by id")
		volume, err := apiClient.FindVolume(id.(string))
		if err != nil {
			return fmt.Errorf("[ERR] failed to retrive volume: %s", err)
		}

		foundVolume = volume
	} else if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the volume by name")
		volume, err := apiClient.FindVolume(name.(string))
		if err != nil {
			return fmt.Errorf("[ERR] failed to retrive volume: %s", err)
		}

		foundVolume = volume
	}

	d.SetId(foundVolume.ID)
	d.Set("name", foundVolume.Name)
	d.Set("size_gb", foundVolume.SizeGigabytes)
	d.Set("mount_point", foundVolume.MountPoint)
	d.Set("created_at", foundVolume.CreatedAt.UTC().String())

	return nil
}
