resource "github_repository" "terraform_provider_mongodb" {
  name         = "terraform-provider-mongodb"
  description  = "Manage MongoDB databases, collections, documents, and indexes using Terraform"
  homepage_url = "https://registry.terraform.io/providers/01Joseph-Hwang10/mongodb/latest"
  topics = [
    "terraform",
    "terraform-provider",
    "mongodb",
    "database",
    "collection",
    "document",
    "index",
  ]

  visibility = "public"

  has_downloads   = true
  has_issues      = true
  has_projects    = true
  has_wiki        = true
  has_discussions = false

  template {
    include_all_branches = false
    owner                = "hashicorp"
    repository           = "terraform-provider-scaffolding-framework"
  }
}

resource "github_actions_secret" "gpg_private_key" {
  secret_name     = "GPG_PRIVATE_KEY"
  repository      = github_repository.terraform_provider_mongodb.name
  plaintext_value = file(var.gpg_private_key)
}

resource "github_actions_secret" "passphrase" {
  secret_name     = "PASSPHRASE"
  repository      = github_repository.terraform_provider_mongodb.name
  plaintext_value = file(var.passphrase)
}
