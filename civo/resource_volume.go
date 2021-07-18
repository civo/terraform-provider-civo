package civo

import (
	"fmt"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Volume resource, with this we can create and manage all volume
func resourceVolume() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A name that you wish to use to refer to this volume",
				ValidateFunc: utils.ValidateName,
			},
			"size_gb": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "A minimum of 1 and a maximum of your available disk space from your quota specifies the size of the volume in gigabytes ",
			},
			"bootable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Mark the volume as bootable",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region for the volume",
			},
			// Computed resource
			"mount_point": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourceVolumeCreate,
		Read:   resourceVolumeRead,
		Update: resourceVolumeUpdate,
		Delete: resourceVolumeDelete,
		//Exists: resourceExistsItem,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// function to create the new volume
func resourceVolumeCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] configuring the volume %s", d.Get("name").(string))
	config := &civogo.VolumeConfig{Name: d.Get("name").(string), SizeGigabytes: d.Get("size_gb").(int), Bootable: d.Get("bootable").(bool)}

	volume, err := apiClient.NewVolume(config)
	if err != nil {
		return fmt.Errorf("[ERR] failed to create a new config: %s", err)
	}

	d.SetId(volume.ID)

	return resourceVolumeRead(d, m)
}

// function to read the volume
func resourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] retrieving the volume %s", d.Id())
	resp, err := apiClient.FindVolume(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERR] failed retrieving the volume: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("size_gb", resp.SizeGigabytes)
	d.Set("bootable", resp.Bootable)
	d.Set("mount_point", resp.MountPoint)
	d.Set("created_at", resp.CreatedAt.UTC().String())

	return nil
}

// function to update the volume
func resourceVolumeUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] retrieving the volume %s", d.Id())
	resp, err := apiClient.FindVolume(d.Id())
	if err != nil {
		return fmt.Errorf("[ERR] failed retrieving the volume: %s", err)
	}

	if d.HasChange("size_gb") {
		if resp.InstanceID != "" {
			_, err := apiClient.DetachVolume(d.Id())
			if err != nil {
				return fmt.Errorf("[WARN] an error occurred while tring to detach volume %s, %s", d.Id(), err)
			}

			time.Sleep(10 * time.Second)

			newSize := d.Get("size_gb").(int)
			_, err = apiClient.ResizeVolume(d.Id(), newSize)
			if err != nil {
				return fmt.Errorf("[ERR] the volume (%s) size not change %s", d.Id(), err)
			}

			time.Sleep(2 * time.Second)

			_, err = apiClient.AttachVolume(d.Id(), resp.InstanceID)
			if err != nil {
				return fmt.Errorf("[ERR] an error occurred while tring to attach the volume %s", d.Id())
			}

		} else {
			newSize := d.Get("size_gb").(int)
			_, err = apiClient.ResizeVolume(d.Id(), newSize)
			if err != nil {
				return fmt.Errorf("[ERR] the volume (%s) size not change %s", d.Id(), err)
			}
		}

	}

	return resourceVolumeRead(d, m)
}

// function to delete the volume
func resourceVolumeDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] deleting the volume %s", d.Id())
	_, err := apiClient.DeleteVolume(d.Id())
	if err != nil {
		return fmt.Errorf("[ERR] an error occurred while tring to delete the volume %s", err)
	}
	return nil
}
