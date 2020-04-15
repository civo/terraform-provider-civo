package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

// Volume resource, with this we can create and manage all volume
func resourceVolume() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name for the volume",
				ValidateFunc: validateName,
			},
			"size_gb": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Size for the volume in GB",
			},
			"bootable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Mark like bootable this volume",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Attach to a instance",
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
			State: schema.ImportStatePassthrough,
		},
	}
}

// function to create the new volume
func resourceVolumeCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] configuring the volume %s", d.Get("name").(string))
	config := &civogo.VolumeConfig{Name: d.Get("name").(string), SizeGigabytes: d.Get("size_gb").(int), Bootable: d.Get("bootable").(bool)}

	volume, err := apiClient.NewVolume(config)
	if err != nil {
		return fmt.Errorf("[ERR] failed to create a new config: %s", err)
	}

	d.SetId(volume.ID)

	if d.Get("instance_id").(string) != "" {
		_, err := apiClient.AttachVolume(d.Id(), d.Get("instance_id").(string))
		if err != nil {
			return fmt.Errorf("[ERR] An error occurred while tring to attach the volume (%s)", d.Id())
		}
	}

	return resourceNetworkRead(d, m)
}

// function to read the volume
func resourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] retrieving the template %s", d.Id())
	resp, err := apiClient.FindVolume(d.Id())
	if err != nil {
		return fmt.Errorf("[ERR] failed retrieving the volume: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("size_gb", resp.SizeGigabytes)
	d.Set("bootable", resp.Bootable)
	d.Set("instance_id", resp.InstanceID)
	d.Set("mount_point", resp.MountPoint)
	d.Set("created_at", resp.CreatedAt.UTC().String())

	return nil
}

// function to update the volume
func resourceVolumeUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	if d.HasChange("instance_id") {
		if d.Get("instance_id").(string) != "" {
			_, err := apiClient.AttachVolume(d.Id(), d.Get("instance_id").(string))
			if err != nil {
				return fmt.Errorf("[WARN] An error occurred while tring to attach the volume %s, %s", d.Id(), err)
			}

			return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				CurrentVolume := civogo.Volume{}

				resp, err := apiClient.ListVolumes()
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("[ERR] failed to get all volume: %s", err))
				}

				for _, volume := range resp {
					if volume.ID == d.Id() {
						CurrentVolume = volume
					}
				}

				if CurrentVolume.MountPoint != "" {
					return resource.RetryableError(fmt.Errorf("expected volume to be mount but the mount point is %s", CurrentVolume.MountPoint))
				}

				return resource.NonRetryableError(resourceVolumeRead(d, m))
			})
		}

		if d.Get("instance_id").(string) == "" {
			_, err := apiClient.DetachVolume(d.Id())
			if err != nil {
				return fmt.Errorf("[ERR] An error occurred while tring to detach volume (%s)", d.Id())
			}

			return resourceVolumeRead(d, m)
		}

	}

	if d.HasChange("size_gb") {
		if d.Get("instance_id").(string) != "" {
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

			_, err = apiClient.AttachVolume(d.Id(), d.Get("instance_id").(string))
			if err != nil {
				return fmt.Errorf("[ERR] an error occurred while tring to attach the volume %s", d.Id())
			}

		}

	}

	return resourceVolumeRead(d, m)
}

// function to delete the volume
func resourceVolumeDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] deleting the volume %s", d.Id())
	_, err := apiClient.DeleteVolume(d.Id())
	if err != nil {
		return fmt.Errorf("[ERR] an error occurred while tring to delete the volume %s", err)
	}
	return nil
}
