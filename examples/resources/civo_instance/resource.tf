
provider "civo" {

  region = "LON1"
}

resource "civo_firewall" "example" {
    name = "example-firewall"
    create_default_rules = true
    network_id = civo_network.example.id

}

resource "civo_network" "example" {
  label = "example-network3"

}

resource "civo_instance" "example" {
    hostname = "example-instance"
    tags = ["nginx"]
    notes = "Created with TF"
    size = "g3.xsmall" # List on CLI: civo instances size
    network_id = civo_network.example.id
    firewall_id = civo_firewall.example.id
    disk_image = "debian-11" # List on CLI: civo diskimage ls
}

# View instances on dashboard: https://dashboard.civo.com/instances
# View instances on CLI: civo instance ls
