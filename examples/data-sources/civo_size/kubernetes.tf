data "civo_size" "kxsmall" {
    filter {
        key = "name"
        values = ["g3.k3s.xsmall"]
    }

    filter {
        key = "type"
        values = ["kubernetes"]
    }
}
