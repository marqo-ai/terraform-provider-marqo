terraform {
  required_providers {
    marqo = {
      source = "registry.terraform.io/marqo/marqo"
    }
  }
}

provider "marqo" {
    host = "https://api.marqo.ai/api/v2"
    api_key = ""
}

data "marqo_read_indices" "example" {}

output "demo_indices" {
  value = data.marqo_read_indices.example
}