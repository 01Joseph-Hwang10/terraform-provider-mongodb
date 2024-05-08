terraform {
  required_providers {
    mongodb = {
      source = "01Joseph-Hwang10/mongodb"
    }
  }
}

provider "mongodb" {
  // You should include valid username and password whose roles have the necessary permissions
  // for the operations you want to perform in the connection string
  //
  // Also, you should attach the options as a query string to the connection string
  // if you want to use it
  uri = "<your-mongodb-connection-string>"
}
