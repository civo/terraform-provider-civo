data "civo_network" "test" {
    label = "test-network"
    region = "LON1"
    cidr_v4        = "10.0.0.0/24"
    nameservers_v4 = ["8.8.8.8", "8.8.4.4", "1.1.1.1"]
}
