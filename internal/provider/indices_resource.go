package provider

import (
	"context"
	"fmt"
	"terraform-provider-marqo/marqo"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &indicesResource{}
	_ resource.ResourceWithConfigure = &indicesResource{}
)

// ManageIndicesResource is a helper function to simplify the provider implementation.
func ManageIndicesResource() resource.Resource {
	return &indicesResource{}
}

// orderResource is the resource implementation.
type indicesResource struct {
	marqoClient *marqo.Client
}

// IndexResourceModel maps the resource schema data.
type IndexResourceModel struct {
	ID        types.String       `tfsdk:"id"`
	IndexName types.String       `tfsdk:"index_name"`
	Settings  IndexSettingsModel `tfsdk:"settings"`
}

type IndexSettingsModel struct {
	Type                         types.String                 `tfsdk:"type"`
	VectorNumericType            types.String                 `tfsdk:"vector_numeric_type"`
	TreatUrlsAndPointersAsImages types.Bool                   `tfsdk:"treat_urls_and_pointers_as_images"`
	Model                        types.String                 `tfsdk:"model"`
	NormalizeEmbeddings          types.Bool                   `tfsdk:"normalize_embeddings"`
	TextPreprocessing            TextPreprocessingModelCreate `tfsdk:"text_preprocessing"`
	ImagePreprocessing           ImagePreprocessingModel      `tfsdk:"image_preprocessing"`
	AnnParameters                AnnParametersModelCreate     `tfsdk:"ann_parameters"`
	FilterStringMaxLength        types.Int64                  `tfsdk:"filter_string_max_length"`
}

type TextPreprocessingModelCreate struct {
	SplitLength  types.Int64  `tfsdk:"split_length"`
	SplitMethod  types.String `tfsdk:"split_method"`
	SplitOverlap types.Int64  `tfsdk:"split_overlap"`
}

type ImagePreprocessingModel struct {
	PatchMethod types.String `tfsdk:"patch_method"`
}

type AnnParametersModelCreate struct {
	SpaceType  types.String    `tfsdk:"space_type"`
	Parameters ParametersModel `tfsdk:"parameters"`
}

type ParametersModel struct {
	EfConstruction types.Int64 `tfsdk:"ef_construction"`
	M              types.Int64 `tfsdk:"m"`
}

// Configure adds the provider configured client to the resource.
func (r *indicesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*marqo.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.marqoClient = client
}

// Metadata returns the resource type name.
func (r *indicesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_index"
}

func (r *indicesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the index.",
			},
			"index_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the index.",
			},
			"settings": schema.SingleNestedAttribute{
				Required:    true,
				Description: "The settings for the index.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
					"vector_numeric_type": schema.StringAttribute{
						Optional: true,
					},
					"treat_urls_and_pointers_as_images": schema.BoolAttribute{
						Optional: true,
					},
					"model": schema.StringAttribute{
						Optional: true,
					},
					"normalize_embeddings": schema.BoolAttribute{
						Optional: true,
					},
					"text_preprocessing": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"split_length":  schema.Int64Attribute{Optional: true},
							"split_method":  schema.StringAttribute{Optional: true},
							"split_overlap": schema.Int64Attribute{Optional: true},
						},
					},
					"image_preprocessing": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"patch_method": schema.StringAttribute{Optional: true},
						},
					},
					"ann_parameters": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"space_type": schema.StringAttribute{Optional: true},
							"parameters": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"ef_construction": schema.Int64Attribute{Optional: true},
									"m":               schema.Int64Attribute{Optional: true},
								},
							},
						},
					},
					"filter_string_max_length": schema.Int64Attribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *indicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IndexResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	indexSettings, err := r.marqoClient.GetIndexSettings(state.IndexName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Index",
			fmt.Sprintf("Could not read index '%s': %s", state.IndexName, err.Error()),
		)
		return
	}

	// Map the fetched index settings to the Terraform state
	state.Settings = IndexSettingsModel{
		Type:                         types.StringValue(indexSettings.Type),
		VectorNumericType:            types.StringValue(indexSettings.VectorNumericType),
		TreatUrlsAndPointersAsImages: types.BoolValue(indexSettings.TreatUrlsAndPointersAsImages),
		Model:                        types.StringValue(indexSettings.Model),
		NormalizeEmbeddings:          types.BoolValue(indexSettings.NormalizeEmbeddings),
		TextPreprocessing: TextPreprocessingModelCreate{
			SplitLength:  types.Int64Value(indexSettings.TextPreprocessing.SplitLength),
			SplitMethod:  types.StringValue(indexSettings.TextPreprocessing.SplitMethod),
			SplitOverlap: types.Int64Value(indexSettings.TextPreprocessing.SplitOverlap),
		},
		ImagePreprocessing: ImagePreprocessingModel{
			PatchMethod: types.StringValue(indexSettings.ImagePreprocessing["patchMethod"].(string)),
		},
		AnnParameters: AnnParametersModelCreate{
			SpaceType: types.StringValue(indexSettings.AnnParameters.SpaceType),
			Parameters: ParametersModel{
				EfConstruction: types.Int64Value(indexSettings.AnnParameters.Parameters.EfConstruction),
				M:              types.Int64Value(indexSettings.AnnParameters.Parameters.M),
			},
		},
		FilterStringMaxLength: types.Int64Value(indexSettings.FilterStringMaxLength),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *indicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model IndexResourceModel
	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//fmt.Println(model.Settings)

	// Construct settings map
	settings := map[string]interface{}{
		"type":                         model.Settings.Type.ValueString(),
		"vectorNumericType":            model.Settings.VectorNumericType.ValueString(),
		"treatUrlsAndPointersAsImages": model.Settings.TreatUrlsAndPointersAsImages.ValueBool(),
		"model":                        model.Settings.Model.ValueString(),
		"normalizeEmbeddings":          model.Settings.NormalizeEmbeddings.ValueBool(),
		"textPreprocessing": map[string]interface{}{
			"splitLength":  model.Settings.TextPreprocessing.SplitLength.ValueInt64(),
			"splitMethod":  model.Settings.TextPreprocessing.SplitMethod.ValueString(),
			"splitOverlap": model.Settings.TextPreprocessing.SplitOverlap.ValueInt64(),
		},
		"imagePreprocessing": map[string]interface{}{
			"patchMethod": model.Settings.ImagePreprocessing.PatchMethod.ValueString(),
		},
		"annParameters": map[string]interface{}{
			"spaceType": model.Settings.AnnParameters.SpaceType.ValueString(),
			"parameters": map[string]interface{}{
				"efConstruction": model.Settings.AnnParameters.Parameters.EfConstruction.ValueInt64(),
				"m":              model.Settings.AnnParameters.Parameters.M.ValueInt64(),
			},
		},
		"filterStringMaxLength": model.Settings.FilterStringMaxLength.ValueInt64(),
	}

	// Remove optional fields if they are not set
	if model.Settings.VectorNumericType.IsNull() {
		delete(settings, "vectorNumericType")
	}
	if model.Settings.TreatUrlsAndPointersAsImages.IsNull() {
		delete(settings, "treatUrlsAndPointersAsImages")
	}
	if model.Settings.Model.IsNull() {
		delete(settings, "model")
	}
	if model.Settings.NormalizeEmbeddings.IsNull() {
		delete(settings, "normalizeEmbeddings")
	}
	if model.Settings.ImagePreprocessing.PatchMethod.IsNull() {
		delete(settings["imagePreprocessing"].(map[string]interface{}), "patchMethod")
	}
	if len(settings["imagePreprocessing"].(map[string]interface{})) == 0 {
		delete(settings, "imagePreprocessing")
	}

	//indexNameAsString := model.IndexName.

	err := r.marqoClient.CreateIndex(model.IndexName.ValueString(), settings)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Create Index", "Could not create index: "+err.Error())
		return
	}

	// Set the index name as the ID in the Terraform state
	model.ID = model.IndexName
	diags = resp.State.Set(ctx, &model)
	resp.Diagnostics.Append(diags...)
}

func (r *indicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model IndexResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.marqoClient.DeleteIndex(model.IndexName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to Delete Index", "Could not delete index: "+err.Error())
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *indicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}
