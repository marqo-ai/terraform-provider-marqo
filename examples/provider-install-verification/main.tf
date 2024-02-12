terraform {
  required_providers {
    marqo = {
      source = "registry.terraform.io/marqo/marqo"
    }
  }
}

provider "marqo" {
    host = "http://localhost:8080"
    api_key = "your"
}

data "marqo_indices" "example" {}
