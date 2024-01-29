package civo

import (
	"fmt"
	"log"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/civo/database"
	"github.com/civo/terraform-provider-civo/civo/disk"
	"github.com/civo/terraform-provider-civo/civo/dns"
	"github.com/civo/terraform-provider-civo/civo/firewall"
	"github.com/civo/terraform-provider-civo/civo/instances"
	"github.com/civo/terraform-provider-civo/civo/ip"
	"github.com/civo/terraform-provider-civo/civo/kubernetes"
	"github.com/civo/terraform-provider-civo/civo/loadbalancer"
	"github.com/civo/terraform-provider-civo/civo/network"
	"github.com/civo/terraform-provider-civo/civo/objectstorage"
	"github.com/civo/terraform-provider-civo/civo/region"
	"github.com/civo/terraform-provider-civo/civo/size"
	"github.com/civo/terraform-provider-civo/civo/ssh"
	"github.com/civo/terraform-provider-civo/civo/volumen"
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
			"civo_disk_image":              disk.DataSourceDiskImage(),
			"civo_kubernetes_version":      kubernetes.DataSourceKubernetesVersion(),
			"civo_kubernetes_cluster":      kubernetes.DataSourceKubernetesCluster(),
			"civo_size":                    size.DataSourceSize(),
			"civo_instances":               instances.DataSourceInstances(),
			"civo_instance":                instances.DataSourceInstance(),
			"civo_dns_domain_name":         dns.DataSourceDNSDomainName(),
			"civo_dns_domain_record":       dns.DataSourceDNSDomainRecord(),
			"civo_network":                 network.DataSourceNetwork(),
			"civo_volume":                  volume.DataSourceVolume(),
			"civo_firewall":                firewall.DataSourceFirewall(),
			"civo_loadbalancer":            loadbalancer.DataSourceLoadBalancer(),
			"civo_ssh_key":                 ssh.DataSourceSSHKey(),
			"civo_object_store":            objectstorage.DataSourceObjectStore(),
			"civo_object_store_credential": objectstorage.DataSourceObjectStoreCredential(),
			"civo_region":                  region.DataSourceRegion(),
			"civo_reserved_ip":             ip.DataSourceReservedIP(),
			"civo_database":                database.DataSourceDatabase(),
			"civo_database_version":        database.DataDatabaseVersion(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"civo_instance":                        instances.ResourceInstance(),
			"civo_instance_reserved_ip_assignment": instances.ResourceInstanceReservedIPAssignment(),
			"civo_network":                         network.ResourceNetwork(),
			"civo_volume":                          volume.ResourceVolume(),
			"civo_volume_attachment":               volume.ResourceVolumeAttachment(),
			"civo_dns_domain_name":                 dns.ResourceDNSDomainName(),
			"civo_dns_domain_record":               dns.ResourceDNSDomainRecord(),
			"civo_firewall":                        firewall.ResourceFirewall(),
			"civo_ssh_key":                         ssh.ResourceSSHKey(),
			"civo_kubernetes_cluster":              kubernetes.ResourceKubernetesCluster(),
			"civo_kubernetes_node_pool":            kubernetes.ResourceKubernetesClusterNodePool(),
			"civo_reserved_ip":                     ip.ResourceReservedIP(),
			"civo_object_store":                    objectstorage.ResourceObjectStore(),
			"civo_object_store_credential":         objectstorage.ResourceObjectStoreCredential(),
			"civo_database":                        database.ResourceDatabase(),
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
