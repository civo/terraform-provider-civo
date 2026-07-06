data "civo_network" "test" {
    label = "test-network"
    region = "LON1"
}

output "network_cidr_v4" {
  value = data.civo_network.test.cidr_v4
}

output "network_nameservers_v4" {
  value = data.civo_network.test.nameservers_v4
}
