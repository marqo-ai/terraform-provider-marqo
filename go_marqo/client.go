package go_marqo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Client struct {
	BaseURL string
	APIKey  string
}

type IndexResponse struct {
	Results []IndexDetail `json:"results"`
}

type IndexDetail struct {
	Created                      string                  `json:"Created"`
	IndexName                    string                  `json:"indexName"`
	NumberOfShards               int64                   `json:"numberOfShards"`
	NumberOfReplicas             int64                   `json:"numberOfReplicas"`
	IndexStatus                  string                  `json:"indexStatus"`
	AllFields                    []AllFieldInput         `json:"allFields"`
	TensorFields                 []string                `json:"tensorFields"`
	NumberOfInferences           int64                   `json:"numberOfInferences"`
	StorageClass                 string                  `json:"storageClass"`
	InferenceType                string                  `json:"inferenceType"`
	DocsCount                    string                  `json:"docs.count"`
	StoreSize                    string                  `json:"store.size"`
	DocsDeleted                  string                  `json:"docs.deleted"`
	SearchQueryTotal             string                  `json:"search.queryTotal"`
	TreatUrlsAndPointersAsImages bool                    `json:"treatUrlsAndPointersAsImages"`
	TreatUrlsAndPointersAsMedia  bool                    `json:"treatUrlsAndPointersAsMedia"`
	MarqoEndpoint                string                  `json:"marqoEndpoint"`
	Type                         string                  `json:"type"`
	VectorNumericType            string                  `json:"vectorNumericType"`
	Model                        string                  `json:"model"`
	NormalizeEmbeddings          bool                    `json:"normalizeEmbeddings"`
	TextPreprocessing            TextPreprocessing       `json:"textPreprocessing"`
	ImagePreprocessing           ImagePreprocessingModel `json:"imagePreprocessing"` // Assuming no specific structure
	VideoPreprocessing           VideoPreprocessingModel `json:"videoPreprocessing"`
	AudioPreprocessing           AudioPreprocessingModel `json:"audioPreprocessing"`
	AnnParameters                AnnParameters           `json:"annParameters"`
	MarqoVersion                 string                  `json:"marqoVersion"`
	FilterStringMaxLength        int64                   `json:"filterStringMaxLength"`
}

type AllFieldInput struct {
	Name            string             `tfsdk:"name"`
	Type            string             `tfsdk:"type"`
	Features        []string           `tfsdk:"features"`
	DependentFields map[string]float64 `tfsdk:"dependentFields"`
}

type ImagePreprocessingModel struct {
	PatchMethod string `json:"patchMethod"`
}

type TextPreprocessing struct {
	SplitLength  int64  `json:"splitLength"`
	SplitMethod  string `json:"splitMethod"`
	SplitOverlap int64  `json:"splitOverlap"`
}

type VideoPreprocessingModel struct {
	SplitLength  int64 `json:"splitLength"`
	SplitOverlap int64 `json:"splitOverlap"`
}

type AudioPreprocessingModel struct {
	SplitLength  int64 `json:"splitLength"`
	SplitOverlap int64 `json:"splitOverlap"`
}

type AnnParameters struct {
	SpaceType  string          `json:"spaceType"`
	Parameters parametersModel `json:"parameters"`
}

// parametersModel maps the parameters part of ANN parameters.
type parametersModel struct {
	EfConstruction int64 `json:"efConstruction"`
	M              int64 `json:"m"`
}

// IndexStats represents the statistics of an index.
type IndexStats struct {
	NumberOfDocuments int64             `json:"numberOfDocuments"`
	NumberOfVectors   int64             `json:"numberOfVectors"`
	Backend           IndexStatsBackend `json:"backend"`
}

// IndexStatsBackend represents the backend statistics of an index.
type IndexStatsBackend struct {
	MemoryUsedPercentage  float64 `json:"memory_used_percentage"`
	StorageUsedPercentage float64 `json:"storage_used_percentage"`
}

// IndexSettings represents the settings of an index.
type IndexSettings struct {
	Type                         string                  `json:"type"`
	VectorNumericType            string                  `json:"vectorNumericType"`
	Model                        string                  `json:"model"`
	NormalizeEmbeddings          bool                    `json:"normalizeEmbeddings"`
	TextPreprocessing            TextPreprocessing       `json:"textPreprocessing"`
	ImagePreprocessing           ImagePreprocessingModel `json:"imagePreprocessing"`
	AnnParameters                AnnParameters           `json:"annParameters"`
	TensorFields                 []string                `json:"tensorFields"`
	AllFields                    []AllFieldInput         `json:"allFields"`
	NumberOfInferences           int64                   `json:"numberOfInferences"`
	InferenceType                string                  `json:"inferenceType"`
	StorageClass                 string                  `json:"storageClass"`
	NumberOfShards               int64                   `json:"numberOfShards"`
	NumberOfReplicas             int64                   `json:"numberOfReplicas"`
	TreatUrlsAndPointersAsImages bool                    `json:"treatUrlsAndPointersAsImages"`
	FilterStringMaxLength        int64                   `json:"filterStringMaxLength"`
}

// NewClient creates and returns a new API client or an error.
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
	// Print the input parameters
	fmt.Println(baseURL)
	fmt.Println(apiKey)

	// Create the client instance
	client := &Client{
		BaseURL: *baseURL,
		APIKey:  *apiKey,
	}

	// Return the client instance and nil for the error
	return client, nil
}

// ListIndices lists all indices.
func (c *Client) ListIndices() ([]IndexDetail, error) {
	url := fmt.Sprintf("%s/indexes", c.BaseURL)
	tflog.Debug(context.Background(), fmt.Sprintf("Sending request to: %s", url))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("X-API-KEY", c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	//tflog.Debug(context.Background(), fmt.Sprintf("Response status: %s", resp.Status))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	//tflog.Debug(context.Background(), fmt.Sprintf("Response body length: %d", len(body)))

	// Log the response body in chunks
	/*
		const chunkSize = 1000
		for i := 0; i < len(body); i += chunkSize {
			end := i + chunkSize
			if end > len(body) {
				end = len(body)
			}
			tflog.Debug(context.Background(), fmt.Sprintf("Response body part %d: %s", i/chunkSize, string(body[i:end])))
		}
	*/

	var response IndexResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	tflog.Debug(context.Background(), fmt.Sprintf("Number of indices in response: %d", len(response.Results)))

	return response.Results, nil
}

// GetIndexSettings fetches settings for a specific index and decodes into IndexSettings model.
func (c *Client) GetIndexSettings(indexName string) (IndexSettings, error) {
	url := fmt.Sprintf("%s/indexes/%s/settings", c.BaseURL, indexName)
	fmt.Println("GetIndexSettings URL: ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return IndexSettings{}, fmt.Errorf("API request error: %s - %v", req.URL.String(), err)
	}

	req.Header.Set("X-API-KEY", c.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return IndexSettings{}, err
	}
	fmt.Println("Settings Response: ", resp)
	defer resp.Body.Close()

	var settings IndexSettings
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return IndexSettings{}, err
	}

	return settings, nil
}

// GetIndexStats fetches stats for a specific index and decodes into IndexStats model.
func (c *Client) GetIndexStats(indexName string) (IndexStats, error) {
	url := fmt.Sprintf("%s/indexes/%s/stats", c.BaseURL, indexName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return IndexStats{}, fmt.Errorf("API request error: %s - %v", req.URL.String(), err)
	}

	req.Header.Set("X-API-KEY", c.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return IndexStats{}, err
	}
	defer resp.Body.Close()

	var stats IndexStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return IndexStats{}, err
	}

	return stats, nil
}

// CreateIndex creates a new index with the given settings.
func (c *Client) CreateIndex(indexName string, settings map[string]interface{}) error {
	url := fmt.Sprintf("%s/indexes/%s", c.BaseURL, indexName)
	//fmt.Printf("%T\n", settings)

	jsonData, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	fmt.Println("Settings: ", settings)
	fmt.Println("Request: ", req)
	//fmt.Println("JSON Body: ", jsonData)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("Response: ", resp)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create index: %s", string(body))
	}

	return nil
}

// DeleteIndex deletes an index by name.
func (c *Client) DeleteIndex(indexName string) error {
	url := fmt.Sprintf("%s/indexes/%s", c.BaseURL, indexName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-API-KEY", c.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete index: %s", string(body))
	}

	return nil
}

// CreateIndex creates a new index with the given settings.
func (c *Client) UpdateIndex(indexName string, settings map[string]interface{}) error {
	url := fmt.Sprintf("%s/indexes/%s", c.BaseURL, indexName)
	//fmt.Printf("%T\n", settings)

	jsonData, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	//fmt.Println("Settings: ", settings)
	//fmt.Println("Request: ", req)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("Response: ", resp)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update index: %s", string(body))
	}

	return nil
}
