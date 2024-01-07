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

func NewIndexesDataSource() resource.DataSource {
	return &IndexesDataSource{}
}

func NewIndexSettingsDataSource() resource.DataSource {
	return &IndexSettingsDataSource{}
}

func NewIndexStatsDataSource() resource.DataSource {
	return &IndexStatsDataSource{}
}


type IndexesDataSource struct {
	providerConfig *ProviderConfiguration
}

type IndexStatsDataSource struct {
	providerConfig *ProviderConfiguration
}

type IndexSettingsDataSource struct {
	providerConfig *ProviderConfiguration
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


func (d *IndexesDataSource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model IndexesDataSourceModel

	url := "https://api.marqo.ai/api/indexes"
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


func (d *IndexSettingsDataSource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model IndexSettingsDataSourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
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


func (d *IndexStatsDataSource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model IndexStatsDataSourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("https://api.marqo.ai/api/indexes/%s/stats", model.IndexName.Value)
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