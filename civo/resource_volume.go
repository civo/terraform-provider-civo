package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

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
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
	}
}

func resourceVolumeCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	config := &civogo.VolumeConfig{Name: d.Get("name").(string), SizeGigabytes: d.Get("size_gb").(int), Bootable: d.Get("bootable").(bool)}

	volume, err := apiClient.NewVolume(config)
	if err != nil {
		return fmt.Errorf("failed to create a new config: %s", err)
	}

	d.SetId(volume.ID)

	if d.Get("instance_id").(string) != "" {
		_, err := apiClient.AttachVolume(d.Id(), d.Get("instance_id").(string))
		if err != nil {
			log.Printf("[WARN] An error occurred while trying to attach the volume (%s)", d.Id())
		}
	}

	return resourceNetworkRead(d, m)
}

func resourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	CurrentVolume := civogo.Volume{}

	resp, err := apiClient.ListVolumes()
	if err != nil {
		return fmt.Errorf("failed to create a new config: %s", err)
	}

	for _, volume := range resp {
		if volume.ID == d.Id() {
			CurrentVolume = volume
		}
	}

	d.Set("name", CurrentVolume.Name)
	d.Set("size_gb", CurrentVolume.SizeGigabytes)
	d.Set("bootable", CurrentVolume.Bootable)
	d.Set("instance_id", CurrentVolume.InstanceID)
	d.Set("mount_point", CurrentVolume.MountPoint)
	d.Set("created_at", CurrentVolume.CreatedAt)

	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	if d.HasChange("instance_id") {
		if d.Get("instance_id").(string) != "" {
			_, err := apiClient.AttachVolume(d.Id(), d.Get("instance_id").(string))
			if err != nil {
				log.Printf("[WARN] An error occurred while trying to attach the volume (%s)", d.Id())
			}

			return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				CurrentVolume := civogo.Volume{}

				resp, err := apiClient.ListVolumes()
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("failed to get all volume: %s", err))
				}

				for _, volume := range resp {
					if volume.ID == d.Id() {
						CurrentVolume = volume
					}
				}

				if CurrentVolume.MountPoint != "" {
					return resource.RetryableError(fmt.Errorf("expected volume to be mount but was in state %s", CurrentVolume.MountPoint))
				}

				return resource.NonRetryableError(resourceVolumeRead(d, m))
			})
		}

		if d.Get("instance_id").(string) == "" {
			_, err := apiClient.DetachVolume(d.Id())
			if err != nil {
				log.Printf("[WARN] An error occurred while trying to detach volume (%s)", d.Id())
			}

			return resourceVolumeRead(d, m)
		}

	}

	if d.HasChange("size_gb") {
		if d.Get("instance_id").(string) != "" {
			_, err := apiClient.DetachVolume(d.Id())
			if err != nil {
				log.Printf("[WARN] An error occurred while trying to detach volume (%s)", d.Id())
			}

			time.Sleep(10 * time.Second)

			newSize := d.Get("size_gb").(int)
			_, err = apiClient.ResizeVolume(d.Id(), newSize)
			if err != nil {
				log.Printf("[INFO] Civo volume (%s) size not change", d.Id())
			}

			time.Sleep(2 * time.Second)

			_, err = apiClient.AttachVolume(d.Id(), d.Get("instance_id").(string))
			if err != nil {
				log.Printf("[WARN] An error occurred while trying to attach the volume (%s)", d.Id())
			}

		}

	}

	return resourceVolumeRead(d, m)
}

func resourceVolumeDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	_, err := apiClient.DeleteVolume(d.Id())
	if err != nil {
		log.Printf("[INFO] Civo volume (%s) was delete", d.Id())
	}
	return nil
}
