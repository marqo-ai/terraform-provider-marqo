terraform {
  required_providers {
    marqo = {
      source = "registry.terraform.io/marqo/marqo"
    }
  }
}

provider "marqo" {}

data "marqo_indices" "example" {}
