package loadbalancer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"net"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceLoadBalancer represent a load balancer in the system
func ResourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Civo load balancer resource. This can be used to create, modify, and delete load balancers.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the load balancer.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"service_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Name of the service associated with the load balancer.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"network_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Network ID associated with the load balancer.",
				ValidateFunc: utils.ValidateUUID,
			},
			"algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "round_robin",
				Description:  "Load balancing algorithm, either 'round_robin' or 'least_connections'.",
				ValidateFunc: validation.StringInSlice([]string{"round_robin", "least_connections"}, false),
			},
			"external_traffic_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "External traffic policy, either 'Cluster' or 'Local'.",
				ValidateFunc: validation.StringInSlice([]string{"Cluster", "Local"}, false),
			},
			"session_affinity": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Session affinity setting, either 'ClientIP' or 'None'.",
				ValidateFunc: validation.StringInSlice([]string{"ClientIP", "None"}, false),
			},
			"session_affinity_config_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Timeout for session affinity in seconds.",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"enable_proxy_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Enable proxy protocol, options are '', 'send-proxy', 'send-proxy-v2'.",
				ValidateFunc: validation.StringInSlice([]string{"", "send-proxy", "send-proxy-v2"}, false),
			},
			"max_concurrent_requests": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "Maximum concurrent requests the load balancer can handle.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"cluster_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "ID of the cluster associated with the load balancer.",
				ValidateFunc: utils.ValidateUUID,
			},
			"firewall_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "ID of the firewall associated with the load balancer.",
				ValidateFunc: utils.ValidateUUID,
			},
			"firewall_rule": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Firewall rules for the load balancer (e.g., 'all', '80,443', '40-80,90-120').",
			},
			"backend": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of backend servers to be load balanced.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "IP address of the backend server.",
							ValidateFunc: validateIPAddress,
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "TCP",
							Description:  "Protocol for backend server communication (TCP or UDP).",
							ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
						},
						"source_port": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Source port for backend server.",
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"target_port": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Target port for backend server.",
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"health_check_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Port used for health checks on the backend server.",
							ValidateFunc: validation.IntBetween(1, 65535),
						},
					},
				},
			},
			"instance_pool": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of instance pools for the load balancer.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of tags for instances in the pool.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"names": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of instance names in the pool.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "TCP",
							Description:  "Protocol for instance pool communication (TCP or UDP).",
							ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
						},
						"source_port": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Source port for instance pool.",
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"target_port": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Target port for instance pool.",
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"health_check": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Health check configuration for instance pool.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:         schema.TypeInt,
										Required:     true,
										Description:  "Port used for health checks.",
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Path used for HTTP health checks.",
									},
								},
							},
						},
					},
				},
			},
		},
		CreateContext: resourceLoadBalancerCreate,
		ReadContext:   resourceLoadBalancerRead,
		UpdateContext: resourceLoadBalancerUpdate,
		DeleteContext: resourceLoadBalancerDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customizeDiffLoadbalancer,
	}
}

// function to create a new load balancer
func resourceLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*civogo.Client)

	// overwrite the region if is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		client.Region = region.(string)
	}

	name := d.Get("name").(string)

	// Check that either backend or instance_pool is provided
	backends, backendSet := d.GetOk("backend")
	instancePools, instancePoolSet := d.GetOk("instance_pool")

	if !backendSet && !instancePoolSet {
		return diag.Errorf("one of 'backend' or 'instance_pool' must be specified")
	}

	// Prepare the load balancer create request
	conf := &civogo.LoadBalancerConfig{
		Name:   name,
		Region: client.Region,
	}

	if v, ok := d.GetOk("service_name"); ok {
		conf.ServiceName = v.(string)
	}

	if v, ok := d.GetOk("firewall_id"); ok {
		conf.FirewallID = v.(string)
	}

	if v, ok := d.GetOk("network_id"); ok {
		conf.NetworkID = v.(string)
	}

	if v, ok := d.GetOk("algorithm"); ok {
		conf.Algorithm = v.(string)
	}

	if v, ok := d.GetOk("external_traffic_policy"); ok {
		conf.ExternalTrafficPolicy = v.(string)
	}

	if v, ok := d.GetOk("session_affinity"); ok {
		conf.SessionAffinity = v.(string)
	}

	if v, ok := d.GetOk("session_affinity_config_timeout"); ok {
		conf.SessionAffinityConfigTimeout = int32(v.(int))
	}

	if v, ok := d.GetOk("enable_proxy_protocol"); ok {
		conf.EnableProxyProtocol = v.(string)
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		conf.ClusterID = v.(string)
	}

	if v, ok := d.GetOk("max_concurrent_requests"); ok {
		num := v.(int)
		conf.MaxConcurrentRequests = &num
	}

	// Set backend configurations if provided
	if backendSet {
		for _, backend := range backends.([]interface{}) {
			b := backend.(map[string]interface{})
			conf.Backends = append(conf.Backends, civogo.LoadBalancerBackendConfig{
				IP:              b["ip"].(string),
				Protocol:        b["protocol"].(string),
				SourcePort:      int32(b["source_port"].(int)),
				TargetPort:      int32(b["target_port"].(int)),
				HealthCheckPort: int32(b["health_check_port"].(int)),
			})
		}
	}

	// Set instance pool configurations if provided
	if instancePoolSet {
		for _, instancePool := range instancePools.([]interface{}) {
			ip := instancePool.(map[string]interface{})

			pool := civogo.LoadBalancerInstancePoolConfig{
				Tags:       convertStringList(ip["tags"].([]interface{})),
				Names:      convertStringList(ip["names"].([]interface{})),
				Protocol:   ip["protocol"].(string),
				SourcePort: int32(ip["source_port"].(int)),
				TargetPort: int32(ip["target_port"].(int)),
				HealthCheck: civogo.HealthCheck{
					Port: int32(ip["health_check"].([]interface{})[0].(map[string]interface{})["port"].(int)),
					Path: ip["health_check"].([]interface{})[0].(map[string]interface{})["path"].(string),
				},
			}
			conf.InstancePools = append(conf.InstancePools, pool)
		}
	}

	// Send create request to API
	loadBalancer, err := client.CreateLoadBalancer(conf)
	if err != nil {
		return diag.Errorf("error creating load balancer: %s", err)
	}

	// Set the ID for the Terraform state
	d.SetId(loadBalancer.ID)

	return resourceLoadBalancerRead(ctx, d, m)
}

// function to read the load balancer
func resourceLoadBalancerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*civogo.Client)

	// overwrite the region if is defined in the datasource
	if region, ok := d.GetOk("region"); ok {
		client.Region = region.(string)
	}

	// Retrieve the load balancer information from the API
	loadBalancer, err := client.GetLoadBalancer(d.Id())
	if err != nil {
		return diag.Errorf("error fetching load balancer: %s", err)
	}

	// Set basic load balancer attributes
	d.Set("name", loadBalancer.Name)
	d.Set("service_name", loadBalancer.ServiceName)
	d.Set("network_id", loadBalancer.NetworkID)
	d.Set("algorithm", loadBalancer.Algorithm)
	d.Set("external_traffic_policy", loadBalancer.ExternalTrafficPolicy)
	d.Set("session_affinity", loadBalancer.SessionAffinity)
	d.Set("session_affinity_config_timeout", loadBalancer.SessionAffinityConfigTimeout)
	d.Set("enable_proxy_protocol", loadBalancer.EnableProxyProtocol)
	d.Set("max_concurrent_requests", loadBalancer.MaxConcurrentRequests)
	d.Set("cluster_id", loadBalancer.ClusterID)
	d.Set("firewall_id", loadBalancer.FirewallID)

	if err := d.Set("backend", flattenLoadBalancerBackend(loadBalancer.Backends)); err != nil {
		return diag.Errorf("error setting backend: %s", err)
	}

	if loadBalancer.InstancePool != nil {
		flattenedInstancePool := flattenInstancePool(loadBalancer.InstancePool)
		if err := d.Set("instance_pool", flattenedInstancePool); err != nil {
			return diag.Errorf("error setting instance_pool: %s", err)
		}
	}

	return nil
}

// function to update the load balancer
func resourceLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*civogo.Client)
	loadBalancerID := d.Id()

	// Initialize an update request
	updateRequest := &civogo.LoadBalancerUpdateConfig{
		Region: client.Region,
		Name:   d.Get("name").(string),
	}

	// Check if any relevant fields have changed and update the request accordingly
	if d.HasChange("service_name") {
		updateRequest.ServiceName = d.Get("service_name").(string)
	}

	if d.HasChange("firewall_id") {
		updateRequest.FirewallID = d.Get("firewall_id").(string)
	}

	if d.HasChange("algorithm") {
		updateRequest.Algorithm = d.Get("algorithm").(string)
	}

	if d.HasChange("external_traffic_policy") {
		updateRequest.ExternalTrafficPolicy = d.Get("external_traffic_policy").(string)
	}

	if d.HasChange("session_affinity") {
		updateRequest.SessionAffinity = d.Get("session_affinity").(string)
	}

	if d.HasChange("session_affinity_config_timeout") {
		updateRequest.SessionAffinityConfigTimeout = int32(d.Get("session_affinity_config_timeout").(int))
	}

	if d.HasChange("enable_proxy_protocol") {
		updateRequest.EnableProxyProtocol = d.Get("enable_proxy_protocol").(string)
	}

	// Send the update request to the Civo API
	_, err := client.UpdateLoadBalancer(loadBalancerID, updateRequest)
	if err != nil {
		return diag.Errorf("error updating load balancer: %s", err)
	}

	return resourceLoadBalancerRead(ctx, d, m)
}

// function to delete the load balancer
func resourceLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] deleting the load balancer %s", d.Id())
	_, err := apiClient.DeleteLoadBalancer(d.Id())
	if err != nil {
		return diag.Errorf("[ERR] an error occurred while trying to delete load balancer %s", d.Id())
	}
	return nil
}

// flattenLoadBalancerBackend converts a slice of LoadBalancerBackend structs into a slice of maps for Terraform state.
func flattenLoadBalancerBackend(backends []civogo.LoadBalancerBackend) []interface{} {
	// Return nil if there are no backends
	if len(backends) == 0 {
		return nil
	}

	// Create a slice to store each backend configuration as a map
	flattenedBackends := make([]interface{}, len(backends))
	for i, backend := range backends {
		// Convert each backend into a map format
		backendMap := map[string]interface{}{
			"ip":                backend.IP,
			"protocol":          backend.Protocol,
			"source_port":       backend.SourcePort,
			"target_port":       backend.TargetPort,
			"health_check_port": backend.HealthCheckPort,
		}
		flattenedBackends[i] = backendMap
	}

	return flattenedBackends
}

// flattenInstancePool converts a slice of InstancePool structs into a slice of maps for Terraform state.
func flattenInstancePool(instancePools []civogo.InstancePool) []interface{} {
	// Return nil if there are no instance pools
	if len(instancePools) == 0 {
		return nil
	}

	// Create a slice to store each instance pool configuration as a map
	flattenedInstancePools := make([]interface{}, len(instancePools))
	for i, pool := range instancePools {
		// Convert HealthCheck struct into a map
		healthCheckMap := map[string]interface{}{
			"port": pool.HealthCheck.Port,
			"path": pool.HealthCheck.Path,
		}

		// Convert each instance pool into a map format
		poolMap := map[string]interface{}{
			"tags":         pool.Tags,
			"names":        pool.Names,
			"protocol":     pool.Protocol,
			"source_port":  pool.SourcePort,
			"target_port":  pool.TargetPort,
			"health_check": []interface{}{healthCheckMap}, // Wrap in a slice to match Terraform's nested object structure
		}
		flattenedInstancePools[i] = poolMap
	}

	return flattenedInstancePools
}

// Helper function for IP address validation
func validateIPAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if ip := net.ParseIP(value); ip == nil {
		errors = append(errors, fmt.Errorf("%q must be a valid IP address", k))
	}
	return
}

// convertStringList converts a list of interface{} to a slice of strings.
func convertStringList(input []interface{}) []string {
	strList := make([]string, len(input))
	for i, v := range input {
		strList[i] = v.(string)
	}
	return strList
}

func customizeDiffLoadbalancer(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if d.Id() != "" && d.HasChange("instance_pool") {
		return fmt.Errorf("the 'instance_pool' field is immutable")
	}
	if d.Id() != "" && d.HasChange("backend") {
		return fmt.Errorf("the 'backend' field is immutable")
	}
	return nil
}
