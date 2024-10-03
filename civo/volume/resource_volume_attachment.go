package volume

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceVolumeAttachment function returns a schema.Resource that represents a Volume Attachment.
// This can be used to create, read, update, and delete operations for a Volume Attachment in the infrastructure.
func ResourceVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Manages volume attachment/detachment to an instance.",
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The ID of target instance for attachment",
			},
			"volume_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The ID of target volume for attachment",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The region for the volume attachment",
			},
		},
		CreateContext: resourceVolumeAttachmentCreate,
		ReadContext:   resourceVolumeAttachmentRead,
		DeleteContext: resourceVolumeAttachmentDelete,
	}
}

// function to create the new volume
func resourceVolumeAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)
	var diags diag.Diagnostics

	// overwrite the region if it's defined
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	instanceID := d.Get("instance_id").(string)
	volumeID := d.Get("volume_id").(string)
	attachAtBoot := d.Get("attach_at_boot").(bool)

	log.Printf("[INFO] retrieving the volume %s", volumeID)
	volume, err := apiClient.FindVolume(volumeID)
	if err != nil {
		return diag.Errorf("[ERR] Error retrieving volume: %s", err)
	}

	if volume.InstanceID == "" || volume.InstanceID != instanceID {

		vuc := civogo.VolumeAttachConfig{
			InstanceID: instanceID,
			Region:     apiClient.Region,
		}

		if attachAtBoot {
			// Notify the terminal
			msg := fmt.Sprintf("To use the volume %s, The instance %s needs to be rebooted", volumeID, instanceID)
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  msg,
			})

			vuc.AttachAtBoot = true
		}

		log.Printf("[INFO] attaching the volume %s to instance %s", volumeID, instanceID)
		_, err := apiClient.AttachVolume(volumeID, vuc)
		if err != nil {
			return diag.Errorf("[ERR] error attaching volume to instance %s", err)
		}
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-%s-", instanceID, volumeID)))

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"attaching"},
		Target:  []string{"attached"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.FindVolume(volumeID)
			if err != nil {
				return 0, "", err
			}
			return resp, resp.Status, nil
		},
		Timeout:        60 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 10,
	}
	_, err = createStateConf.WaitForStateContext(context.Background())
	if err != nil {
		return diag.Errorf("error waiting for volume (%s) to be attached: %s", d.Id(), err)
	}

	return resourceVolumeAttachmentRead(ctx, d, m)
}

// function to read the volume
func resourceVolumeAttachmentRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it's defined
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	instanceID := d.Get("instance_id").(string)
	volumeID := d.Get("volume_id").(string)

	log.Printf("[INFO] retrieving the volume %s", volumeID)
	resp, err := apiClient.FindVolume(volumeID)
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}

		return diag.Errorf("[ERR] failed retrieving the volume: %s", err)
	}

	if resp.InstanceID == "" || resp.InstanceID != instanceID {
		log.Printf("[DEBUG] Volume Attachment (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return nil
}

// function to delete the volume
func resourceVolumeAttachmentDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it's defined
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	volumeID := d.Get("volume_id").(string)

	log.Printf("[INFO] Detaching the volume %s", d.Id())
	_, err := apiClient.DetachVolume(volumeID)
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while trying to detach the volume %s", err)
	}
	return nil
}
