terraform {
  required_providers {
    marqo-terraform = {
      source = "registry.terraform.io/marqo/marqo-terraform"
    }
  }
}

provider "marqo-terraform" {
  host    = "https://api.marqo.ai/api/v2"
  api_key = var.marqo_api_key
}

resource "marqo-terraform_index" "example" {
  index_name = "example_index_dependent"
  settings = {
    type                = "structured"
    vector_numeric_type = "float"
    all_fields = [
      { "name" : "text_field", "type" : "text", "features" : ["lexical_search"] },
      { "name" : "image_field", "type" : "image_pointer" },
      {
        "name" : "multimodal_field",
        "type" : "multimodal_combination",
        "dependent_fields" : {
          "image_field" : 0.8,
          "text_field" : 0.1
        },
      },
    ],
    tensor_fields        = ["multimodal_field"],
    model                = "open_clip/ViT-L-14/laion2b_s32b_b82k"
    normalize_embeddings = true
    inference_type       = "marqo.CPU.small"
    text_preprocessing = {
      split_length  = 2
      split_method  = "sentence"
      split_overlap = 0
    }
    image_preprocessing = {
      patch_method = null
    }
    ann_parameters = {
      space_type = "prenormalized-angular"
      parameters = {
        ef_construction = 512
        m               = 16
      }
    }
  }
}

output "created_index" {
  value = marqo-terraform_index.example
}

variable "marqo_api_key" {
  type        = string
  description = "Marqo API key"
}