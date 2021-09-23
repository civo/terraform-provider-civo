package civo

import (
	"fmt"
	"log"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Volume resource, with this we can create and manage all volume
func resourceVolume() *schema.Resource {
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
				Description: "A minimum of 1 and a maximum of your available disk space from your quota specifies the size of the volume in gigabytes ",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region for the volume, if not declare we use the region in declared in the provider.",
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
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of this resource.",
			},
		},
		Create: resourceVolumeCreate,
		Read:   resourceVolumeRead,
		Update: resourceVolumeUpdate,
		Delete: resourceVolumeDelete,
		//Exists: resourceExistsItem,
		Importer: &schema.ResourceImporter{
			State: resourceVolumeImport,
		},
	}
}

// function to create the new volume
func resourceVolumeCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] configuring the volume %s", d.Get("name").(string))
	config := &civogo.VolumeConfig{
		Name:          d.Get("name").(string),
		SizeGigabytes: d.Get("size_gb").(int),
		NetworkID:     d.Get("network_id").(string),

		// if "region" is set at provider level, use it
		Region: apiClient.Region,
	}

	// if "region" is set in configuration file, use it
	if region, ok := d.GetOk("region"); ok {
		currentRegion := region.(string)
		apiClient.Region = currentRegion
		config.Region = currentRegion
	}

	_, err := apiClient.FindNetwork(config.NetworkID)
	if err != nil {
		return fmt.Errorf("[ERR] Unable to find network ID %q in %q region", config.NetworkID, config.Region)
	}

	volume, err := apiClient.NewVolume(config)
	if err != nil {
		return fmt.Errorf("[ERR] failed to create a new volume: %s", err)
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
	d.Set("network_id", resp.NetworkID)
	d.Set("size_gb", resp.SizeGigabytes)
	d.Set("mount_point", resp.MountPoint)

	return nil
}

// function to update the volume
func resourceVolumeUpdate(d *schema.ResourceData, m interface{}) error {
	/*
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
	*/

	if d.HasChange("size_gb") {
		return fmt.Errorf("[ERR] Resize operation is not available at this moment - we are working to re-enable it soon")

		/*
			if resp.InstanceID != "" {
				_, err := apiClient.DetachVolume(d.Id())
				if err != nil {
					return fmt.Errorf("[WARN] an error occurred while trying to detach volume %s, %s", d.Id(), err)
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
		*/
	}

	if d.HasChange("network_id") {
		return fmt.Errorf("[ERR] Network change for volume is not supported at this moment")
	}

	if d.HasChange("name") {
		return fmt.Errorf("[ERR] Name change for volume is not supported at this moment")
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
		apiClient.Region = currentRegion

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
