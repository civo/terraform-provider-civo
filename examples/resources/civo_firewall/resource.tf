# Create a network
resource "civo_network" "custom_net" {
    label = "my-custom-network"
}

# Create a firewall
resource "civo_firewall" "www" {
  name = "www"
  network_id = civo_network.custom_net.id
}
