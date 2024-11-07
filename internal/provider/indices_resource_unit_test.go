package provider

import (
	"context"
	"fmt"
	"marqo/go_marqo"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMarqoClient is a mock implementation of the Marqo client
type MockMarqoClient struct {
	mock.Mock
}

func (m *MockMarqoClient) CreateIndex(indexName string, settings map[string]interface{}) error {
	args := m.Called(indexName, settings)
	return args.Error(0)
}

func (m *MockMarqoClient) DeleteIndex(indexName string) error {
	args := m.Called(indexName)
	return args.Error(0)
}

func (m *MockMarqoClient) ListIndices() ([]go_marqo.IndexDetail, error) {
	args := m.Called()
	return args.Get(0).([]go_marqo.IndexDetail), args.Error(1)
}

func (m *MockMarqoClient) UpdateIndex(indexName string, settings map[string]interface{}) error {
	args := m.Called(indexName, settings)
	return args.Error(0)
}

// Helper function to create a test IndexResourceModel
func createTestIndexModel() IndexResourceModel {
	return IndexResourceModel{
		IndexName: types.StringValue("test-index"),
		Settings: IndexSettingsModel{
			Type:               types.StringValue("unstructured"),
			NumberOfInferences: types.Int64Value(1),
			StorageClass:       types.StringValue("marqo.basic"),
			NumberOfShards:     types.Int64Value(1),
			NumberOfReplicas:   types.Int64Value(0),
			Model:              types.StringValue("test-model"),
			InferenceType:      types.StringValue("CPU"),
		},
	}
}

func TestIndicesResource_Create(t *testing.T) {
	mockClient := new(MockMarqoClient)
	r := &indicesResource{
		marqoClient: mockClient,
	}

	// Setup test data
	model := createTestIndexModel()

	// Setup expectations
	mockClient.On("CreateIndex", "test-index", mock.Anything).Return(nil)

	// Create context and request
	ctx := context.Background()
	req := resource.CreateRequest{}
	resp := resource.CreateResponse{}

	// Set the plan in the request
	req.Plan.Set(ctx, &model)

	// Call Create
	r.Create(ctx, req, &resp)

	// Assertions
	assert.False(t, resp.Diagnostics.HasError(), "Create should not return an error")
	mockClient.AssertExpectations(t)

	// Test error case
	mockClient.On("CreateIndex", "test-index", mock.Anything).Return(fmt.Errorf("creation failed"))
	r.Create(ctx, req, &resp)
	assert.True(t, resp.Diagnostics.HasError(), "Create should return an error when client fails")
}

func TestIndicesResource_Read(t *testing.T) {
	mockClient := new(MockMarqoClient)
	r := &indicesResource{
		marqoClient: mockClient,
	}

	// Setup test data
	model := createTestIndexModel()
	indices := []go_marqo.IndexDetail{
		{
			IndexName:        "test-index",
			Type:             "unstructured",
			NumberOfShards:   1,
			NumberOfReplicas: 0,
			StorageClass:     "BASIC",
			InferenceType:    "CPU",
		},
	}

	// Setup expectations
	mockClient.On("ListIndices").Return(indices, nil)

	// Create context and request
	ctx := context.Background()
	req := resource.ReadRequest{}
	resp := resource.ReadResponse{}

	// Set the state in the request
	req.State.Set(ctx, &model)

	// Call Read
	r.Read(ctx, req, &resp)

	// Assertions
	assert.False(t, resp.Diagnostics.HasError(), "Read should not return an error")
	mockClient.AssertExpectations(t)

	// Test error case
	mockClient.On("ListIndices").Return([]go_marqo.IndexDetail{}, fmt.Errorf("read failed"))
	r.Read(ctx, req, &resp)
	assert.True(t, resp.Diagnostics.HasError(), "Read should return an error when client fails")
}

func TestIndicesResource_Update(t *testing.T) {
	mockClient := new(MockMarqoClient)
	r := &indicesResource{
		marqoClient: mockClient,
	}

	// Setup test data
	model := createTestIndexModel()

	// Setup expectations
	mockClient.On("UpdateIndex", "test-index", mock.Anything).Return(nil)

	// Create context and request
	ctx := context.Background()
	req := resource.UpdateRequest{}
	resp := resource.UpdateResponse{}

	// Set the plan in the request
	req.Plan.Set(ctx, &model)

	// Call Update
	r.Update(ctx, req, &resp)

	// Assertions
	assert.False(t, resp.Diagnostics.HasError(), "Update should not return an error")
	mockClient.AssertExpectations(t)

	// Test error case
	mockClient.On("UpdateIndex", "test-index", mock.Anything).Return(fmt.Errorf("update failed"))
	r.Update(ctx, req, &resp)
	assert.True(t, resp.Diagnostics.HasError(), "Update should return an error when client fails")
}

func TestIndicesResource_Delete(t *testing.T) {
	mockClient := new(MockMarqoClient)
	r := &indicesResource{
		marqoClient: mockClient,
	}

	// Setup test data
	model := createTestIndexModel()

	// Setup expectations
	mockClient.On("DeleteIndex", "test-index").Return(nil)

	// Create context and request
	ctx := context.Background()
	req := resource.DeleteRequest{}
	resp := resource.DeleteResponse{}

	// Set the state in the request
	req.State.Set(ctx, &model)

	// Call Delete
	r.Delete(ctx, req, &resp)

	// Assertions
	assert.False(t, resp.Diagnostics.HasError(), "Delete should not return an error")
	mockClient.AssertExpectations(t)

	// Test error case
	mockClient.On("DeleteIndex", "test-index").Return(fmt.Errorf("deletion failed"))
	r.Delete(ctx, req, &resp)
	assert.True(t, resp.Diagnostics.HasError(), "Delete should return an error when client fails")
}

func TestValidateAndConstructAllFields(t *testing.T) {
	tests := []struct {
		name    string
		input   []AllFieldInput
		want    []map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid input with all fields",
			input: []AllFieldInput{
				{
					Name:     types.StringValue("test_field"),
					Type:     types.StringValue("text"),
					Features: []types.String{types.StringValue("lexical_search")},
					DependentFields: map[string]types.Float64{
						"field1": types.Float64Value(0.8),
						"field2": types.Float64Value(0.2),
					},
				},
			},
			want: []map[string]interface{}{
				{
					"name":     "test_field",
					"type":     "text",
					"features": []string{"lexical_search"},
					"dependentFields": map[string]float64{
						"field1": 0.8,
						"field2": 0.2,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing required fields",
			input: []AllFieldInput{
				{
					Name: types.StringNull(),
					Type: types.StringValue("text"),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateAndConstructAllFields(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAndConstructAllFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateAndConstructAllFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatesAreEqual(t *testing.T) {
	tests := []struct {
		name     string
		existing *IndexResourceModel
		desired  *IndexResourceModel
		want     bool
	}{
		{
			name: "identical states",
			existing: &IndexResourceModel{
				IndexName: types.StringValue("test"),
				Settings: IndexSettingsModel{
					Type:               types.StringValue("unstructured"),
					NumberOfInferences: types.Int64Value(1),
				},
			},
			desired: &IndexResourceModel{
				IndexName: types.StringValue("test"),
				Settings: IndexSettingsModel{
					Type:               types.StringValue("unstructured"),
					NumberOfInferences: types.Int64Value(1),
				},
			},
			want: true,
		},
		{
			name: "different states",
			existing: &IndexResourceModel{
				IndexName: types.StringValue("test"),
				Settings: IndexSettingsModel{
					Type:               types.StringValue("unstructured"),
					NumberOfInferences: types.Int64Value(1),
				},
			},
			desired: &IndexResourceModel{
				IndexName: types.StringValue("test"),
				Settings: IndexSettingsModel{
					Type:               types.StringValue("unstructured"),
					NumberOfInferences: types.Int64Value(2),
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := statesAreEqual(tt.existing, tt.desired)
			if got != tt.want {
				t.Errorf("statesAreEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToInt64(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  types.Int64
	}{
		{
			name:  "valid integer",
			input: "123",
			want:  types.Int64Value(123),
		},
		{
			name:  "invalid integer",
			input: "abc",
			want:  types.Int64Null(),
		},
		{
			name:  "empty string",
			input: "",
			want:  types.Int64Null(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringToInt64(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}
