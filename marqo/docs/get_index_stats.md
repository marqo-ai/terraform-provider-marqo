# Stats
Stats give information about indexes, including the number of documents and vectors in the index.
Note that the number of vectors is not necessarily the same as the number of documents, as a document can have multiple vectors or zero vectors.

***

Please also be aware that the number of documents and the number of vectors may not be immediately consistent during indexing processes.
However, these values will be eventually consistent with each other after indexing processes have finished.


## Stats Object

```json
{
  "numberOfDocuments": 4,
  "numberOfVectors": 4,
  "backend": {
    "memoryUsedPercentage": 0.73484113083,
    "storageUsedPercentage": 37.01321365493
  }
}
```

| Name                    | Type    | Description                      
|:------------------------| :------ |:---------------------------------|
| **`numberOfDocuments`** | Integer | number of documents in the index |
| **`numberOfVectors`**   | Integer | number of vectors in the index   |


## Get index stats

Get statistics and information about an index 
```http
GET /indexes/{index_name}/stats
```

### Path parameters

| Name              | Type   | Description                                                               |
| :---------------- | :----- | :------------------|
| **`index_name`**  | String  |  name of the index |

### Example

=== "cURL"
    ```bash
    curl -XGET http://localhost:8882/indexes/my-first-index/stats
    ```
=== "Python"
    ```Python
    results = client.index("my-first-index").get_stats()
    ```
 

#### Response: `200`
```json
{
  "numberOfDocuments": 4,
  "numberOfVectors": 4,
  "backend": {
    "memoryUsedPercentage": 0.73484113083,
    "storageUsedPercentage": 37.01321365493
  }
}
```
