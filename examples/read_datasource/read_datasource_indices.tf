terraform {
  required_providers {
    marqo = {
      source = "registry.terraform.io/marqo/marqo"
    }
  }
}

provider "marqo" {
    host = "https://api.marqo.ai/api/v2"
    api_key = "/xDdPdfdVdmuKv0Xc3F9pzOFUmKDmNOVsqQwHPILbAb4dzNHlrjrzk0bsZl7+DFw"
}

data "marqo_read_indices" "example" {}

output "demo_indices" {
  value = data.marqo_read_indices.example
}