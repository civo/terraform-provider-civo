data "civo_size" "small" {
    filter {
        key = "name"
        values = ["db.small"]
        match_by = "re"
    }
    filter {
        key = "type"
        values = ["database"]
    }
}

resource "civo_database" "custom_database" {
    name = "custom_database"
    size = element(data.civo_size.small.sizes, 0).name
    nodes = 2
}