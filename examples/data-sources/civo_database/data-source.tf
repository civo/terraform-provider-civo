data "civo_database" "test" {
    name = "test-database"
    region = "LON1"
    size = element(data.civo_database.size.small.sizes, 0).name
    nodes = 2
}