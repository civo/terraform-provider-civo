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

# Create a firewall
resource "civo_firewall" "my-firewall" {
    name = "my-firewall"
}

# Create a firewall rule
resource "civo_firewall_rule" "kubernetes" {
    firewall_id = civo_firewall.my-firewall.id
    protocol = "tcp"
    start_port = "6443"
    end_port = "6443"
    cidr = ["0.0.0.0/0"]
    direction = "ingress"
    label = "kubernetes-api-server"
}

# Create a cluster
resource "civo_kubernetes_cluster" "my-cluster" {
    name = "my-cluster"
    applications = "Portainer,Linkerd:Linkerd & Jaeger"
    firewall_id = civo_firewall.my-firewall.id
    node_pool {
        label = "front-end" // Optional
        size = element(data.civo_size.xsmall.sizes, 0).name
        node_count = 3
    }
}
