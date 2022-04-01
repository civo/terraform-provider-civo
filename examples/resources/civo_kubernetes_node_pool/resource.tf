# Query xsmall instance size
data "civo_instances_size" "xsmall" {
    filter {
        key = "type"
        values = ["kubernetes"]
    }

    sort {
        key = "ram"
        direction = "asc"
    }
}

# Create a cluster
resource "civo_kubernetes_cluster" "my-cluster" {
    name = "my-cluster"
    applications = "Portainer,Linkerd:Linkerd & Jaeger"
    firewall_id = civo_firewall.my-firewall.id
    pools {
        size = element(data.civo_instances_size.xsmall.sizes, 0).name
        node_count = 3
    }
}

# Add a node pool
resource "civo_kubernetes_node_pool" "front-end" {
   cluster_id = civo_kubernetes_cluster.my-cluster.id
   node_count = 1 // Optional
   size = element(data.civo_instances_size.xsmall.sizes, 0).name // Optional
   region = "LON1"
}
