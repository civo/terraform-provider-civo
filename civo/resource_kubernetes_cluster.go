package civo

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				Deprecated:   "This field will be deprecated in the next major release, please use the 'pools' field instead",
				Description:  "The number of instances to create (optional, the default at the time of writing is 3)",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"target_nodes_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Deprecated:  "This field will be deprecated in the next major release, please use the 'pools' field instead",
				Description: "The size of each node (optional, the default is currently g4s.kube.medium)",
			},
			"kubernetes_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The version of k3s to install (optional, the default is currently the latest available)",
			},
			"cni": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The cni for the k3s to install (the default is `flannel`) valid options are `cilium` or `flannel`",
				ValidateFunc: utils.ValidateCNIName,
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
			"cluster_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of cluster to create, valid options are `k3s` or `talos` the default is `k3s`",
			},
			// Computed resource
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
		CreateContext: resourceKubernetesClusterCreate,
		ReadContext:   resourceKubernetesClusterRead,
		UpdateContext: resourceKubernetesClusterUpdate,
		DeleteContext: resourceKubernetesClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

// schema for the node pool in the cluster
func nodePoolSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"label": {
					Type:             schema.TypeString,
					Optional:         true,
					Computed:         true,
					ValidateDiagFunc: utils.ValidateNameOnlyContainsAlphanumericCharacters,
					Description:      "Node pool label, if you don't provide one, we will generate one for you",
				},
				"node_count": {
					Type:         schema.TypeInt,
					Required:     true,
					Description:  "Number of nodes in the nodepool",
					ValidateFunc: validation.IntAtLeast(1),
				},
				"size": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Size of the nodes in the nodepool",
				},
				"instance_names": {
					Type:        schema.TypeList,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Instance names in the nodepool",
				},
				"public_ip_node_pool": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "Node pool belongs to the public ip node pool",
				},
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
func resourceKubernetesClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
			return diag.Errorf("[ERR] failed to get the default network: %s", err)
		}
		config.NetworkID = defaultNetwork.ID
	}

	if attr, ok := d.GetOk("kubernetes_version"); ok {
		config.KubernetesVersion = attr.(string)
	}

	if attr, ok := d.GetOk("tags"); ok {
		config.Tags = attr.(string)
	} else {
		config.Tags = ""
	}

	if attr, ok := d.GetOk("cni"); ok {
		config.CNIPlugin = attr.(string)
	}

	if attr, ok := d.GetOk("applications"); ok {
		if utils.CheckAPPName(attr.(string), apiClient) {
			config.Applications = attr.(string)
		} else {
			return diag.Errorf("[ERR] the app that tries to install is not valid: %s", attr.(string))
		}
	} else {
		config.Applications = ""
	}

	if attr, ok := d.GetOk("cluster_type"); ok {
		config.ClusterType = attr.(string)
	}

	if attr, ok := d.GetOk("firewall_id"); ok {
		firewallID := attr.(string)
		firewall, err := apiClient.FindFirewall(firewallID)
		if err != nil {
			return diag.Errorf("[ERR] unable to find firewall - %s", err)
		}

		if firewall.NetworkID != config.NetworkID {
			return diag.Errorf("[ERR] firewall %s is not part of network %s", firewall.ID, config.NetworkID)
		}

		config.InstanceFirewall = firewallID
	}

	pools := expandNodePools(d.Get("pools").([]interface{}))
	config.Pools = pools

	log.Printf("[INFO] creating a new kubernetes cluster %s", d.Get("name").(string))
	log.Printf("[INFO] kubernertes config %+v", config)
	resp, err := apiClient.NewKubernetesClusters(config)
	if err != nil {
		return diag.Errorf("[ERR] failed to create the kubernetes cluster: %s", err)
	}

	d.SetId(resp.ID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"BUILDING", "AVAILABLE", "UPGRADING", "SCALING"},
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
		return diag.Errorf("error waiting for cluster (%s) to be created: %s", d.Id(), err)
	}

	return resourceKubernetesClusterRead(ctx, d, m)
}

// function to read the kubernetes cluster
func resourceKubernetesClusterRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.Errorf("[ERR] failed to find the kubernetes cluster: %s", err)
	}

	d.Set("name", resp.Name)
	d.Set("region", apiClient.Region)
	d.Set("network_id", resp.NetworkID)
	d.Set("num_target_nodes", resp.NumTargetNode)
	d.Set("target_nodes_size", resp.TargetNodeSize)
	d.Set("kubernetes_version", resp.KubernetesVersion)
	d.Set("cluster_type", resp.ClusterType)
	d.Set("cni", resp.CNIPlugin)
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

	if err := d.Set("pools", flattenNodePool(resp)); err != nil {
		return diag.Errorf("[ERR] error retrieving the pool for kubernetes cluster error: %#v", err)
	}

	if err := d.Set("installed_applications", flattenInstalledApplication(resp.InstalledApplications)); err != nil {
		return diag.Errorf("[ERR] error retrieving the installed application for kubernetes cluster error: %#v", err)
	}

	return nil
}

// function to update the kubernetes cluster
func resourceKubernetesClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	config := &civogo.KubernetesClusterConfig{}

	if d.HasChange("network_id") {
		return diag.Errorf("[ERR] Network change (%q) for existing cluster is not available at this moment", "network_id")
	}

	if d.HasChange("firewall_id") {
		return diag.Errorf("[ERR] Firewall change (%q) for existing cluster is not available at this moment", "firewall_id")
	}

	// Update the node pool if necessary
	if !d.HasChange("pools") {
		return resourceKubernetesClusterRead(ctx, d, m)
	}

	if d.HasChange("pools") {
		old, new := d.GetChange("pools")
		oldPool := old.([]interface{})[0].(map[string]interface{})
		newPool := new.([]interface{})[0].(map[string]interface{})

		// if the size is different, then return and error as we can't change the size of a pool
		if oldPool["size"].(string) != newPool["size"].(string) {
			return diag.Errorf("[ERR] Size change (%q) for existing cluster is not available at this moment", "size")
		}

		config.Region = apiClient.Region
		kubernetesCluster, err := apiClient.FindKubernetesCluster(d.Id())
		if err != nil {
			return diag.Errorf("[ERR] failed to find the kubernetes cluster: %s", err)
		}

		targetNodePool := ""
		nodePools := []civogo.KubernetesClusterPoolConfig{}
		for _, v := range kubernetesCluster.Pools {
			nodePools = append(nodePools, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
			if targetNodePool == "" && v.ID == newPool["label"].(string) {
				targetNodePool = v.ID
			}
		}

		nodePools = updateNodePool(nodePools, targetNodePool, newPool["node_count"].(int))
		config.Pools = nodePools
	}

	if d.HasChange("kubernetes_version") {
		// config.KubernetesVersion = d.Get("kubernetes_version").(string)
		// config.Region = apiClient.Region
		return diag.Errorf("[ERR] Kubernetes version upgrade (%q attribute) is not supported yet", "kubernetes_version")
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
	_, err := apiClient.UpdateKubernetesCluster(d.Id(), config)
	if err != nil {
		return diag.Errorf("[ERR] failed to update kubernetes cluster: %s", err)
	}

	err = waitForKubernetesNodePoolCreate(apiClient, d, d.Id())
	if err != nil {
		return diag.Errorf("Error updating Kubernetes node pool: %s", err)
	}

	return resourceKubernetesClusterRead(ctx, d, m)
}

// function to delete the kubernetes cluster
func resourceKubernetesClusterDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	log.Printf("[INFO] deleting the kubernetes cluster %s", d.Id())
	_, err := apiClient.DeleteKubernetesCluster(d.Id())
	if err != nil {
		return diag.Errorf("[INFO] an error occurred while trying to delete the kubernetes cluster %s", err)
	}

	return nil
}

// function to flatten all instances inside the cluster
func flattenNodePool(cluster *civogo.KubernetesCluster) []interface{} {
	if cluster.Pools == nil {
		return nil
	}

	flattenedPool := make([]interface{}, 0)

	poolInstanceNames := make([]string, 0)
	poolInstanceNames = append(poolInstanceNames, cluster.Pools[0].InstanceNames...)

	rawPool := map[string]interface{}{
		"label":               cluster.Pools[0].ID,
		"node_count":          cluster.Pools[0].Count,
		"size":                cluster.Pools[0].Size,
		"instance_names":      poolInstanceNames,
		"public_ip_node_pool": cluster.Pools[0].PublicIPNodePool,
	}

	flattenedPool = append(flattenedPool, rawPool)

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

// exapandNodePools function to expand the node pools
func expandNodePools(nodePools []interface{}) []civogo.KubernetesClusterPoolConfig {
	expandedNodePools := make([]civogo.KubernetesClusterPoolConfig, 0, len(nodePools))
	for _, rawPool := range nodePools {
		pool := rawPool.(map[string]interface{})

		poolID := uuid.NewString()
		if pool["label"].(string) != "" {
			poolID = pool["label"].(string)
		}

		cr := civogo.KubernetesClusterPoolConfig{
			ID:    poolID,
			Size:  pool["size"].(string),
			Count: pool["node_count"].(int),
		}

		if pool["public_ip_node_pool"].(bool) {
			cr.PublicIPNodePool = pool["public_ip_node_pool"].(bool)
		}

		expandedNodePools = append(expandedNodePools, cr)
	}

	return expandedNodePools
}
