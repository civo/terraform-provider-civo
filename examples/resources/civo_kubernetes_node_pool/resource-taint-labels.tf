# Query xsmall instance size
data "civo_size" "xsmall" {
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
        size = element(data.civo_size.xsmall.sizes, 0).name
        node_count = 3
    }
}

# Add a node pool
resource "civo_kubernetes_node_pool" "back-end" {
   cluster_id = civo_kubernetes_cluster.my-cluster.id
   label = "back-end" // Optional
   node_count = 1 // Optional
   size = element(data.civo_size.xsmall.sizes, 0).name // Optional
   region = "LON1"

   labels = {
     service  = "backend"
     priority = "high"
   }

  taint {
    key    = "workloadKind"
    value  = "database"
    effect = "NoSchedule"
  }

  taint {
    key    = "workloadKind"
    value  = "frontend"
    effect = "NoSchedule"
  }
}
