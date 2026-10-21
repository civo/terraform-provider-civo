data "civo_kubernetes_clusters" "lon1" {
  region = "LON1"
  filter {
    key    = "status"
    values = ["ACTIVE"]
  }
}

output "cluster_endpoints" {
  value = {
    for cluster in data.civo_kubernetes_clusters.lon1.clusters :
    cluster.name => cluster.api_endpoint
  }
}
