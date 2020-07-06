package civo

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider Civo cloud provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CIVO_TOKEN", ""),
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
		},
		ResourcesMap: map[string]*schema.Resource{
			"civo_instance":           resourceInstance(),
			"civo_network":            resourceNetwork(),
			"civo_volume":             resourceVolume(),
			"civo_volume_attachment":  resourceVolumeAttachment(),
			"civo_dns_domain_name":    resourceDNSDomainName(),
			"civo_dns_domain_record":  resourceDNSDomainRecord(),
			"civo_firewall":           resourceFirewall(),
			"civo_firewall_rule":      resourceFirewallRule(),
			"civo_loadbalancer":       resourceLoadBalancer(),
			"civo_ssh_key":            resourceSSHKey(),
			"civo_template":           resourceTemplate(),
			"civo_snapshot":           resourceSnapshot(),
			"civo_kubernetes_cluster": resourceKubernetesCluster(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// Provider configuration
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	if token, ok := d.GetOk("token"); ok {
		client, _ := civogo.NewClient(token.(string))
		return client, nil
	}

	return nil, fmt.Errorf("[ERR] token not found")

}
