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
    num_target_nodes = 1
    target_nodes_size = element(data.civo_instances_size.xsmall.sizes, 0).name
}

# Add a node pool
resource "civo_kubernetes_node_pool" "front-end" {
   cluster_id = civo_kubernetes_cluster.my-cluster.id
   num_target_nodes = 1
   region = "LON1"
}
