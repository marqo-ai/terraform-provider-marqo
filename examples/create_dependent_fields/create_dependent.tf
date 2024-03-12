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

resource "marqo_index" "example" {
  index_name = "example_index_dependent"
  settings = {
    type = "structured"
    vector_numeric_type = "float"
    all_fields = [
        {"name": "text_field", "type": "text", "features": ["lexical_search"]},
        {"name": "image_field", "type": "image_pointer"},
        {
            "name": "multimodal_field",
            "type": "multimodal_combination",
            "dependent_fields": {
                "image_field": 0.8,
                "text_field": 0.1
            },
        },
    ],
    tensor_fields = ["multimodal_field"],
    treat_urls_and_pointers_as_images = true
    model = "open_clip/ViT-L-14/laion2b_s32b_b82k"
    normalize_embeddings = true
    inference_type = "marqo.CPU.small"
    text_preprocessing = {
      split_length = 2
      split_method = "sentence"
      split_overlap = 0
    }
    image_preprocessing = {
      patch_method = null
    }
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