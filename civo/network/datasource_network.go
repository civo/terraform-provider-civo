package network

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
				AtLeastOneOf: []string{"id", "label"},
			},
			"label": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				AtLeastOneOf: []string{"id", "label"},
				Description:  "The label of an existing network",
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The region an existing network lives in; used to scope the lookup by id or label",
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
			"cidr_v4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CIDR block for the network",
			},
			"nameservers_v4": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of nameservers for the network",
			},
		},
	}
}

func dataSourceNetworkRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient = utils.RegionalClient(apiClient, region.(string))
	}

	var foundNetwork *civogo.Network

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the network by id")
		network, err := apiClient.FindVPCNetwork(id.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive network: %s", err)
		}

		foundNetwork = network
	} else if label, ok := d.GetOk("label"); ok {
		log.Printf("[INFO] Getting the network by label")
		network, err := apiClient.FindVPCNetwork(label.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive network: %s", err)
		}

		foundNetwork = network
	}

	if foundNetwork == nil {
		return diag.Errorf("[ERR] network not found; specify a valid id or label")
	}

	d.SetId(foundNetwork.ID)
	d.Set("name", foundNetwork.Name)
	d.Set("label", foundNetwork.Label)
	d.Set("region", apiClient.Region)
	d.Set("default", foundNetwork.Default)
	d.Set("cidr_v4", foundNetwork.CIDR)
	d.Set("nameservers_v4", foundNetwork.NameserversV4)

	return nil
}
