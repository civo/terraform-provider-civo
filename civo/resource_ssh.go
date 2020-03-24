package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceSSHKey() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "a string that will be the reference for the SSH key.",
				ValidateFunc: validateName,
			},
			"public_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "a string containing the SSH public key.",
				Sensitive:   true,
				ForceNew:    true,
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "a string containing the SSH finger print.",
			},
		},
		Create: resourceSSHKeyCreate,
		Read:   resourceSSHKeyRead,
		Update: resourceSSHKeyUpdate,
		Delete: resourceSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceSSHKeyCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	sshKey, err := apiClient.NewSSHKey(d.Get("name").(string), d.Get("public_key").(string))
	if err != nil {
		fmt.Errorf("[WARN] failed to create a new ssh key: %s", err)
		return err
	}

	d.SetId(sshKey.ID)

	return resourceSSHKeyRead(d, m)
}

func resourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	sshKey, err := apiClient.FindSSHKey(d.Id())
	if err != nil {
		if sshKey != nil {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[WARN] error retrieving ssh key: %s", err)
	}

	d.Set("name", sshKey.Name)
	d.Set("fingerprint", sshKey.Fingerprint)

	return nil
}

func resourceSSHKeyUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	if d.HasChange("name") {
		if d.Get("name").(string) != "" {
			_, err := apiClient.UpdateSSHKey(d.Get("name").(string), d.Id())
			if err != nil {
				log.Printf("[WARN] an error occurred while trying to rename the ssh key (%s)", d.Id())
			}
		}
	}

	return resourceSSHKeyRead(d, m)
}

func resourceSSHKeyDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	_, err := apiClient.DeleteSSHKey(d.Id())
	if err != nil {
		log.Printf("[INFO] civo ssh key (%s) was delete", d.Id())
	}
	return nil
}
