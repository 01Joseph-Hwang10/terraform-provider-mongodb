resource "mongodb_database" "default" {
  name          = "default"
  force_destroy = false
}

resource "mongodb_database_collection" "users" {
  database      = mongodb_database.default.name
  name          = "users"
  force_destroy = false
}
