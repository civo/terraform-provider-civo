package civo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Kubernetes Cluster resource, with this you can manage all cluster from terraform
func resourceKubernetesClusterNodePool() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Civo Kubernetes node pool resource. While the default node pool must be defined in the `civo_kubernetes_cluster` resource, this resource can be used to add additional ones to a cluster.",
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The ID of your cluster",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"label": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: utils.ValidateNameOnlyContainsAlphanumericCharacters,
				Description:      "Node pool label, if you don't provide one, we will generate one for you",
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The region of the node pool, has to match that of the cluster",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"num_target_nodes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Deprecated:  "This field is deprecated, please use `node_count` instead",
				Description: "the number of instances to create (optional, the default at the time of writing is 3)",
			},
			"node_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "the number of instances to create (optional, the default at the time of writing is 3)",
			},
			"target_nodes_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated, please use `size` instead",
				Description: "the size of each node (optional, the default is currently g4s.kube.medium)",
			},
			"size": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "the size of each node (optional, the default is currently g4s.kube.medium)",
			},
			"instance_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance names in the nodepool",
			},
		},
		CreateContext: resourceKubernetesClusterNodePoolCreate,
		ReadContext:   resourceKubernetesClusterNodePoolRead,
		UpdateContext: resourceKubernetesClusterNodePoolUpdate,
		DeleteContext: resourceKubernetesClusterNodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKubernetesClusterNodePoolImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

// function to create a new cluster
func resourceKubernetesClusterNodePoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	clusterID := d.Get("cluster_id").(string)

	count := 0
	if attr, ok := d.GetOk("node_count"); ok {
		count = attr.(int)
	} else {
		count = 3
	}

	size := ""
	if attr, ok := d.GetOk("size"); ok {
		size = attr.(string)
	} else {
		size = "g4s.kube.medium"
	}

	nodePoolLabel := uuid.NewString()
	if attr, ok := d.GetOk("label"); ok {
		nodePoolLabel = attr.(string)
	}

	log.Printf("[INFO] getting kubernetes cluster %s in the region %s", clusterID, apiClient.Region)
	getKubernetesCluster, err := apiClient.GetKubernetesCluster(clusterID)
	if err != nil {
		return diag.Errorf("[INFO] error getting kubernetes cluster: %s", clusterID)
	}

	newPool := []civogo.KubernetesClusterPoolConfig{}
	for _, v := range getKubernetesCluster.Pools {
		newPool = append(newPool, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
	}

	newPool = append(newPool, civogo.KubernetesClusterPoolConfig{ID: nodePoolLabel, Count: count, Size: size})

	log.Printf("[INFO] configuring kubernetes cluster %s to add pool %s", getKubernetesCluster.ID, nodePoolLabel)
	config := &civogo.KubernetesClusterConfig{
		Pools:  newPool,
		Region: apiClient.Region,
	}

	log.Printf("[INFO] Creating a new kubernetes cluster pool %s", nodePoolLabel)
	_, err = apiClient.UpdateKubernetesCluster(getKubernetesCluster.ID, config)
	if err != nil {
		return diag.Errorf("[ERR] failed to create the kubernetes cluster: %s", err)
	}

	d.SetId(nodePoolLabel)

	err = waitForKubernetesNodePoolCreate(apiClient, d, clusterID)
	if err != nil {
		return diag.Errorf("Error creating Kubernetes node pool: %s", err)
	}

	return resourceKubernetesClusterNodePoolRead(ctx, d, m)

}

// function to read the kubernetes cluster
func resourceKubernetesClusterNodePoolRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)
	clusterID := d.Get("cluster_id").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var poolIndex int

	log.Printf("[INFO] retrieving the kubernetes cluster %s", clusterID)
	resp, err := apiClient.GetKubernetesCluster(clusterID)
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("[ERR] failed to find the kubernetes cluster: %s", err)
	}

	for i, v := range resp.Pools {
		if v.ID == d.Id() {
			poolIndex = i
		}
	}

	d.SetId(resp.Pools[poolIndex].ID)
	d.Set("cluster_id", resp.ID)
	d.Set("node_count", resp.Pools[poolIndex].Count)
	d.Set("size", resp.Pools[poolIndex].Size)

	poolInstanceNames := make([]string, 0)
	poolInstanceNames = append(poolInstanceNames, resp.Pools[poolIndex].InstanceNames...)

	d.Set("instance_names", poolInstanceNames)

	return diags
}

// function to update the kubernetes cluster
func resourceKubernetesClusterNodePoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	old, new := d.GetChange("size")
	if old != new {
		return diag.Errorf("[ERR] Size change (%q) for existing pool is not available at this moment", "size")
	}

	clusterID := d.Get("cluster_id").(string)
	count := 0

	if d.HasChange("node_count") {
		count = d.Get("node_count").(int)
	}

	getKubernetesCluster, err := apiClient.GetKubernetesCluster(clusterID)
	if err != nil {
		return diag.Errorf("[INFO] error getting kubernetes cluster: %s", clusterID)
	}

	nodePool := []civogo.KubernetesClusterPoolConfig{}
	for _, v := range getKubernetesCluster.Pools {
		nodePool = append(nodePool, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
	}

	nodePool = updateNodePool(nodePool, d.Id(), count)

	log.Printf("[INFO] configuring kubernetes cluster %s to add pool", d.Id())
	config := &civogo.KubernetesClusterConfig{
		Pools:  nodePool,
		Region: apiClient.Region,
	}

	log.Printf("[INFO] updating the kubernetes cluster %s", d.Id())
	_, err = apiClient.UpdateKubernetesCluster(getKubernetesCluster.ID, config)
	if err != nil {
		return diag.Errorf("[ERR] failed to update kubernetes cluster: %s", err)
	}

	err = waitForKubernetesNodePoolCreate(apiClient, d, clusterID)
	if err != nil {
		return diag.Errorf("Error updating Kubernetes node pool: %s", err)
	}

	return resourceKubernetesClusterNodePoolRead(ctx, d, m)
}

// function to delete the kubernetes cluster
func resourceKubernetesClusterNodePoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	clusterID := d.Get("cluster_id").(string)
	getKubernetesCluster, err := apiClient.GetKubernetesCluster(clusterID)
	if err != nil {
		return diag.Errorf("[INFO] error getting kubernetes cluster: %s", clusterID)
	}

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	nodePool := []civogo.KubernetesClusterPoolConfig{}
	for _, v := range getKubernetesCluster.Pools {
		nodePool = append(nodePool, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
	}

	nodePool = removeNodePool(nodePool, d.Id())
	config := &civogo.KubernetesClusterConfig{
		Pools:  nodePool,
		Region: apiClient.Region,
	}

	log.Printf("[INFO] deleting the kubernetes cluster %s", d.Id())
	_, err = apiClient.UpdateKubernetesCluster(getKubernetesCluster.ID, config)
	if err != nil {
		return diag.Errorf("[INFO] an error occurred while trying to delete the kubernetes cluster pool %s", err)
	}

	err = waitForKubernetesNodePoolDelete(apiClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceKubernetesClusterNodePoolRead(ctx, d, m)
}

// custom import to able to add a node pool to the terraform
func resourceKubernetesClusterNodePoolImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	apiClient := m.(*civogo.Client)
	regions, err := apiClient.ListRegions()
	if err != nil {
		return nil, err
	}

	clusterID, nodePoolID, err := utils.ResourceCommonParseID(d.Id())
	if err != nil {
		return nil, err
	}

	poolFound := false
	for _, region := range regions {
		if poolFound {
			break
		}

		currentRegionCode := region.Code
		apiClient.Region = currentRegionCode
		log.Printf("[INFO] Retriving the node pool %s from region %s", nodePoolID, currentRegionCode)
		resp, err := apiClient.GetKubernetesCluster(clusterID)
		if err != nil {
			continue // move on and find in another region
		}

		for _, v := range resp.Pools {
			if v.ID == nodePoolID {
				poolFound = true
				d.SetId(v.ID)
				d.Set("cluster_id", resp.ID)
				d.Set("label", v.ID)
				d.Set("node_count", v.Count)
				d.Set("size", v.Size)
				d.Set("region", currentRegionCode)
			}
		}
	}

	if !poolFound {
		return nil, fmt.Errorf("[ERR] Node pool %s not found", nodePoolID)
	}

	return []*schema.ResourceData{d}, nil
}

// RemoveNodePool is a utility function to remove node pool from a kuberentes cluster
func removeNodePool(s []civogo.KubernetesClusterPoolConfig, id string) []civogo.KubernetesClusterPoolConfig {
	for k, v := range s {
		if strings.Contains(v.ID, id) {
			s[len(s)-1], s[k] = s[k], s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	return s
}

// UpdateNodePool is a utility function to update node pool from a kuberentes cluster
func updateNodePool(s []civogo.KubernetesClusterPoolConfig, id string, count int) []civogo.KubernetesClusterPoolConfig {
	for k, v := range s {
		if strings.Contains(v.ID, id) {
			s[k].Count = count
			break
		}
	}
	return s
}

// waitForKubernetesNodePoolCreate is a utility function to wait for a node pool to be created
func waitForKubernetesNodePoolCreate(client *civogo.Client, d *schema.ResourceData, clusterID string) error {
	var (
		tickerInterval        = 10 * time.Second
		timeoutSeconds        = d.Timeout(schema.TimeoutCreate).Seconds()
		timeout               = int(timeoutSeconds / tickerInterval.Seconds())
		n                     = 0
		totalRequiredInstance = 0
		totalRunningInstance  = 0
		ticker                = time.NewTicker(tickerInterval)
		nodePoolID            = d.Id()
	)

	for range ticker.C {
		cluster, err := client.GetKubernetesCluster(clusterID)
		if err != nil {
			ticker.Stop()
			return fmt.Errorf("error trying to read cluster state: %s", err)
		}

		for _, v := range cluster.RequiredPools {
			if v.ID == nodePoolID {
				totalRequiredInstance = v.Count
				break
			}
		}

		for _, v := range cluster.Pools {
			if v.ID == nodePoolID {
				totalRunningInstance = v.Count
				break
			}
		}

		allRunning := totalRunningInstance == totalRequiredInstance
		for _, n := range cluster.Pools {
			if n.ID == nodePoolID {
				for _, node := range n.Instances {
					if node.Status != "ACTIVE" {
						allRunning = false
					}
				}
			}
		}

		if allRunning {
			ticker.Stop()
			return nil
		}

		if n > timeout {
			ticker.Stop()
			break
		}

		n++
	}

	return fmt.Errorf("timeout waiting to create nodepool %s", nodePoolID)
}

// waitForKubernetesNodePoolDelete is a utility function to wait for a node pool to be deleted
func waitForKubernetesNodePoolDelete(client *civogo.Client, d *schema.ResourceData) error {
	var (
		tickerInterval = 10 * time.Second
		timeoutSeconds = d.Timeout(schema.TimeoutDelete).Seconds()
		timeout        = int(timeoutSeconds / tickerInterval.Seconds())
		n              = 0
		ticker         = time.NewTicker(tickerInterval)
		clusterID      = d.Get("cluster_id").(string)
		nodePoolID     = d.Id()
	)

	for range ticker.C {

		cluster, err := client.GetKubernetesCluster(clusterID)
		if err != nil {
			ticker.Stop()
			return fmt.Errorf("error trying to read cluster state: %s", err)
		}

		nodePool := []civogo.KubernetesClusterPoolConfig{}
		for _, v := range cluster.Pools {
			nodePool = append(nodePool, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
		}

		if !utils.InPool(nodePoolID, nodePool) {
			ticker.Stop()
			return nil
		}

		if n > timeout {
			ticker.Stop()
			break
		}

		n++
	}

	return fmt.Errorf("timeout waiting to delete nodepool")
}
