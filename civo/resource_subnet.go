package civo

import (
	"context"
	"log"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The resource Subnet represent a Subnet inside the cloud
func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Civo Subnet resource. This can be used to create, modify, and delete Subnets.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateName,
				Description:  "The name of the subnet",
			},
			"networkID": {
				Type:        schema.TypeString,
				Description: "The network ID of an existing subnet",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of an existing subnet",
			},
			"subnetSize": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The size of an existing subnet",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "It is the default subnet",
			},
		},
		CreateContext: resourceSubnetCreate,
		ReadContext:   resourceSubnetRead,
		UpdateContext: resourceSubnetUpdate,
		DeleteContext: resourceSubnetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// function to create a new Subnet
func resourceSubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] creating the new Subnet %s", d.Get("label").(string))
	subnet, err := apiClient.NewSubnet(d.Get("name").(string))
	if err != nil {
		return diag.Errorf("[ERR] failed to create a new Subnet: %s", err)
	}

	d.SetId(subnet.ID)

	return resourceSubnetRead(ctx, d, m)
}

// function to read a Subnet
func resourceSubnetRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	CurrentSubnet := civogo.Subnet{}

	log.Printf("[INFO] retriving the Subnet %s", d.Id())
	resp, err := apiClient.ListSubnets()
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}

		return diag.Errorf("[ERR] failed to list the Subnet: %s", err)
	}

	for _, net := range resp {
		if net.ID == d.Id() {
			CurrentSubnet = net
		}
	}

	d.Set("name", CurrentSubnet.Name)
	d.Set("networkID", CurrentSubnet.NetworkID)
	d.Set("subnetSize", CurrentSubnet.SubnetSize)
	d.Set("status", CurrentSubnet.Status)
	d.Set("default", CurrentSubnet.Default)
	return nil
}

// function to update the Subnet
func resourceSubnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	if d.HasChange("label") {
		log.Printf("[INFO] updating the Subnet %s", d.Id())
		_, err := apiClient.RenameSubnet(d.Get("name").(string), d.Id())
		if err != nil {
			return diag.Errorf("[ERR] An error occurred while rename the Subnet %s", d.Id())
		}
		return resourceSubnetRead(ctx, d, m)
	}
	return resourceSubnetRead(ctx, d, m)
}

// function to delete a Subnet
func resourceSubnetDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] deleting the Subnet %s", d.Id())
	_, err := apiClient.DeleteSubnet(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while tring to delete the Subnet %s", d.Id())
	}
	return nil
}
