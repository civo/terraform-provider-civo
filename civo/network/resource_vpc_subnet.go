package network

import (
	"context"
	"log"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceVPCSubnet function returns a schema.Resource that represents a VPC Subnet.
// This can be used to create, read, and delete operations for a VPC Subnet in the infrastructure.
func ResourceVPCSubnet() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Civo VPC subnet resource. This can be used to create and delete subnets within a VPC network.",
		CreateContext: resourceVPCSubnetCreate,
		ReadContext:   resourceVPCSubnetRead,
		DeleteContext: resourceVPCSubnetDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name for the VPC subnet",
				ValidateFunc: utils.ValidateName,
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the VPC network this subnet belongs to",
			},
			"region": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				Description:      "The region of the subnet",
				DiffSuppressFunc: utils.IgnoreCaseDiff,
			},
			"subnet_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The size of the subnet",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the subnet",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVPCSubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	networkID := d.Get("network_id").(string)
	subnetConfig := civogo.SubnetConfig{
		Name: d.Get("name").(string),
	}

	log.Printf("[INFO] creating VPC subnet %s in network %s", subnetConfig.Name, networkID)
	subnet, err := apiClient.CreateVPCSubnet(networkID, subnetConfig)
	if err != nil {
		return diag.Errorf("[ERR] failed to create VPC subnet: %s", err)
	}

	d.SetId(subnet.ID)

	return resourceVPCSubnetRead(ctx, d, m)
}

func resourceVPCSubnetRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	networkID := d.Get("network_id").(string)

	log.Printf("[INFO] retrieving VPC subnet %s", d.Id())
	subnet, err := apiClient.GetVPCSubnet(networkID, d.Id())
	if err != nil {
		if subnet == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("[ERR] failed to get VPC subnet: %s", err)
	}

	d.Set("name", subnet.Name)
	d.Set("network_id", subnet.NetworkID)
	d.Set("subnet_size", subnet.SubnetSize)
	d.Set("status", subnet.Status)
	d.Set("region", apiClient.Region)

	return nil
}

func resourceVPCSubnetDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	networkID := d.Get("network_id").(string)

	log.Printf("[INFO] deleting VPC subnet %s from network %s", d.Id(), networkID)
	_, err := apiClient.DeleteVPCSubnet(networkID, d.Id())
	if err != nil {
		return diag.Errorf("[ERR] failed to delete VPC subnet: %s", err)
	}

	return nil
}
