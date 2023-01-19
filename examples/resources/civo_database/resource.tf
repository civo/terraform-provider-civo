resource "civo_database" "custom_database" {
    name = "test_database"
    size = element(data.civo_database.size.small.sizes, 0).name
    nodes = 2
}