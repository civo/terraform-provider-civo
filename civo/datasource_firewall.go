package civo

import (
	"fmt"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Data source to get from the api a specific firewall
// using the id or the label
func dataSourceFirewall() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Retrieve information about a firewall for use in other resources.",
			"This data source provides all of the firewall's properties as configured on your Civo account.",
			"Firewalls may be looked up by id or name, and you can optionally pass region if you want to make a lookup for an expecific firewall inside that region.",
		}, "\n\n"),
		Read: dataSourceFirewallRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The name of the firewall",
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The region where the firewall is",
			},
			// Computed resource
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the associated network",
			},
		},
	}
}

func dataSourceFirewallRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	var foundFirewall *civogo.Firewall

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the firewall by id")
		firewall, err := apiClient.FindFirewall(id.(string))
		if err != nil {
			return fmt.Errorf("[ERR] failed to retrive firewall: %s", err)
		}

		foundFirewall = firewall
	} else if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the firewall by name")
		firewall, err := apiClient.FindFirewall(name.(string))
		if err != nil {
			return fmt.Errorf("[ERR] failed to retrive firewall: %s", err)
		}

		foundFirewall = firewall
	}

	d.SetId(foundFirewall.ID)
	d.Set("name", foundFirewall.Name)
	d.Set("network_id", foundFirewall.NetworkID)
	d.Set("region", apiClient.Region)

	return nil
}
