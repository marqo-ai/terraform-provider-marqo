package marqo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	BaseURL         string
	APIKey          string
	ReturnTelemetry bool
}

func NewClient(baseURL, apiKey string, returnTelemetry bool) *Client {
	return &Client{
		BaseURL:         baseURL,
		APIKey:          apiKey,
		ReturnTelemetry: returnTelemetry,
	}
}

func (c *Client) CreateIndex(indexName string, settings map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/indexes/%s", c.BaseURL, indexName)
	jsonData, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) DeleteIndex(indexName string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/indexes/%s", c.BaseURL, indexName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetIndex(indexName string) (*Index, error) {
	url := fmt.Sprintf("%s/indexes/%s/stats", c.BaseURL, indexName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var index Index
	if err := json.Unmarshal(body, &index); err != nil {
		return nil, err
	}

	return &index, nil
}
