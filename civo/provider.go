package civo

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	"github.com/civo/terraform-provider-civo/civo/volume"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				Type:             schema.TypeString,
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc("CIVO_TOKEN", ""),
				Description:      "This is the Civo API token. Alternatively, this can also be specified using `CIVO_TOKEN` environment variable.",
				Deprecated:       "",
				ValidateDiagFunc: validateTokenUsage,
			},
			"credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CIVO_CREDENTIAL_FILE", ""),
				Description: "Path to the Civo credentials file. Can be specified using CIVO_CREDENTIAL_FILE environment variable.",
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

	if token, ok := getToken(d); ok {
		tokenValue = token.(string)
	} else {
		return nil, fmt.Errorf("[ERR] No token configuration found in $CIVO_TOKEN or ~/.civo.json. Please go to https://dashboard.civo.com/security to fetch one")
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

func getToken(d *schema.ResourceData) (interface{}, bool) {
	var exists = true

	// Gets you the token atrribute value or falls back to reading CIVO_TOKEN environment variable
	if token, ok := d.GetOk("token"); ok {
		return token, exists
	}

	// Check for credentials file specified in provider config
	if credFile, ok := d.GetOk("credentials_file"); ok {
		token, err := readTokenFromFile(credFile.(string))
		if err == nil {
			return token, exists
		}
	}

	// Check for default CLI config file
	homeDir, err := getHomeDir()
	if err == nil {
		token, err := readTokenFromFile(filepath.Join(homeDir, ".civo.json"))
		if err == nil {
			return token, exists
		}
	}

	return nil, !exists

}

var getHomeDir = func() (string, error) {
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}
	// Fall back to os.UserHomeDir() if HOME is not set
	return os.UserHomeDir()
}

func readTokenFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return "", err
	}

	var config struct {
		APIKeys map[string]string `json:"apikeys"`
		Meta    struct {
			CurrentAPIKey string `json:"current_apikey"`
		} `json:"meta"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return "", err
	}

	// Get the current API key name
	currentKeyName := config.Meta.CurrentAPIKey

	// Fetch the corresponding token
	token, ok := config.APIKeys[currentKeyName]

	if !ok {
		return "", fmt.Errorf("API key '%s' not found", currentKeyName)
	}

	return token, nil
}

func validateTokenUsage(v interface{}, path cty.Path) diag.Diagnostics {
	val := v.(string)

	// Ensures warning is not shown when "CIVO_TOKEN" environment variable is set.
	if token := os.Getenv("CIVO_TOKEN"); token != "" {
		val = ""
	}
	var diags diag.Diagnostics

	if val != "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deprecated Attribute Usage",
			Detail:   "The \"token\" attribute is deprecated. Please use the CIVO_TOKEN environment variable or the credentials_file attribute instead.",
		})
	}

	return diags
}
