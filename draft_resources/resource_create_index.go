package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CreateIndexResource struct {
	client *http.Client
}

// NewCreateIndexResource returns a new CreateIndexResource.
func NewCreateIndexResource() resource.Resource {
	return &CreateIndexResource{}
}

type CreateIndexResourceModel struct {
	IndexName  types.String `tfsdk:"index_name"`
	Settings   types.Map    `tfsdk:"settings"`
}

func (r *CreateIndexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_create_index"
}

func (r *CreateIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = &types.Schema{
		Attributes: map[string]types.Attribute{
			"index_name": {
				Type:     types.StringType,
				Required: true,
			},
			"settings": {
				Type:     types.MapType,
				Optional: true,
			},
		},
	}
}

func (r *CreateIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model CreateIndexResourceModel
	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	indexName := model.IndexName.Value
	settings, _ := json.Marshal(model.Settings)

	url := fmt.Sprintf("https://api.marqo.ai/indexes/%s", indexName)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(settings))
	if err != nil {
		resp.Diagnostics.AddError("Error creating HTTP request", fmt.Sprintf("Failed to create HTTP request: %s", err.Error()))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	if r.client == nil {
		resp.Diagnostics.AddError("Client Error", "HTTP client is not initialized")
		return
	}

	httpResp, err := r.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create index, got error: %s", err))
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned non-200 status code: %d", httpResp.StatusCode))
		return
	}


	// Assuming the index creation is successful and doesn't return any specific ID
	model.Id = types.String{Value: indexName}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
