# Create Structured Index

Structured indexes in Marqo are tailored for datasets with a defined schema and are particularly effective for complex queries like sorting, grouping, and filtering. They are designed for fast, in-memory operations.

By default, the settings for structured index look like this. Settings can be set as the index is created.

***

```http
POST /indexes/{index_name}
```

Create index with (optional) settings.
This endpoint accepts the `application/json` content type.

## Path parameters

| Name                                        | Type   | Description       |
|:--------------------------------------------|:-------|:------------------|
| <div class="no-wrap">**`index_name`**</div> | String | name of the index |

## Body Parameters

The settings for the index are represented as a nested JSON object that contains the default settings for the index. The
parameters are as follows:

| Name                                                 | Type       | Default value                  | Description                                                                                                        |
|:-----------------------------------------------------|:-----------|:-------------------------------|:-------------------------------------------------------------------------------------------------------------------|
| **`allFields`**                                      | List       | `-`                            | list of fields that might be indexed or queried                                                                    |
| **`tensorFields`**                                   | List       | `[]`                           | list of fields that are treated as tensors                                                                         |
| **`model`**                                          | String     | `hf/e5-base-v2` | The model to use to vectorise doc content in `add_documents()` calls for the index                                 |
| **`modelProperties`**                                | Dictionary | `""`                           | The model properties object corresponding to `model` (for custom models)                                           |
| <div class="no-wrap">**`normalizeEmbeddings`**</div> | Boolean    | `true`                         | Normalize the embeddings to have unit length                                                                       |
| **`textPreprocessing`**                              | Dictionary | `""`                           | The text preprocessing object                                                                                      |
| **`imagePreprocessing`**                             | Dictionary | `""`                           | The image preprocessing object                                                                                     |
| **`annParameters`**                                  | Dictionary | `""`                           | The ANN algorithm parameter object                                                                                 |
| **`type`**                                           | String     | `-`                            | Type of the index. The default value is `unstructured`, but for the structured index this needs to be `structured` |
| **`vectorNumericType`**                              | String     | `float`                     | Numeric type for vector encoding                                                                                   |

## Text Preprocessing Object

The `textPreprocessing` object contains the specifics of how you want the index to preprocess text. The parameters are
as follows:

| Name                                          | Type    | Default value | Description                                                                                                                                                                                                                |
|:----------------------------------------------|:--------|:--------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **`splitLength`**                             | Integer | `2`           | The length of the chunks after splitting by split_method                                                                                                                                                                   |
| <div class="no-wrap">**`splitOverlap`**</div> | Integer | `0`           | The length of overlap between adjacent chunks                                                                                                                                                                              |
| **`splitMethod`**                             | String  | `sentence`    | The method by which text is chunked (`character`, `word`, `sentence`, or `passage`)                                                                                                                                        |
| **`overrideTextChunkPrefix`**                 | String  | `null`        | A string to be added to the start of all text chunks in documents before vectorisation. Only affects vectors generated. Text itself will not be stored. Overrides `modelProperties`-level prefix.                         |
| **`overrideTextQueryPrefix`**                 | String  | `null`        | A string to be added to the start of all search text queries before vectorisation. Only affects vectors generated. Text itself will not be returned or used for lexical search. Overrides `modelProperties`-level prefix. |

## Image Preprocessing Object

The `imagePreprocessing` object contains the specifics of how you want the index to preprocess images. The parameters
are as follows:

| Name                                         | Type   | Default value | Description                                                  |
|:---------------------------------------------|:-------|:--------------|:-------------------------------------------------------------|
| <div class="no-wrap">**`patchMethod`**</div> | String | `null`        | The method by which images are chunked (`simple` or `frcnn`) |

## ANN Algorithm Parameter object

The `annParameters` object contains hyperparameters for the approximate nearest neighbour algorithm used for tensor
storage within Marqo. The parameters are as follows:

| Name                                       | Type   | Default value | Description                                                                                                |
|:-------------------------------------------|:-------|:--------------|:-----------------------------------------------------------------------------------------------------------|
| <div class="no-wrap">**`spaceType`**</div> | String | `prenormalized-anglar` | The function used to measure the distance between two points in ANN (`l1`, `l2`, `linf`, or `prenormalized-anglar`. |
| **`parameters`**                           | Dict   | `""`          | The hyperparameters for the ANN method (which is always `hnsw` for Marqo).                                 |

## HNSW Method Parameters Object

`parameters` can have the following values:

| Name                                            | Type | Default value | Description                                                                                                                                                                                             |
|:------------------------------------------------|:-----|:--------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| <div class="no-wrap">**`efConstruction`**</div> | int  | `512`         | The size of the dynamic list used during k-NN graph creation. Higher values lead to a more accurate graph but slower indexing speed. It is recommended to keep this between 2 and 800 (maximum is 4096) |
| **`m`**                                         | int  | `16`          | The number of bidirectional links that the plugin creates for each new element. Increasing and decreasing this value can have a large impact on memory consumption. Keep this value between 2 and 100.  |

## Model Properties Object

This flexible object, used by `modelProperties` is used to set up models that aren't available in Marqo by default (
models available by default are listed [here](../../Guides/Models-Reference/dense_retrieval.md)).
The structure of this object will vary depending on the model.

For Open CLIP models, see [here](../../Guides/Models-Reference/dense_retrieval.md#generic-clip-models)
for `modelProperties` format and example usage.

For Generic SBERT models, see [here](../../Guides/Models-Reference/dense_retrieval.md#generic-sbert-models)
for `modelProperties` format and example usage.

## Sample structured index settings

Here's a sample settings for a structured index using marqo client:

```python
import marqo

settings = {
    "type": "structured",
    "vectorNumericType": "float",
    "model": "open_clip/ViT-B-32/laion2b_s34b_b79k",
    "normalizeEmbeddings": True,
    "textPreprocessing": {
        "splitLength": 2,
        "splitOverlap": 0,
        "splitMethod": "sentence",
    },
    "imagePreprocessing": {"patchMethod": None},
    "allFields": [
        {"name": "text_field", "type": "text", "features": ["lexical_search"]},
        {"name": "caption", "type": "text", "features": ["lexical_search", "filter"]},
        {"name": "tags", "type": "array<text>", "features": ["filter"]},
        {"name": "image_field", "type": "image_pointer"},
        {"name": "my_int", "type": "int", "features": ["score_modifier"]},
        {
            "name": "multimodal_field",
            "type": "multimodal_combination",
            "dependent_fields": {"image_field": 0.8, "text_field": 0.1},
        },
    ],
    "tensorFields": ["multimodal_field"],
    "annParameters": {
        "spaceType": "prenormalized-angular",
        "parameters": {"efConstruction": 512, "m": 16},
    },
}

mq = marqo.Client(url="http://localhost:8882")

mq.create_index("my-first-structured-index", settings_dict=settings)
```