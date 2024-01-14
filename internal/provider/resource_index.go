package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &IndexResource{}

func NewIndexResource() resource.Resource {
	return &IndexResource{}
}

type IndexResource struct {
	providerConfig *ProviderConfiguration
}

type IndexResourceModel struct {
	Name               types.String `tfsdk:"name"`
	IndexDefaults      types.Map    `tfsdk:"index_defaults"`
	NumberOfShards     types.Int64  `tfsdk:"number_of_shards"`
	NumberOfReplicas   types.Int64  `tfsdk:"number_of_replicas"`
	InferenceType      types.String `tfsdk:"inference_type"`
	StorageClass       types.String `tfsdk:"storage_class"`
	NumberOfInferences types.Int64  `tfsdk:"number_of_inferences"`
}

func (r *IndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IndexResourceModel

	// Retrieve the provider configuration
	providerConfig := req.ProviderData.(*ProviderConfiguration)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Construct the request body from the data model
	requestBody, err := json.Marshal(map[string]interface{}{
		"index_defaults":       data.IndexDefaults.Value,
		"number_of_shards":     data.NumberOfShards.Value,
		"number_of_replicas":   data.NumberOfReplicas.Value,
		"inference_type":       data.InferenceType.Value,
		"storage_class":        data.StorageClass.Value,
		"number_of_inferences": data.NumberOfInferences.Value,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating request body", err.Error())
		return
	}

	// Use the API key from provider configuration
	apiKey := providerConfig.APIKey
	url := "https://api.marqo.ai/api/indexes/" + data.Name.Value

	// Create the HTTP request for creating the index
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create request", err.Error())
		return
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("x-api-key", apiKey)

	// Send the request
	httpResponse, err := providerConfig.APIClient.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error sending request to API", err.Error())
		return
	}
	defer httpResponse.Body.Close()

	// Handle the response
	if httpResponse.StatusCode != http.StatusCreated {
		// Read the response body for detailed error message
		bodyBytes, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			resp.Diagnostics.AddError("Error reading response body", err.Error())
			return
		}

		resp.Diagnostics.AddError("API returned non-OK status", fmt.Sprintf("API Error: %s", string(bodyBytes)))
		return
	}

	// Optionally, you can decode the response body if the API returns useful data
	// var result YourAPIResponseStruct
	// err = json.NewDecoder(httpResponse.Body).Decode(&result)
	// if err != nil {
	//     resp.Diagnostics.AddError("Error decoding response body", err.Error())
	//     return
	// }

	// Update the state with the created resource ID or other necessary information
	// Assuming the API returns the ID or name of the created index in the response,
	// or you can use the name provided in the request
	// data.Id = types.String{Value: data.Name.Value}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve the provider configuration
	providerConfig := req.ProviderData.(*ProviderConfiguration)

	var data IndexResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Name.Value
	url := fmt.Sprintf("https://api.marqo.ai/api/indexes/%s", id)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create request", err.Error())
		return
	}
	httpReq.Header.Add("x-api-key", providerConfig.APIKey)

	httpResponse, err := providerConfig.APIClient.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error sending request to API", err.Error())
		return
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			resp.Diagnostics.AddError("Error reading response body", err.Error())
			return
		}
		resp.Diagnostics.AddError("API returned non-OK status", fmt.Sprintf("API Error: %s", string(bodyBytes)))
		return
	}

	// The resource is successfully deleted. Terraform will automatically remove it from the state.
}