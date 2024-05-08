resource "mongodb_database" "default" {
  name          = "default"
  force_destroy = false
}

resource "mongodb_database_collection" "users" {
  database      = mongodb_database.default.name
  name          = "users"
  force_destroy = false
}

resource "mongodb_database_index" "user_age_index" {
  database      = mongodb_database.default.name
  collection    = mongodb_database_collection.users.name
  field         = "age"
  direction     = 1
  unique        = false
  force_destroy = false
}
