package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

// Data source to get from the api a specific Load Balancer
// using the id or the hostname of the Load Balancer
func dataSourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
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
			},
			// Computed resource
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tls_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tls_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_request_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fail_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_conns": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ignore_invalid_backend_tls": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"backend": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLoadBalancerRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

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
