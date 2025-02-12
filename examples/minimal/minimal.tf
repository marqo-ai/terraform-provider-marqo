terraform {
  required_providers {
    marqo = {
      source  = "marqo-ai/marqo"
      version = "1.2.1"
    }
  }
}

provider "marqo" {
  host    = "https://controller.marqo-staging.com/api/v2"
  api_key = var.marqo_api_key
}

resource "marqo_index" "example" {
  index_name = "optional1"
  settings = {
    type                 = "unstructured"
    model                = "open_clip/ViT-L-14/laion2b_s32b_b82k"
    inference_type       = "marqo.GPU"
    number_of_inferences = 2
    number_of_replicas   = 0
    number_of_shards     = 1
    storage_class        = "marqo.basic"
  }
}

output "created_index" {
  value = marqo_index.example
}

variable "marqo_api_key" {
  type        = string
  description = "Marqo API key"
}
