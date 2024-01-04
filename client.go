package marqo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MarqoClient struct {
	BaseURL string
}

func NewMarqoClient(baseURL string) *MarqoClient {
	return &MarqoClient{
		BaseURL: baseURL,
	}
}

// Structs for Marqo index creation
type IndexSettings struct {
	IndexDefaults IndexDefaults `json:"index_defaults"`
	NumberOfShards int `json:"number_of_shards"`
	NumberOfReplicas int `json:"number_of_replicas"`
}

type IndexDefaults struct {
	TreatUrlsAndPointersAsImages bool `json:"treat_urls_and_pointers_as_images"`
	Model string `json:"model"`
	//  other fields?
}

// CreateIndex creates a new index in Marqo
func (c *MarqoClient) CreateIndex(indexName string, settings IndexSettings) error {
	url := fmt.Sprintf("%s/indexes/%s", c.BaseURL, indexName)
	jsonData, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and check the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("error creating index: %s", string(body))
	}

	return nil
}

// DeleteIndex deletes an existing index in Marqo
func (c *MarqoClient) DeleteIndex(indexName string) error {
	url := fmt.Sprintf("%s/indexes/%s", c.BaseURL, indexName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and check the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("error deleting index: %s", string(body))
	}

	return nil
}

// ListIndexes lists all indexes in Marqo
func (c *MarqoClient) ListIndexes() ([]string, error) {
	url := fmt.Sprintf("%s/indexes", c.BaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error listing indexes: %s", string(body))
	}

	var result struct {
		Results []struct {
			IndexName string `json:"index_name"`
		} `json:"results"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var indexes []string
	for _, r := range result.Results {
		indexes = append(indexes, r.IndexName)
	}
	return indexes, nil
}

// UpdateDocuments updates documents in a specified index in Marqo
func (c *MarqoClient) UpdateDocuments(indexName string, documents []map[string]interface{}) error {
	url := fmt.Sprintf("%s/indexes/%s/documents", c.BaseURL, indexName)
	jsonData, err := json.Marshal(map[string]interface{}{
		"documents": documents,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("error updating documents: %s", string(body))
	}

	return nil
}

// DeleteDocuments deletes specified documents from an index in Marqo
func (c *MarqoClient) DeleteDocuments(indexName string, documentIDs []string) error {
	url := fmt.Sprintf("%s/indexes/%s/documents/delete-batch", c.BaseURL, indexName)
	jsonData, err := json.Marshal(documentIDs)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("error deleting documents: %s", string(body))
	}

	return nil
}