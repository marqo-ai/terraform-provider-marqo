---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "marqo_read_indices Data Source - terraform-provider-marqo"
subcategory: ""
description: |-
  
---

# marqo_read_indices (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The unique identifier for the resource.

### Read-Only

- `items` (Attributes List) (see [below for nested schema](#nestedatt--items))
- `last_updated` (String) The last time the resource was updated.

<a id="nestedatt--items"></a>
### Nested Schema for `items`

Optional:

- `all_fields` (Attributes List) (see [below for nested schema](#nestedatt--items--all_fields))
- `tensor_fields` (List of String)

Read-Only:

- `ann_parameters` (Attributes) (see [below for nested schema](#nestedatt--items--ann_parameters))
- `created` (String) The creation date of the index
- `docs_count` (String) The number of documents in the index
- `docs_deleted` (String) The number of documents deleted from the index
- `filter_string_max_length` (String) The filter string max length
- `image_preprocessing` (Attributes) (see [below for nested schema](#nestedatt--items--image_preprocessing))
- `index_name` (String) The name of the index
- `index_status` (String) The status of the index
- `inference_type` (String) The type of inference used by the index
- `marqo_endpoint` (String) The Marqo endpoint used by the index
- `marqo_version` (String) The version of Marqo used by the index
- `model` (String) The model used by the index
- `normalize_embeddings` (Boolean) Indicates if embeddings should be normalized
- `number_of_inferences` (String) The number of inferences made by the index
- `number_of_replicas` (String) The number of replicas for the index
- `number_of_shards` (String) The number of shards for the index
- `search_query_total` (String) The total number of search queries made on the index
- `storage_class` (String) The storage class of the index
- `store_size` (String) The size of the index storage
- `text_preprocessing` (Attributes) (see [below for nested schema](#nestedatt--items--text_preprocessing))
- `treat_urls_and_pointers_as_images` (Boolean) Indicates if URLs and pointers should be treated as images
- `type` (String) The type of the index
- `vector_numeric_type` (String) The numeric type of the vector

<a id="nestedatt--items--all_fields"></a>
### Nested Schema for `items.all_fields`

Optional:

- `dependent_fields` (Map of Number)
- `features` (List of String)
- `name` (String)
- `type` (String)


<a id="nestedatt--items--ann_parameters"></a>
### Nested Schema for `items.ann_parameters`

Read-Only:

- `parameters` (Attributes) (see [below for nested schema](#nestedatt--items--ann_parameters--parameters))
- `space_type` (String) The space type for ANN parameters

<a id="nestedatt--items--ann_parameters--parameters"></a>
### Nested Schema for `items.ann_parameters.parameters`

Read-Only:

- `ef_construction` (String) The efConstruction parameter for ANN
- `m` (String) The m parameter for ANN



<a id="nestedatt--items--image_preprocessing"></a>
### Nested Schema for `items.image_preprocessing`

Read-Only:

- `patch_method` (String) The patch method for image preprocessing


<a id="nestedatt--items--text_preprocessing"></a>
### Nested Schema for `items.text_preprocessing`

Read-Only:

- `split_length` (String) The split length for text preprocessing
- `split_method` (String) The split method for text preprocessing
- `split_overlap` (String) The split overlap for text preprocessing
