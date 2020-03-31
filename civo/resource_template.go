package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

func resourceTemplate() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "This is a unqiue, alphanumerical, short, human readable code for the template (required).",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "This is a short human readable name for the template (optional).",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "This is the ID of a bootable volume, either owned by you or global" +
					"(optional; but must be specified if no image_id is specified).",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "This is the Image ID of any default template or the ID of another template," +
					"either owned by you or global (optional; but must be specified if no volume_id is specified).",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"short_description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "A one line description of the template (optional)e",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "A multi-line description of the template, in Markdown format (optional).",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"default_username": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The default username to suggest that the user creates (optional: defaults to civo).",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"cloud_config": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Commonly referred to as 'user-data', this is a customisation script that is run after" +
					"the instance is first booted. We recommend using cloud-config as it's a great distribution-agnostic" +
					"way of configuring cloud servers. If you put $INITIAL_USER in your script, this will automatically" +
					"be replaced by the initial user chosen when creating the instance, $INITIAL_PASSWORD will be" +
					"replaced with the random password generated by the system, $HOSTNAME is the fully qualified" +
					"domain name of the instance and $SSH_KEY will be the content of the SSH public key." +
					"(this is technically optional, but you won't really be able to use instances without it -" +
					"see our learn guide on templates for more information)",
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		Create: resourceTemplateCreate,
		Read:   resourceTemplateRead,
		Update: resourceTemplateUpdate,
		Delete: resourceTemplateDelete,
	}
}

func resourceTemplateCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	_, hasCode := d.GetOk("code")
	_, hasVolumeID := d.GetOk("volume_id")
	_, hasImageID := d.GetOk("image_id")

	if !hasCode {
		return fmt.Errorf("`code` must be assigned")
	}

	if !hasVolumeID && !hasImageID {
		return fmt.Errorf("`volume_id` or `image_id` must be assigned")
	}

	config := &civogo.Template{
		Code: d.Get("code").(string),
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("short_description"); ok {
		config.ShortDescription = attr.(string)
	}

	if attr, ok := d.GetOk("description"); ok {
		config.Description = attr.(string)
	}

	if attr, ok := d.GetOk("volume_id"); ok {
		config.VolumeID = attr.(string)
	}

	if attr, ok := d.GetOk("image_id"); ok {
		config.ImageID = attr.(string)
	}

	if attr, ok := d.GetOk("default_username"); ok {
		config.DefaultUsername = attr.(string)
	}

	if attr, ok := d.GetOk("cloud_config"); ok {
		config.CloudConfig = attr.(string)
	}

	resp, err := apiClient.NewTemplate(config)
	if err != nil {
		fmt.Errorf("[WARN] failed to create template: %s", err)
		return err
	}

	d.SetId(resp.ID)

	return resourceTemplateRead(d, m)
}

func resourceTemplateRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	resp, err := apiClient.GetTemplateByCode(d.Get("code").(string))
	if err != nil {
		fmt.Errorf("[WARN] failed to create template: %s", err)
		return err
	}

	d.Set("code", resp.Code)
	d.Set("name", resp.Name)
	d.Set("volume_id", resp.VolumeID)
	d.Set("image_id", resp.ImageID)
	d.Set("short_description", resp.ShortDescription)
	d.Set("description", resp.Description)
	d.Set("default_username", resp.DefaultUsername)
	d.Set("cloud_config", resp.CloudConfig)

	return nil
}

func resourceTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	conf := &civogo.Template{
		Name:             d.Get("name").(string),
		ShortDescription: d.Get("short_description").(string),
		Description:      d.Get("description").(string),
		DefaultUsername:  d.Get("default_username").(string),
		CloudConfig:      d.Get("cloud_config").(string),
		VolumeID:         d.Get("volume_id").(string),
		ImageID:          d.Get("image_id").(string),
	}

	_, err := apiClient.UpdateTemplate(d.Id(), conf)
	if err != nil {
		fmt.Errorf("[WARN] failed to update template: %s", err)
		return err
	}

	return resourceTemplateRead(d, m)
}

func resourceTemplateDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	_, err := apiClient.DeleteTemplate(d.Id())
	if err != nil {
		log.Printf("[INFO] civo template (%s) was delete", d.Id())
	}

	return nil
}
