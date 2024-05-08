data "mongodb_database" "default" {
  name = "default"
}

data "mongodb_database_collection" "users" {
  database = data.mongodb_database.default.name
  name     = "users"
}


data "mongodb_database_document" "first_user" {
  database    = data.mongodb_database.default.name
  collection  = data.mongodb_database_collection.users.name
  document_id = "<stringified-mongodb-object-id>"
}

// Example usage of the data source: writing the document to a local file
resource "local_file" "first_user_document" {
  content  = data.mongodb_database_document.first_user.document
  filename = "${path.module}/first-user.json"
}
