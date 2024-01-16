package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ReadIndicesStatusesResource struct {
	client *http.Client
}

func NewReadIndicesStatusesResource() resource.Resource {
	return &ReadIndicesStatusesResource{}
}

type ReadIndicesStatusesResourceModel struct {
	Indices []IndexStatus `tfsdk:"indices"`
}

type IndexStatus struct {
	IndexName          types.String `tfsdk:"index_name"`
	NumberOfShards     types.Int64  `tfsdk:"number_of_shards"`
	NumberOfReplicas   types.Int64  `tfsdk:"number_of_replicas"`
}

func (r *ReadIndicesStatusesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_read_indices_statuses"
}

func (r *ReadIndicesStatusesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = &types.Schema{
		Attributes: map[string]types.Attribute{
			"indices": {
				Type: types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]types.Attribute{
					"index_name": {
						Type:     types.StringType,
						Computed: true,
					},
					"number_of_shards": {
						Type:     types.Int64Type,
						Computed: true,
					},
					"number_of_replicas": {
						Type:     types.Int64Type,
						Computed: true,
					},
				}}},
				Computed: true,
			},
		},
	}
}

func (r *ReadIndicesStatusesResource) Read(ctx context.Context, req resource.ReadRequest, resp*resource.ReadResponse) {
	var model ReadIndicesStatusesResourceModel

	url := "https://api.marqo.ai/indexes"
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error creating request", fmt.Sprintf("Failed to create a request for reading indices statuses: %s", err.Error()))
		return
	}

	// Note: Add authentication header if required

	httpResp, err := r.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error sending request to Marqo API", fmt.Sprintf("Failed to send request: %s", err.Error()))
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unexpected API response status: %s", httpResp.Status))
		return
	}

	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError("Error reading response body", fmt.Sprintf("Failed to read response body: %s", err.Error()))
		return
	}

	var indexesResponse struct {
		Results []struct {
			IndexName string `json:"index_name"`
		} `json:"results"`
	}

	err = json.Unmarshal(bodyBytes, &indexesResponse)
	if err != nil {
		resp.Diagnostics.AddError("Error decoding response body", fmt.Sprintf("Failed to decode response body: %s", err.Error()))
		return
	}

	for _, index := range indexesResponse.Results {
		status := IndexStatus{
			IndexName: types.String{Value: index.IndexName},
			// Populate other fields like NumberOfShards and NumberOfReplicas as needed
		}
		model.Indices = append(model.Indices, status)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)

}