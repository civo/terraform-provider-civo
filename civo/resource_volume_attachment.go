package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
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

	instanceId := d.Get("instance_id").(string)
	volumeId := d.Get("volume_id").(string)

	log.Printf("[INFO] retrieving the volume %s", volumeId)
	volume, err := apiClient.FindVolume(volumeId)
	if err != nil {
		return fmt.Errorf("[ERR] Error retrieving volume: %s", err)
	}

	if volume.InstanceID == "" || volume.InstanceID != instanceId {
		// Only one volume can be attached at one time to a single droplet.
		log.Printf("[INFO] attaching the volume %s to instance %s", volumeId, instanceId)
		_, err := apiClient.AttachVolume(volumeId, instanceId)
		if err != nil {
			return fmt.Errorf("[ERR] error attaching volume to instance %s", err)
		}
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-%s-", instanceId, volumeId)))

	return resourceNetworkRead(d, m)
}

// function to read the volume
func resourceVolumeAttachmentRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	instanceId := d.Get("instance_id").(string)
	volumeId := d.Get("volume_id").(string)

	log.Printf("[INFO] retrieving the volume %s", volumeId)
	resp, err := apiClient.FindVolume(volumeId)
	if err != nil {
		if resp != nil {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERR] failed retrieving the volume: %s", err)
	}

	if resp.InstanceID == "" || resp.InstanceID != instanceId {
		log.Printf("[DEBUG] Volume Attachment (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return nil
}

// function to delete the volume
func resourceVolumeAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	volumeId := d.Get("volume_id").(string)

	log.Printf("[INFO] Detaching the volume %s", d.Id())
	_, err := apiClient.DetachVolume(volumeId)
	if err != nil {
		return fmt.Errorf("[ERR] an error occurred while tring to detach the volume %s", err)
	}
	return nil
}
