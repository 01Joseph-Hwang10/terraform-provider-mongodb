variable "github_owner" {
  type        = string
  description = <<EOT
    Variable for GitHub owner.

    This represents what account or organization the repository will be created under.
  EOT
}

variable "gpg_private_key" {
  type        = string
  description = "Path to the GPG private key file."
}

variable "passphrase" {
  type        = string
  description = "Path to the passphrase file."
}
