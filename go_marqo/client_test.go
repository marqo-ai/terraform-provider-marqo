package go_marqo_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"marqo/go_marqo"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	baseURL := "http://example.com"
	apiKey := "test-api-key"

	client, err := go_marqo.NewClient(&baseURL, &apiKey)
	assert.NoError(t, err)
	assert.Equal(t, baseURL, client.BaseURL)
	assert.Equal(t, apiKey, client.APIKey)
}

func TestListIndices(t *testing.T) {
	// Create a test server to mock the API response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/indexes", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results": [{"indexName": "test-index"}]}`))
	}))
	defer server.Close()

	client := &go_marqo.Client{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	}

	indices, err := client.ListIndices()
	assert.NoError(t, err)
	assert.Len(t, indices, 1)
	assert.Equal(t, "test-index", indices[0].IndexName)
}

func TestGetIndexSettings(t *testing.T) {
	// Create a test server to mock the API response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/indexes/test-index/settings", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"type": "test-type"}`))
	}))
	defer server.Close()

	client := &go_marqo.Client{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	}

	settings, err := client.GetIndexSettings("test-index")
	assert.NoError(t, err)
	assert.Equal(t, "test-type", settings.Type)
}

func TestCreateIndex(t *testing.T) {
	// Create a test server to mock the API response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/indexes/test-index", r.URL.Path)

		// Verify headers
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		// Read and verify request body
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		// Verify JSON payload
		var settings map[string]interface{}
		err = json.Unmarshal(body, &settings)
		assert.NoError(t, err)
		assert.Equal(t, "cpu", settings["type"])

		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Index created successfully"}`))
	}))
	defer server.Close()

	client := &go_marqo.Client{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	}

	// Test settings
	settings := map[string]interface{}{
		"type": "cpu",
	}

	// Test the CreateIndex function
	err := client.CreateIndex("test-index", settings)
	assert.NoError(t, err)
}
