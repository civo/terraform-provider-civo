package civo

import (
	// "context"
	"fmt"
	"log"
	"strings"

	// "time"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/google/uuid"

	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The region of the node pool, has to match that of the cluster",
				ValidateFunc: validation.StringIsNotEmpty,
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
				Description: "the size of each node (optional, the default is currently g4s.kube.medium)",
			},
		},
		Create: resourceKubernetesClusterNodePoolCreate,
		Read:   resourceKubernetesClusterNodePoolRead,
		Update: resourceKubernetesClusterNodePoolUpdate,
		Delete: resourceKubernetesClusterNodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKubernetesClusterNodePoolImport,
		},
	}
}

// function to create a new cluster
func resourceKubernetesClusterNodePoolCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	cluster_id := d.Get("cluster_id").(string)

	count := 0
	if attr, ok := d.GetOk("num_target_nodes"); ok {
		count = attr.(int)
	} else {
		count = 3
	}

	size := ""
	if attr, ok := d.GetOk("target_nodes_size"); ok {
		size = attr.(string)
	} else {
		size = "g4s.kube.medium"
	}

	log.Printf("[INFO] getting kubernetes cluster %s in the region %s", cluster_id, apiClient.Region)
	getKubernetesCluster, err := apiClient.GetKubernetesCluster(cluster_id)
	if err != nil {
		return err
	}

	newPool := []civogo.KubernetesClusterPoolConfig{}
	for _, v := range getKubernetesCluster.Pools {
		newPool = append(newPool, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
	}

	poolID := uuid.NewString()
	newPool = append(newPool, civogo.KubernetesClusterPoolConfig{ID: poolID, Count: count, Size: size})

	log.Printf("[INFO] configuring kubernetes cluster %s to add pool %s", cluster_id, poolID[:6])
	config := &civogo.KubernetesClusterConfig{
		Pools: newPool,
	}

	log.Printf("[INFO] Creating a new kubernetes cluster pool %s", poolID[:6])
	_, err = apiClient.UpdateKubernetesCluster(cluster_id, config)
	if err != nil {
		return fmt.Errorf("[ERR] failed to create the kubernetes cluster: %s", err)
	}

	d.SetId(poolID)

	return resourceKubernetesClusterNodePoolRead(d, m)

}

// function to read the kubernetes cluster
func resourceKubernetesClusterNodePoolRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)
	cluster_id := d.Get("cluster_id").(string)

	log.Printf("[INFO] retrieving the kubernetes cluster %s", cluster_id)
	resp, err := apiClient.GetKubernetesCluster(cluster_id)
	if err != nil {
		if resp == nil {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERR] failed to find the kubernetes cluster: %s", err)
	}

	d.Set("cluster_id", resp.ID)
	for _, v := range resp.Pools {
		if v.ID == d.Id() {
			d.Set("num_target_nodes", v.Count)
			d.Set("target_nodes_size", v.Size)
		}
	}

	return nil
}

// function to update the kubernetes cluster
func resourceKubernetesClusterNodePoolUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	cluster_id := d.Get("cluster_id").(string)
	count := 0

	if d.HasChange("num_target_nodes") {
		count = d.Get("num_target_nodes").(int)
	}

	getKubernetesCluster, err := apiClient.GetKubernetesCluster(cluster_id)
	if err != nil {
		return err
	}

	nodePool := []civogo.KubernetesClusterPoolConfig{}
	for _, v := range getKubernetesCluster.Pools {
		nodePool = append(nodePool, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
	}

	nodePool = updateNodePool(nodePool, d.Id(), count)

	log.Printf("[INFO] configuring kubernetes cluster %s to add pool", d.Id()[:6])
	config := &civogo.KubernetesClusterConfig{
		Pools: nodePool,
	}

	log.Printf("[INFO] updating the kubernetes cluster %s", d.Id())
	_, err = apiClient.UpdateKubernetesCluster(cluster_id, config)
	if err != nil {
		return fmt.Errorf("[ERR] failed to update kubernetes cluster: %s", err)
	}

	return resourceKubernetesClusterNodePoolRead(d, m)
}

// function to delete the kubernetes cluster
func resourceKubernetesClusterNodePoolDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*civogo.Client)

	cluster_id := d.Get("cluster_id").(string)
	getKubernetesCluster, err := apiClient.GetKubernetesCluster(cluster_id)
	if err != nil {
		return err
	}

	nodePool := []civogo.KubernetesClusterPoolConfig{}
	for _, v := range getKubernetesCluster.Pools {
		nodePool = append(nodePool, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
	}

	nodePool = removeNodePool(nodePool, d.Id())
	config := &civogo.KubernetesClusterConfig{
		Pools: nodePool,
	}

	log.Printf("[INFO] deleting the kubernetes cluster %s", d.Id())
	_, err = apiClient.UpdateKubernetesCluster(cluster_id, config)
	if err != nil {
		return fmt.Errorf("[INFO] an error occurred while tring to delete the kubernetes cluster pool %s", err)
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
		resp, err := apiClient.GetKubernetesCluster(clusterID)
		if err != nil {
			continue // move on and find in another region
		}

		for _, v := range resp.Pools {
			if v.ID == nodePoolID {
				poolFound = true
				d.SetId(v.ID)
				d.Set("cluster_id", resp.ID)
				d.Set("region", currentRegionCode)
				d.Set("num_target_nodes", v.Count)
				d.Set("target_nodes_size", v.Size)
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
	key := 0
	for k, v := range s {
		if strings.Contains(v.ID, id) {
			key = k
			break
		}
	}

	s[len(s)-1], s[key] = s[key], s[len(s)-1]
	return s[:len(s)-1]
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
