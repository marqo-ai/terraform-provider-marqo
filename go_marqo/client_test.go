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
		_, err := w.Write([]byte(`{"results": [{"indexName": "test-index"}]}`))
		if err != nil {
			t.Fatal(err)
		}
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
		_, err := w.Write([]byte(`{"type": "test-type"}`))
		if err != nil {
			t.Fatal(err)
		}
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
		_, err = w.Write([]byte(`{"message": "Index created successfully"}`))
		if err != nil {
			t.Fatal(err)
		}
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

func TestGetIndexStats(t *testing.T) {
	// Create a test server to mock the API response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/indexes/test-index/stats", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"numberOfDocuments": 100,
			"numberOfVectors": 200,
			"backend": {
				"memory_used_percentage": 75.5,
				"storage_used_percentage": 60.2
			}
		}`))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	client := &go_marqo.Client{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	}

	stats, err := client.GetIndexStats("test-index")
	assert.NoError(t, err)
	assert.Equal(t, int64(100), stats.NumberOfDocuments)
	assert.Equal(t, int64(200), stats.NumberOfVectors)
	assert.Equal(t, 75.5, stats.Backend.MemoryUsedPercentage)
	assert.Equal(t, 60.2, stats.Backend.StorageUsedPercentage)
}

func TestDeleteIndex(t *testing.T) {
	// Create a test server to mock the API response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/indexes/test-index", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message": "Index deleted successfully"}`))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	client := &go_marqo.Client{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	}

	err := client.DeleteIndex("test-index")
	assert.NoError(t, err)
}

func TestUpdateIndex(t *testing.T) {
	// Create a test server to mock the API response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/indexes/test-index", r.URL.Path)

		// Verify headers
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-KEY"))

		// Read and verify request body
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		// Verify JSON payload
		var settings map[string]interface{}
		err = json.Unmarshal(body, &settings)
		assert.NoError(t, err)
		assert.Equal(t, "gpu", settings["type"])

		// Return success response
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(`{"message": "Index updated successfully"}`))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	client := &go_marqo.Client{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	}

	// Test settings
	settings := map[string]interface{}{
		"type": "gpu",
	}

	// Test the UpdateIndex function
	err := client.UpdateIndex("test-index", settings)
	assert.NoError(t, err)
}
