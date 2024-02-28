
# Delete index
Delete an index. 

***

**Note: This operation cannot be undone, and the deleted index can't be recovered**

```http
DELETE /indexes/{index_name}
```

## Example
=== "cURL"
    ```shell
    curl -XDELETE http://localhost:8882/indexes/my-first-index
    ```
=== "Python"
    ```Python
    results = mq.index("my-first-index").delete()
    ```

### Response: `200 OK`
```
{"acknowledged": true}
```
