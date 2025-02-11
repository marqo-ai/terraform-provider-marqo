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
	_ resource.Resource              = &indicesResource{}
	_ resource.ResourceWithConfigure = &indicesResource{}
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
			return &IndexResourceModel{
				//ID:        types.StringValue(indexDetail.IndexName),
				IndexName:     types.StringValue(indexDetail.IndexName),
				MarqoEndpoint: types.StringValue(indexDetail.MarqoEndpoint),
				Timeouts:      existingTimeouts,
				Settings: IndexSettingsModel{
					Type:                         types.StringValue(indexDetail.Type),
					VectorNumericType:            types.StringValue(indexDetail.VectorNumericType),
					TreatUrlsAndPointersAsImages: types.BoolValue(indexDetail.TreatUrlsAndPointersAsImages),
					TreatUrlsAndPointersAsMedia:  types.BoolValue(indexDetail.TreatUrlsAndPointersAsMedia),
					Model:                        types.StringValue(indexDetail.Model),
					ModelProperties:              convertModelPropertiesToResource(&indexDetail.ModelProperties),
					AllFields:                    ConvertMarqoAllFieldInputs(indexDetail.AllFields),
					TensorFields:                 indexDetail.TensorFields,
					NormalizeEmbeddings:          types.BoolValue(indexDetail.NormalizeEmbeddings),
					InferenceType:                types.StringValue(indexDetail.InferenceType),
					NumberOfInferences:           types.Int64Value(indexDetail.NumberOfInferences),
					StorageClass:                 types.StringValue(indexDetail.StorageClass),
					NumberOfShards:               types.Int64Value(indexDetail.NumberOfShards),
					NumberOfReplicas:             types.Int64Value(indexDetail.NumberOfReplicas),
					TextPreprocessing: &TextPreprocessingModelCreate{
						SplitLength:  types.Int64Value(indexDetail.TextPreprocessing.SplitLength),
						SplitMethod:  types.StringValue(indexDetail.TextPreprocessing.SplitMethod),
						SplitOverlap: types.Int64Value(indexDetail.TextPreprocessing.SplitOverlap),
					},
					ImagePreprocessing: &ImagePreprocessingModel{
						PatchMethod: types.StringValue(indexDetail.ImagePreprocessing.PatchMethod),
					},
					VideoPreprocessing: &VideoPreprocessingModelCreate{
						SplitLength:  types.Int64Value(indexDetail.VideoPreprocessing.SplitLength),
						SplitOverlap: types.Int64Value(indexDetail.VideoPreprocessing.SplitOverlap),
					},
					AudioPreprocessing: &AudioPreprocessingModelCreate{
						SplitLength:  types.Int64Value(indexDetail.AudioPreprocessing.SplitLength),
						SplitOverlap: types.Int64Value(indexDetail.AudioPreprocessing.SplitOverlap),
					},
					AnnParameters: &AnnParametersModelCreate{
						SpaceType: types.StringValue(indexDetail.AnnParameters.SpaceType),
						Parameters: ParametersModel{
							EfConstruction: types.Int64Value(indexDetail.AnnParameters.Parameters.EfConstruction),
							M:              types.Int64Value(indexDetail.AnnParameters.Parameters.M),
						},
					},
					FilterStringMaxLength: types.Int64Value(indexDetail.FilterStringMaxLength),
				},
			}, true
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

		// Ignore these fields for structured indexes
		if newState.Settings.Type.ValueString() == "structured" {
			newState.Settings.FilterStringMaxLength = types.Int64Null()
			newState.Settings.TreatUrlsAndPointersAsImages = types.BoolNull()
			newState.Settings.TreatUrlsAndPointersAsMedia = types.BoolNull()
		}

		// Handle image_preprocessing
		if state.Settings.ImagePreprocessing == nil {
			newState.Settings.ImagePreprocessing = nil
		} else if newState.Settings.ImagePreprocessing != nil &&
			newState.Settings.ImagePreprocessing.PatchMethod.ValueString() == "" {
			newState.Settings.ImagePreprocessing.PatchMethod = types.StringNull()
		}

		// preserve the video/audio preprocessing from current state since api does not return them
		if state.Settings.VideoPreprocessing != nil {
			newState.Settings.VideoPreprocessing = state.Settings.VideoPreprocessing
		}
		if state.Settings.AudioPreprocessing != nil {
			newState.Settings.AudioPreprocessing = state.Settings.AudioPreprocessing
		}

		// Then handle zero values
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

		// Handle model properties
		if newState.Settings.ModelProperties.IsEmpty() {
			newState.Settings.ModelProperties = nil
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

// waitForIndexStatus waits for an index to reach a target status or be deleted
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
		"type":                         model.Settings.Type.ValueString(),
		"vectorNumericType":            model.Settings.VectorNumericType.ValueString(),
		"treatUrlsAndPointersAsImages": model.Settings.TreatUrlsAndPointersAsImages.ValueBool(),
		"treatUrlsAndPointersAsMedia":  model.Settings.TreatUrlsAndPointersAsMedia.ValueBool(),
		"model":                        model.Settings.Model.ValueString(),
		"modelProperties":              model.Settings.ModelProperties,
		"normalizeEmbeddings":          model.Settings.NormalizeEmbeddings.ValueBool(),
		"allFields":                    convertAllFieldsToMap(model.Settings.AllFields),
		"tensorFields":                 model.Settings.TensorFields,
		"inferenceType":                model.Settings.InferenceType.ValueString(),
		"numberOfInferences":           model.Settings.NumberOfInferences.ValueInt64(),
		"storageClass":                 model.Settings.StorageClass.ValueString(),
		"numberOfShards":               model.Settings.NumberOfShards.ValueInt64(),
		"numberOfReplicas":             model.Settings.NumberOfReplicas.ValueInt64(),
		"filterStringMaxLength":        model.Settings.FilterStringMaxLength.ValueInt64(),
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
		settings["imagePreprocessing"] = map[string]interface{}{
			"patchMethod": model.Settings.ImagePreprocessing.PatchMethod.ValueString(),
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
	if model.Settings.VectorNumericType.IsNull() {
		delete(settings, "vectorNumericType")
	}
	if model.Settings.TreatUrlsAndPointersAsImages.IsNull() {
		delete(settings, "treatUrlsAndPointersAsImages")
	}
	if model.Settings.TreatUrlsAndPointersAsMedia.IsNull() {
		delete(settings, "treatUrlsAndPointersAsMedia")
	}
	if model.Settings.Model.IsNull() {
		delete(settings, "model")
	}
	if model.Settings.NormalizeEmbeddings.IsNull() {
		delete(settings, "normalizeEmbeddings")
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
	if model.Settings.InferenceType.IsNull() {
		delete(settings, "inferenceType")
	}
	if model.Settings.NumberOfInferences.IsNull() {
		delete(settings, "numberOfInferences")
	}
	if model.Settings.StorageClass.IsNull() {
		settings["storageClass"] = "marqo.basic"
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
	if model.Settings.FilterStringMaxLength.IsNull() {
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

			// Set state to the (potentially updated) existing index
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
			fmt.Sprintf("Index %s deletion did not complete: %s", indexName, err))
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

	indexName := model.IndexName.ValueString()

	// Check current index status before attempting update
	indices, err := r.marqoClient.ListIndices()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to List Indices",
			fmt.Sprintf("Could not check index status: %s", err.Error()))
		return
	}

	var currentStatus string
	indexFound := false
	for _, index := range indices {
		if index.IndexName == indexName {
			currentStatus = index.IndexStatus
			indexFound = true
			break
		}
	}

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
			fmt.Sprintf("Index %s update did not complete: %s", indexName, err))
		return
	}

	// Do final read to get the complete state
	// Preserve computed/meta fields from current state
	var state IndexResourceModel
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
