data "civo_instances" "small-size" {
    region = "LON1"
    filter {
        key = "size"
        values = [g3.small]
    }
}
