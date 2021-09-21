package civo

import (
	"fmt"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Data source to get from the api a specific network
// using the id or the label
func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Retrieve information about a network for use in other resources.",
			"This data source provides all of the network's properties as configured on your Civo account.",
			"Networks may be looked up by id or label, and you can optionally pass region if you want to make a lookup for an expecific network inside that region.",
		}, "\n\n"),
		Read: dataSourceNetworkRead,
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

func dataSourceNetworkRead(d *schema.ResourceData, m interface{}) error {
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
			return fmt.Errorf("[ERR] failed to retrive network: %s", err)
		}

		foundNetwork = network
	} else if label, ok := d.GetOk("label"); ok {
		log.Printf("[INFO] Getting the network by label")
		network, err := apiClient.FindNetwork(label.(string))
		if err != nil {
			return fmt.Errorf("[ERR] failed to retrive network: %s", err)
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
