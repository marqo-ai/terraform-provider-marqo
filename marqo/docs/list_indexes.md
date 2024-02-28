# List Indexes

***

```http
GET /indexes
```
List indexes

## Example
=== "cURL"
    ```shell
    curl http://localhost:8882/indexes
    ```
=== "Python"
    ```Python
    mq.get_indexes()
    ```

### Response: `200 OK`
```
{
  "results": [
    {
      "index_name": "Book Collection"
    },
    {
      "index_name": "Animal facts"
    }
  ]
}
```