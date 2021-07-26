package civo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Kubernetes Cluster resource, with this you can manage all cluster from terraform
func resourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "a name for your cluster, must be unique within your account (required)",
				ValidateFunc: utils.ValidateNameSize,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region for the cluster, if not declare we use the region in declared in the provider",
			},
			"network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The network for the cluster, if not declare we use the default one",
			},
			"num_target_nodes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "the number of instances to create (optional, the default at the time of writing is 3)",
			},
			"target_nodes_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "the size of each node (optional, the default is currently g2.k3s.medium)",
			},
			"kubernetes_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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
					"Application names are case-sensitive; the available applications can be listed with the civo CLI:" +
					"'civo kubernetes applications ls'." +
					"If you want to remove a default installed application, prefix it with a '-', e.g. -Traefik.",
			},
			// Computed resource
			"instances":              instanceSchema(),
			"installed_applications": applicationSchema(),
			"pools":                  nodePoolSchema(),
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
		Create: resourceKubernetesClusterCreate,
		Read:   resourceKubernetesClusterRead,
		Update: resourceKubernetesClusterUpdate,
		Delete: resourceKubernetesClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// schema for the instances
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
func nodePoolSchema() *schema.Schema {
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
				"instances": instanceSchema(),
			},
		},
	}
}

// schema for the application in the cluster
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

// function to create a new cluster
func resourceKubernetesClusterCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] configuring a new kubernetes cluster %s", d.Get("name").(string))

	config := &civogo.KubernetesClusterConfig{
		Region:      apiClient.Region,
		NodeDestroy: "",
	}

	if name, ok := d.GetOk("name"); ok {
		config.Name = name.(string)
	} else {
		config.Name = utils.RandomName()
	}

	if networtID, ok := d.GetOk("network_id"); ok {
		config.NetworkID = networtID.(string)
	} else {
		defaultNetwork, err := apiClient.GetDefaultNetwork()
		if err != nil {
			return fmt.Errorf("[ERR] failed to get the default network: %s", err)
		}
		config.NetworkID = defaultNetwork.ID
	}

	if attr, ok := d.GetOk("num_target_nodes"); ok {
		config.NumTargetNodes = attr.(int)
	} else {
		config.NumTargetNodes = 3
	}

	if attr, ok := d.GetOk("target_nodes_size"); ok {
		config.TargetNodesSize = attr.(string)
	} else {
		config.TargetNodesSize = "g3.k3s.small"
	}

	if attr, ok := d.GetOk("kubernetes_version"); ok {
		config.KubernetesVersion = attr.(string)
	}

	if attr, ok := d.GetOk("tags"); ok {
		config.Tags = attr.(string)
	} else {
		config.Tags = ""
	}

	if attr, ok := d.GetOk("applications"); ok {
		if utils.CheckAPPName(attr.(string), apiClient) {
			config.Applications = attr.(string)
		} else {
			return fmt.Errorf("[ERR] the app that tries to install is not valid: %s", attr.(string))
		}
	} else {
		config.Applications = ""
	}

	log.Printf("[INFO] creating a new kubernetes cluster %s", d.Get("name").(string))
	log.Printf("[INFO] kubernertes config %+v", config)
	resp, err := apiClient.NewKubernetesClusters(config)
	if err != nil {
		return fmt.Errorf("[ERR] failed to create the kubernetes cluster: %s", err)
	}

	d.SetId(resp.ID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"BUILDING"},
		Target:  []string{"ACTIVE"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.GetKubernetesCluster(d.Id())
			if err != nil {
				return 0, "", err
			}
			return resp, resp.Status, nil
		},
		Timeout:        60 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 10,
	}
	_, err = createStateConf.WaitForStateContext(context.Background())
	if err != nil {
		return fmt.Errorf("error waiting for cluster (%s) to be created: %s", d.Id(), err)
	}

	return resourceKubernetesClusterRead(d, m)

}

// function to read the kubernetes cluster
func resourceKubernetesClusterRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] retrieving the kubernetes cluster %s", d.Id())
	resp, err := apiClient.GetKubernetesCluster(d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERR] failed to find the kubernetes cluster: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("region", apiClient.Region)
	d.Set("network_id", resp.NetworkID)
	d.Set("num_target_nodes", resp.NumTargetNode)
	d.Set("target_nodes_size", resp.TargetNodeSize)
	d.Set("kubernetes_version", resp.KubernetesVersion)
	d.Set("tags", strings.Join(resp.Tags, ", "))
	d.Set("status", resp.Status)
	d.Set("ready", resp.Ready)
	d.Set("kubeconfig", resp.KubeConfig)
	d.Set("api_endpoint", resp.APIEndPoint)
	d.Set("master_ip", resp.MasterIP)
	d.Set("dns_entry", resp.DNSEntry)
	// d.Set("built_at", resp.BuiltAt.UTC().String())
	d.Set("created_at", resp.CreatedAt.UTC().String())

	if err := d.Set("instances", flattenInstances(resp.Instances)); err != nil {
		return fmt.Errorf("[ERR] error retrieving the instances for kubernetes cluster error: %#v", err)
	}

	if err := d.Set("pools", flattenNodePool(resp)); err != nil {
		return fmt.Errorf("[ERR] error retrieving the pool for kubernetes cluster error: %#v", err)
	}

	if err := d.Set("installed_applications", flattenInstalledApplication(resp.InstalledApplications)); err != nil {
		return fmt.Errorf("[ERR] error retrieving the installed application for kubernetes cluster error: %#v", err)
	}

	return nil
}

// function to update the kubernetes cluster
func resourceKubernetesClusterUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	config := &civogo.KubernetesClusterConfig{}

	if d.HasChange("num_target_nodes") {
		config.NumTargetNodes = d.Get("num_target_nodes").(int)
		config.Region = apiClient.Region
	}

	if d.HasChange("kubernetes_version") {
		config.KubernetesVersion = d.Get("kubernetes_version").(string)
		config.Region = apiClient.Region
	}

	if d.HasChange("applications") {
		config.Applications = d.Get("applications").(string)
		config.Region = apiClient.Region
	}

	if d.HasChange("name") {
		config.Applications = d.Get("name").(string)
		config.Region = apiClient.Region
	}

	if d.HasChange("tags") {
		config.Tags = d.Get("tags").(string)
	}

	log.Printf("[INFO] updating the kubernetes cluster %s", d.Id())
	_, err := apiClient.UpdateKubernetesCluster(d.Id(), config)
	if err != nil {
		return fmt.Errorf("[ERR] failed to update kubernetes cluster: %s", err)
	}

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"SCALING"},
		Target:  []string{"ACTIVE"},
		Refresh: func() (interface{}, string, error) {
			resp, err := apiClient.GetKubernetesCluster(d.Id())
			if err != nil {
				return 0, "", err
			}
			return resp, resp.Status, nil
		},
		Timeout:        60 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 10,
	}
	_, err = createStateConf.WaitForStateContext(context.Background())
	if err != nil {
		return fmt.Errorf("error waiting for cluster (%s) to be created: %s", d.Id(), err)
	}

	return resourceKubernetesClusterRead(d, m)
}

// function to delete the kubernetes cluster
func resourceKubernetesClusterDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] deleting the kubernetes cluster %s", d.Id())
	_, err := apiClient.DeleteKubernetesCluster(d.Id())
	if err != nil {
		return fmt.Errorf("[INFO] an error occurred while tring to delete the kubernetes cluster %s", err)
	}

	return nil
}

// function to flatten all instances inside the cluster
func flattenInstances(instances []civogo.KubernetesInstance) []interface{} {
	if instances == nil {
		return nil
	}

	flattenedInstances := make([]interface{}, 0)
	for _, instance := range instances {
		rawInstance := map[string]interface{}{
			"hostname":  instance.Hostname,
			"cpu_cores": instance.CPUCores,
			"ram_mb":    instance.RAMMegabytes,
			"disk_gb":   instance.DiskGigabytes,
			"status":    instance.Status,
		}

		flattenedInstances = append(flattenedInstances, rawInstance)
	}

	return flattenedInstances
}

// function to flatten all instances inside the cluster
func flattenNodePool(cluster *civogo.KubernetesCluster) []interface{} {
	if cluster.Pools == nil {
		return nil
	}

	flattenedPool := make([]interface{}, 0)
	for _, pool := range cluster.Pools {
		flattenedPoolInstance := make([]interface{}, 0)
		for _, v := range pool.Instances {

			instanceData := searchInstance(cluster.Instances, v.Hostname)

			rawPoolInstance := map[string]interface{}{
				"hostname":  v.Hostname,
				"size":      pool.Size,
				"cpu_cores": instanceData.CPUCores,
				"ram_mb":    instanceData.RAMMegabytes,
				"disk_gb":   instanceData.DiskGigabytes,
				"status":    v.Status,
				"tags":      v.Tags,
			}
			flattenedPoolInstance = append(flattenedPoolInstance, rawPoolInstance)
		}

		instanceName := append(pool.InstanceNames, pool.InstanceNames...)

		rawPool := map[string]interface{}{
			"id":             pool.ID,
			"count":          pool.Count,
			"size":           pool.Size,
			"instance_names": instanceName,
			"instances":      flattenedPoolInstance,
		}

		flattenedPool = append(flattenedPool, rawPool)
	}

	return flattenedPool
}

// function to flatten all applications inside the cluster
func flattenInstalledApplication(apps []civogo.KubernetesInstalledApplication) []interface{} {
	if apps == nil {
		return nil
	}

	flattenedInstalledApplication := make([]interface{}, 0)
	for _, app := range apps {
		rawInstalledApplication := map[string]interface{}{
			"application": app.Name,
			"version":     app.Version,
			"installed":   app.Installed,
			"category":    app.Category,
		}

		flattenedInstalledApplication = append(flattenedInstalledApplication, rawInstalledApplication)
	}

	return flattenedInstalledApplication
}

func searchInstance(instanceList []civogo.KubernetesInstance, hostname string) civogo.KubernetesInstance {
	for _, v := range instanceList {
		if strings.Contains(v.Hostname, hostname) {
			return v
		}
	}
	return civogo.KubernetesInstance{}
}
