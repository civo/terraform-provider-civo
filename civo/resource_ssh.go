package civo

import (
	"context"
	"log"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Ssh resource, with this we can create and manage all Snapshot
func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Civo SSH key resource to allow you to manage SSH keys for instance access. Keys created with this resource can be referenced in your instance configuration via their ID.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "a string that will be the reference for the SSH key.",
				ValidateFunc: utils.ValidateName,
			},
			"public_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "a string containing the SSH public key.",
				ForceNew:    true,
			},
			// Computed resource
			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "a string containing the SSH finger print.",
			},
		},
		CreateContext: resourceSSHKeyCreate,
		ReadContext:   resourceSSHKeyRead,
		UpdateContext: resourceSSHKeyUpdate,
		DeleteContext: resourceSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// function to create a new ssh key
func resourceSSHKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] creating the new ssh key %s", d.Get("name").(string))
	sshKey, err := apiClient.NewSSHKey(d.Get("name").(string), d.Get("public_key").(string))
	if err != nil {
		return diag.Errorf("[ERR] failed to create a new ssh key: %s", err)
	}

	d.SetId(sshKey.ID)

	return resourceSSHKeyRead(ctx, d, m)
}

// function to read a ssh key
func resourceSSHKeyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] retrieving the new ssh key %s", d.Get("name").(string))
	sshKey, err := apiClient.FindSSHKey(d.Id())
	if err != nil {
		if sshKey == nil {
			d.SetId("")
			return nil
		}

		return diag.Errorf("[ERR] error retrieving ssh key: %s", err)
	}

	d.Set("name", sshKey.Name)
	d.Set("fingerprint", sshKey.Fingerprint)

	return nil
}

// function to update the ssh key
func resourceSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	if d.HasChange("name") {
		if d.Get("name").(string) != "" {
			log.Printf("[INFO] updating the ssh key %s", d.Get("name").(string))
			_, err := apiClient.UpdateSSHKey(d.Get("name").(string), d.Id())
			if err != nil {
				return diag.Errorf("[ERR] an error occurred while trying to rename the ssh key %s", d.Id())
			}
		}
	}

	return resourceSSHKeyRead(ctx, d, m)
}

// function to delete the ssh key
func resourceSSHKeyDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] deleting the ssh key %s", d.Id())
	_, err := apiClient.DeleteSSHKey(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while trying to delete the ssh key %s", d.Id())
	}
	return nil
}
