package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTemplateRead,
		Schema: map[string]*schema.Schema{
			"code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "code of the image",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the image",
			},
			"volume_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "volume_id of the image",
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "image_id of the image",
			},
			"short_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "short_description of the image",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "description of the image",
			},
			"default_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "default_username of the image",
			},
			"cloud_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "cloud_config of the image",
			},
		},
	}
}

func dataSourceTemplateRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	code, hasCode := d.GetOk("code")

	if !hasCode {
		return fmt.Errorf("`code` must be assigned")
	}

	if hasCode {
		image, err := apiClient.GetTemplateByCode(code.(string))
		if err != nil {
			fmt.Errorf("[ERR] failed to retrive template: %s", err)
			return err
		}

		d.SetId(image.ID)
		d.Set("code", image.Code)
		d.Set("name", image.Name)
		d.Set("volume_id", image.VolumeID)
		d.Set("image_id", image.ImageID)
		d.Set("short_description", image.ShortDescription)
		d.Set("description", image.Description)
		d.Set("default_username", image.DefaultUsername)
		d.Set("cloud_config", image.CloudConfig)
	}

	return nil
}
