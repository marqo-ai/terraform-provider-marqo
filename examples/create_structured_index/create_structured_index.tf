terraform {
  required_providers {
    marqo = {
      source  = "marqo/marqo"
      version = "1.2.0"
    }
  }
}

provider "marqo" {
  host    = "https://api.marqo.ai/api/v2"
  api_key = var.marqo_api_key
}

resource "marqo_index" "example" {
  index_name = "example_index_dependent_1"
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
          "imageField" : 0.8,
          "textField" : 0.1
        },
      },
    ],
    number_of_inferences = 1
    storage_class        = "marqo.basic"
    number_of_replicas   = 0
    number_of_shards     = 2
    tensor_fields        = ["multimodal_field"],
    model                = "open_clip/ViT-L-14/laion2b_s32b_b82k"
    normalize_embeddings = true
    inference_type       = "marqo.CPU.large"
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
  value = marqo_index.example
}

variable "marqo_api_key" {
  type        = string
  description = "Marqo API key"
}
