data "civo_size" "dbxsmall" {
    filter {
        key = "name"
        values = ["g3.db.xsmall"]
    }

    filter {
        key = "type"
        values = ["database"]
    }
}