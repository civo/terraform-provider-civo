package network

import (
	"context"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceVPCSubnet function returns a schema.Resource that represents a VPC Subnet.
// This can be used to query and retrieve details about a specific VPC Subnet in the infrastructure using its id or name.
func DataSourceVPCSubnet() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Retrieve information about a VPC subnet for use in other resources.",
			"This data source provides all of the subnet's properties as configured on your Civo account.",
			"Subnets may be looked up by id or name, and require the network_id.",
		}, "\n\n"),
		ReadContext: dataSourceVPCSubnetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The ID of the VPC subnet",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The name of the VPC subnet",
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the VPC network this subnet belongs to",
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The region of the subnet",
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
	}
}

func dataSourceVPCSubnetRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	networkID := d.Get("network_id").(string)

	var searchBy string
	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the VPC subnet by id")
		searchBy = id.(string)
	} else if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the VPC subnet by name")
		searchBy = name.(string)
	}

	subnet, err := apiClient.FindVPCSubnet(searchBy, networkID)
	if err != nil {
		return diag.Errorf("[ERR] failed to retrieve VPC subnet: %s", err)
	}

	d.SetId(subnet.ID)
	d.Set("name", subnet.Name)
	d.Set("network_id", subnet.NetworkID)
	d.Set("subnet_size", subnet.SubnetSize)
	d.Set("status", subnet.Status)
	d.Set("region", apiClient.Region)

	return nil
}
