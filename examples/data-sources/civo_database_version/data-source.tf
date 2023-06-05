# Example for mysql
data civo_database_version "mysql" {
  filter {
        key = "engine"
        values = ["mysql"]
    }
}

# Example for postgresql
data civo_database_version "mysql" {
  filter {
        key = "engine"
        values = ["postgresql"]
    }
}

# To use this data source, make sure you have a database cluster created.
resource "civo_database" "custom_database" {
    name = "custom_database"
    size = element(data.civo_size.small.sizes, 0).name
    nodes = 2
    engine = element(data.civo_database_version.mysql.versions, 0).engine
    version = element(data.civo_database_version.mysql.versions, 0).version
}