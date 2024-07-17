terraform {
  required_providers {
    marqo = {
      source = "registry.terraform.io/marqo/marqo"
    }
  }
}

provider "marqo" {
  host    = "https://api.marqo.ai"
  api_key = var.marqo_api_key
}

data "marqo_indices" "example" {}

variable "marqo_api_key" {
  type        = string
  description = "Marqo API key"
}