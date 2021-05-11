package civo

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	_ "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Provider Civo cloud provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CIVO_TOKEN", ""),
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				DefaultFunc: schema.EnvDefaultFunc("CIVO_REGION", ""),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"civo_template":           dataSourceTemplate(),
			"civo_kubernetes_version": dataSourceKubernetesVersion(),
			"civo_kubernetes_cluster": dataSourceKubernetesCluster(),
			"civo_instances_size":     dataSourceInstancesSize(),
			"civo_instances":          dataSourceInstances(),
			"civo_instance":           dataSourceInstance(),
			"civo_dns_domain_name":    dataSourceDNSDomainName(),
			"civo_dns_domain_record":  dataSourceDNSDomainRecord(),
			"civo_network":            dataSourceNetwork(),
			"civo_volume":             dataSourceVolume(),
			"civo_loadbalancer":       dataSourceLoadBalancer(),
			"civo_ssh_key":            dataSourceSSHKey(),
			"civo_snapshot":           dataSourceSnapshot(),
			"civo_region":             dataSourceRegion(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"civo_instance":             resourceInstance(),
			"civo_network":              resourceNetwork(),
			"civo_volume":               resourceVolume(),
			"civo_volume_attachment":    resourceVolumeAttachment(),
			"civo_dns_domain_name":      resourceDNSDomainName(),
			"civo_dns_domain_record":    resourceDNSDomainRecord(),
			"civo_firewall":             resourceFirewall(),
			"civo_firewall_rule":        resourceFirewallRule(),
			"civo_loadbalancer":         resourceLoadBalancer(),
			"civo_ssh_key":              resourceSSHKey(),
			"civo_template":             resourceTemplate(),
			"civo_snapshot":             resourceSnapshot(),
			"civo_kubernetes_cluster":   resourceKubernetesCluster(),
			"civo_kubernetes_node_pool": resourceKubernetesClusterNodePool(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// Provider configuration
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var regionValue, tokenValue string

	if region, ok := d.GetOk("region"); ok {
		regionValue = region.(string)
	}

	if token, ok := d.GetOk("token"); ok {
		tokenValue = token.(string)
	} else {
		return nil, fmt.Errorf("[ERR] token not found")
	}

	client, err := civogo.NewClient(tokenValue, regionValue)
	if err != nil {
		return nil, err
	}
	return client, nil

}
