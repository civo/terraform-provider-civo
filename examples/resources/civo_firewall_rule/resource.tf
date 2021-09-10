# Query small instance size
data "civo_instances_size" "small" {
    filter {
        key = "name"
        values = ["g3.small"]
        match_by = "re"
    }

    filter {
        key = "type"
        values = ["instance"]
    }

}

# Query instance template
data "civo_template" "debian" {
   filter {
        key = "name"
        values = ["debian-10"]
   }
}

# Create a new instance
resource "civo_instance" "foo" {
    hostname = "foo.com"
    size = element(data.civo_instances_size.small.sizes, 0).name
    template = element(data.civo_template.debian.templates, 0).id
}

# Create a network
resource "civo_network" "custom_net" {
    label = "my-custom-network"
}

# Create a firewall
resource "civo_firewall" "custom_firewall" {
  name = "my-custom-firewall"
  network_id = civo_network.custom_net.id
}

# Create a firewall rule and only allow
# connections from instance we created above
resource "civo_firewall_rule" "custom_port" {
    firewall_id = civo_firewall.custom_firewall.id
    protocol = "tcp"
    start_port = "3000"
    end_port = "3000"
    cidr = [format("%s/%s",civo_instance.foo.public_ip,"32")]
    direction = "ingress"
    label = "custom-application"
    depends_on = [civo_firewall.custom_firewall]
}
