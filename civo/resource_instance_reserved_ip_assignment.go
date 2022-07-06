package civo

import (
	"context"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The instance reserved ip assignment resource schema definition
// represent the instance reserved ip assignment resource
func resourceInstanceReservedIPAssignment() *schema.Resource {
	return &schema.Resource{
		Description: "The instance reserved ip assignment resource schema definition",
		Schema: map[string]*schema.Schema{
			"reserved_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The reserved ip id",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The instance id",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region of the ip",
			},
		},
		CreateContext: resourceInstanceReservedIPCreate,
		ReadContext:   resourceInstanceReservedIPRead,
		DeleteContext: resourceInstanceReservedIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

// function to create a instance
func resourceInstanceReservedIPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	// We check if the instance is valid and if it is not we return an error
	instance, err := apiClient.GetInstance(d.Get("instance_id").(string))
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while tring to get instance %s", d.Get("instance_id").(string))
	}

	// We check if the reserved ip is valid and if it is not we return an error
	reservedIP, err := apiClient.FindIP(d.Get("reserved_ip_id").(string))
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while tring to get reserved ip %s", d.Get("reserved_ip_id").(string))
	}

	if reservedIP.AssignedTo.ID != "" {
		return diag.Errorf("[ERR] the reserved ip %s is already assigned to an instance", reservedIP.ID)
	}

	// We send to assign the reserved ip to the instance
	log.Printf("[INFO] assigning the reserved ip %s to the instance %s", d.Get("reserved_ip_id").(string), d.Get("instance_id").(string))

	assignedTo := &civogo.AssignedTo{
		ID:   instance.ID,
		Type: "instance",
		Name: instance.Hostname,
	}
	_, err = apiClient.AssignIP(reservedIP.ID, assignedTo)
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while tring to assign reserved ip %s to instance %s", d.Get("reserved_ip_id").(string), d.Get("instance_id").(string))
	}

	d.SetId(resource.UniqueId())

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"ASSIGNED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.GetInstance(instance.ID)
			if err != nil {
				return 0, "", err
			}
			if resp.PublicIP != reservedIP.IP {
				return 0, "PENDING", nil
			}
			return resp, "ASSIGNED", nil
		},
		Timeout:        60 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}
	_, err = createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for ip be assingne (%s) to the instance: %s", d.Id(), err)
	}

	return resourceInstanceReservedIPRead(ctx, d, m)

}

// function to read the instance
func resourceInstanceReservedIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	instanceID := d.Get("instance_id").(string)
	reservedID := d.Get("reserved_ip_id").(string)

	// We check if the reserved ip is valid and if it is not we return an error
	reservedIP, err := apiClient.FindIP(reservedID)
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while tring to get reserved ip %s", reservedID)
	}

	if reservedIP.AssignedTo.ID != instanceID {
		d.SetId("")
		return diag.Errorf("[ERR] the reserved ip %s is not assigned to the instance %s", reservedIP.ID, instanceID)
	}

	return nil
}

// function to delete instance
func resourceInstanceReservedIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	// We check if the reserved ip is valid and if it is not we return an error
	log.Printf("[INFO] unassign the ip (%s) from the instance", d.Id())
	_, err := apiClient.UnassignIP(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while tring to unassign the ip %s", d.Id())
	}

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"DONE"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.FindIP(d.Id())
			if err != nil {
				return 0, "", err
			}
			if resp.AssignedTo.ID != "" {
				return 0, "PENDING", nil
			}
			return resp, "DONE", nil
		},
		Timeout:        60 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}
	_, err = createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for ip be unassign (%s)", d.Id(), err)
	}

	return nil
}
