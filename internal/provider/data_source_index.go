package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type IndexesDataSource struct {
	providerConfig *ProviderConfiguration
}

type IndexStatsDataSource struct {
	providerConfig *ProviderConfiguration
}

type IndexSettingsDataSource struct {
	providerConfig *ProviderConfiguration
}

func NewIndexesDataSource() datasource.DataSource {
	return &IndexesDataSource{}
}

func NewIndexSettingsDataSource() datasource.DataSource {
	return &IndexSettingsDataSource{}
}

func NewIndexStatsDataSource() datasource.DataSource {
	return &IndexStatsDataSource{}
}

func (d *IndexSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_indexes"
}

func (d *IndexesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_indexes"
}

func (d *IndexStatsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_indexes"
}

func (d *IndexesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Data source for retrieving a list of indexes from Marqo API",

        Attributes: map[string]schema.Attribute{
            "indexes": schema.ListAttribute{
                MarkdownDescription: "List of index names",
                Computed:            true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *IndexStatsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Data source for retrieving statistics of a specific index from Marqo API",

        Attributes: map[string]schema.Attribute{
            "index_name": schema.StringAttribute{
                MarkdownDescription: "Name of the index",
                Required:            true,
            },
            "number_of_documents": schema.Int64Attribute{
                MarkdownDescription: "Number of documents in the index",
                Computed:            true,
            },
            "number_of_vectors": schema.Int64Attribute{
                MarkdownDescription: "Number of vectors in the index",
                Computed:            true,
            },
        },
    }
}

func (d *IndexSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Data source for retrieving settings of a specific index from Marqo API",

        Attributes: map[string]schema.Attribute{
            "index_name": schema.StringAttribute{
                MarkdownDescription: "Name of the index",
                Required:            true,
            },
			// Implement nested attributes here
        },
    }
}

type IndexesDataSourceModel struct {
	Indexes []types.String `tfsdk:"indexes"`
}

type IndexStatsDataSourceModel struct {
	IndexName        types.String `tfsdk:"index_name"`
	NumberOfDocuments types.Int64  `tfsdk:"number_of_documents"`
	NumberOfVectors   types.Int64  `tfsdk:"number_of_vectors"`
}

type IndexSettingsDataSourceModel struct {
	IndexName   types.String `tfsdk:"index_name"`
	Settings    types.Map    `tfsdk:"settings"`
}

func (d *IndexesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model IndexesDataSourceModel

	url := "https://api.marqo.ai/api/indexes"
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create request", err.Error())
		return
	}

	providerConfig := req.Provider.(*ProviderConfiguration)
	//providerConfig := req.ProviderData.(*ProviderConfiguration)
	httpReq.Header.Add("x-api-key", providerConfig.APIKey)

	httpResp, err := providerConfig.APIClient.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error sending request to API", err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			resp.Diagnostics.AddError("Error reading response body", err.Error())
			return
		}
		resp.Diagnostics.AddError("API returned non-OK status", fmt.Sprintf("API Error: %s", string(bodyBytes)))
		return
	}

	var apiResponse struct {
		Results []struct {
			IndexName string `json:"index_name"`
		} `json:"results"`
	}
	err = json.NewDecoder(httpResp.Body).Decode(&apiResponse)
	if err != nil {
		resp.Diagnostics.AddError("Error decoding response body", err.Error())
		return
	}

	for _, index := range apiResponse.Results {
		model.Indexes = append(model.Indexes, types.String{Value: index.IndexName})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}


func (d *IndexSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model IndexSettingsDataSourceModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("https://api.marqo.ai/api/indexes/%s/settings", model.IndexName.Value)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create request", err.Error())
		return
	}
	
	providerConfig := req.ProviderData.(*ProviderConfiguration)
	httpReq.Header.Add("x-api-key", providerConfig.APIKey)

	httpResp, err := providerConfig.APIClient.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error sending request to API", err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			resp.Diagnostics.AddError("Error reading response body", err.Error())
			return
		}
		resp.Diagnostics.AddError("API returned non-OK status", fmt.Sprintf("API Error: %s", string(bodyBytes)))
		return
	}

	var settings map[string]interface{}
	err = json.NewDecoder(httpResp.Body).Decode(&settings)
	if err != nil {
		resp.Diagnostics.AddError("Error decoding response body", err.Error())
		return
	}

	model.Settings = types.Map{Value: settings}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}


func (d *IndexStatsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model IndexStatsDataSourceModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("https://api.marqo.ai/api/indexes/%s/stats", model.IndexName.Value)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create request", err.Error())
		return
	}

	//providerConfig := req.Provider.(*ProviderConfiguration)
	providerConfig := req.ProviderData.(*ProviderConfiguration)
	httpReq.Header.Add("x-api-key", providerConfig.APIKey)

	httpResp, err := providerConfig.APIClient.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error sending request to API", err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			resp.Diagnostics.AddError("Error reading response body", err.Error())
			return
		}
		resp.Diagnostics.AddError("API returned non-OK status", fmt.Sprintf("API Error: %s", string(bodyBytes)))
		return
	}

	var stats struct {
		NumberOfDocuments int64 `json:"numberOfDocuments"`
		NumberOfVectors   int64 `json:"numberOfVectors"`
	}
	err = json.NewDecoder(httpResp.Body).Decode(&stats)
	if err != nil {
		resp.Diagnostics.AddError("Error decoding response body", err.Error())
		return
	}

	model.NumberOfDocuments = types.Int64{Value: stats.NumberOfDocuments}
	model.NumberOfVectors = types.Int64{Value: stats.NumberOfVectors}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}