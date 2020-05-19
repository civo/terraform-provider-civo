package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

// Data source to get from the api a specific domain
// using the id or the name of the domain
func dataSourceVolume() *schema.Resource {
	return &schema.Resource{
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
			},
			// Computed resource
			"size_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bootable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mount_point": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	var foundVolume *civogo.Volume

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the volume by id")
		volume, err := apiClient.FindVolume(id.(string))
		if err != nil {
			fmt.Errorf("[ERR] failed to retrive volume: %s", err)
			return err
		}

		foundVolume = volume
	} else if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the volume by name")
		volume, err := apiClient.FindVolume(name.(string))
		if err != nil {
			fmt.Errorf("[ERR] failed to retrive volume: %s", err)
			return err
		}

		foundVolume = volume
	}

	d.SetId(foundVolume.ID)
	d.Set("name", foundVolume.Name)
	d.Set("size_gb", foundVolume.SizeGigabytes)
	d.Set("bootable", foundVolume.Bootable)
	d.Set("mount_point", foundVolume.MountPoint)
	d.Set("created_at", foundVolume.CreatedAt.UTC().String())

	return nil
}
