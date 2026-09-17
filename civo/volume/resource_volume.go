package volume

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceVolume function returns a schema.Resource that represents a Volume.
// This can be used to create, read, update, and delete operations for a Volume in the infrastructure.
func ResourceVolume() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Civo volume which can be attached to an instance in order to provide expanded storage.",
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
				Description: "A minimum of 1 and a maximum of your available disk space from your quota specifies the size of the volume in gigabytes. Increases are applied in place, on an attached volume too when its volume type supports online expansion; decreases are rejected at plan time.",
			},
			"region": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "The region for the volume, if not declare we use the region in declared in the provider.",
				DiffSuppressFunc: utils.IgnoreCaseDiff,
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network that the volume belongs to",
			},
			// Computed resource
			"mount_point": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The mount point of the volume (from instance's perspective)",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the volume",
			},
		},
		CreateContext: resourceVolumeCreate,
		ReadContext:   resourceVolumeRead,
		UpdateContext: resourceVolumeUpdate,
		DeleteContext: resourceVolumeDelete,
		CustomizeDiff: customizeDiffVolume,
		Importer: &schema.ResourceImporter{
			State: resourceVolumeImport,
		},
	}
}

// customizeDiffVolume rejects a size decrease at plan time. Volume size is additive only: the API
// refuses a shrink, so failing the plan is friendlier than failing the apply.
func customizeDiffVolume(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.Id() == "" || !diff.HasChange("size_gb") {
		return nil
	}
	oldSize, newSize := diff.GetChange("size_gb")
	if newSize.(int) < oldSize.(int) {
		return fmt.Errorf("size_gb cannot be decreased (from %d to %d): volume size is additive only", oldSize.(int), newSize.(int))
	}
	return nil
}

// function to create the new volume
func resourceVolumeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	apiClient, err := utils.RegionalClient(apiClient, utils.ResolveRegion(d, utils.NetworkRef("network_id")))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] configuring the volume %s", d.Get("name").(string))

	config := &civogo.VolumeConfig{
		Name:          d.Get("name").(string),
		SizeGigabytes: d.Get("size_gb").(int),
		NetworkID:     d.Get("network_id").(string),
		Region:        apiClient.Region,
	}

	if v, ok := d.GetOk("volume_type"); ok {
		config.VolumeType = v.(string)
	}

	_, err = apiClient.FindNetwork(config.NetworkID)
	if err != nil {
		return diag.Errorf("[ERR] Unable to find network ID %q in %q region", config.NetworkID, config.Region)
	}

	volume, err := apiClient.NewVolume(config)
	if err != nil {
		return diag.Errorf("[ERR] failed to create a new volume: %s", err)
	}

	d.SetId(volume.ID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"creating"},
		Target:  []string{"available"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.FindVolume(d.Id())
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
		return diag.Errorf("error waiting for volume (%s) to be created: %s", d.Id(), err)
	}

	return resourceVolumeRead(ctx, d, m)
}

// function to read the volume
func resourceVolumeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	apiClient, err := utils.RegionalClient(apiClient, utils.ResolveRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] retrieving the volume %s", d.Id())
	resp, err := apiClient.FindVolume(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("[ERR] failed retrieving the volume: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("network_id", resp.NetworkID)
	d.Set("size_gb", resp.SizeGigabytes)
	d.Set("mount_point", resp.MountPoint)
	d.Set("volume_type", resp.VolumeType)

	return nil
}

// function to update the volume
func resourceVolumeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*civogo.Client)

	apiClient, err := utils.RegionalClient(apiClient, utils.ResolveRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] retrieving the volume %s", d.Id())
	resp, err := apiClient.FindVolume(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] failed retrieving the volume: %s", err)
	}

	if d.HasChange("size_gb") {
		// Resize in place, attached or not. The API grows an attached volume online when its volume
		// type supports expansion, and refuses with a clear message when it does not (code
		// "volume_onlie_resize": detach the volume, or move it to an expandable volume type). That
		// verdict belongs to the user, so it is surfaced rather than worked around: the previous
		// detach / resize / re-attach sequence took the instance's storage offline for every resize
		// and re-attached hotplugged volumes as boot attachments.
		newSize := d.Get("size_gb").(int)
		log.Printf("[INFO] resizing the volume %s to %dGB (attached to %q)", d.Id(), newSize, resp.InstanceID)
		if _, err := apiClient.ResizeVolume(d.Id(), newSize); err != nil {
			return diag.Errorf("[ERR] failed to resize the volume %s to %dGB: %s", d.Id(), newSize, err)
		}

		// The API reports the requested size as soon as the resize is admitted, so a delivered
		// resize cannot be told apart from a pending one from here; wait only for the volume to
		// leave the transitional states before reading it back. Any other state (available,
		// attached, attaching, ...) is a settled one for this purpose.
		err = resource.RetryContext(ctx, 60*time.Minute, func() *resource.RetryError {
			volume, err := apiClient.FindVolume(d.Id())
			if err != nil {
				return resource.NonRetryableError(err)
			}
			switch volume.Status {
			case "resizing", "migrating":
				return resource.RetryableError(fmt.Errorf("volume %s is still %s", d.Id(), volume.Status))
			}
			return nil
		})
		if err != nil {
			return diag.Errorf("[ERR] error waiting for the volume %s to finish resizing: %s", d.Id(), err)
		}
	}

	if d.HasChange("network_id") {
		return diag.Errorf("[ERR] Network change for volume is not supported at this moment")
	}

	if d.HasChange("name") {
		return diag.Errorf("[ERR] Name change for volume is not supported at this moment")
	}

	return resourceVolumeRead(ctx, d, m)
}

// function to delete the volume
func resourceVolumeDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	apiClient, err := utils.RegionalClient(apiClient, utils.ResolveRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] deleting the volume %s", d.Id())
	_, err = apiClient.DeleteVolume(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while trying to delete the volume %s", err)
	}
	return nil
}

// custom import to able to import a volume
func resourceVolumeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	apiClient := m.(*civogo.Client)
	regions, err := apiClient.ListRegions()
	if err != nil {
		return nil, err
	}

	volumeFound := false
	for _, region := range regions {
		if volumeFound {
			break
		}

		currentRegion := region.Code
		apiClient, err = utils.RegionalClient(apiClient, utils.WithRegion(currentRegion))
		if err != nil {
			return nil, err
		}

		volumes, err := apiClient.ListVolumes()
		if err != nil {
			return nil, err
		}

		for _, volume := range volumes {
			if volume.ID == d.Id() {
				volumeFound = true
				d.SetId(volume.ID)
				d.Set("name", volume.Name)
				d.Set("network_id", volume.NetworkID)
				d.Set("region", currentRegion)
				d.Set("size_gb", volume.SizeGigabytes)
				d.Set("mount_point", volume.MountPoint)
			}
		}
	}

	if !volumeFound {
		return nil, fmt.Errorf("[ERR] Volume %s not found", d.Id())
	}

	return []*schema.ResourceData{d}, nil
}
