# Settings
Get settings of an index. For a conceptual overview of index settings, refer to our [Index API Reference](../Indexes/create_structured_index.md).

***

## Get index settings

```http
GET /indexes/{index_name}/settings
```

### Path parameters

| Name              | Type   | Description                                                               |
| :---------------- | :----- | :------------------|
| **`index_name`**  | String  |  name of the index |

### Example

=== "cURL"
    ```bash
    curl -XGET http://localhost:8882/indexes/my-first-index/settings
    ```
=== "Python"
    ```Python
    results = mq.index("my-first-index").get_settings()
    ```

#### Response: `200`
```json
{
  "annParameters": {
    "parameters": {
      "efConstruction": 512,
      "m": 16
    },
    "spaceType": "prenormalized-angular"
  },
  "filterStringMaxLength": 20,
  "imagePreprocessing": {},
  "model": "hf/e5-base-v2",
  "normalizeEmbeddings": true,
  "textPreprocessing": {
    "splitLength": 2,
    "splitMethod": "sentence",
    "splitOverlap": 0
  },
  "treatUrlsAndPointersAsImages": false,
  "type": "unstructured",
  "vectorNumericType": "float"
}
```
