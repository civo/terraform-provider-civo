package civo

import (
	"context"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The Object Store resource represents an ObjectStore object
// and with it you can handle the Object Stores created with Terraform.
func resourceObjectStore() *schema.Resource {
	return &schema.Resource{
		Description: "Provides an Object Store resource. This can be used to create, modify, and delete object stores.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateNameSize,
				Description:  "The name of the Object Store. Must be unique.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region for the Object Store, if not declared we use the region as declared in the provider (Defaults to LON1)",
			},
			"max_size_gb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     500,
				Description: "The maximum size of the Object Store. Default is 500GB.",
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The access key ID from the Object Store credential. If this is not set, a new credential will be created.",
			},
			"bucket_url": {
				Type:        schema.TypeString,
				Description: "The endpoint of the Object Store. It is generated by the provider.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the Object Store.",
			},
		},
		CreateContext: resourceObjectStoreCreate,
		ReadContext:   resourceObjectStoreRead,
		UpdateContext: resourceObjectStoreUpdate,
		DeleteContext: resourceObjectStoreDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

// Function to create an Object Store
func resourceObjectStoreCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] configuring the Object Store %s", d.Get("name").(string))
	config := &civogo.CreateObjectStoreRequest{
		Name:      d.Get("name").(string),
		MaxSizeGB: int64(d.Get("max_size_gb").(int)),
		Region:    apiClient.Region,
	}

	if AccessKeyID, ok := d.GetOk("access_key_id"); ok {
		config.AccessKeyID = AccessKeyID.(string)
	}

	log.Printf("[INFO] creating the Object Store %s", d.Get("name").(string))
	store, err := apiClient.NewObjectStore(config)
	if err != nil {
		return diag.Errorf("[ERR] failed to create Object Store: %s", err)
	}

	d.SetId(store.ID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"creating"},
		Target:  []string{"ready"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.GetObjectStore(d.Id())
			if err != nil {
				return 0, "", err
			}
			return resp, resp.Status, nil
		},
		Timeout:        60 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}
	_, err = createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for Object Store (%s) to be created: %s", d.Id(), err)
	}

	return resourceObjectStoreRead(ctx, d, m)

}

// Function to read Object Store
func resourceObjectStoreRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] retriving the Object Store %s", d.Id())
	resp, err := apiClient.GetObjectStore(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}

		return diag.Errorf("[ERR] failed to retrive the Object Store: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("max_size_gb", resp.MaxSize)
	d.Set("region", apiClient.Region)
	d.Set("access_key_id", resp.OwnerInfo.AccessKeyID)
	d.Set("bucket_url", resp.BucketURL)
	d.Set("status", resp.Status)

	return nil
}

// Function to update the Object Store
func resourceObjectStoreUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	_, err := apiClient.FindObjectStore(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] failed to find Object Store: %s", err)
	}

	config := &civogo.UpdateObjectStoreRequest{
		Region: apiClient.Region,
	}

	if d.HasChange("max_size_gb") {
		config.MaxSizeGB = int64(d.Get("max_size_gb").(int))
	}

	log.Printf("[INFO] updating the Object Store %s", d.Id())
	_, err = apiClient.UpdateObjectStore(d.Id(), config)
	if err != nil {
		return diag.Errorf("[ERR] failed to update Object Store: %s", err)
	}

	return resourceObjectStoreRead(ctx, d, m)
}

// Function to delete an Object Store
func resourceObjectStoreDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if it is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] deleting the Object Store %s", d.Id())
	_, err := apiClient.DeleteObjectStore(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while tring to delete the Object Store %s", d.Id())
	}
	return nil
}
