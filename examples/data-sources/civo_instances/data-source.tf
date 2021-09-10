data "civo_instances" "small-size" {
    region = "NYC1"
    filter {
        key = "size"
        values = [g3.small]
    }
}
