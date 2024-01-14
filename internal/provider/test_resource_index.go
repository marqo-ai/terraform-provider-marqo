package provider

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "testing"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockAPIClient simulates an HTTP client for testing
type MockAPIClient struct {
    mock.Mock
}

func (m *MockAPIClient) Do(req *http.Request) (*http.Response, error) {
    args := m.Called(req)
    return args.Get(0).(*http.Response), args.Error(1)
}

func TestIndexResource_Create(t *testing.T) {
    // Mocking the HTTP client
    client := new(MockAPIClient)
    providerConfig := &ProviderConfiguration{APIClient: client}

    // Mock response for the API call
    mockResp := &http.Response{
        StatusCode: http.StatusCreated,
        Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{}`))),
    }
    client.On("Do", mock.Anything).Return(mockResp, nil)

    // Set up the resource
    resource := NewIndexResource().(*IndexResource)
    resource.providerConfig = providerConfig

    // Create request and response objects
    ctx := context.TODO()
    req := resource.CreateRequest{}
    resp := resource.CreateResponse{}

    // Perform the test
    resource.Create(ctx, req, &resp)

    // Assert no errors in diagnostics
    assert.False(t, resp.Diagnostics.HasError())
}

func TestIndexResource_Read(t *testing.T) {
    // Mocking the HTTP client
    client := new(MockAPIClient)
	providerConfig := &ProviderConfiguration{APIClient: client}

	// Mock response for the list indexes API call
	listIndexesResponse := `{
		"results": [
			{"index_name": "test-index"}
		]
	}`
	mockListIndexesResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(listIndexesResponse))),
	}

	// Mock response for the get settings API call
	getIndexSettingsResponse := `{
		"number_of_shards": 3,
		"number_of_replicas": 2
	}`
	mockGetIndexSettingsResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(getIndexSettingsResponse))),
	}

	client.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/api/indexes"
	})).Return(mockListIndexesResp, nil)

	client.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/api/indexes/test-index/settings"
	})).Return(mockGetIndexSettingsResp, nil)

	// Set up the resource
	resource := NewIndexResource().(*IndexResource)
	resource.providerConfig = providerConfig

	// Create request and response objects
	ctx := context.TODO()
	req := resource.ReadRequest{}
	resp := resource.ReadResponse{}

	// Perform the test
	resource.Read(ctx, req, &resp)

	// Assert no errors in diagnostics
	assert.False(t, resp.Diagnostics.HasError())

	// Assert the resource data
	var data IndexResourceModel
	err := resp.State.Get(ctx, &data)
	assert.NoError(t, err)
	assert.Equal(t, "test-index", data.Name.Value)
	assert.Equal(t, int64(3), data.NumberOfShards.Value)
	assert.Equal(t, int64(2), data.NumberOfReplicas.Value)
	}
   
func TestIndexResource_Delete(t *testing.T) {
	// Mocking the HTTP client
	client := new(MockAPIClient)
	providerConfig := &ProviderConfiguration{APIClient: client}

	// Mock response for the API call
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{}`))),
	}
	client.On("Do", mock.Anything).Return(mockResp, nil)

	// Set up the resource
	resource := NewIndexResource().(*IndexResource)
	resource.providerConfig = providerConfig

	// Create request and response objects
	ctx := context.TODO()
	req := resource.DeleteRequest{}
	resp := resource.DeleteResponse{}

	// Perform the test
	resource.Delete(ctx, req, &resp)

	// Assert no errors in diagnostics
	assert.False(t, resp.Diagnostics.HasError())
}

