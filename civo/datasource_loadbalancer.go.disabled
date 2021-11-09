package civo

import (
	"fmt"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Data source to get from the api a specific Load Balancer
// using the id or the hostname of the Load Balancer
func dataSourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Get information on a load balancer for use in other resources. This data source provides all of the load balancers properties as configured on your Civo account.",
			"An error will be raised if the provided load balancer name does not exist in your Civo account.",
		}, "\n\n"),
		Read: dataSourceLoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "hostname"},
			},
			"hostname": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "hostname"},
				Description:  "The hostname of the load balancer",
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The region where load balancer is running",
			},
			// Computed resource
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol used",
			},
			"tls_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If is set will be returned",
			},
			"tls_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If is set will be returned",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port set in the configuration",
			},
			"max_request_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The max request size set in the configuration",
			},
			"policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy set in the load balancer",
			},
			"health_check_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path to check the health of the backend",
			},
			"fail_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The wait time until the backend is marked as a failure",
			},
			"max_conns": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "How many concurrent connections can each backend handle",
			},
			"ignore_invalid_backend_tls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Should self-signed/invalid certificates be ignored from the backend servers",
			},
			"backend": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance ID",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol used in the configuration",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port set in the configuration",
						},
					},
				},
			},
		},
	}
}

func dataSourceLoadBalancerRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	var searchBy string

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the LoadBalancer by id")
		searchBy = id.(string)
	} else if hostname, ok := d.GetOk("hostname"); ok {
		log.Printf("[INFO] Getting the LoadBalancer by hostname")
		searchBy = hostname.(string)
	}

	lb, err := apiClient.FindLoadBalancer(searchBy)
	if err != nil {
		return fmt.Errorf("[ERR] failed to retrive LoadBalancer: %s", err)
	}

	d.SetId(lb.ID)
	d.Set("hostname", lb.Hostname)
	d.Set("protocol", lb.Protocol)
	d.Set("tls_certificate", lb.TLSCertificate)
	d.Set("tls_key", lb.TLSKey)
	d.Set("port", lb.Port)
	d.Set("max_request_size", lb.MaxRequestSize)
	d.Set("policy", lb.Policy)
	d.Set("health_check_path", lb.HealthCheckPath)
	d.Set("fail_timeout", lb.FailTimeout)
	d.Set("max_conns", lb.MaxConns)
	d.Set("ignore_invalid_backend_tls", lb.IgnoreInvalidBackendTLS)

	if err := d.Set("backend", flattenLoadBalancerBackend(lb.Backends)); err != nil {
		return fmt.Errorf("[ERR] error retrieving the backend for load balancer error: %#v", err)
	}

	return nil
}
