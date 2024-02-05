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

# Create a cluster with labels and taints
resource "civo_kubernetes_cluster" "my-cluster" {
    name = "my-cluster"
    applications = "Portainer,Linkerd:Linkerd & Jaeger"
    firewall_id = civo_firewall.my-firewall.id
    
    pools {
        size = element(data.civo_size.xsmall.sizes, 0).name
        node_count = 3

        labels = {
          service  = "backend"
          priority = "high"
        }

        taint {
          key    = "workloadKind"
          value  = "database"
          effect = "NoSchedule"
        }
    }
}

