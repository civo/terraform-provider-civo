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

// Data source to get from the api a specific subnet
// using the id or the label
func dataSourceSubnet() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Retrieve information about a subnet for use in other resources.",
			"This data source provides all of the subnet's properties as configured on your Civo account.",
			"Subnets may be looked up by id or label, and you can optionally pass region if you want to make a lookup for an expecific subnet inside that region.",
		}, "\n\n"),
		ReadContext: dataSourceSubnetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
				AtLeastOneOf: []string{"id", "label", "region"},
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the subnet",
			},
			"networkID": {
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The network ID of an existing subnet",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The status of an existing subnet",
			},
			"subnetSize": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The size of an existing subnet",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "It is the default subnet",
			},
		},
	}
}

func dataSourcesubnetRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	var foundSubnet *civogo.Subnet

	if networkID, ok := d.GetOk("networkID"); ok {
		log.Printf("[INFO] Getting the subnet by id")
		subnet, err := apiClient.FindNetwork(networkID.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive subnet: %s", err)
		}

		foundSubnet = subnet
	}

	d.SetId(foundSubnet.ID)
	d.Set("name", foundSubnet.Name)
	d.Set("networkID", foundSubnet.NetworkID)
	d.Set("status", foundSubnet.Status)
	d.Set("subnetSize", foundSubnet.SubnetSize)
	d.Set("default", foundSubnet.Default)

	return nil
}
