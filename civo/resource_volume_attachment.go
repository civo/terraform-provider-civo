package civo

import (
	"fmt"
	"log"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Volume resource, with this we can create and manage all volume
func resourceVolumeAttachment() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"volume_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
		Create: resourceVolumeAttachmentCreate,
		Read:   resourceVolumeAttachmentRead,
		Delete: resourceVolumeAttachmentDelete,
	}
}

// function to create the new volume
func resourceVolumeAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	instanceID := d.Get("instance_id").(string)
	volumeID := d.Get("volume_id").(string)

	log.Printf("[INFO] retrieving the volume %s", volumeID)
	volume, err := apiClient.FindVolume(volumeID)
	if err != nil {
		return fmt.Errorf("[ERR] Error retrieving volume: %s", err)
	}

	if volume.InstanceID == "" || volume.InstanceID != instanceID {
		// Only one volume can be attached at one time to a single droplet.
		log.Printf("[INFO] attaching the volume %s to instance %s", volumeID, instanceID)
		_, err := apiClient.AttachVolume(volumeID, instanceID)
		if err != nil {
			return fmt.Errorf("[ERR] error attaching volume to instance %s", err)
		}
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-%s-", instanceID, volumeID)))

	return resourceVolumeAttachmentRead(d, m)
}

// function to read the volume
func resourceVolumeAttachmentRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	instanceID := d.Get("instance_id").(string)
	volumeID := d.Get("volume_id").(string)

	log.Printf("[INFO] retrieving the volume %s", volumeID)
	resp, err := apiClient.FindVolume(volumeID)
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERR] failed retrieving the volume: %s", err)
	}

	if resp.InstanceID == "" || resp.InstanceID != instanceID {
		log.Printf("[DEBUG] Volume Attachment (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return nil
}

// function to delete the volume
func resourceVolumeAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	volumeID := d.Get("volume_id").(string)

	log.Printf("[INFO] Detaching the volume %s", d.Id())
	_, err := apiClient.DetachVolume(volumeID)
	if err != nil {
		return fmt.Errorf("[ERR] an error occurred while tring to detach the volume %s", err)
	}
	return nil
}
