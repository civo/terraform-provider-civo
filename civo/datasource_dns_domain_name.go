package civo

import (
	"fmt"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// dataSourceDnsDomainName data source to get from the api a specific domain
// using the id or the name of the domain
func dataSourceDNSDomainName() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Get information on a domain. This data source provides the name and the id.",
			"An error will be raised if the provided domain name is not in your Civo account.",
		}, "\n\n"),
		Read: dataSourceDNSDomainNameRead,
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
				Description:  "The name of the domain",
			},
		},
	}
}

func dataSourceDNSDomainNameRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	var foundDomain *civogo.DNSDomain

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the domain by id")
		domain, err := apiClient.FindDNSDomain(id.(string))
		if err != nil {
			return fmt.Errorf("[ERR] failed to retrive domain: %s", err)
		}

		foundDomain = domain
	} else if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the domain by name")
		image, err := apiClient.FindDNSDomain(name.(string))
		if err != nil {
			return fmt.Errorf("[ERR] failed to retrive domain: %s", err)
		}

		foundDomain = image
	}

	d.SetId(foundDomain.ID)
	d.Set("name", foundDomain.Name)

	return nil
}
