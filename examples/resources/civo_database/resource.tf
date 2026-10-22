data "civo_size" "small" {
  filter {
    key      = "name"
    values   = ["db.small"]
    match_by = "re"
  }
  filter {
    key    = "type"
    values = ["database"]
  }
}

data "civo_database_version" "mysql" {
  filter {
    key    = "engine"
    values = ["mysql"]
  }
}

resource "civo_database" "custom_database" {
  name    = "custom_database"
  size    = element(data.civo_size.small.sizes, 0).name
  nodes   = 2
  engine  = element(data.civo_database_version.mysql.versions, 0).engine
  version = element(data.civo_database_version.mysql.versions, 0).version
}

# Grow the storage beyond the plan, enable extensions and set server parameters.
# Storage can only be increased; the GB above the plan are billed as an add-on.
resource "civo_database" "tuned_database" {
  name       = "tuned_database"
  size       = "g3.db.small"
  nodes      = 3
  engine     = "PostgreSQL"
  version    = "17"
  storage_gb = 120
  extensions = ["pg_cron", "postgis"]
  database_parameters = {
    max_connections = "200"
  }
}

# Fork an existing database as it was at a point in time. engine and version
# must match the source; size and nodes may differ.
resource "civo_database" "forked_database" {
  name               = "forked_database"
  size               = civo_database.tuned_database.size
  nodes              = 1
  engine             = civo_database.tuned_database.engine
  version            = civo_database.tuned_database.version
  source_database_id = civo_database.tuned_database.id
  source_target_time = "2026-10-21T10:00:00Z"
}
