data "civo_size" "kfsmall" {
    filter {
        key = "name"
        values = ["g3.kf.small"]
    }

    filter {
        key = "type"
        values = ["kfcluster"]
    }
}