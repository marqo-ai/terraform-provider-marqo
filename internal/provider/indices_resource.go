package provider

import (
	"context"
	"fmt"
	"marqo/go_marqo"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &indicesResource{}
	_ resource.ResourceWithConfigure   = &indicesResource{}
	_ resource.ResourceWithImportState = &indicesResource{}
)

// ManageIndicesResource is a helper function to simplify the provider implementation.
func ManageIndicesResource() resource.Resource {
	return &indicesResource{}
}

// orderResource is the resource implementation.
type indicesResource struct {
	marqoClient *go_marqo.Client
}

// IndexResourceModel maps the resource schema data.
type IndexResourceModel struct {
	IndexName     types.String       `tfsdk:"index_name"`
	Settings      IndexSettingsModel `tfsdk:"settings"`
	MarqoEndpoint types.String       `tfsdk:"marqo_endpoint"`
	Timeouts      *timeouts          `tfsdk:"timeouts"`
}

type timeouts struct {
	Create types.String `tfsdk:"create"`
	Update types.String `tfsdk:"update"`
	Delete types.String `tfsdk:"delete"`
}

type IndexSettingsModel struct {
	Type                         types.String                   `tfsdk:"type"`
	VectorNumericType            types.String                   `tfsdk:"vector_numeric_type"`
	NumberOfInferences           types.Int64                    `tfsdk:"number_of_inferences"`
	AllFields                    []AllFieldInput                `tfsdk:"all_fields"`
	TensorFields                 []string                       `tfsdk:"tensor_fields"`
	InferenceType                types.String                   `tfsdk:"inference_type"`
	StorageClass                 types.String                   `tfsdk:"storage_class"`
	NumberOfShards               types.Int64                    `tfsdk:"number_of_shards"`
	NumberOfReplicas             types.Int64                    `tfsdk:"number_of_replicas"`
	TreatUrlsAndPointersAsImages types.Bool                     `tfsdk:"treat_urls_and_pointers_as_images"`
	TreatUrlsAndPointersAsMedia  types.Bool                     `tfsdk:"treat_urls_and_pointers_as_media"`
	Model                        types.String                   `tfsdk:"model"`
	ModelProperties              *ModelPropertiesModelCreate    `tfsdk:"model_properties"`
	NormalizeEmbeddings          types.Bool                     `tfsdk:"normalize_embeddings"`
	TextPreprocessing            *TextPreprocessingModelCreate  `tfsdk:"text_preprocessing"`
	ImagePreprocessing           *ImagePreprocessingModel       `tfsdk:"image_preprocessing"`
	VideoPreprocessing           *VideoPreprocessingModelCreate `tfsdk:"video_preprocessing"`
	AudioPreprocessing           *AudioPreprocessingModelCreate `tfsdk:"audio_preprocessing"`
	AnnParameters                *AnnParametersModelCreate      `tfsdk:"ann_parameters"`
	FilterStringMaxLength        types.Int64                    `tfsdk:"filter_string_max_length"`
}

type ModelPropertiesModelCreate struct {
	Name             types.String        `tfsdk:"name"`
	Dimensions       types.Int64         `tfsdk:"dimensions"`
	Type             types.String        `tfsdk:"type"`
	Tokens           types.Int64         `tfsdk:"tokens"`
	ModelLocation    *ModelLocationModel `tfsdk:"model_location"`
	Url              types.String        `tfsdk:"url"`
	TrustRemoteCode  types.Bool          `tfsdk:"trust_remote_code"`
	IsMarqtunedModel types.Bool          `tfsdk:"is_marqtuned_model"`
}

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

type VideoPreprocessingModelCreate struct {
	SplitLength  types.Int64 `tfsdk:"split_length"`
	SplitOverlap types.Int64 `tfsdk:"split_overlap"`
}

type AudioPreprocessingModelCreate struct {
	SplitLength  types.Int64 `tfsdk:"split_length"`
	SplitOverlap types.Int64 `tfsdk:"split_overlap"`
}

type AnnParametersModelCreate struct {
	SpaceType  types.String     `tfsdk:"space_type"`
	Parameters *ParametersModel `tfsdk:"parameters"`
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

	client, ok := req.ProviderData.(*go_marqo.Client)

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
			"index_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the index.",
			},
			"marqo_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "The Marqo endpoint used by the index",
			},
			"timeouts": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"create": schema.StringAttribute{
						Optional:    true,
						Description: "Time to wait for index to be ready (e.g., '30m', '1h'). Default is 30m.",
					},
					"update": schema.StringAttribute{
						Optional:    true,
						Description: "Time to wait for index to be updated (e.g., '30m', '1h'). Default is 30m.",
					},
					"delete": schema.StringAttribute{
						Optional:    true,
						Description: "Time to wait for index to be deleted (e.g., '15m', '1h'). Default is 15m.",
					},
				},
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
						Required: true,
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
									ElementType: types.Float64Type,
								},
							},
						},
					},
					"tensor_fields": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
					"inference_type": schema.StringAttribute{
						Required: true,
					},
					"storage_class": schema.StringAttribute{
						Required: true,
					},
					"number_of_shards": schema.Int64Attribute{
						Required: true,
					},
					"number_of_replicas": schema.Int64Attribute{
						Required: true,
					},
					"treat_urls_and_pointers_as_images": schema.BoolAttribute{
						Optional: true,
					},
					"treat_urls_and_pointers_as_media": schema.BoolAttribute{
						Optional: true,
					},
					"model": schema.StringAttribute{
						Required: true,
					},
					"model_properties": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"name":       schema.StringAttribute{Optional: true},
							"dimensions": schema.Int64Attribute{Optional: true},
							"type":       schema.StringAttribute{Optional: true},
							"tokens":     schema.Int64Attribute{Optional: true},
							"model_location": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"s3": schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{
											"bucket": schema.StringAttribute{Optional: true},
											"key":    schema.StringAttribute{Optional: true},
										},
									},
									"hf": schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{
											"repo_id":  schema.StringAttribute{Optional: true},
											"filename": schema.StringAttribute{Optional: true},
										},
									},
									"auth_required": schema.BoolAttribute{Optional: true},
								},
							},
							"url":                schema.StringAttribute{Optional: true},
							"trust_remote_code":  schema.BoolAttribute{Optional: true},
							"is_marqtuned_model": schema.BoolAttribute{Optional: true},
						},
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
					"video_preprocessing": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"split_length":  schema.Int64Attribute{Optional: true},
							"split_overlap": schema.Int64Attribute{Optional: true},
						},
					},
					"audio_preprocessing": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"split_length":  schema.Int64Attribute{Optional: true},
							"split_overlap": schema.Int64Attribute{Optional: true},
						},
					},
					"ann_parameters": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"space_type": schema.StringAttribute{
								Optional: true,
							},
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
		if field.Name.IsNull() || field.Type.IsNull() {
			return nil, fmt.Errorf("each field must have a name and type")
		}
		fieldMap := map[string]interface{}{
			"name":            field.Name.ValueString(),
			"type":            field.Type.ValueString(),
			"features":        []string{},
			"dependentFields": map[string]float64{},
		}

		if len(field.Features) > 0 {
			features := []string{}
			for _, feature := range field.Features {
				features = append(features, feature.ValueString())
			}
			fieldMap["features"] = features
		}

		if len(field.DependentFields) > 0 {
			dependentFields := make(map[string]float64)
			for key, value := range field.DependentFields {
				dependentFields[key] = value.ValueFloat64()
			}
			fieldMap["dependentFields"] = dependentFields
		}

		allFields = append(allFields, fieldMap)
	}
	return allFields, nil
}

// Utility function to convert []AllFieldInput to a format suitable for settings map.
func convertAllFieldsToMap(allFieldsInput []AllFieldInput) []map[string]interface{} {
	allFields := []map[string]interface{}{}
	for _, field := range allFieldsInput {
		fieldMap := map[string]interface{}{
			"name":            field.Name.ValueString(),
			"type":            field.Type.ValueString(),
			"features":        []string{},
			"dependentFields": map[string]float64{},
		}

		if len(field.Features) > 0 {
			features := []string{}
			for _, feature := range field.Features {
				features = append(features, feature.ValueString())
			}
			fieldMap["features"] = features
		}

		if len(field.DependentFields) > 0 {
			dependentFieldsMap := make(map[string]float64)
			for key, value := range field.DependentFields {
				dependentFieldsMap[key] = value.ValueFloat64()
			}
			fieldMap["dependentFields"] = dependentFieldsMap
		}

		allFields = append(allFields, fieldMap)
	}
	return allFields
}

func convertModelLocationToAPI(modelLocation *ModelLocationModel) map[string]interface{} {
	if modelLocation == nil {
		return nil
	}

	result := map[string]interface{}{
		"authRequired": modelLocation.AuthRequired.ValueBool(),
	}

	if modelLocation.S3 != nil {
		result["s3"] = map[string]interface{}{
			"bucket": modelLocation.S3.Bucket.ValueString(),
			"key":    modelLocation.S3.Key.ValueString(),
		}
	}

	if modelLocation.Hf != nil {
		result["hf"] = map[string]interface{}{
			"repoId":   modelLocation.Hf.RepoId.ValueString(),
			"filename": modelLocation.Hf.Filename.ValueString(),
		}
	}

	return result
}

func (m *ModelPropertiesModelCreate) IsEmpty() bool {
	if m == nil {
		return true
	}
	return m.Name.IsNull() &&
		m.Dimensions.IsNull() &&
		m.Type.IsNull() &&
		m.Tokens.IsNull() &&
		m.Url.IsNull() &&
		(m.ModelLocation == nil || m.ModelLocation.IsEmpty())
}

func (m *ModelLocationModel) IsEmpty() bool {
	if m == nil {
		return true
	}
	return m.AuthRequired.IsNull() &&
		(m.S3 == nil || (m.S3.Bucket.IsNull() && m.S3.Key.IsNull())) &&
		(m.Hf == nil || (m.Hf.RepoId.IsNull() && m.Hf.Filename.IsNull()))
}

func convertModelPropertiesToResource(props *go_marqo.ModelProperties) *ModelPropertiesModelCreate {
	if props == nil {
		return nil
	}

	model := &ModelPropertiesModelCreate{}

	// Convert only non-empty values
	if props.Name != "" {
		model.Name = types.StringValue(props.Name)
	}
	if props.Dimensions != 0 {
		model.Dimensions = types.Int64Value(props.Dimensions)
	}
	if props.Type != "" {
		model.Type = types.StringValue(props.Type)
	}
	if props.Tokens != 0 {
		model.Tokens = types.Int64Value(props.Tokens)
	}
	if props.Url != "" {
		model.Url = types.StringValue(props.Url)
	}
	if props.TrustRemoteCode {
		model.TrustRemoteCode = types.BoolValue(true)
	}
	if props.IsMarqtunedModel {
		model.IsMarqtunedModel = types.BoolValue(true)
	}
	// Only convert ModelLocation if it has non-null values
	if loc := convertModelLocation(props.ModelLocation); loc != nil {
		model.ModelLocation = loc
	}

	// Only return the model if it's not empty.
	if model.IsEmpty() {
		return nil
	}

	return model
}

func (r *indicesResource) findAndCreateState(indices []go_marqo.IndexDetail, indexName string, existingTimeouts *timeouts) (*IndexResourceModel, bool) {
	for _, indexDetail := range indices {
		if indexDetail.IndexName == indexName {
			// Create a new model with proper null handling
			model := &IndexResourceModel{
				IndexName:     types.StringValue(indexDetail.IndexName),
				MarqoEndpoint: types.StringValue(indexDetail.MarqoEndpoint),
				Timeouts:      existingTimeouts,
				Settings: IndexSettingsModel{
					Type:               types.StringValue(indexDetail.Type),
					InferenceType:      types.StringValue(indexDetail.InferenceType),
					NumberOfInferences: types.Int64Value(indexDetail.NumberOfInferences),
					StorageClass:       types.StringValue(indexDetail.StorageClass),
					NumberOfShards:     types.Int64Value(indexDetail.NumberOfShards),
					NumberOfReplicas:   types.Int64Value(indexDetail.NumberOfReplicas),
					Model:              types.StringValue(indexDetail.Model),
				},
			}

			// Handle Bool fields
			if indexDetail.TreatUrlsAndPointersAsImages == nil {
				model.Settings.TreatUrlsAndPointersAsImages = types.BoolNull()
			} else {
				model.Settings.TreatUrlsAndPointersAsImages = types.BoolValue(*indexDetail.TreatUrlsAndPointersAsImages)
			}
			if indexDetail.TreatUrlsAndPointersAsMedia == nil {
				model.Settings.TreatUrlsAndPointersAsMedia = types.BoolNull()
			} else {
				model.Settings.TreatUrlsAndPointersAsMedia = types.BoolValue(*indexDetail.TreatUrlsAndPointersAsMedia)
			}
			if indexDetail.NormalizeEmbeddings == nil {
				model.Settings.NormalizeEmbeddings = types.BoolNull()
			} else {
				model.Settings.NormalizeEmbeddings = types.BoolValue(*indexDetail.NormalizeEmbeddings)
			}

			// Handle optional string fields
			if indexDetail.VectorNumericType != "" {
				model.Settings.VectorNumericType = types.StringValue(indexDetail.VectorNumericType)
			} else {
				model.Settings.VectorNumericType = types.StringNull()
			}

			// Handle optional numeric fields
			if indexDetail.FilterStringMaxLength > 0 {
				model.Settings.FilterStringMaxLength = types.Int64Value(indexDetail.FilterStringMaxLength)
			} else {
				model.Settings.FilterStringMaxLength = types.Int64Null()
			}

			// Handle model properties
			if !reflect.DeepEqual(indexDetail.ModelProperties, go_marqo.ModelProperties{}) {
				model.Settings.ModelProperties = convertModelPropertiesToResource(&indexDetail.ModelProperties)
			} else {
				model.Settings.ModelProperties = nil
			}

			// Handle AllFields
			if len(indexDetail.AllFields) > 0 {
				model.Settings.AllFields = ConvertMarqoAllFieldInputs(indexDetail.AllFields)
			} else {
				model.Settings.AllFields = nil
			}

			// Handle TensorFields
			if len(indexDetail.TensorFields) > 0 {
				model.Settings.TensorFields = indexDetail.TensorFields
			} else {
				model.Settings.TensorFields = nil
			}

			// Handle TextPreprocessing
			if indexDetail.TextPreprocessing == (go_marqo.TextPreprocessing{}) {
				model.Settings.TextPreprocessing = nil
			} else {
				model.Settings.TextPreprocessing = &TextPreprocessingModelCreate{
					SplitLength:  types.Int64Value(indexDetail.TextPreprocessing.SplitLength),
					SplitMethod:  types.StringValue(indexDetail.TextPreprocessing.SplitMethod),
					SplitOverlap: types.Int64Value(indexDetail.TextPreprocessing.SplitOverlap),
				}
			}

			// Handle ImagePreprocessing
			if indexDetail.ImagePreprocessing.PatchMethod == "" {
				model.Settings.ImagePreprocessing = nil
			} else {
				model.Settings.ImagePreprocessing = &ImagePreprocessingModel{
					PatchMethod: types.StringValue(indexDetail.ImagePreprocessing.PatchMethod),
				}
			}

			// Handle VideoPreprocessing
			if indexDetail.VideoPreprocessing.SplitLength > 0 || indexDetail.VideoPreprocessing.SplitOverlap > 0 {
				model.Settings.VideoPreprocessing = &VideoPreprocessingModelCreate{
					SplitLength:  types.Int64Value(indexDetail.VideoPreprocessing.SplitLength),
					SplitOverlap: types.Int64Value(indexDetail.VideoPreprocessing.SplitOverlap),
				}
			} else {
				model.Settings.VideoPreprocessing = nil
			}

			// Handle AudioPreprocessing
			if indexDetail.AudioPreprocessing.SplitLength > 0 || indexDetail.AudioPreprocessing.SplitOverlap > 0 {
				model.Settings.AudioPreprocessing = &AudioPreprocessingModelCreate{
					SplitLength:  types.Int64Value(indexDetail.AudioPreprocessing.SplitLength),
					SplitOverlap: types.Int64Value(indexDetail.AudioPreprocessing.SplitOverlap),
				}
			} else {
				model.Settings.AudioPreprocessing = nil
			}

			// Handle AnnParameters
			if indexDetail.AnnParameters.SpaceType != "" ||
				indexDetail.AnnParameters.Parameters.EfConstruction > 0 ||
				indexDetail.AnnParameters.Parameters.M > 0 {
				model.Settings.AnnParameters = &AnnParametersModelCreate{
					SpaceType: types.StringValue(indexDetail.AnnParameters.SpaceType),
					Parameters: &ParametersModel{
						EfConstruction: types.Int64Value(indexDetail.AnnParameters.Parameters.EfConstruction),
						M:              types.Int64Value(indexDetail.AnnParameters.Parameters.M),
					},
				}
			} else {
				model.Settings.AnnParameters = nil
			}

			return model, true
		}
	}
	return nil, false
}

func (r *indicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Initialize the state variable based on the IndexResourceModel
	var state IndexResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if this is likely an import operation
	isImport := state.MarqoEndpoint.IsNull() &&
		(state.Settings.Type.ValueString() == "" ||
			state.Settings.InferenceType.ValueString() == "" ||
			state.Settings.NumberOfInferences.ValueInt64() == 0)

	if isImport {
		tflog.Info(ctx, fmt.Sprintf("Detected import operation for index %s", state.IndexName.ValueString()))
	}

	tflog.Debug(ctx, "Calling marqo client ListIndices")
	indices, err := r.marqoClient.ListIndices()
	if err != nil {
		resp.Diagnostics.AddError("Failed to List Indices", fmt.Sprintf("Could not list indices: %s", err.Error()))
		return
	}

	newState, found := r.findAndCreateState(indices, state.IndexName.ValueString(), state.Timeouts)

	// Handle inference_type field
	if newState != nil {
		inferenceTypeMap := map[string]string{
			"CPU":       "marqo.CPU.large", // verify this
			"CPU.SMALL": "marqo.CPU.small",
			"CPU.LARGE": "marqo.CPU.large",
			"GPU":       "marqo.GPU",
		}

		storaceClassMap := map[string]string{
			"BASIC":       "marqo.basic",
			"BALANCED":    "marqo.balanced",
			"PERFORMANCE": "marqo.performance",
		}

		if !newState.Settings.InferenceType.IsNull() {
			currentValue := newState.Settings.InferenceType.ValueString()
			if mappedValue, exists := inferenceTypeMap[currentValue]; exists {
				newState.Settings.InferenceType = types.StringValue(mappedValue)
			}
		}

		if !newState.Settings.StorageClass.IsNull() {
			currentValue := newState.Settings.StorageClass.ValueString()
			if mappedValue, exists := storaceClassMap[currentValue]; exists {
				newState.Settings.StorageClass = types.StringValue(mappedValue)
			}
		}

		// marqo doesn't return timeouts, so we maintain the existing state
		newState.Timeouts = state.Timeouts

		// Special handling for import case - if this is a new import (state has empty values)
		// we need to ensure consistent null values
		if isImport {
			// For import operations, we want to set all optional fields to null
			// unless they have explicit non-default values

			// Handle optional scalar fields
			if newState.Settings.FilterStringMaxLength.ValueInt64() == 0 {
				newState.Settings.FilterStringMaxLength = types.Int64Null()
			}

			// Set empty string vector_numeric_type to null
			if newState.Settings.VectorNumericType.ValueString() == "" {
				newState.Settings.VectorNumericType = types.StringNull()
			}

			// Set empty objects to null
			if newState.Settings.AllFields != nil && len(newState.Settings.AllFields) == 0 {
				newState.Settings.AllFields = nil
			}
			if newState.Settings.TensorFields != nil && len(newState.Settings.TensorFields) == 0 {
				newState.Settings.TensorFields = nil
			}

			// Handle optional nested objects
			if newState.Settings.TextPreprocessing != nil &&
				newState.Settings.TextPreprocessing.SplitLength.ValueInt64() == 0 &&
				newState.Settings.TextPreprocessing.SplitMethod.ValueString() == "" &&
				newState.Settings.TextPreprocessing.SplitOverlap.ValueInt64() == 0 {
				newState.Settings.TextPreprocessing = nil
			}

			if newState.Settings.ImagePreprocessing != nil &&
				newState.Settings.ImagePreprocessing.PatchMethod.ValueString() == "" {
				newState.Settings.ImagePreprocessing.PatchMethod = types.StringNull()
			}

			if newState.Settings.VideoPreprocessing != nil &&
				newState.Settings.VideoPreprocessing.SplitLength.ValueInt64() == 0 &&
				newState.Settings.VideoPreprocessing.SplitOverlap.ValueInt64() == 0 {
				newState.Settings.VideoPreprocessing = nil
			}

			if newState.Settings.AudioPreprocessing != nil &&
				newState.Settings.AudioPreprocessing.SplitLength.ValueInt64() == 0 &&
				newState.Settings.AudioPreprocessing.SplitOverlap.ValueInt64() == 0 {
				newState.Settings.AudioPreprocessing = nil
			}

			if newState.Settings.AnnParameters != nil &&
				newState.Settings.AnnParameters.SpaceType.ValueString() == "" &&
				newState.Settings.AnnParameters.Parameters != nil &&
				newState.Settings.AnnParameters.Parameters.EfConstruction.ValueInt64() == 0 &&
				newState.Settings.AnnParameters.Parameters.M.ValueInt64() == 0 {
				newState.Settings.AnnParameters = nil
			}

			if newState.Settings.ModelProperties != nil && newState.Settings.ModelProperties.IsEmpty() {
				newState.Settings.ModelProperties = nil
			}

			// Image Preprocessing
			if state.Settings.ImagePreprocessing == nil {
				newState.Settings.ImagePreprocessing = nil
			} else {
				if newState.Settings.ImagePreprocessing == nil {
					newState.Settings.ImagePreprocessing = &ImagePreprocessingModel{}
				}
				if newState.Settings.ImagePreprocessing.PatchMethod.ValueString() == "" {
					newState.Settings.ImagePreprocessing.PatchMethod = types.StringNull()
				}

			}

			// preserve the video/audio preprocessing from current state since api does not return them
			if state.Settings.VideoPreprocessing != nil {
				newState.Settings.VideoPreprocessing = state.Settings.VideoPreprocessing
			} else if newState.Settings.VideoPreprocessing != nil {
				// If not in state but returned by API with zero values, set to null
				// Doesnt do anything now because API doesn't return video/audio preprocessing
				if newState.Settings.VideoPreprocessing.SplitLength.ValueInt64() == 0 &&
					newState.Settings.VideoPreprocessing.SplitOverlap.ValueInt64() == 0 {
					newState.Settings.VideoPreprocessing = nil
				}
			}

			if state.Settings.AudioPreprocessing != nil {
				newState.Settings.AudioPreprocessing = state.Settings.AudioPreprocessing
			} else if newState.Settings.AudioPreprocessing != nil {
				// If not in state but returned by API with zero values, set to null
				// Doesnt do anything now because API doesn't return video/audio preprocessing
				if newState.Settings.AudioPreprocessing.SplitLength.ValueInt64() == 0 &&
					newState.Settings.AudioPreprocessing.SplitOverlap.ValueInt64() == 0 {
					newState.Settings.AudioPreprocessing = nil
				}
			}
		} else {
			// For non-import operations, preserve values from the existing state

			// Handle AllFields
			if state.Settings.AllFields == nil {
				newState.Settings.AllFields = nil
			} else if len(newState.Settings.AllFields) == 0 {
				newState.Settings.AllFields = []AllFieldInput{}
			} else {
				// Ensure features and dependent_fields are always set
				for i := range newState.Settings.AllFields {
					if len(newState.Settings.AllFields[i].Features) == 0 {
						newState.Settings.AllFields[i].Features = nil
					}
					if len(newState.Settings.AllFields[i].DependentFields) == 0 {
						newState.Settings.AllFields[i].DependentFields = nil
					}
				}
			}

			// Handle TensorFields
			if len(state.Settings.TensorFields) == 0 {
				newState.Settings.TensorFields = nil
			}

			// Handle optional nested objects
			// If these fields are not set in the state, set them to null
			if state.Settings.TextPreprocessing == nil {
				newState.Settings.TextPreprocessing = nil
			}
			if state.Settings.AnnParameters == nil {
				newState.Settings.AnnParameters = nil
			}

			// Handle optional scalar fields
			// If these fields are null in the state, keep them null
			if state.Settings.FilterStringMaxLength.IsNull() {
				newState.Settings.FilterStringMaxLength = types.Int64Null()
			}

			// For boolean fields, if they're explicitly set in the configuration,
			// use those values; otherwise keep them null
			if !state.Settings.NormalizeEmbeddings.IsNull() {
				newState.Settings.NormalizeEmbeddings = state.Settings.NormalizeEmbeddings
			} else {
				newState.Settings.NormalizeEmbeddings = types.BoolNull()
			}
			if !state.Settings.TreatUrlsAndPointersAsImages.IsNull() {
				newState.Settings.TreatUrlsAndPointersAsImages = state.Settings.TreatUrlsAndPointersAsImages
			} else {
				newState.Settings.TreatUrlsAndPointersAsImages = types.BoolNull()
			}
			if !state.Settings.TreatUrlsAndPointersAsMedia.IsNull() {
				newState.Settings.TreatUrlsAndPointersAsMedia = state.Settings.TreatUrlsAndPointersAsMedia
			} else {
				newState.Settings.TreatUrlsAndPointersAsMedia = types.BoolNull()
			}

			if state.Settings.VectorNumericType.IsNull() {
				newState.Settings.VectorNumericType = types.StringNull()
			}

			// Handle image_preprocessing
			if state.Settings.ImagePreprocessing == nil {
				newState.Settings.ImagePreprocessing = nil
			} else {
				// Ensure we always have an ImagePreprocessing object if it was in the config
				if newState.Settings.ImagePreprocessing == nil {
					newState.Settings.ImagePreprocessing = &ImagePreprocessingModel{}
				}

				if newState.Settings.ImagePreprocessing.PatchMethod.ValueString() == "" {
					newState.Settings.ImagePreprocessing.PatchMethod = types.StringNull()
				}

			}

			// preserve the video/audio preprocessing from current state since api does not return them
			if state.Settings.VideoPreprocessing != nil {
				newState.Settings.VideoPreprocessing = state.Settings.VideoPreprocessing
			} else if newState.Settings.VideoPreprocessing != nil {
				// If not in state but returned by API with zero values, set to null
				if newState.Settings.VideoPreprocessing.SplitLength.ValueInt64() == 0 &&
					newState.Settings.VideoPreprocessing.SplitOverlap.ValueInt64() == 0 {
					newState.Settings.VideoPreprocessing = nil
				}
			}

			if state.Settings.AudioPreprocessing != nil {
				newState.Settings.AudioPreprocessing = state.Settings.AudioPreprocessing
			} else if newState.Settings.AudioPreprocessing != nil {
				// If not in state but returned by API with zero values, set to null
				if newState.Settings.AudioPreprocessing.SplitLength.ValueInt64() == 0 &&
					newState.Settings.AudioPreprocessing.SplitOverlap.ValueInt64() == 0 {
					newState.Settings.AudioPreprocessing = nil
				}
			}
		}

		// Ignore these fields for structured indexes
		if newState.Settings.Type.ValueString() == "structured" {
			newState.Settings.FilterStringMaxLength = types.Int64Null()
			newState.Settings.TreatUrlsAndPointersAsImages = types.BoolNull()
			newState.Settings.TreatUrlsAndPointersAsMedia = types.BoolNull()
		}

		// Remove null fields
		if newState.Settings.InferenceType.IsNull() {
			newState.Settings.InferenceType = types.StringNull()
		}
	}

	// if index no longer exists in cloud, delete the state
	if !found {
		resp.Diagnostics.AddWarning("Resource Not Found", "The specified index does not exist in the cloud. The state will be deleted.")
		// Then Totally Remove from terraform resources
		resp.State.RemoveResource(ctx)
		return
	}

	// Set the updated state
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Standalone function to compare states.
func statesAreEqual(existing *IndexResourceModel, desired *IndexResourceModel) bool {
	// A deep comparison between existing and desired states
	return reflect.DeepEqual(existing.Settings, desired.Settings)
}

// waitForIndexStatus waits for an index to reach a target status or be deleted.
func (r *indicesResource) waitForIndexStatus(ctx context.Context, indexName string, targetStatus string, timeoutDuration time.Duration, isDelete bool) error {
	timeout := time.After(timeoutDuration)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	start := time.Now()

	// For delete operations, check if index is in READY state first
	if isDelete {
		indices, err := r.marqoClient.ListIndices()
		if err != nil {
			return fmt.Errorf("error checking index status before deletion: %v", err)
		}

		for _, index := range indices {
			if index.IndexName == indexName {
				if index.IndexStatus != "READY" && index.IndexStatus != "DELETING" {
					return fmt.Errorf("cannot delete index %s: index is in %s state, must be in READY state",
						indexName, index.IndexStatus)
				}
				break
			}
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Waiting up to %v for index %s to reach %s...",
		timeoutDuration,
		indexName,
		func() string {
			if isDelete {
				return "deletion"
			}
			return fmt.Sprintf("status %s", targetStatus)
		}()))

	for {
		select {
		case <-timeout:
			indices, err := r.marqoClient.ListIndices()
			if err != nil {
				return fmt.Errorf("timeout checking final status: %v", err)
			}

			var currentStatus string
			exists := false
			for _, index := range indices {
				if index.IndexName == indexName {
					currentStatus = index.IndexStatus
					exists = true
					break
				}
			}

			if isDelete {
				if !exists {
					return nil
				}
				return fmt.Errorf("index %s still exists after %v (status: %s)",
					indexName, timeoutDuration, currentStatus)
			}

			if !exists {
				return fmt.Errorf("index %s no longer exists while waiting for status %s",
					indexName, targetStatus)
			}

			return fmt.Errorf("timeout waiting for index %s to reach status %s after %v - current status is %s",
				indexName, targetStatus, timeoutDuration, currentStatus)

		case <-ticker.C:
			indices, err := r.marqoClient.ListIndices()
			if err != nil {
				tflog.Error(ctx, fmt.Sprintf("Error checking index status: %s", err))
				continue
			}

			// For delete operations, we check if the index no longer exists
			if isDelete {
				indexExists := false
				var currentStatus string
				for _, index := range indices {
					if index.IndexName == indexName {
						indexExists = true
						currentStatus = index.IndexStatus
						break
					}
				}
				if !indexExists {
					tflog.Info(ctx, fmt.Sprintf("Index %s has been successfully deleted (total time: %v)",
						indexName, time.Since(start)))
					return nil
				}
				tflog.Info(ctx, fmt.Sprintf("Index %s still exists with status %s, continuing to wait... (elapsed: %v)",
					indexName, currentStatus, time.Since(start)))
				continue
			}

			// For create/update operations, we check for the target status
			for _, index := range indices {
				if index.IndexName == indexName {
					tflog.Info(ctx, fmt.Sprintf("Index %s status: %s (elapsed: %v)",
						indexName, index.IndexStatus, time.Since(start)))

					if index.IndexStatus == targetStatus {
						tflog.Info(ctx, fmt.Sprintf("Index %s has reached status %s (total time: %v)",
							indexName, targetStatus, time.Since(start)))
						return nil
					} else if index.IndexStatus == "FAILED" {
						return fmt.Errorf("index %s reached FAILED status while waiting for %s",
							indexName, targetStatus)
					}
					break
				}
			}
		}
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
		"type":                  model.Settings.Type.ValueString(),
		"vectorNumericType":     model.Settings.VectorNumericType.ValueString(),
		"model":                 model.Settings.Model.ValueString(),
		"modelProperties":       model.Settings.ModelProperties,
		"allFields":             convertAllFieldsToMap(model.Settings.AllFields),
		"tensorFields":          model.Settings.TensorFields,
		"inferenceType":         model.Settings.InferenceType.ValueString(),
		"numberOfInferences":    model.Settings.NumberOfInferences.ValueInt64(),
		"storageClass":          model.Settings.StorageClass.ValueString(),
		"numberOfShards":        model.Settings.NumberOfShards.ValueInt64(),
		"numberOfReplicas":      model.Settings.NumberOfReplicas.ValueInt64(),
		"filterStringMaxLength": model.Settings.FilterStringMaxLength.ValueInt64(),
	}

	if !model.Settings.TreatUrlsAndPointersAsImages.IsNull() {
		settings["treatUrlsAndPointersAsImages"] = model.Settings.TreatUrlsAndPointersAsImages.ValueBool()
	}
	if !model.Settings.TreatUrlsAndPointersAsMedia.IsNull() {
		settings["treatUrlsAndPointersAsMedia"] = model.Settings.TreatUrlsAndPointersAsMedia.ValueBool()
	}
	if !model.Settings.NormalizeEmbeddings.IsNull() {
		settings["normalizeEmbeddings"] = model.Settings.NormalizeEmbeddings.ValueBool()
	}

	// Optional dictionary fields
	if model.Settings.ModelProperties != nil {
		modelPropertiesMap := map[string]interface{}{
			"name":             model.Settings.ModelProperties.Name.ValueString(),
			"dimensions":       model.Settings.ModelProperties.Dimensions.ValueInt64(),
			"type":             model.Settings.ModelProperties.Type.ValueString(),
			"tokens":           model.Settings.ModelProperties.Tokens.ValueInt64(),
			"url":              model.Settings.ModelProperties.Url.ValueString(),
			"trustRemoteCode":  model.Settings.ModelProperties.TrustRemoteCode.ValueBool(),
			"isMarqtunedModel": model.Settings.ModelProperties.IsMarqtunedModel.ValueBool(),
		}

		if model.Settings.ModelProperties.ModelLocation != nil {
			modelPropertiesMap["modelLocation"] = convertModelLocationToAPI(model.Settings.ModelProperties.ModelLocation)
		}

		settings["modelProperties"] = modelPropertiesMap
	}
	if model.Settings.TextPreprocessing != nil {
		settings["textPreprocessing"] = map[string]interface{}{
			"splitLength":  model.Settings.TextPreprocessing.SplitLength.ValueInt64(),
			"splitMethod":  model.Settings.TextPreprocessing.SplitMethod.ValueString(),
			"splitOverlap": model.Settings.TextPreprocessing.SplitOverlap.ValueInt64(),
		}
	}
	if model.Settings.ImagePreprocessing != nil {
		if model.Settings.ImagePreprocessing.PatchMethod.IsNull() {
			settings["imagePreprocessing"] = nil

		} else {
			settings["imagePreprocessing"] = map[string]interface{}{
				"patchMethod": model.Settings.ImagePreprocessing.PatchMethod.ValueString(),
			}
		}
	}
	if model.Settings.VideoPreprocessing != nil {
		settings["videoPreprocessing"] = map[string]interface{}{
			"splitLength":  model.Settings.VideoPreprocessing.SplitLength.ValueInt64(),
			"splitOverlap": model.Settings.VideoPreprocessing.SplitOverlap.ValueInt64(),
		}
	}
	if model.Settings.AudioPreprocessing != nil {
		settings["audioPreprocessing"] = map[string]interface{}{
			"splitLength":  model.Settings.AudioPreprocessing.SplitLength.ValueInt64(),
			"splitOverlap": model.Settings.AudioPreprocessing.SplitOverlap.ValueInt64(),
		}
	}
	if model.Settings.AnnParameters != nil {
		settings["annParameters"] = map[string]interface{}{
			"spaceType": model.Settings.AnnParameters.SpaceType.ValueString(),
			"parameters": map[string]interface{}{
				"efConstruction": model.Settings.AnnParameters.Parameters.EfConstruction.ValueInt64(),
				"m":              model.Settings.AnnParameters.Parameters.M.ValueInt64(),
			},
		}
	}

	// Remove optional fields if they are not set
	if model.Settings.TreatUrlsAndPointersAsImages.IsNull() {
		delete(settings, "treatUrlsAndPointersAsImages")
	}
	if model.Settings.TreatUrlsAndPointersAsMedia.IsNull() {
		delete(settings, "treatUrlsAndPointersAsMedia")
	}
	if model.Settings.NormalizeEmbeddings.IsNull() {
		delete(settings, "normalizeEmbeddings")
	}
	if model.Settings.Model.IsNull() {
		delete(settings, "model")
	}
	if model.Settings.ModelProperties == nil ||
		(model.Settings.ModelProperties.Name.IsNull() &&
			model.Settings.ModelProperties.Dimensions.IsNull() &&
			model.Settings.ModelProperties.Type.IsNull() &&
			model.Settings.ModelProperties.Tokens.IsNull() &&
			model.Settings.ModelProperties.Url.IsNull() &&
			model.Settings.ModelProperties.TrustRemoteCode.IsNull()) {
		delete(settings, "modelProperties")
	}
	if model.Settings.TextPreprocessing == nil ||
		(model.Settings.TextPreprocessing.SplitLength.IsNull() &&
			model.Settings.TextPreprocessing.SplitMethod.IsNull() &&
			model.Settings.TextPreprocessing.SplitOverlap.IsNull()) {
		delete(settings, "textPreprocessing")
	}

	if model.Settings.VideoPreprocessing == nil ||
		(model.Settings.VideoPreprocessing.SplitLength.IsNull() &&
			model.Settings.VideoPreprocessing.SplitOverlap.IsNull()) {
		delete(settings, "videoPreprocessing")
	}
	if model.Settings.AudioPreprocessing == nil ||
		(model.Settings.AudioPreprocessing.SplitLength.IsNull() &&
			model.Settings.AudioPreprocessing.SplitOverlap.IsNull()) {
		delete(settings, "audioPreprocessing")
	}
	if model.Settings.ImagePreprocessing == nil ||
		(model.Settings.ImagePreprocessing.PatchMethod.IsNull()) {
		delete(settings, "imagePreprocessing")
	}
	if model.Settings.AnnParameters == nil ||
		(model.Settings.AnnParameters.SpaceType.IsNull() &&
			model.Settings.AnnParameters.Parameters.EfConstruction.IsNull() &&
			model.Settings.AnnParameters.Parameters.M.IsNull()) {
		delete(settings, "annParameters")
	}
	if model.Settings.VectorNumericType.IsNull() || model.Settings.VectorNumericType.ValueString() == "" {
		delete(settings, "vectorNumericType")
	}
	if len(model.Settings.AllFields) == 0 {
		delete(settings, "allFields")
	}
	if len(model.Settings.TensorFields) == 0 {
		delete(settings, "tensorFields")
	}
	if model.Settings.FilterStringMaxLength.IsNull() || model.Settings.FilterStringMaxLength.IsUnknown() {
		delete(settings, "filterStringMaxLength")
	}

	tflog.Debug(ctx, "Creating index with settings: %#v", settings)

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

	// Parse timeout duration
	timeoutDuration := 30 * time.Minute // default timeout
	if model.Timeouts != nil && model.Timeouts.Create.ValueString() != "" {
		parsedTimeout, err := time.ParseDuration(model.Timeouts.Create.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Timeout Duration",
				fmt.Sprintf("Could not parse timeout duration: %s. Expected format: '30m', '1h', etc.", err),
			)
			return
		}
		timeoutDuration = parsedTimeout
		tflog.Info(ctx, fmt.Sprintf("Using configured timeout of %v", timeoutDuration))
	}

	indexName := model.IndexName.ValueString()
	err := r.marqoClient.CreateIndex(indexName, settings)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			tflog.Info(ctx, fmt.Sprintf("Index %s already exists. Checking if it needs to be updated.", indexName))

			indices, err := r.marqoClient.ListIndices()
			if err != nil {
				resp.Diagnostics.AddError("Failed to List Indices", fmt.Sprintf("Could not list indices: %s", err.Error()))
				return
			}

			existingState, found := r.findAndCreateState(indices, indexName, model.Timeouts)
			if !found {
				resp.Diagnostics.AddError("Failed to Find Index", fmt.Sprintf("Index %s not found after creation", indexName))
				return
			}

			// Compare existing state with desired state
			if !statesAreEqual(existingState, &model) {
				// Attempt to update the existing index
				err = r.marqoClient.UpdateIndex(indexName, settings)
				if err != nil {
					resp.Diagnostics.AddError("Failed to Update Existing Index",
						fmt.Sprintf("Index %s exists but couldn't be updated to match the configuration: %s", indexName, err.Error()))
					return
				}
				tflog.Info(ctx, fmt.Sprintf("Index %s updated to match configuration.", indexName))
			} else {
				tflog.Info(ctx, fmt.Sprintf("Existing index %s matches configuration. No update needed.", indexName))
			}

			// Set state to the (potentially updated) existing index with preserved null values
			diags = resp.State.Set(ctx, existingState)
			resp.Diagnostics.Append(diags...)
			resp.Diagnostics.AddWarning(fmt.Sprintf("Index %s already existed and has been imported into Terraform state.", indexName),
				"Any differences between the existing index and your configuration have been resolved by updating the index.")
			return
		}

		resp.Diagnostics.AddError("Failed to Create Index", "Could not create index: "+err.Error())
		return
	}

	// Set initial state
	model.MarqoEndpoint = types.StringValue("pending")
	diags = resp.State.Set(ctx, &model)
	resp.Diagnostics.Append(diags...)

	// Wait for the index to be ready
	err = r.waitForIndexStatus(ctx, indexName, "READY", timeoutDuration, false)
	if err != nil {
		// If waiting failed, attempt to clean up the index
		deleteErr := r.marqoClient.DeleteIndex(indexName)
		if deleteErr != nil {
			resp.Diagnostics.AddError(
				"Index Creation Failed and Cleanup Failed",
				fmt.Sprintf("Index %s creation failed: %s\n\nAttempted cleanup also failed: %s\nManual cleanup may be required.",
					indexName, err, deleteErr),
			)
		} else {
			resp.Diagnostics.AddError(
				"Index Creation Failed",
				fmt.Sprintf("Index %s creation failed and was automatically cleaned up.\nError: %s",
					indexName, err),
			)
		}
		return
	}

	// Do final read to get the complete state
	readResp := resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{State: resp.State}, &readResp)
	if readResp.Diagnostics.HasError() {
		resp.Diagnostics.Append(readResp.Diagnostics...)
		return
	}

	// Update the response state with the read state
	resp.State = readResp.State
}

func (r *indicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model IndexResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	indexName := model.IndexName.ValueString()

	// Default timeout of 15 minutes for deletion
	timeoutDuration := 15 * time.Minute
	if model.Timeouts != nil && model.Timeouts.Delete.ValueString() != "" {
		parsedTimeout, err := time.ParseDuration(model.Timeouts.Delete.ValueString())
		if err == nil {
			timeoutDuration = parsedTimeout
			tflog.Info(ctx, fmt.Sprintf("Using configured delete timeout of %v", timeoutDuration))
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Invalid delete timeout duration: %s, using default of 15m", err))
		}
	}

	// Attempt to delete the index
	err := r.marqoClient.DeleteIndex(indexName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Index",
			fmt.Sprintf("Could not delete index %s.\n"+
				"Error details: %s", indexName, err.Error()))
		return
	}

	// Wait for the index to be deleted
	err = r.waitForIndexStatus(ctx, indexName, "", timeoutDuration, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Timeout Waiting for Index Deletion",
			fmt.Sprintf("Index %s deletion did not complete within the timeout period: %s", indexName, err))
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *indicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model IndexResourceModel
	var state IndexResourceModel

	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that only allowed fields are being modified
	// The only modifiable fields are:
	// - inference_type
	// - number_of_inferences
	// - number_of_replicas (can only go up)
	// - number_of_shards (can only go up)

	// Check for changes in non-modifiable fields
	if model.Settings.Type.ValueString() != state.Settings.Type.ValueString() {
		resp.Diagnostics.AddError(
			"Cannot Modify Index Type",
			fmt.Sprintf("The index type cannot be modified from '%s' to '%s'. You must destroy and recreate the index to change this field.",
				state.Settings.Type.ValueString(), model.Settings.Type.ValueString()))
		return
	}

	if model.Settings.Model.ValueString() != state.Settings.Model.ValueString() {
		resp.Diagnostics.AddError(
			"Cannot Modify Index Model",
			fmt.Sprintf("The model cannot be modified from '%s' to '%s'. You must destroy and recreate the index to change this field.",
				state.Settings.Model.ValueString(), model.Settings.Model.ValueString()))
		return
	}

	if model.Settings.StorageClass.ValueString() != state.Settings.StorageClass.ValueString() {
		resp.Diagnostics.AddError(
			"Cannot Modify Storage Class",
			fmt.Sprintf("The storage class cannot be modified from '%s' to '%s'. You must destroy and recreate the index to change this field.",
				state.Settings.StorageClass.ValueString(), model.Settings.StorageClass.ValueString()))
		return
	}

	// Check for changes in other fields that are not modifiable
	// This is not an exhaustive list, but covers the most common fields
	if !reflect.DeepEqual(model.Settings.AllFields, state.Settings.AllFields) {
		resp.Diagnostics.AddError(
			"Cannot Modify All Fields",
			"The all_fields configuration cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	if !reflect.DeepEqual(model.Settings.TensorFields, state.Settings.TensorFields) {
		resp.Diagnostics.AddError(
			"Cannot Modify Tensor Fields",
			"The tensor_fields configuration cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	// Check for changes in nested objects
	if !reflect.DeepEqual(model.Settings.TextPreprocessing, state.Settings.TextPreprocessing) {
		resp.Diagnostics.AddError(
			"Cannot Modify Text Preprocessing",
			"The text_preprocessing configuration cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	if !reflect.DeepEqual(model.Settings.ImagePreprocessing, state.Settings.ImagePreprocessing) {
		resp.Diagnostics.AddError(
			"Cannot Modify Image Preprocessing",
			"The image_preprocessing configuration cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	if !reflect.DeepEqual(model.Settings.VideoPreprocessing, state.Settings.VideoPreprocessing) {
		resp.Diagnostics.AddError(
			"Cannot Modify Video Preprocessing",
			"The video_preprocessing configuration cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	if !reflect.DeepEqual(model.Settings.AudioPreprocessing, state.Settings.AudioPreprocessing) {
		resp.Diagnostics.AddError(
			"Cannot Modify Audio Preprocessing",
			"The audio_preprocessing configuration cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	if !reflect.DeepEqual(model.Settings.AnnParameters, state.Settings.AnnParameters) {
		resp.Diagnostics.AddError(
			"Cannot Modify ANN Parameters",
			"The ann_parameters configuration cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	if !reflect.DeepEqual(model.Settings.ModelProperties, state.Settings.ModelProperties) {
		resp.Diagnostics.AddError(
			"Cannot Modify Model Properties",
			"The model_properties configuration cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	// Check for changes in scalar fields
	if !model.Settings.VectorNumericType.IsNull() && !state.Settings.VectorNumericType.IsNull() &&
		model.Settings.VectorNumericType.ValueString() != state.Settings.VectorNumericType.ValueString() {
		resp.Diagnostics.AddError(
			"Cannot Modify Vector Numeric Type",
			"The vector_numeric_type field cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	if !model.Settings.FilterStringMaxLength.IsNull() && !state.Settings.FilterStringMaxLength.IsNull() &&
		model.Settings.FilterStringMaxLength.ValueInt64() != state.Settings.FilterStringMaxLength.ValueInt64() {
		resp.Diagnostics.AddError(
			"Cannot Modify Filter String Max Length",
			"The filter_string_max_length field cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	// Check for changes in boolean fields
	if !model.Settings.TreatUrlsAndPointersAsImages.IsNull() && !state.Settings.TreatUrlsAndPointersAsImages.IsNull() &&
		model.Settings.TreatUrlsAndPointersAsImages.ValueBool() != state.Settings.TreatUrlsAndPointersAsImages.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot Modify treat_urls_and_pointers_as_images",
			"This field cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	if !model.Settings.TreatUrlsAndPointersAsMedia.IsNull() && !state.Settings.TreatUrlsAndPointersAsMedia.IsNull() &&
		model.Settings.TreatUrlsAndPointersAsMedia.ValueBool() != state.Settings.TreatUrlsAndPointersAsMedia.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot Modify treat_urls_and_pointers_as_media",
			"This field cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	if !model.Settings.NormalizeEmbeddings.IsNull() && !state.Settings.NormalizeEmbeddings.IsNull() &&
		model.Settings.NormalizeEmbeddings.ValueBool() != state.Settings.NormalizeEmbeddings.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot Modify normalize_embeddings",
			"This field cannot be modified. You must destroy and recreate the index to change this field.")
		return
	}

	indexName := model.IndexName.ValueString()

	// Check current index status before attempting update
	indices, err := r.marqoClient.ListIndices()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to List Indices",
			fmt.Sprintf("Could not check index status: %s", err.Error()),
		)
		return
	}

	var currentStatus string
	var currentIndex *go_marqo.IndexDetail
	indexFound := false
	for _, index := range indices {
		if index.IndexName == indexName {
			currentStatus = index.IndexStatus
			currentIndex = &index
			indexFound = true
			break
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Current index state before update - Name: %s, Status: %s, Shards: %d, Replicas: %d, Storage: %s, InferenceType: %s, NumberOfInferences: %d",
		currentIndex.IndexName,
		currentIndex.IndexStatus,
		currentIndex.NumberOfShards,
		currentIndex.NumberOfReplicas,
		currentIndex.StorageClass,
		currentIndex.InferenceType,
		currentIndex.NumberOfInferences))

	if !indexFound {
		resp.Diagnostics.AddError(
			"Index Not Found",
			fmt.Sprintf("Index %s does not exist", indexName))
		return
	}

	if currentStatus != "READY" {
		resp.Diagnostics.AddError(
			"Index Not Ready",
			fmt.Sprintf("Cannot update index %s: current status is %s, must be READY", indexName, currentStatus))
		return
	}

	// Add detailed logging for storage class
	tflog.Debug(ctx, fmt.Sprintf("Storage class from current index: '%s'", currentIndex.StorageClass))

	// Map storage class values similar to inference type
	storageClassMap := map[string]string{
		"BASIC":             "marqo.basic",
		"BALANCED":          "marqo.balanced",
		"PERFORMANCE":       "marqo.performance",
		"marqo.basic":       "marqo.basic",
		"marqo.balanced":    "marqo.balanced",
		"marqo.performance": "marqo.performance",
	}

	mappedStorageClass := currentIndex.StorageClass
	if mapped, exists := storageClassMap[currentIndex.StorageClass]; exists {
		mappedStorageClass = mapped
		tflog.Debug(ctx, fmt.Sprintf("Mapped storage class from '%s' to '%s'", currentIndex.StorageClass, mappedStorageClass))
	} else {
		tflog.Warn(ctx, fmt.Sprintf("No mapping found for storage class '%s'", currentIndex.StorageClass))
	}

	// Construct settings map with only the modifiable settings
	settings := map[string]interface{}{
		"inferenceType":      model.Settings.InferenceType.ValueString(),
		"numberOfInferences": model.Settings.NumberOfInferences.ValueInt64(),
		"numberOfShards":     model.Settings.NumberOfShards.ValueInt64(),
		"numberOfReplicas":   model.Settings.NumberOfReplicas.ValueInt64(),
		"storageClass":       mappedStorageClass,
		"type":               currentIndex.Type,
	}

	// Log raw values before any processing
	tflog.Debug(ctx, fmt.Sprintf("Raw values from current index - Storage: '%s', Type: '%s', Shards: %d, Replicas: %d",
		currentIndex.StorageClass,
		currentIndex.Type,
		currentIndex.NumberOfShards,
		currentIndex.NumberOfReplicas))

	// Log the desired changes
	tflog.Debug(ctx, fmt.Sprintf("Desired changes - InferenceType: %s, NumberOfInferences: %d, NumberOfShards: %d, NumberOfReplicas: %d",
		model.Settings.InferenceType.ValueString(),
		model.Settings.NumberOfInferences.ValueInt64(),
		model.Settings.NumberOfShards.ValueInt64(),
		model.Settings.NumberOfReplicas.ValueInt64()))

	// Validate that shards and replicas are only increasing
	if model.Settings.NumberOfShards.ValueInt64() < currentIndex.NumberOfShards {
		resp.Diagnostics.AddError(
			"Invalid Shard Count",
			fmt.Sprintf("Cannot decrease number of shards from %d to %d. Shards can only be increased.",
				currentIndex.NumberOfShards, model.Settings.NumberOfShards.ValueInt64()))
		return
	}

	if model.Settings.NumberOfReplicas.ValueInt64() < currentIndex.NumberOfReplicas {
		resp.Diagnostics.AddError(
			"Invalid Replica Count",
			fmt.Sprintf("Cannot decrease number of replicas from %d to %d. Replicas can only be increased.",
				currentIndex.NumberOfReplicas, model.Settings.NumberOfReplicas.ValueInt64()))
		return
	}

	// Check for drift between configuration and actual infrastructure
	// Instead of silently using API values, detect drift and provide clear error messages
	if currentIndex.NumberOfShards > model.Settings.NumberOfShards.ValueInt64() {
		resp.Diagnostics.AddError(
			"Drift Detected: Shard Count",
			fmt.Sprintf("The API reports %d shards, but your configuration specifies %d shards.\n\n"+
				"Shards can only be increased, not decreased. To resolve this:\n"+
				"1. Update your configuration to match the current state: number_of_shards = %d\n"+
				"2. Run terraform apply again\n\n"+
				"This ensures your Terraform configuration accurately reflects the actual infrastructure.",
				currentIndex.NumberOfShards, model.Settings.NumberOfShards.ValueInt64(), currentIndex.NumberOfShards))
		return
	}

	if currentIndex.NumberOfReplicas > model.Settings.NumberOfReplicas.ValueInt64() {
		resp.Diagnostics.AddError(
			"Drift Detected: Replica Count",
			fmt.Sprintf("The API reports %d replicas, but your configuration specifies %d replicas.\n\n"+
				"Replicas can only be increased, not decreased. To resolve this:\n"+
				"1. Update your configuration to match the current state: number_of_replicas = %d\n"+
				"2. Run terraform apply again\n\n"+
				"This ensures your Terraform configuration accurately reflects the actual infrastructure.",
				currentIndex.NumberOfReplicas, model.Settings.NumberOfReplicas.ValueInt64(), currentIndex.NumberOfReplicas))
		return
	}

	// Check if the number of inferences has been changed outside of Terraform
	if !model.Settings.NumberOfInferences.IsNull() &&
		currentIndex.NumberOfInferences != model.Settings.NumberOfInferences.ValueInt64() {
		tflog.Info(ctx, fmt.Sprintf("API reports different number of inferences (%d) than plan (%d), will update to plan value",
			currentIndex.NumberOfInferences, model.Settings.NumberOfInferences.ValueInt64()))
	}

	// Check if the inference type has been changed outside of Terraform
	if !model.Settings.InferenceType.IsNull() &&
		currentIndex.InferenceType != model.Settings.InferenceType.ValueString() &&
		!strings.Contains(model.Settings.InferenceType.ValueString(), currentIndex.InferenceType) &&
		!strings.Contains(currentIndex.InferenceType, model.Settings.InferenceType.ValueString()) {
		tflog.Info(ctx, fmt.Sprintf("API reports different inference type (%s) than plan (%s), will update to plan value",
			currentIndex.InferenceType, model.Settings.InferenceType.ValueString()))
	}

	if model.Settings.InferenceType.IsNull() {
		delete(settings, "inferenceType")
	}

	if model.Settings.NumberOfInferences.IsNull() {
		delete(settings, "numberOfInferences")
	}

	if model.Settings.NumberOfShards.IsNull() {
		delete(settings, "numberOfShards")
	}
	if model.Settings.NumberOfReplicas.IsNull() {
		delete(settings, "numberOfReplicas")
	}

	tflog.Debug(ctx, fmt.Sprintf("Final update settings being sent: %+v", settings))

	// Default timeout of 30 minutes for update
	timeoutDuration := 30 * time.Minute
	if model.Timeouts != nil && model.Timeouts.Update.ValueString() != "" {
		parsedTimeout, err := time.ParseDuration(model.Timeouts.Update.ValueString())
		if err == nil {
			timeoutDuration = parsedTimeout
			tflog.Info(ctx, fmt.Sprintf("Using configured update timeout of %v", timeoutDuration))
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Invalid update timeout duration: %s, using default of 30m", err))
		}
	}

	// Attempt to update the index
	err = r.marqoClient.UpdateIndex(indexName, settings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Index",
			fmt.Sprintf("Could not update index %s.\n"+
				"Error details: %s", indexName, err.Error()))
		return
	}

	// Wait for the index to be ready
	err = r.waitForIndexStatus(ctx, indexName, "READY", timeoutDuration, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Timeout Waiting for Index Update",
			fmt.Sprintf("Index %s update did not complete within the timeout period: %s", indexName, err))
		return
	}

	// Do final read to get the complete state
	// Preserve computed/meta fields from current state
	model.MarqoEndpoint = state.MarqoEndpoint
	model.Timeouts = state.Timeouts
	readResp := resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{State: resp.State}, &readResp)
	if readResp.Diagnostics.HasError() {
		resp.Diagnostics.Append(readResp.Diagnostics...)
		return
	}

	// Update the response state with the read state
	diags = resp.State.Set(ctx, &model)
	resp.Diagnostics.Append(diags...)
	resp.State = readResp.State
}

// ImportState imports an existing index into Terraform state.
func (r *indicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Extract the index name from the import ID
	indexName := req.ID

	tflog.Info(ctx, fmt.Sprintf("Importing index %s", indexName))

	// List all indices to find the one we're importing
	indices, err := r.marqoClient.ListIndices()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Listing Indices During Import",
			fmt.Sprintf("Could not list indices: %s", err.Error()),
		)
		return
	}

	// Find the index in the list
	var foundIndex bool
	for _, index := range indices {
		if index.IndexName == indexName {
			foundIndex = true
			break
		}
	}

	if !foundIndex {
		resp.Diagnostics.AddError(
			"Index Not Found",
			fmt.Sprintf("Index with name %s was not found", indexName),
		)
		return
	}

	// Create a minimal state with the index name and empty settings
	// The Read method will be called after this to populate the full state
	initialState := IndexResourceModel{
		IndexName: types.StringValue(indexName),
		Settings: IndexSettingsModel{
			Type:               types.StringValue(""),
			InferenceType:      types.StringValue(""),
			NumberOfInferences: types.Int64Value(0),
			StorageClass:       types.StringValue(""),
			NumberOfShards:     types.Int64Value(0),
			NumberOfReplicas:   types.Int64Value(0),
			Model:              types.StringValue(""),
			// Set boolean fields to null by default
			NormalizeEmbeddings:          types.BoolNull(),
			TreatUrlsAndPointersAsImages: types.BoolNull(),
			TreatUrlsAndPointersAsMedia:  types.BoolNull(),
			// Initialize preprocessing fields to null
			ImagePreprocessing: nil,
			VideoPreprocessing: &VideoPreprocessingModelCreate{
				SplitLength:  types.Int64Null(),
				SplitOverlap: types.Int64Null(),
			},
			AudioPreprocessing: &AudioPreprocessingModelCreate{
				SplitLength:  types.Int64Null(),
				SplitOverlap: types.Int64Null(),
			},
			// Other fields will be populated by the Read method
		},
	}

	// Set the entire state at once
	diags := resp.State.Set(ctx, &initialState)
	resp.Diagnostics.Append(diags...)
}
