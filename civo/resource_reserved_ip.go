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

func resourceReservedIP() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Civo reserved IP to represent a publicly-accessible static IP addresses that can be mapped to one of your Instancesor Load Balancer.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name for the ip address",
				ValidateFunc: utils.ValidateName,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region of the ip",
			},
			// Computed resource
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP Address of the resource",
			},
		},
		CreateContext: resourceReservedIPCreate,
		ReadContext:   resourceReservedIPRead,
		UpdateContext: resourceReservedIPUpdate,
		DeleteContext: resourceReservedIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// function to create a new IP resource
func resourceReservedIPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] creating the new ip address %s", d.Get("name").(string))
	newIP := &civogo.CreateIPRequest{
		Name:   d.Get("name").(string),
		Region: apiClient.Region,
	}
	ipAddress, err := apiClient.NewIP(newIP)
	if err != nil {
		return diag.Errorf("[ERR] failed to create a new ip address: %s", err)
	}

	d.SetId(ipAddress.ID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"BUILDING"},
		Target:  []string{"ACTIVE"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.FindIP(d.Id())
			if err != nil {
				return 0, "", err
			}
			if resp.IP == "" {
				return 0, "BUILDING", nil
			}
			return resp, "ACTIVE", nil
		},
		Timeout:        60 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 10,
	}
	_, err = createStateConf.WaitForStateContext(context.Background())
	if err != nil {
		return diag.Errorf("error waiting for ip resource (%s) to be created: %s", d.Id(), err)
	}

	return resourceReservedIPRead(ctx, d, m)
}

// function to read a the IP resource
func resourceReservedIPRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] retriving the ip address %s", d.Id())
	resp, err := apiClient.FindIP(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}

		return diag.Errorf("[ERR] failed to get the ips: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("region", apiClient.Region)
	d.Set("ip", resp.IP)

	return nil
}

// function to update the IP resource
func resourceReservedIPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	if d.HasChange("name") {
		log.Printf("[INFO] updating the iop name %s", d.Id())
		ipUpdate := &civogo.UpdateIPRequest{
			Name: d.Get("name").(string),
		}
		_, err := apiClient.UpdateIP(d.Id(), ipUpdate)
		if err != nil {
			return diag.Errorf("[ERR] An error occurred while rename the ip resource %s", d.Id())
		}
		return resourceReservedIPRead(ctx, d, m)
	}
	return resourceReservedIPRead(ctx, d, m)
}

// function to delete a network
func resourceReservedIPDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] deleting the ip resource %s", d.Id())
	_, err := apiClient.DeleteIP(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while trying to delete the ip resource %s", d.Id())
	}
	return nil
}
