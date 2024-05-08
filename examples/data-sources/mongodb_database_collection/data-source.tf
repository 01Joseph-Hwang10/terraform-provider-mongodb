data "mongodb_database" "default" {
  name = "default"
}

data "mongodb_database_collection" "users" {
  database = data.mongodb_database.default.name
  name     = "users"
}
