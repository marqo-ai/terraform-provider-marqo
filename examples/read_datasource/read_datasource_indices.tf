terraform {
  required_providers {
    marqo = {
      source = "registry.terraform.io/marqo/marqo"
    }
  }
}

provider "marqo" {
    host = "https://api.marqo.ai/api/v2"
    api_key = var.marqo_api_key
}

data "marqo_read_indices" "example" {
  id = 1
}

output "indices_in_marqo_cloud" {
  value = data.marqo_read_indices.example
}

variable "marqo_api_key" {
  type = string
  description = "Marqo API key"
}