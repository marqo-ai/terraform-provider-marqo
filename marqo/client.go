package marqo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	BaseURL string
	APIKey  string
}

// NewClient creates and returns a new API client or an error
func NewClient(baseURL, apiKey *string) (*Client, error) {
	// Validate the input parameters
	if baseURL == nil || *baseURL == "" {
		return nil, errors.New("baseURL is required but was not provided")
	}
	if apiKey == nil || *apiKey == "" {
		return nil, errors.New("apiKey is required but was not provided")
	}

	//
	// TO IMPLEMENT:
	// - Translate is_marqo_cloud = False
	//    if url is not None:
	//	if url.lower().startswith(os.environ.get("MARQO_CLOUD_URL", "https://api.marqo.ai")):
	//		instance_mappings = MarqoCloudInstanceMappings(control_base_url=url, api_key=api_key)
	//		is_marqo_cloud = True
	//	else:
	//		instance_mappings = DefaultInstanceMappings(url, main_user, main_password)
	//

	// Create the client instance
	client := &Client{
		BaseURL: *baseURL,
		APIKey:  *apiKey,
	}

	// Return the client instance and nil for the error
	return client, nil
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

// Might need to make types.go file to define the types that are returned from the API
