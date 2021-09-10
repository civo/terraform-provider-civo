data "civo_kubernetes_version" "stable" {
    filter {
        key = "type"
        values = ["stable"]
    }
}
