# Marqo Terraform Provider

This page will walk you through using OpenTofu to manage your Marqo Cloud resources.

OpenTofu, a fork of Terraform, is an open-source infrastructure as code tool that enables you to safely and predictably create, change, and delete cloud resources. 
For more information on opentofu and terraform, please visit the following links: [OpenTofu](https://opentofu.org) and [Terraform](https://www.terraform.io/).

The marqo opentofu provider is located at [`registry.opentofu.org/marqo-ai/marqo`](https://github.com/opentofu/registry/blob/main/providers/m/marqo-ai/marqo.json).

If you wish to use the terraform provider instead, please note the following
- replace the provider source with `marqo-ai/marqo`
- replace all `tofu` commands in the guide below with `terraform`

Detailed documentation can be found in marqodocs [here](https://docs.marqo.ai/2.11/Cloud-Reference/opentofu_provider/).

---

## Installation Instructions

1. Install Opentofu by following the instructions on the [Opentofu website](https://opentofu.org/docs/intro/install/). Alternatively, install Terraform by following the instructions on the [Terraform website](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli).
2. Create a Marqo configuration (a `.tf` file with your marqo API and endpoint details) in a new directory
3. Run `tofu init` in the directory you created the configuration in to initialize the configuration
4. At this point, you should be able to run `tofu plan` and `tofu apply` to create and manage your Marqo resources.

Some common commands are:
- `tofu plan` - checks whether the configuration is valid and what actions will be taken
- `tofu apply` - creates or updates the resources in your account
- `tofu destroy` - deletes the resources in your account

See the [Opentofu documentation](https://opentofu.org/docs/intro/) for more information on how to use Opentofu.

## Overview of Features

The marqo opentofu provider supports the following:

- A datasource called `marqo_read_indices` that allows you to read all of your marqo indexes in your account.
- A resource called `marqo_index` that allows you to create and manage a marqo index.

## Sample Configuration

For both of the examples below, create a file within each configuration directory named `terraform.tfvars` containing your api key as follows

```python
marqo_api_key = "<KEY>"
```

Note that the host must be set to `"https://api.marqo.ai/api/v2"`

### Reading All Indexes in Your Account (datasource)

```terraform
terraform {
  required_providers {
    marqo = {
      source = "registry.opentofu.org/marqo-ai/marqo"
      version = "1.0.1"
    }
  }
}

provider "marqo" {
  host    = "https://api.marqo.ai/api/v2"
  api_key = var.marqo_api_key
}

data "marqo_read_indices" "example" {
  id = 1
}

output "indices_in_marqo_cloud" {
  value = data.marqo_read_indices.example
}

variable "marqo_api_key" {
  type        = string
  description = "Marqo API key"
}
```

### Creating and Managing a Structured Index (resource)

```terraform
terraform {
  required_providers {
    marqo = {
      source = "registry.opentofu.org/marqo-ai/marqo"
      version = "1.0.1"
    }
  }
}

provider "marqo" {
  host    = "https://api.marqo.ai/api/v2"
  api_key = var.marqo_api_key
}

resource "marqo_index" "example" {
  index_name = "example_index_dependent_2"
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
  value = marqo_index.example
}

variable "marqo_api_key" {
  type        = string
  description = "Marqo API key"
}
```
