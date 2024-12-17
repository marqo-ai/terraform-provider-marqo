terraform {
  required_providers {
    marqo = {
      source = "marqo/marqo"
      version = "1.2.0"
    }
  }
}

provider "marqo" {
  host    = "https://api.marqo.ai/api/v2"
  api_key = var.marqo_api_key
}

resource "marqo_index" "example" {
  index_name = "example_index_languagebind"
  settings = {
    type                = "structured"
    vector_numeric_type = "float"
    all_fields = [
      { "name" : "textField", "type" : "text", "features" : ["lexical_search"] },
      { "name" : "videoField", "type" : "video_pointer" },
      { "name" : "audioField", "type" : "audio_pointer" },
      { "name" : "imageField", "type" : "image_pointer" },
      {
        "name" : "multimodal_field",
        "type" : "multimodal_combination",
        "dependent_fields" : {
          "imageField" : 0.8,
          "textField" : 0.1,
          "videoField" : 0.1,
          "audioField" : 0.1
        },
      },
    ],
    number_of_inferences = 1
    storage_class        = "marqo.basic"
    number_of_replicas   = 0
    number_of_shards     = 1
    tensor_fields        = ["multimodal_field", "textField", "videoField", "audioField", "imageField"],
    model                = "LanguageBind/Video_V1.5_FT_Audio_FT_Image"
    normalize_embeddings = true
    inference_type       = "marqo.GPU"
    text_preprocessing = {
      split_length  = 2
      split_method  = "sentence"
      split_overlap = 0
    }
    image_preprocessing = {
      patch_method = null
    }
    video_preprocessing = {
      split_length  = 5
      split_overlap = 1
    }
    audio_preprocessing = {
      split_length  = 5
      split_overlap = 1
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
