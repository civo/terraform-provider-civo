package civo

import (
	"context"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Data source to get from the api a specific ObjectStore
// using the id or the name
func dataSourceObjectStore() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Get information of an Object Store for use in other resources. This data source provides all of the Object Store's properties as configured on your Civo account.",
			"Note: This data source returns a single Object Store. When specifying a name, an error will be raised if more than one Object Stores with the same name found.",
		}, "\n\n"),
		ReadContext: dataSourceObjectStoreRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the Object Store",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The name of the Object Store",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region of an existing Object Store",
			},
			"generated_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The generated name of the Object Store",
			},
			"max_size_gb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     500,
				Description: "The maximum size of the Object Store",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Description: "The endpoint of the Object Store",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the Object Store",
			},
		},
	}
}

func dataSourceObjectStoreRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	var foundStore *civogo.ObjectStore

	if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the Object Store by name")
		store, err := apiClient.FindObjectStore(name.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive Object Store: %s", err)
		}

		foundStore = store
	}

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the Object Store by name")
		store, err := apiClient.FindObjectStore(id.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive Object Store: %s", err)
		}

		foundStore = store
	}

	maxSize, err := utils.StringToInt(foundStore.MaxSize)
	if err != nil {
		return diag.Errorf("[ERR] failed to convert the max size to int: %s", err)
	}

	d.SetId(foundStore.ID)
	d.Set("name", foundStore.Name)
	d.Set("generated_name", foundStore.GeneratedName)
	d.Set("max_size_gb", maxSize)
	d.Set("endpoint", foundStore.ObjectStoreEndpoint)
	d.Set("status", foundStore.Status)

	return nil
}
