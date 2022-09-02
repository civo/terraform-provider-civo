# Create a network
resource "civo_network" "custom_net" {
  label = "my-custom-network"
}

# Create a firewall
resource "civo_firewall" "www" {
  name       = "www"
  network_id = civo_network.custom_net.id
}

# Create a firewall with the default rules
resource "civo_firewall" "www" {
  name                 = "www"
  network_id           = civo_network.custom_net.id
  create_default_rules = true
}

# Create a firewall withouth the default rules but with a custom rule
resource "civo_firewall" "www" {
  name                 = "www"
  network_id           = civo_network.custom_net.id
  create_default_rules = false
  ingress_rule {
    label      = "k8s"
    protocol   = "tcp"
    port_range = "6443"
    cidr       = ["192.168.1.1/32", "192.168.10.4/32", "192.168.10.10/32"]
    action     = "allow"
  }

  ingress_rule {
    label      = "ssh"
    protocol   = "tcp"
    port_range = "22"
    cidr       = ["192.168.1.1/32", "192.168.10.4/32", "192.168.10.10/32"]
    action     = "allow"
  }

  egress_rule {
    label      = "all"
    protocol   = "tcp"
    port_range = "1-65535"
    cidr       = ["0.0.0.0/0"]
    action     = "allow"
  }
}
