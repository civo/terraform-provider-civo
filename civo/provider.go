package civo

import (
	"fmt"
	"log"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// ProviderVersion is the version of the provider to set in the User-Agent header
	ProviderVersion = "dev"

	// ProdAPI is the Base URL for CIVO Production API
	ProdAPI = "https://api.civo.com"
)

// Provider Civo cloud provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CIVO_TOKEN", ""),
				Description: "This is the Civo API token. Alternatively, this can also be specified using `CIVO_TOKEN` environment variable.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CIVO_REGION", ""),
				Description: "If region is not set, then no region will be used and them you need expensify in every resource even if you expensify here you can overwrite in a resource.",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CIVO_API_URL", ProdAPI),
				Description: "The Base URL to use for CIVO API.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			// "civo_template":           dataSourceTemplate(),
			"civo_disk_image":              dataSourceDiskImage(),
			"civo_kubernetes_version":      dataSourceKubernetesVersion(),
			"civo_kubernetes_cluster":      dataSourceKubernetesCluster(),
			"civo_size":                    dataSourceSize(),
			"civo_instances":               dataSourceInstances(),
			"civo_instance":                dataSourceInstance(),
			"civo_dns_domain_name":         dataSourceDNSDomainName(),
			"civo_dns_domain_record":       dataSourceDNSDomainRecord(),
			"civo_network":                 dataSourceNetwork(),
			"civo_volume":                  dataSourceVolume(),
			"civo_firewall":                dataSourceFirewall(),
			"civo_loadbalancer":            dataSourceLoadBalancer(),
			"civo_ssh_key":                 dataSourceSSHKey(),
			"civo_object_store":            dataSourceObjectStore(),
			"civo_object_store_credential": dataSourceObjectStoreCredential(),
			"civo_region":                  dataSourceRegion(),
			"civo_reserved_ip":             dataSourceReservedIP(),
			"civo_database":                dataSourceDatabase(),
			"civo_database_version":        dataDatabaseVersion(),
			// "civo_snapshot":           dataSourceSnapshot(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"civo_instance":                        resourceInstance(),
			"civo_network":                         resourceNetwork(),
			"civo_volume":                          resourceVolume(),
			"civo_volume_attachment":               resourceVolumeAttachment(),
			"civo_dns_domain_name":                 resourceDNSDomainName(),
			"civo_dns_domain_record":               resourceDNSDomainRecord(),
			"civo_firewall":                        resourceFirewall(),
			"civo_firewall_rule":                   resourceFirewallRule(),
			"civo_ssh_key":                         resourceSSHKey(),
			"civo_kubernetes_cluster":              resourceKubernetesCluster(),
			"civo_kubernetes_node_pool":            resourceKubernetesClusterNodePool(),
			"civo_reserved_ip":                     resourceReservedIP(),
			"civo_instance_reserved_ip_assignment": resourceInstanceReservedIPAssignment(),
			"civo_object_store":                    resourceObjectStore(),
			"civo_object_store_credential":         resourceObjectStoreCredential(),
			"civo_database":                        resourceDatabase(),
			// "civo_loadbalancer":         resourceLoadBalancer(),
			// "civo_template": resourceTemplate(),
			// "civo_snapshot":             resourceSnapshot(),

		},
		ConfigureFunc: providerConfigure,
	}
}

// Provider configuration
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var regionValue, tokenValue, apiURL string
	var client *civogo.Client
	var err error

	if region, ok := d.GetOk("region"); ok {
		regionValue = region.(string)
	}

	if token, ok := d.GetOk("token"); ok {
		tokenValue = token.(string)
	} else {
		return nil, fmt.Errorf("[ERR] token not found")
	}

	if apiEndpoint, ok := d.GetOk("api_endpoint"); ok {
		apiURL = apiEndpoint.(string)
	} else {
		apiURL = ProdAPI
	}
	client, err = civogo.NewClientWithURL(tokenValue, apiURL, regionValue)
	if err != nil {
		return nil, err
	}

	userAgent := &civogo.Component{
		Name:    "terraform-provider-civo",
		Version: ProviderVersion,
	}
	client.SetUserAgent(userAgent)

	log.Printf("[DEBUG] Civo API URL: %s\n", apiURL)
	return client, nil
}
