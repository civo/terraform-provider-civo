package kubernetes

import (
	"fmt"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/terraform-provider-civo/internal/datalist"
	"github.com/civo/terraform-provider-civo/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceKubernetesClusters Data source to get and filter all Kubernetes clusters in a region
func DataSourceKubernetesClusters() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		Description: strings.Join([]string{
			"Get information on Kubernetes clusters for use in other resources, with the ability to filter and sort the results. If no filters are specified, all clusters in the region will be returned.",
			"Note: You can use the `civo_kubernetes_cluster` data source to obtain metadata about a single cluster, including its kubeconfig, if you already know the id or name.",
		}, "\n\n"),
		RecordSchema: kubernetesClustersSchema(),
		ExtraQuerySchema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If used, all clusters will be from the provided region",
			},
		},
		ResultAttributeName: "clusters",
		FlattenRecord:       flattenDataSourceKubernetesClusters,
		GetRecords:          getDataSourceKubernetesClusters,
	}

	return datalist.NewResource(dataListConfig)
}

// kubernetesClusterRecord pairs a cluster with the region it was listed from,
// because the API does not return a region per cluster
type kubernetesClusterRecord struct {
	cluster civogo.KubernetesCluster
	region  string
}

func getDataSourceKubernetesClusters(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	apiClient := m.(*civogo.Client)

	region, ok := extra["region"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to find `region` key from query data")
	}

	if region != "" {
		var err error
		apiClient, err = utils.RegionalClient(apiClient, utils.WithRegion(region))
		if err != nil {
			return nil, err
		}
	}

	// ListKubernetesClusters walks every page of results, so large accounts are not truncated
	clusters, err := apiClient.ListKubernetesClusters()
	if err != nil {
		return nil, fmt.Errorf("[ERR] error retrieving kubernetes clusters: %s", err)
	}

	// With no region on the data source or the provider, the API answers from the
	// account's default region, so look that up to report it on each cluster. The
	// region is informational, so a failed lookup leaves it empty instead of failing the read.
	listedRegion := apiClient.Region
	if listedRegion == "" && len(clusters.Items) > 0 {
		if defaultRegion, err := apiClient.GetDefaultRegion(); err != nil {
			log.Printf("[WARN] could not look up the default region for civo_kubernetes_clusters: %s", err)
		} else {
			listedRegion = defaultRegion.Code
		}
	}

	records := make([]interface{}, 0, len(clusters.Items))
	for _, cluster := range clusters.Items {
		records = append(records, kubernetesClusterRecord{cluster: cluster, region: listedRegion})
	}

	return records, nil
}

func flattenDataSourceKubernetesClusters(record, _ interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
	r := record.(kubernetesClusterRecord)
	cluster := r.cluster

	// datalist filters set fields as a *schema.Set, so tags can't be passed as a plain []string
	tags := schema.NewSet(schema.HashString, nil)
	for _, tag := range cluster.Tags {
		tags.Add(tag)
	}

	flattened := map[string]interface{}{
		"id":                     cluster.ID,
		"name":                   cluster.Name,
		"region":                 r.region,
		"cluster_type":           cluster.ClusterType,
		"kubernetes_version":     cluster.KubernetesVersion,
		"cni":                    cluster.CNIPlugin,
		"volume_type":            cluster.VolumeType,
		"tags":                   tags,
		"status":                 cluster.Status,
		"ready":                  cluster.Ready,
		"api_endpoint":           cluster.APIEndPoint,
		"master_ip":              cluster.MasterIP,
		"dns_entry":              cluster.DNSEntry,
		"network_id":             cluster.NetworkID,
		"firewall_id":            cluster.FirewallID,
		"created_at":             cluster.CreatedAt.UTC().String(),
		"pools":                  flattenDataSourceNodePool(&cluster),
		"installed_applications": flattenInstalledApplication(cluster.InstalledApplications),
	}

	return flattened, nil
}

// kubernetesClustersSchema mirrors the civo_kubernetes_cluster data source, minus the kubeconfig
// (a credential, only needed for one cluster), the deprecated fields, and `applications`, which
// that data source never sets (installed_applications carries the same information)
func kubernetesClustersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "The ID of the Kubernetes cluster",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the Kubernetes cluster",
		},
		"region": {
			Type:        schema.TypeString,
			Description: "The region where the cluster is running",
		},
		"cluster_type": {
			Type:        schema.TypeString,
			Description: "The type of the cluster, `k3s` or `talos`",
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Description: "The version of Kubernetes",
		},
		"cni": {
			Type:        schema.TypeString,
			Description: "The CNI plugin used by the cluster",
		},
		"volume_type": {
			Type:        schema.TypeString,
			Description: "The volume type used for the Kubernetes nodes",
		},
		"tags": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "A list of tags",
		},
		"status": {
			Type:        schema.TypeString,
			Description: "The status of the Kubernetes cluster",
		},
		"ready": {
			Type:        schema.TypeBool,
			Description: "If the Kubernetes cluster is ready",
		},
		"api_endpoint": {
			Type:        schema.TypeString,
			Description: "The base URL of the API server on the Kubernetes master node",
		},
		"master_ip": {
			Type:        schema.TypeString,
			Description: "The IP of the Kubernetes master node",
		},
		"dns_entry": {
			Type:        schema.TypeString,
			Description: "The unique DNS entry for the cluster, which points to the master",
		},
		"network_id": {
			Type:        schema.TypeString,
			Description: "The ID of the network the cluster is in",
		},
		"firewall_id": {
			Type:        schema.TypeString,
			Description: "The ID of the firewall used by the cluster",
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "The date the Kubernetes cluster was created",
		},
		"pools": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: nodePoolSchema(false),
			},
			Description: "The node pools of the cluster",
		},
		"installed_applications": installedApplicationsSchema(),
	}
}

func installedApplicationsSchema() *schema.Schema {
	s := dataSourceApplicationSchema()
	s.Description = "The marketplace applications installed on the cluster"
	return s
}
