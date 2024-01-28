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

// DataSourceNetwork function returns a schema.Resource that represents a Network.
// This can be used to query and retrieve details about a specific Network in the infrastructure using its id or label.
func DataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Retrieve information about a network for use in other resources.",
			"This data source provides all of the network's properties as configured on your Civo account.",
			"Networks may be looked up by id or label, and you can optionally pass region if you want to make a lookup for a specific network inside that region.",
		}, "\n\n"),
		ReadContext: dataSourceNetworkRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				AtLeastOneOf: []string{"id", "label", "region"},
			},
			"label": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				AtLeastOneOf: []string{"id", "label", "region"},
				Description:  "The label of an existing network",
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				AtLeastOneOf: []string{"id", "label", "region"},
				Description:  "The region of an existing network",
			},
			// Computed resource
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the network",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If is the default network",
			},
		},
	}
}

func dataSourceNetworkRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	var foundNetwork *civogo.Network

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the network by id")
		network, err := apiClient.FindNetwork(id.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive network: %s", err)
		}

		foundNetwork = network
	} else if label, ok := d.GetOk("label"); ok {
		log.Printf("[INFO] Getting the network by label")
		network, err := apiClient.FindNetwork(label.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive network: %s", err)
		}

		foundNetwork = network
	}

	d.SetId(foundNetwork.ID)
	d.Set("name", foundNetwork.Name)
	d.Set("label", foundNetwork.Label)
	d.Set("region", apiClient.Region)
	d.Set("default", foundNetwork.Default)

	return nil
}
