package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeleteIndexResource struct {
	client *http.Client
}

func NewDeleteIndexResource() resource.Resource {
	return &DeleteIndexResource{}
}

type DeleteIndexResourceModel struct {
	IndexName types.String `tfsdk:"index_name"`
}

func (r *DeleteIndexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delete_index"
}

func (r *DeleteIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = &types.Schema{
		Attributes: map[string]types.Attribute{
			"index_name": {
				Type:     types.StringType,
				Required: true,
			},
		},
	}
}

func (r *DeleteIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model DeleteIndexResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	indexName := model.IndexName.Value
	url := fmt.Sprintf("https://api.marqo.ai/indexes/%s", indexName)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error creating request", fmt.Sprintf("Failed to create a request for deleting an index: %s", err.Error()))
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
	// The index is successfully deleted. Terraform will automatically remove it from the state.
}