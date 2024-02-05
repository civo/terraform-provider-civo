package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceKubernetesClusterNodePool function returns a schema.Resource that represents a node pool in a Kubernetes cluster.
// This can be used to create, read, update, and delete operations for a node pool in a Kubernetes cluster from terraform.
func ResourceKubernetesClusterNodePool() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Civo Kubernetes node pool resource. While the default node pool must be defined in the `civo_kubernetes_cluster` resource, this resource can be used to add additional ones to a cluster.",
		Schema:        nodePoolSchema(true),
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

	// We check if the cluster exists before creating the node pool or made any process
	log.Printf("[INFO] getting kubernetes cluster %s in the region %s", clusterID, apiClient.Region)
	getKubernetesCluster, err := apiClient.GetKubernetesCluster(clusterID)
	if err != nil {
		return diag.Errorf("[INFO] error getting kubernetes cluster: %s", clusterID)
	}

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

	nodePoolLabes := map[string]string{}
	if attr, ok := d.GetOk("labels"); ok {
		for k, v := range attr.(map[string]interface{}) {
			nodePoolLabes[k] = v.(string)
		}
	}

	nodePoolTains := []corev1.Taint{}
	if attr, ok := d.GetOk("taint"); ok {
		for _, v := range attr.(*schema.Set).List() {
			taint := v.(map[string]interface{})
			nodePoolTains = append(nodePoolTains, corev1.Taint{
				Key:    taint["key"].(string),
				Value:  taint["value"].(string),
				Effect: corev1.TaintEffect(taint["effect"].(string)),
			})
		}
	}

	newPool := &civogo.KubernetesClusterPoolUpdateConfig{
		ID:     nodePoolLabel,
		Count:  count,
		Size:   size,
		Labels: nodePoolLabes,
		Taints: nodePoolTains,
		Region: apiClient.Region,
	}

	if value, ok := d.GetOk("public_ip_node_pool"); ok {
		newPool.PublicIPNodePool = value.(bool)
	}

	log.Printf("[INFO] configuring kubernetes cluster %s to add pool %s", getKubernetesCluster.ID, nodePoolLabel)
	log.Printf("[INFO] Creating a new kubernetes cluster pool %s", nodePoolLabel)
	_, err = apiClient.CreateKubernetesClusterPool(getKubernetesCluster.ID, newPool)
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

	log.Printf("[INFO] retrieving the kubernetes cluster %s", clusterID)
	resp, err := apiClient.GetKubernetesCluster(clusterID)
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("[ERR] failed to find the kubernetes cluster: %s", err)
	}

	log.Printf("[INFO] retrieving the kubernetes cluster pool %s", d.Id())
	respPool, err := apiClient.GetKubernetesClusterPool(clusterID, d.Id())
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("[ERR] failed to find the kubernetes cluster pool: %s", err)
	}

	d.SetId(respPool.ID)
	d.Set("cluster_id", resp.ID)
	d.Set("node_count", respPool.Count)
	d.Set("size", respPool.Size)

	if respPool.PublicIPNodePool {
		d.Set("public_ip_node_pool", respPool.PublicIPNodePool)
	}

	poolInstanceNames := make([]string, 0)
	poolInstanceNames = append(poolInstanceNames, respPool.InstanceNames...)

	d.Set("instance_names", poolInstanceNames)

	if len(respPool.Labels) > 0 {
		d.Set("labels", respPool.Labels)
	}

	if len(respPool.Taints) > 0 {
		taints := make([]map[string]interface{}, 0)
		for _, taint := range respPool.Taints {
			taints = append(taints, map[string]interface{}{
				"key":    taint.Key,
				"value":  taint.Value,
				"effect": string(taint.Effect),
			})
		}
		d.Set("taint", taints)
	}

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
	poolUpdate := &civogo.KubernetesClusterPoolUpdateConfig{
		Region: apiClient.Region,
	}

	getKubernetesCluster, err := apiClient.GetKubernetesCluster(clusterID)
	if err != nil {
		return diag.Errorf("[INFO] error getting kubernetes cluster: %s", clusterID)
	}

	if d.HasChange("node_count") {
		poolUpdate.Count = d.Get("node_count").(int)
	}

	if d.HasChange("labels") {
		if attr, ok := d.GetOk("labels"); ok {
			nodePoolLabels := make(map[string]string)
			for k, v := range attr.(map[string]interface{}) {
				if s, ok := v.(string); ok {
					nodePoolLabels[k] = s
				}
			}
			poolUpdate.Labels = nodePoolLabels
		} else {
			poolUpdate.Labels = nil
		}
	}

	if d.HasChange("taint") {
		if attr, ok := d.GetOk("taint"); ok {
			log.Printf("[INFO] dentro de tains")
			nodePoolTains := []corev1.Taint{}
			for _, v := range attr.(*schema.Set).List() {
				taint := v.(map[string]interface{})
				nodePoolTains = append(nodePoolTains, corev1.Taint{
					Key:    taint["key"].(string),
					Value:  taint["value"].(string),
					Effect: corev1.TaintEffect(taint["effect"].(string)),
				})
			}
			poolUpdate.Taints = nodePoolTains
		} else {
			poolUpdate.Taints = []corev1.Taint{}
		}
	}

	log.Printf("[INFO] updating the kubernetes cluster pool %s", d.Id())
	_, err = apiClient.UpdateKubernetesClusterPool(getKubernetesCluster.ID, d.Id(), poolUpdate)
	if err != nil {
		return diag.Errorf("[ERR] failed to update kubernetes cluster pool: %s", err)
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

	log.Printf("[INFO] deleting the kubernetes cluster %s", d.Id())
	_, err = apiClient.DeleteKubernetesClusterPool(getKubernetesCluster.ID, d.Id())
	if err != nil {
		return diag.Errorf("[INFO] an error occurred while trying to delete the kubernetes cluster pool %s", err)
	}

	// Add retry logic here to delete the node pool
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete)-time.Minute, func() *retry.RetryError {
		_, err := apiClient.GetKubernetesClusterPool(getKubernetesCluster.ID, d.Id())
		if err != nil {
			if errors.Is(err, civogo.DatabaseClusterPoolNotFoundError) {
				log.Printf("[INFO] kubernetes node pool %s deleted", d.Id())
				return nil
			}
			log.Printf("[INFO] error trying to read kubernetes cluster pool: %s", err)
			return retry.NonRetryableError(fmt.Errorf("error waiting for Kubernetes node pool to be deleted: %s", err))
		}
		log.Printf("[INFO] kubernetes node pool %s still exists", d.Id())
		return retry.RetryableError(fmt.Errorf("kubernetes node pool still exists"))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
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
		respPool, err := apiClient.GetKubernetesClusterPool(clusterID, nodePoolID)
		if err != nil {
			continue
		}

		if respPool.ID == nodePoolID {
			poolFound = true
			d.SetId(respPool.ID)
			d.Set("cluster_id", clusterID)
			d.Set("label", respPool.ID)
			d.Set("node_count", respPool.Count)
			d.Set("size", respPool.Size)
			d.Set("region", currentRegionCode)
			if respPool.PublicIPNodePool {
				d.Set("public_ip_node_pool", respPool.PublicIPNodePool)
			}
		}
	}

	if !poolFound {
		return nil, fmt.Errorf("[ERR] Node pool %s not found", nodePoolID)
	}

	return []*schema.ResourceData{d}, nil
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
