resource "mongodb_database" "default" {
  name          = "default"
  force_destroy = false
}
