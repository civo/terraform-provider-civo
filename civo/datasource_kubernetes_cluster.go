package civo

import (
	"fmt"
	"log"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Data source to get from the api a specific instance
// using the id or the hostname
func dataSourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKubernetesClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"num_target_nodes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"target_nodes_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kubernetes_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"applications": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instances":              dataSourceInstanceSchema(),
			"installed_applications": dataSourceApplicationSchema(),
			"pools":                  dataSourcenodePoolSchema(),
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
			"master_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_entry": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// schema for the instances
func dataSourceInstanceSchema() *schema.Schema {
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
				"cpu_cores": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"ram_mb": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"disk_gb": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"status": {
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

// schema for the node pool in the cluster
func dataSourcenodePoolSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"count": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"size": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"instance_names": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"instances": dataSourceInstanceSchema(),
			},
		},
	}
}

// schema for the application in the cluster
func dataSourceApplicationSchema() *schema.Schema {
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

func dataSourceKubernetesClusterRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	var foundCluster *civogo.KubernetesCluster

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the kubernetes Cluster by id")
		kubeCluster, err := apiClient.FindKubernetesCluster(id.(string))
		if err != nil {
			return fmt.Errorf("[ERR] failed to retrive kubernetes cluster: %s", err)
		}

		foundCluster = kubeCluster
	} else if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the kubernetes Cluster by name")
		kubeCluster, err := apiClient.FindKubernetesCluster(name.(string))
		if err != nil {
			return fmt.Errorf("[ERR] failed to retrive kubernetes cluster: %s", err)
		}

		foundCluster = kubeCluster
	}

	d.SetId(foundCluster.ID)
	d.Set("name", foundCluster.Name)
	d.Set("num_target_nodes", foundCluster.NumTargetNode)
	d.Set("target_nodes_size", foundCluster.TargetNodeSize)
	d.Set("kubernetes_version", foundCluster.KubernetesVersion)
	d.Set("tags", foundCluster.Tags)
	d.Set("status", foundCluster.Status)
	d.Set("ready", foundCluster.Ready)
	d.Set("kubeconfig", foundCluster.KubeConfig)
	d.Set("api_endpoint", foundCluster.APIEndPoint)
	d.Set("master_ip", foundCluster.MasterIP)
	d.Set("dns_entry", foundCluster.DNSEntry)
	d.Set("created_at", foundCluster.CreatedAt.UTC().String())

	if err := d.Set("pools", flattenNodePool(foundCluster)); err != nil {
		return fmt.Errorf("[ERR] error retrieving the pools for kubernetes cluster error: %#v", err)
	}

	if err := d.Set("instances", flattenInstances(foundCluster.Instances)); err != nil {
		return fmt.Errorf("[ERR] error retrieving the instances for kubernetes cluster error: %#v", err)
	}

	if err := d.Set("installed_applications", flattenInstalledApplication(foundCluster.InstalledApplications)); err != nil {
		return fmt.Errorf("[ERR] error retrieving the installed application for kubernetes cluster error: %#v", err)
	}

	return nil
}
