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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Kubernetes Cluster resource, with this you can manage all cluster from terraform
func resourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Civo Kubernetes cluster resource. This can be used to create, delete, and modify clusters.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Name for your cluster, must be unique within your account",
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
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "The number of instances to create (optional, the default at the time of writing is 3)",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"target_nodes_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The size of each node (optional, the default is currently g3.k3s.medium)",
			},
			"kubernetes_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The version of k3s to install (optional, the default is currently the latest available)",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Space separated list of tags, to be used freely as required",
			},
			"applications": {
				Type:     schema.TypeString,
				Optional: true,
				Description: strings.Join([]string{
					"Comma separated list of applications to install.",
					"Spaces within application names are fine, but shouldn't be either side of the comma.",
					"Application names are case-sensitive; the available applications can be listed with the Civo CLI:",
					"'civo kubernetes applications ls'.",
					"If you want to remove a default installed application, prefix it with a '-', e.g. -Traefik.",
					"For application that supports plans, you can use 'app_name:app_plan' format e.g. 'Linkerd:Linkerd & Jaeger' or 'MariaDB:5GB'.",
				}, " "),
			},
			"firewall_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The existing firewall ID to use for this cluster",
			},
			// Computed resource
			"instances":              instanceSchema(),
			"installed_applications": applicationSchema(),
			"pools":                  nodePoolSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the cluster",
			},
			"ready": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When cluster is ready, this will return `true`",
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The kubeconfig of the cluster",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API server endpoint of the cluster",
			},
			"master_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the master node",
			},
			"dns_entry": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The DNS name of the cluster",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the cluster was created",
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
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Instance's hostname",
				},
				"size": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Instance's size",
				},
				"cpu_cores": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Instance's CPU cores",
				},
				"ram_mb": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Instance's RAM (MB)",
				},
				"disk_gb": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Instance's disk (GB)",
				},
				"status": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Instance's status",
				},
				"tags": {
					Type:        schema.TypeSet,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Instance's tags",
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
				"count": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Number of nodes in the nodepool",
				},
				"size": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Size of the nodes in the nodepool",
				},
				"instance_names": {
					Type:        schema.TypeSet,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Instance names in the nodepool",
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
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of application",
				},
				"version": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Version of application",
				},
				"installed": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Application installation status (`true` if installed)",
				},
				"category": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Category of the application",
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
		config.TargetNodesSize = "g3.k3s.medium"
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

	if attr, ok := d.GetOk("firewall_id"); ok {
		firewallID := attr.(string)
		firewall, err := apiClient.FindFirewall(firewallID)
		if err != nil {
			return fmt.Errorf("[ERR] unable to find firewall - %s", err)
		}

		if firewall.NetworkID != config.NetworkID {
			return fmt.Errorf("[ERR] firewall %s is not part of network %s", firewall.ID, config.NetworkID)
		}

		config.InstanceFirewall = firewallID
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
	d.Set("tags", strings.Join(resp.Tags, " ")) // space separated tags
	d.Set("status", resp.Status)
	d.Set("ready", resp.Ready)
	d.Set("kubeconfig", resp.KubeConfig)
	d.Set("api_endpoint", resp.APIEndPoint)
	d.Set("master_ip", resp.MasterIP)
	d.Set("dns_entry", resp.DNSEntry)
	// d.Set("built_at", resp.BuiltAt.UTC().String())
	d.Set("created_at", resp.CreatedAt.UTC().String())
	d.Set("firewall_id", resp.FirewallID)

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

	if d.HasChange("network_id") {
		return fmt.Errorf("[ERR] Network change (%q) for existing cluster is not available at this moment", "network_id")
	}

	if d.HasChange("firewall_id") {
		return fmt.Errorf("[ERR] Firewall change (%q) for existing cluster is not available at this moment", "firewall_id")
	}

	if d.HasChange("target_nodes_size") {
		errMsg := []string{
			"[ERR] Unable to update 'target_nodes_size' after creation.",
			"Please create a new node pool with the new node size.",
		}
		return fmt.Errorf(strings.Join(errMsg, " "))
	}

	if d.HasChange("num_target_nodes") {
		numTargetNodes := d.Get("num_target_nodes").(int)

		config.Region = apiClient.Region
		kubernetesCluster, err := apiClient.FindKubernetesCluster(d.Id())
		if err != nil {
			return err
		}

		targetNodePool := ""
		nodePools := []civogo.KubernetesClusterPoolConfig{}
		for _, v := range kubernetesCluster.Pools {
			nodePools = append(nodePools, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})

			if targetNodePool == "" && v.Size == d.Get("target_nodes_size").(string) {
				targetNodePool = v.ID
			}
		}

		nodePools = updateNodePool(nodePools, targetNodePool, numTargetNodes)
		config.Pools = nodePools
	}

	if d.HasChange("kubernetes_version") {
		// config.KubernetesVersion = d.Get("kubernetes_version").(string)
		// config.Region = apiClient.Region
		return fmt.Errorf("[ERR] Kubernetes version upgrade (%q attribute) is not supported yet", "kubernetes_version")
	}

	if d.HasChange("applications") {
		config.Applications = d.Get("applications").(string)
		config.Region = apiClient.Region
	}

	if d.HasChange("name") {
		config.Name = d.Get("name").(string)
		config.Region = apiClient.Region
	}

	if d.HasChange("tags") {
		config.Tags = d.Get("tags").(string)
	}

	log.Printf("[INFO] updating the kubernetes cluster %s", d.Id())
	log.Printf("[DEBUG] KubernetesClusterConfig: %+v\n", config)
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
