package civo

import (
	"fmt"
	"log"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// This resource represent a load balancer in the system
func resourceLoadBalancer() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "the hostname to receive traffic for, e.g. www.example.com (optional: sets hostname to loadbalancer-uuid.civo.com if blank)",
				ValidateFunc: utils.ValidateName,
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "either http or https. If you specify https then you must also provide the next two fields, the default is http",
				ValidateFunc: validation.StringInSlice([]string{
					"http",
					"https",
				}, false),
			},
			"tls_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "if your protocol is https then you should send the TLS certificate in Base64-encoded PEM format",
			},
			"tls_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "if your protocol is https then you should send the TLS private key in Base64-encoded PEM format",
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "you can listen on any port, the default is 80 to match the default protocol of http," +
					"if not you must specify it here (commonly 80 for HTTP or 443 for HTTPS)",
				ValidateFunc: validation.NoZeroValues,
			},
			"max_request_size": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "the size in megabytes of the maximum request content that will be accepted, defaults to 20",
				ValidateFunc: validation.IntAtLeast(20),
			},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
				Description: "one of: least_conn (sends new requests to the least busy server), " +
					"random (sends new requests to a random backend), round_robin (sends new requests to the next backend in order), " +
					"ip_hash (sends requests from a given IP address to the same backend), default is random",
				ValidateFunc: validation.StringInSlice([]string{
					"least_conn",
					"random",
					"round_robin",
					"ip_hash",
				}, false),
			},
			"health_check_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "what URL should be used on the backends to determine if it's OK (2xx/3xx status), defaults to /",
			},
			"fail_timeout": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "how long to wait in seconds before determining a backend has failed, defaults to 30",
				ValidateFunc: validation.IntAtLeast(30),
			},
			"max_conns": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "how many concurrent connections can each backend handle, defaults to 10",
				ValidateFunc: validation.IntAtLeast(10),
			},
			"ignore_invalid_backend_tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "should self-signed/invalid certificates be ignored from the backend servers, defaults to true",
			},
			"backend": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "a list of backend instances, each containing an instance_id, protocol (http or https) and port",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: utils.ValidateName,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http",
								"https",
							}, false),
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
					},
				},
			},
		},
		Create: resourceLoadBalancerCreate,
		Read:   resourceLoadBalancerRead,
		Update: resourceLoadBalancerUpdate,
		Delete: resourceLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// function to create a new load balancer
func resourceLoadBalancerCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] configuring the load balancer %s", d.Get("hostname").(string))
	conf := &civogo.LoadBalancerConfig{
		Hostname:       d.Get("hostname").(string),
		Protocol:       d.Get("protocol").(string),
		Port:           d.Get("port").(int),
		MaxRequestSize: d.Get("max_request_size").(int),
		Policy:         d.Get("policy").(string),
		MaxConns:       d.Get("max_conns").(int),
		Backends:       expandLoadBalancerBackend(d.Get("backend").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("tls_certificate"); ok {
		conf.TLSCertificate = v.(string)
	}

	if v, ok := d.GetOk("tls_key"); ok {
		conf.TLSKey = v.(string)
	}

	if v, ok := d.GetOk("health_check_path"); ok {
		conf.HealthCheckPath = v.(string)
	}

	if v, ok := d.GetOk("fail_timeout"); ok {
		conf.FailTimeout = v.(int)
	}

	if v, ok := d.GetOk("ignore_invalid_backend_tls"); ok {
		conf.IgnoreInvalidBackendTLS = v.(bool)
	}

	log.Printf("[INFO] creating the load balancer %s", d.Get("hostname").(string))
	lb, err := apiClient.CreateLoadBalancer(conf)
	if err != nil {
		return fmt.Errorf("[ERR] failed to create a new load balancer: %s", err)
	}

	d.SetId(lb.ID)

	return resourceLoadBalancerRead(d, m)
}

// function to read the load balancer
func resourceLoadBalancerRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] retrieving the load balancer %s", d.Id())
	resp, err := apiClient.FindLoadBalancer(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERR] error retrieving load balancer: %s", err)
	}

	d.Set("hostname", resp.Hostname)
	d.Set("protocol", resp.Protocol)
	d.Set("tls_certificate", resp.TLSCertificate)
	d.Set("tls_key", resp.TLSKey)
	d.Set("port", resp.Port)
	d.Set("max_request_size", resp.MaxRequestSize)
	d.Set("policy", resp.Policy)
	d.Set("health_check_path", resp.HealthCheckPath)
	d.Set("fail_timeout", resp.FailTimeout)
	d.Set("max_conns", resp.MaxConns)
	d.Set("ignore_invalid_backend_tls", resp.IgnoreInvalidBackendTLS)

	if err := d.Set("backend", flattenLoadBalancerBackend(resp.Backends)); err != nil {
		return fmt.Errorf("[ERR] error retrieving the backend for load balancer error: %#v", err)
	}

	return nil
}

// function to update the load balancer
func resourceLoadBalancerUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] configuring the load balancer to update %s", d.Id())
	conf := &civogo.LoadBalancerConfig{
		Hostname:       d.Get("hostname").(string),
		Protocol:       d.Get("protocol").(string),
		Port:           d.Get("port").(int),
		MaxRequestSize: d.Get("max_request_size").(int),
		Policy:         d.Get("policy").(string),
		MaxConns:       d.Get("max_conns").(int),
		FailTimeout:    d.Get("fail_timeout").(int),
		Backends:       expandLoadBalancerBackend(d.Get("backend").(*schema.Set).List()),
	}

	if d.HasChange("tls_certificate") {
		conf.TLSCertificate = d.Get("tls_certificate").(string)
	}

	if d.HasChange("tls_key") {
		conf.TLSKey = d.Get("tls_key").(string)
	}

	if d.HasChange("health_check_path") {
		conf.HealthCheckPath = d.Get("health_check_path").(string)
	}

	if d.HasChange("ignore_invalid_backend_tls") {
		conf.IgnoreInvalidBackendTLS = d.Get("ignore_invalid_backend_tls").(bool)
	}

	log.Printf("[INFO] updating the load balancer %s", d.Id())
	_, err := apiClient.UpdateLoadBalancer(d.Id(), conf)
	if err != nil {
		return fmt.Errorf("[ERR] failed to update load balancer: %s", err)
	}

	return resourceLoadBalancerRead(d, m)

}

// function to delete the load balancer
func resourceLoadBalancerDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	log.Printf("[INFO] deleting the load balancer %s", d.Id())
	_, err := apiClient.DeleteLoadBalancer(d.Id())
	if err != nil {
		return fmt.Errorf("[ERR] an error occurred while tring to delete load balancer %s", d.Id())
	}
	return nil
}

// function to expand the load balancer backend to send to the api
func expandLoadBalancerBackend(backend []interface{}) []civogo.LoadBalancerBackendConfig {
	expandedBackend := make([]civogo.LoadBalancerBackendConfig, 0, len(backend))
	for _, rawBackend := range backend {

		rule := rawBackend.(map[string]interface{})

		r := civogo.LoadBalancerBackendConfig{
			Protocol:   rule["protocol"].(string),
			InstanceID: rule["instance_id"].(string),
			Port:       rule["port"].(int),
		}

		expandedBackend = append(expandedBackend, r)
	}
	return expandedBackend
}

// function to flatten the load balancer backend when is coming from the api
func flattenLoadBalancerBackend(backend []civogo.LoadBalancerBackend) []interface{} {
	if backend == nil {
		return nil
	}

	flattenedBackend := make([]interface{}, len(backend))
	for i, back := range backend {
		instanceID := back.InstanceID
		protocol := back.Protocol
		port := back.Port

		rawRule := map[string]interface{}{
			"instance_id": instanceID,
			"protocol":    protocol,
			"port":        port,
		}

		flattenedBackend[i] = rawRule
	}

	return flattenedBackend
}
