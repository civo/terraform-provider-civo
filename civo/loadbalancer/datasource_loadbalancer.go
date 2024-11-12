package loadbalancer

import (
	"context"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceLoadBalancer function returns a schema.Resource that represents a Load Balancer.
// This can be used to query and retrieve details about a specific Load Balancer in the infrastructure using its id or hostname.
func DataSourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Get information on a load balancer for use in other resources. This data source provides all of the load balancers properties as configured on your Civo account.",
			"An error will be raised if the provided load balancer name does not exist in your Civo account.",
		}, "\n\n"),
		ReadContext: dataSourceLoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the load balancer to retrieve (You can find this id from service annotations 'kubernetes.civo.com/loadbalancer-id')",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the load balancer (You can find this name from service annotations 'kubernetes.civo.com/loadbalancer-name')",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region of the load balancer, if you declare this field, the datasource will use this value instead of the one defined in the provider",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public ip of the load balancer",
			},
			"algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The algorithm used by the load balancer",
			},
			"external_traffic_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The external traffic policy of the load balancer",
			},
			"session_affinity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The session affinity of the load balancer",
			},
			"session_affinity_config_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The session affinity config timeout of the load balancer",
			},
			"enable_proxy_protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The enabled proxy protocol of the load balancer",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private ip of the load balancer",
			},
			"firewall_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The firewall id of the load balancer",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cluster id of the load balancer",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the load balancer",
			},
			"backends": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ip of the backend",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol of the backend",
						},
						"source_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The source port of the backend",
						},
						"target_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The target port of the backend",
						},
						"health_check_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The health check port of the backend",
						},
					},
				},
			},
		},
	}
}

func dataSourceLoadBalancerRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	var searchBy string

	if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the LoadBalancer by name")
		searchBy = name.(string)
	} else if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the LoadBalancer by id")
		searchBy = id.(string)
	}

	lb, err := apiClient.FindLoadBalancer(searchBy)
	if err != nil {
		return diag.Errorf("[ERR] failed to retrive LoadBalancer: %s", err)
	}

	d.SetId(lb.ID)
	d.Set("name", lb.Name)
	d.Set("public_ip", lb.PublicIP)
	d.Set("algorithm", lb.Algorithm)
	d.Set("external_traffic_policy", lb.ExternalTrafficPolicy)
	d.Set("session_affinity", lb.SessionAffinity)
	d.Set("session_affinity_config_timeout", lb.SessionAffinityConfigTimeout)
	d.Set("enable_proxy_protocol", lb.EnableProxyProtocol)
	d.Set("private_ip", lb.PrivateIP)
	d.Set("firewall_id", lb.FirewallID)
	d.Set("cluster_id", lb.ClusterID)
	d.Set("state", lb.State)

	if err := d.Set("backends", flattenLoadBalancerBackend(lb.Backends)); err != nil {
		return diag.Errorf("[ERR] error retrieving the backends for load balancer error: %#v", err)
	}

	return nil
}
