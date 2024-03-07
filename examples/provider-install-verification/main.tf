terraform {
  required_providers {
    marqo = {
      source = "registry.terraform.io/marqo/marqo"
    }
  }
}

provider "marqo" {
    host = "https://api.marqo.ai"
    api_key = "mXU94cDBGV32u+Ha8A/DKw9sUl6ldz8uce8E9JmhpiCniq3pjhWdAxlIv6Iog8eU"
}

data "marqo_indices" "example" {}
