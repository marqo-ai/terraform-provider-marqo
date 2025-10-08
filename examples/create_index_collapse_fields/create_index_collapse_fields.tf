terraform {
  required_providers {
    marqo = {
      source = "registry.terraform.io/marqo/marqo"
    }
  }
}

provider "marqo" {
  host    = "https://api.marqo.ai/api/v2"
  api_key = var.marqo_api_key
}

resource "marqo_index" "c_f_test" {
  index_name = "%s"
  timeouts = {
    create = "60m"
    update = "60m"
    delete = "45m"
  }
  settings = {
    type                 = "unstructured"
    model                = "open_clip/ViT-L-14/laion2b_s32b_b82k"
    inference_type       = "marqo.CPU.large"
    number_of_inferences = 1
    number_of_replicas   = 0
    number_of_shards     = 1
    storage_class        = "marqo.basic"
    collapse_fields = [
      {
        name       = "product_family_2"
        min_groups = 200
      }
    ]
  }
}

output "created_index" {
  value = marqo_index.example
}

variable "marqo_api_key" {
  type        = string
  description = "Marqo API key"
}