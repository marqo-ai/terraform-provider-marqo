package provider

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-marqo/marqo"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

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
	NumberOfInferences           types.Int64                  `tfsdk:"number_of_inferences"`
	AllFields                    []AllFieldInput              `tfsdk:"all_fields"`
	TensorFields                 []string                     `tfsdk:"tensor_fields"`
	InferenceType                types.String                 `tfsdk:"inference_type"`
	StorageClass                 types.String                 `tfsdk:"storage_class"`
	NumberOfShards               types.Int64                  `tfsdk:"number_of_shards"`
	NumberOfReplicas             types.Int64                  `tfsdk:"number_of_replicas"`
	TreatUrlsAndPointersAsImages types.Bool                   `tfsdk:"treat_urls_and_pointers_as_images"`
	Model                        types.String                 `tfsdk:"model"`
	NormalizeEmbeddings          types.Bool                   `tfsdk:"normalize_embeddings"`
	TextPreprocessing            TextPreprocessingModelCreate `tfsdk:"text_preprocessing"`
	ImagePreprocessing           ImagePreprocessingModel      `tfsdk:"image_preprocessing"`
	AnnParameters                AnnParametersModelCreate     `tfsdk:"ann_parameters"`
	FilterStringMaxLength        types.Int64                  `tfsdk:"filter_string_max_length"`
}

//             "dependentFields": {"image_field": 0.8, "text_field": 0.1},

type AllFieldInput struct {
	Name            types.String             `tfsdk:"name"`
	Type            types.String             `tfsdk:"type"`
	Features        []types.String           `tfsdk:"features"`
	DependentFields map[string]types.Float64 `tfsdk:"dependent_fields"`
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
					"number_of_inferences": schema.Int64Attribute{
						Optional: true,
					},
					"all_fields": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{Optional: true},
								"type": schema.StringAttribute{Optional: true},
								"features": schema.ListAttribute{
									Optional:    true,
									ElementType: types.StringType,
								},
								// Sample:  "dependentFields": {"image_field": 0.8, "text_field": 0.1},
								"dependent_fields": schema.MapAttribute{
									Optional:    true,
									ElementType: types.Int64Type,
								},
							},
						},
					},
					"tensor_fields": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},

					"inference_type": schema.StringAttribute{
						Optional: true,
					},
					"storage_class": schema.StringAttribute{
						Optional: true,
					},
					"number_of_shards": schema.Int64Attribute{
						Optional: true,
					},
					"number_of_replicas": schema.Int64Attribute{
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

// Utility function to convert standard Go string to types.Int64 .
func StringToInt64(str string) types.Int64 {
	intVal, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		// Handle the error appropriately. Here, we return a Null types.Int64 to indicate failure.
		return types.Int64Null()
	}
	return types.Int64Value(intVal)
}

// validateAndConstructAllFields validates the allFields input and constructs the corresponding setting.
func validateAndConstructAllFields(allFieldsInput []AllFieldInput) ([]map[string]interface{}, error) {
	var allFields []map[string]interface{}
	for _, field := range allFieldsInput {
		// Basic validation example. Expand based on full Marqo documentation requirements.
		if field.Name.IsNull() || field.Type.IsNull() {
			return nil, fmt.Errorf("each field must have a name and type")
		}
		fieldMap := map[string]interface{}{
			"name":     field.Name.ValueString(),
			"type":     field.Type.ValueString(),
			"features": []string{}, // Convert types.String to string
		}
		// Assert the type of "features" before appending
		features, ok := fieldMap["features"].([]string)
		if !ok {
			// Handle the error
			features = []string{}
		}
		for _, feature := range field.Features {
			features = append(features, feature.ValueString())
		}
		fieldMap["features"] = features
		//if len(field.DependentFields) > 0 {
		//	dependentFields := make(map[string]float64)
		//	for key, value := range field.DependentFields {
		//		dependentFields[key] = value.ValueFloat64()
		//	}
		//	fieldMap["dependent_fields"] = dependentFields
		//}
		allFields = append(allFields, fieldMap)
	}
	return allFields, nil
}

// Utility function to convert []AllFieldInput to a format suitable for settings map.
func convertAllFieldsToMap(allFieldsInput []AllFieldInput) []map[string]interface{} {
	allFields := []map[string]interface{}{}
	for _, field := range allFieldsInput {
		fieldMap := map[string]interface{}{
			"name": field.Name.ValueString(),
			"type": field.Type.ValueString(),
			// Add other necessary fields from AllFieldInput struct
		}
		// Assuming Features is a slice of types.String and needs conversion
		features := []string{}
		for _, feature := range field.Features {
			features = append(features, feature.ValueString())
		}
		fieldMap["features"] = features

		// Convert DependentFields if necessary
		//dependentFieldsMap := make(map[string]float64)
		//for key, value := range field.DependentFields {
		//	dependentFieldsMap[key] = value.ValueFloat64()
		//}
		//if len(dependentFieldsMap) > 0 {
		//	fieldMap["dependent_fields"] = dependentFieldsMap
		//}

		allFields = append(allFields, fieldMap)
	}
	return allFields
}

// constructTensorFields constructs the tensorFields setting from the input.
func constructTensorFields(tensorFieldsInput []string) ([]string, error) {
	// Blank for now
	return tensorFieldsInput, nil
}

func (r *indicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Initialize the state variable based on the IndexResourceModel
	var state IndexResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Calling marqo client ListIndices")
	indices, err := r.marqoClient.ListIndices()
	if err != nil {
		resp.Diagnostics.AddError("Failed to List Indices", fmt.Sprintf("Could not list indices: %s", err.Error()))
		return
	}

	for _, indexDetail := range indices {
		if indexDetail.IndexName == state.IndexName.ValueString() {
			// Update the state with the details from the indexDetail
			state.Settings = IndexSettingsModel{
				Type:                         types.StringValue(indexDetail.Type),
				VectorNumericType:            types.StringValue(indexDetail.VectorNumericType),
				TreatUrlsAndPointersAsImages: types.BoolValue(indexDetail.TreatUrlsAndPointersAsImages),
				Model:                        types.StringValue(indexDetail.Model),
				AllFields:                    ConvertMarqoAllFieldInputs(indexDetail.AllFields),
				TensorFields:                 indexDetail.TensorFields,
				NormalizeEmbeddings:          types.BoolValue(indexDetail.NormalizeEmbeddings),
				InferenceType:                types.StringValue(indexDetail.InferenceType),
				NumberOfInferences:           StringToInt64(indexDetail.NumberOfInferences),
				StorageClass:                 types.StringValue(indexDetail.StorageClass),
				NumberOfShards:               StringToInt64(indexDetail.NumberOfShards),
				NumberOfReplicas:             StringToInt64(indexDetail.NumberOfReplicas),
				TextPreprocessing: TextPreprocessingModelCreate{
					SplitLength:  StringToInt64(indexDetail.TextPreprocessing.SplitLength),
					SplitMethod:  types.StringValue(indexDetail.TextPreprocessing.SplitMethod),
					SplitOverlap: StringToInt64(indexDetail.TextPreprocessing.SplitOverlap),
				},
				ImagePreprocessing: ImagePreprocessingModel{
					PatchMethod: types.StringValue(indexDetail.ImagePreprocessing.PatchMethod),
				},
				AnnParameters: AnnParametersModelCreate{
					SpaceType: types.StringValue(indexDetail.AnnParameters.SpaceType),
					Parameters: ParametersModel{
						EfConstruction: StringToInt64(indexDetail.AnnParameters.Parameters.EfConstruction),
						M:              StringToInt64(indexDetail.AnnParameters.Parameters.M),
					},
				},
				FilterStringMaxLength: StringToInt64(indexDetail.FilterStringMaxLength),
			}
			fmt.Print("tensorFields: ", indexDetail.TensorFields)
			break
		}
	}

	// implement deletion of state if resource no longer exists in cloud

	// Set the updated state
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

	// Construct settings map
	settings := map[string]interface{}{
		"type":                         model.Settings.Type.ValueString(),
		"vectorNumericType":            model.Settings.VectorNumericType.ValueString(),
		"treatUrlsAndPointersAsImages": model.Settings.TreatUrlsAndPointersAsImages.ValueBool(),
		"model":                        model.Settings.Model.ValueString(),
		"normalizeEmbeddings":          model.Settings.NormalizeEmbeddings.ValueBool(),
		"allFields":                    convertAllFieldsToMap(model.Settings.AllFields),
		"tensorFields":                 model.Settings.TensorFields,
		"inferenceType":                model.Settings.InferenceType.ValueString(),
		"numberOfInferences":           model.Settings.NumberOfInferences.ValueInt64(),
		"storageClass":                 model.Settings.StorageClass.ValueString(),
		"numberOfShards":               model.Settings.NumberOfShards.ValueInt64(),
		"numberOfReplicas":             model.Settings.NumberOfReplicas.ValueInt64(),
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
	if imagePreprocessing, ok := settings["imagePreprocessing"].(map[string]interface{}); ok {
		if model.Settings.ImagePreprocessing.PatchMethod.IsNull() {
			delete(imagePreprocessing, "patchMethod")
		}
		if len(imagePreprocessing) == 0 {
			delete(settings, "imagePreprocessing")
		}
	}
	if model.Settings.InferenceType.IsNull() {
		delete(settings, "inferenceType")
	}
	if model.Settings.NumberOfInferences.IsNull() {
		delete(settings, "numberOfInferences")
	}
	if model.Settings.StorageClass.IsNull() {
		delete(settings, "storageClass")
	}
	if model.Settings.NumberOfShards.IsNull() {
		delete(settings, "numberOfShards")
	}
	if model.Settings.NumberOfReplicas.IsNull() {
		delete(settings, "numberOfReplicas")
	}
	if len(model.Settings.AllFields) == 0 {
		delete(settings, "allFields")
	}
	if len(model.Settings.TensorFields) == 0 {
		delete(settings, "tensorFields")
	}

	//indexNameAsString := model.IndexName.

	// Adjust settings for structured index
	if model.Settings.Type.ValueString() == "structured" {
		allFields, err := validateAndConstructAllFields(model.Settings.AllFields)
		if err != nil {
			resp.Diagnostics.AddError("Invalid allFields", "Error validating allFields: "+err.Error())
			return
		}
		settings["allFields"] = allFields

		//if len(model.Settings.TensorFields) > 0 {
		//	tensorFields, err := constructTensorFields(model.Settings.TensorFields)
		//	if err != nil {
		//		resp.Diagnostics.AddError("Invalid tensorFields", "Error validating tensorFields: "+err.Error())
		//		return
		//	}
		//	settings["tensorFields"] = tensorFields
		//}
	}

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
	var model IndexResourceModel
	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Construct settings map
	settings := map[string]interface{}{
		"inferenceType":      model.Settings.InferenceType.ValueString(),
		"numberOfInferences": model.Settings.NumberOfInferences.ValueInt64(),
	}

	if model.Settings.InferenceType.IsNull() {
		delete(settings, "inferenceType")
	}
	if model.Settings.NumberOfInferences.IsNull() {
		delete(settings, "numberOfInferences")
	}

	//indexNameAsString := model.IndexName.

	err := r.marqoClient.UpdateIndex(model.IndexName.ValueString(), settings)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Update Index", "Could not create index: "+err.Error())
		return
	}

	// Set the index name as the ID in the Terraform state
	model.ID = model.IndexName
	diags = resp.State.Set(ctx, &model)
	resp.Diagnostics.Append(diags...)
}