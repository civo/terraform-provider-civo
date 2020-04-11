package civo

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "a name for your cluster, must be unique within your account (required)",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"num_target_nodes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "the number of instances to create (optional, the default at the time of writing is 3)",
			},
			"target_nodes_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "g2.small",
				Description: "the size of each node (optional, the default is currently g2.small)",
			},
			"kubernetes_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1.0.0",
				Description: "the version of k3s to install (optional, the default is currently the latest available)",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "a space separated list of tags, to be used freely as required (optional)",
			},
			"applications": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "a comma separated list of applications to install." +
					"Spaces within application names are fine, but shouldn't be either side of the comma." +
					"If you want to remove a default installed application, prefix it with a '-', e.g. -traefik.",
			},
			// Computed resource
			"instances":              instanceSchema(),
			"installed_applications": applicationSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ready": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"kubeconfig": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_entry": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"built_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create: resourceKubernetesClusterCreate,
		Read:   resourceKubernetesClusterRead,
		Update: resourceKubernetesClusterUpdate,
		Delete: resourceKubernetesClusterDelete,
	}
}

func instanceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hostname": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"size": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"region": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"created_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"firewall_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"public_ip": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"tags": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func applicationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"application": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"version": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"installed": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"category": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func resourceKubernetesClusterCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	config := &civogo.KubernetesClusterConfig{
		Name:            d.Get("name").(string),
		TargetNodesSize: d.Get("target_nodes_size").(string),
	}

	if attr, ok := d.GetOk("num_target_nodes"); ok {
		config.NumTargetNodes = attr.(int)
	}

	if attr, ok := d.GetOk("kubernetes_version"); ok {
		config.KubernetesVersion = attr.(string)
	}

	if attr, ok := d.GetOk("tags"); ok {
		config.Tags = attr.(string)
	}

	if attr, ok := d.GetOk("applications"); ok {
		config.Applications = attr.(string)
	}

	resp, err := apiClient.NewKubernetesClusters(config)
	if err != nil {
		return fmt.Errorf("[WARN] failed to create the kubernets cluster: %s", err)
	}

	d.SetId(resp.ID)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := apiClient.FindKubernetesCluster(d.Id())
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[WARN] error geting kubernetes cluster: %s", err))
		}

		if resp.Ready != true {
			return resource.RetryableError(fmt.Errorf("[WARN] waiting for the kubernets cluster to be created but the status is %s", resp.Status))
		}

		return resource.NonRetryableError(resourceKubernetesClusterRead(d, m))
	})
}

func resourceKubernetesClusterRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	resp, err := apiClient.FindKubernetesCluster(d.Id())
	if err != nil {
		return fmt.Errorf("[WARN] failed to find the kubernets cluster: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("num_target_nodes", resp.NumTargetNode)
	d.Set("target_nodes_size", resp.TargetNodeSize)
	d.Set("kubernetes_version", resp.KubernetesVersion)
	d.Set("tags", resp.Tags)
	d.Set("status", resp.Status)
	d.Set("ready", resp.Ready)
	d.Set("kubeconfig", resp.KubeConfig)
	d.Set("api_endpoint", resp.APIEndPoint)
	d.Set("dns_entry", resp.DNSEntry)
	d.Set("built_at", resp.BuiltAt.UTC().String())
	d.Set("created_at", resp.CreatedAt.UTC().String())

	if err := d.Set("instances", flattenInstances(resp.Instances)); err != nil {
		return fmt.Errorf("[WARN] error retrieving the instances for kubernetes cluster error: %#v", err)
	}

	if err := d.Set("installed_applications", flattenInstalledApplication(resp.InstalledApplications)); err != nil {
		return fmt.Errorf("[WARN] error retrieving the installed application for kubernetes cluster error: %#v", err)
	}

	return nil
}

func resourceKubernetesClusterUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	config := &civogo.KubernetesClusterConfig{}

	if d.HasChange("num_target_nodes") || d.HasChange("kubernetes_version") || d.HasChange("applications") || d.HasChange("name") {
		config.Name = d.Get("name").(string)
		config.NumTargetNodes = d.Get("num_target_nodes").(int)
		config.KubernetesVersion = d.Get("kubernetes_version").(string)
		config.Applications = d.Get("applications").(string)
	}

	_, err := apiClient.UpdateKubernetesCluster(d.Id(), config)
	if err != nil {
		return fmt.Errorf("[WARN] failed to update kubernetes cluster: %s", err)
	}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := apiClient.FindKubernetesCluster(d.Id())
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[WARN] error geting kubernetes cluster: %s", err))
		}

		if resp.Status != "ACTIVE" {
			return resource.RetryableError(fmt.Errorf("[WARN] waiting for the kubernets cluster to be created but the status is %s", resp.Status))
		}

		return resource.NonRetryableError(resourceKubernetesClusterRead(d, m))
	})
}

func resourceKubernetesClusterDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	_, err := apiClient.DeleteKubernetesCluster(d.Id())
	if err != nil {
		return fmt.Errorf("[INFO] kubernets cluster (%s) was delete", d.Id())
	}

	return nil
}

func flattenInstances(instances []civogo.KubernetesInstance) []interface{} {
	if instances == nil {
		return nil
	}

	flattenedInstances := make([]interface{}, 0)
	for _, instance := range instances {
		rawInstance := map[string]interface{}{
			"hostname":    instance.Hostname,
			"size":        instance.Size,
			"region":      instance.Region,
			"status":      instance.Status,
			"firewall_id": instance.FirewallID,
			"public_ip":   instance.PublicIP,
			"tags":        instance.Tags,
			"created_at":  instance.CreatedAt.UTC().String(),
		}

		flattenedInstances = append(flattenedInstances, rawInstance)
	}

	return flattenedInstances
}

func flattenInstalledApplication(apps []civogo.KubernetesInstalledApplication) []interface{} {
	if apps == nil {
		return nil
	}

	flattenedInstalledApplication := make([]interface{}, 0)
	for _, app := range apps {
		rawInstalledApplication := map[string]interface{}{
			"application": app.Application,
			"version":     app.Version,
			"installed":   app.Installed,
			"category":    app.Category,
		}

		flattenedInstalledApplication = append(flattenedInstalledApplication, rawInstalledApplication)
	}

	return flattenedInstalledApplication
}
