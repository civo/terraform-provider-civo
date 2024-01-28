package objectstorage

import (
	"context"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceObjectStoreCredential function returns a schema.Resource that represents an Object Store Credential.
// This can be used to query and retrieve details about a specific Object Store Credential in the infrastructure using its id or name.
func DataSourceObjectStoreCredential() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Get information of an Object Store Credential for use in other resources. This data source provides all of the Object Store Credential's properties as configured on your Civo account.",
			"Note: This data source returns a single Object Store Credential. When specifying a name, an error will be raised if more than one Object Store Credentials with the same name found.",
		}, "\n\n"),
		ReadContext: dataSourceObjectStoreCredentialRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the Object Store Credential",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The name of the Object Store Credential",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region of an existing Object Store",
			},
			// Computed values
			"access_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The access key id of the Object Store Credential",
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret access key of the Object Store Credential",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the Object Store Credential",
			},
		},
	}
}

func dataSourceObjectStoreCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	var foundStoreCredential *civogo.ObjectStoreCredential

	if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the Object Store Credential by name")
		storeCredential, err := apiClient.FindObjectStoreCredential(name.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive Object Store Credential: %s", err)
		}

		foundStoreCredential = storeCredential
	}

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the Object Store Credential by name")
		storeCredential, err := apiClient.FindObjectStoreCredential(id.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive Object Store Credential: %s", err)
		}

		foundStoreCredential = storeCredential
	}

	d.SetId(foundStoreCredential.ID)
	d.Set("name", foundStoreCredential.Name)
	d.Set("region", apiClient.Region)
	d.Set("access_key_id", foundStoreCredential.AccessKeyID)
	d.Set("secret_access_key", foundStoreCredential.SecretAccessKeyID)
	d.Set("status", foundStoreCredential.Status)

	return nil
}
