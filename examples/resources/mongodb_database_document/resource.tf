resource "mongodb_database" "default" {
  name          = "default"
  force_destroy = false
}

resource "mongodb_database_collection" "users" {
  database      = mongodb_database.default.name
  name          = "users"
  force_destroy = false
}

resource "mongodb_database_document" "first_user" {
  database   = mongodb_database.default.name
  collection = mongodb_database_collection.users.name
  document = jsonencode({
    name = "John Doe"
    age  = 25
  })
}
