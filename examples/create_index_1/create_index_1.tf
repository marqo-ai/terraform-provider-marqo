terraform {
  required_providers {
    marqo = {
      #source = "registry.terraform.io/marqo/marqo"
      source = "registry.terraform.io/marqo/marqo"
    }
  }
}

provider "marqo" {
    host = "https://api.marqo.ai/api/v2"
    api_key = var.marqo_api_key
}

resource "marqo_index" "example" {
  index_name = "example_index_1"
  settings = {
    type = "unstructured"
    vector_numeric_type = "float"
    treat_urls_and_pointers_as_images = true
    model = "open_clip/ViT-L-14/laion2b_s32b_b82k"
    normalize_embeddings = true
    inference_type = "marqo.CPU.small"
    all_fields = []
    number_of_inferences = 1
    number_of_replicas = 0
    number_of_shards = 1
    storage_class = "marqo.basic"
    text_preprocessing = {
      split_length = 2
      split_method  = "sentence"
      split_overlap = 0
    }
    image_preprocessing = {} 
    ann_parameters = {
      space_type = "prenormalized-angular"
      parameters = {
        ef_construction = 512
        m = 16
      }
    }
    filter_string_max_length = 20
  }
}

output "created_index" {
  value = marqo_index.example
}

variable "marqo_api_key" {
  type = string
  description = "Marqo API key"
}