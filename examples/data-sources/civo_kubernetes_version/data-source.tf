data "civo_kubernetes_version" "talos" {
  filter {
    key    = "type"
    values = ["talos"]
  }
}

data "civo_kubernetes_version" "k3s" {
  filter {
    key    = "type"
    values = ["k3s"]
  }
}
