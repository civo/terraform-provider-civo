# Get network
data "civo_network" "default_network" {
    label = "Default"
}

# Create volume
resource "civo_volume" "db" {
    name = "backup-data"
    size_gb = 5
    network_id = data.civo_network.default_network.id
    depends_on = [
      data.civo_network.default_network
    ]
}
