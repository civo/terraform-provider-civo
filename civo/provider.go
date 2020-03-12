package civo

import (
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CIVO_TOKEN", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"civo_instance":          resourceInstance(),
			"civo_network":           resourceNetwork(),
			"civo_volume":            resourceVolume(),
			"civo_dns_domain_name":   resourceDnsDomainName(),
			"civo_dns_domain_record": resourceDnsDomainRecord(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	token := d.Get("token").(string)
	client, _ := civogo.NewClient(token)
	return client, nil
}
