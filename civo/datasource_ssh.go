package civo

import (
	"context"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Data source to get from the api a specific ssh key
// using the id or the name
func dataSourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Get information on a SSH key. This data source provides the name, and fingerprint as configured on your Civo account.",
			"An error will be raised if the provided SSH key name does not exist in your Civo account.",
		}, "\n\n"),
		ReadContext: dataSourceSSHKeyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The name of the SSH key",
			},
			// Computed resource
			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fingerprint of the public key of the SSH key",
			},
		},
	}
}

func dataSourceSSHKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	var searchBy string

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the ssh key by id")
		searchBy = id.(string)
	} else if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the ssh key by label")
		searchBy = name.(string)
	}

	sshKey, err := apiClient.FindSSHKey(searchBy)
	if err != nil {
		return diag.Errorf("[ERR] failed to retrive network: %s", err)
	}

	d.SetId(sshKey.ID)
	d.Set("name", sshKey.Name)
	d.Set("fingerprint", sshKey.Fingerprint)

	return nil
}
