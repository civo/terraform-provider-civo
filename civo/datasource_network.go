package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

// Data source to get from the api a specific network
// using the id or the label
func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "label"},
			},
			"label": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "label"},
			},
			// Computed resource
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	var foundNetwork *civogo.Network

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the network by id")
		network, err := apiClient.FindNetwork(id.(string))
		if err != nil {
			fmt.Errorf("[ERR] failed to retrive network: %s", err)
			return err
		}

		foundNetwork = network
	} else if label, ok := d.GetOk("label"); ok {
		log.Printf("[INFO] Getting the network by label")
		network, err := apiClient.FindNetwork(label.(string))
		if err != nil {
			fmt.Errorf("[ERR] failed to retrive network: %s", err)
			return err
		}

		foundNetwork = network
	}

	d.SetId(foundNetwork.ID)
	d.Set("name", foundNetwork.Name)
	d.Set("label", foundNetwork.Label)
	d.Set("region", foundNetwork.Region)
	d.Set("default", foundNetwork.Default)
	d.Set("cidr", foundNetwork.CIDR)

	return nil
}
