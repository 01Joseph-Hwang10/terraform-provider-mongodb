# Terraform Provider MongoDB

`01Joseph-Hwang10/terraform-provider-mongodb` allows you to manage MongoDB databases, collections, documents, and indexes.

## Quick Example

In this example, we will create a database, a collection, and a document in MongoDB.

First, add the provider to your Terraform configuration:

```terraform
terraform {
  required_providers {
    mongodb = {
      source = "01Joseph-Hwang10/mongodb"
    }
  }
}

provider "mongodb" {
  uri = "<your-mongodb-connection-string>"
}
```

Then, create a database, a collection, and a document:

```terraform
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
```

Finally, run `terraform apply` to create the database, collection, and document in MongoDB.

## API Documentation

See the [API Documentation](./docs/index.md) for more information.

## Contributing

Any contribution is welcome! Check out [CONTRIBUTING.md](https://github.com/01Joseph-Hwang10/terraform-provider-mongodb/blob/main/.github/CONTRIBUTING.md) and [CODE_OF_CONDUCT.md](https://github.com/01Joseph-Hwang10/terraform-provider-mongodb/blob/main/.github/CODE_OF_CONDUCT.md) for more information on how to get started.

## License

`terraform-provider-mongodb` is licensed under a [MIT License](https://github.com/01Joseph-Hwang10/terraform-provider-mongodb/blob/main/LICENSE).