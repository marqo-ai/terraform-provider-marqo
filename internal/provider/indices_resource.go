package provider

import (
	"context"
	"fmt"
	"marqo/go_marqo"
	"reflect"
	"strconv"
	"strings"

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
	//ID        types.String       `tfsdk:"id"`
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
	TreatUrlsAndPointersAsMedia  types.Bool                   `tfsdk:"treat_urls_and_pointers_as_media"`
	Model                        types.String                 `tfsdk:"model"`
	NormalizeEmbeddings          types.Bool                   `tfsdk:"normalize_embeddings"`
	TextPreprocessing            TextPreprocessingModelCreate `tfsdk:"text_preprocessing"`
	ImagePreprocessing           ImagePreprocessingModel      `tfsdk:"image_preprocessing"`
	VideoPreprocessing           *VideoPreprocessingModel     `tfsdk:"video_preprocessing"`
	AudioPreprocessing           *AudioPreprocessingModel     `tfsdk:"audio_preprocessing"`
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

type VideoPreprocessingModel struct {
	SplitLength  types.Int64 `tfsdk:"split_length"`
	SplitOverlap types.Int64 `tfsdk:"split_overlap"`
}

type AudioPreprocessingModel struct {
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
			//"id": schema.StringAttribute{
			//	Computed:    true,
			//	Description: "The unique identifier for the index.",
			//},
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

// constructTensorFields constructs the tensorFields setting from the input.
//func constructTensorFields(tensorFieldsInput []string) ([]string, error) {
//	// Blank for now
//	return tensorFieldsInput, nil
//}

func (r *indicesResource) findAndCreateState(indices []go_marqo.IndexDetail, indexName string) (*IndexResourceModel, bool) {
	for _, indexDetail := range indices {
		if indexDetail.IndexName == indexName {
			return &IndexResourceModel{
				//ID:        types.StringValue(indexDetail.IndexName),
				IndexName: types.StringValue(indexDetail.IndexName),
				Settings: IndexSettingsModel{
					Type:                         types.StringValue(indexDetail.Type),
					VectorNumericType:            types.StringValue(indexDetail.VectorNumericType),
					TreatUrlsAndPointersAsImages: types.BoolValue(indexDetail.TreatUrlsAndPointersAsImages),
					TreatUrlsAndPointersAsMedia:  types.BoolValue(indexDetail.TreatUrlsAndPointersAsMedia),
					Model:                        types.StringValue(indexDetail.Model),
					AllFields:                    ConvertMarqoAllFieldInputs(indexDetail.AllFields),
					TensorFields:                 indexDetail.TensorFields,
					NormalizeEmbeddings:          types.BoolValue(indexDetail.NormalizeEmbeddings),
					InferenceType:                types.StringValue(indexDetail.InferenceType),
					NumberOfInferences:           types.Int64Value(indexDetail.NumberOfInferences),
					StorageClass:                 types.StringValue(indexDetail.StorageClass),
					NumberOfShards:               types.Int64Value(indexDetail.NumberOfShards),
					NumberOfReplicas:             types.Int64Value(indexDetail.NumberOfReplicas),
					TextPreprocessing: TextPreprocessingModelCreate{
						SplitLength:  types.Int64Value(indexDetail.TextPreprocessing.SplitLength),
						SplitMethod:  types.StringValue(indexDetail.TextPreprocessing.SplitMethod),
						SplitOverlap: types.Int64Value(indexDetail.TextPreprocessing.SplitOverlap),
					},
					ImagePreprocessing: ImagePreprocessingModel{
						PatchMethod: types.StringValue(indexDetail.ImagePreprocessing.PatchMethod),
					},
					VideoPreprocessing: &VideoPreprocessingModel{
						SplitLength:  types.Int64Value(indexDetail.VideoPreprocessing.SplitLength),
						SplitOverlap: types.Int64Value(indexDetail.VideoPreprocessing.SplitOverlap),
					},
					AudioPreprocessing: &AudioPreprocessingModel{
						SplitLength:  types.Int64Value(indexDetail.AudioPreprocessing.SplitLength),
						SplitOverlap: types.Int64Value(indexDetail.AudioPreprocessing.SplitOverlap),
					},
					AnnParameters: AnnParametersModelCreate{
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

	newState, found := r.findAndCreateState(indices, state.IndexName.ValueString())

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

		// Ensure features and dependent_fields are always set
		for i := range newState.Settings.AllFields {
			if len(newState.Settings.AllFields[i].Features) == 0 {
				newState.Settings.AllFields[i].Features = nil
			}
			if len(newState.Settings.AllFields[i].DependentFields) == 0 {
				newState.Settings.AllFields[i].DependentFields = nil
			}
		}

		// Ignore these fields for structured indexes
		if newState.Settings.Type.ValueString() == "structured" {
			newState.Settings.FilterStringMaxLength = types.Int64Null()
			newState.Settings.TreatUrlsAndPointersAsImages = types.BoolNull()
			newState.Settings.TreatUrlsAndPointersAsMedia = types.BoolNull()
		}

		// Handle image_preprocessing.patch_method
		if newState.Settings.ImagePreprocessing.PatchMethod.ValueString() == "" {
			newState.Settings.ImagePreprocessing.PatchMethod = types.StringNull()
		}

		// Remove null fields
		if newState.Settings.InferenceType.IsNull() {
			newState.Settings.InferenceType = types.StringNull()
		}
	}

	// if index no longer exists in cloud, delete the state
	if !found {

		resp.Diagnostics.AddWarning("Resource Not Found", "test The specified index does not exist in the cloud. The state will be deleted.")
		//state = IndexResourceModel{}
		// Then Totally Remove from terraform resources
		resp.State.RemoveResource(ctx)
		//resp.State.Set(ctx, &IndexResourceModel{})
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
	// Implement a deep comparison between existing and desired states
	// This is a basic implementation - you may need to adjust based on your specific needs
	return reflect.DeepEqual(existing.Settings, desired.Settings)
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
	if model.Settings.TextPreprocessing.SplitLength.IsNull() &&
		model.Settings.TextPreprocessing.SplitMethod.IsNull() &&
		model.Settings.TextPreprocessing.SplitOverlap.IsNull() {
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
		// Set storageClass to marqo.basic
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
		if strings.Contains(err.Error(), "already exists") {
			tflog.Info(ctx, fmt.Sprintf("Index %s already exists. Checking if it needs to be updated.", model.IndexName.ValueString()))

			indices, err := r.marqoClient.ListIndices()
			if err != nil {
				resp.Diagnostics.AddError("Failed to List Indices", fmt.Sprintf("Could not list indices: %s", err.Error()))
				return
			}

			existingState, found := r.findAndCreateState(indices, model.IndexName.ValueString())
			if !found {
				resp.Diagnostics.AddError("Failed to Find Index", fmt.Sprintf("Index %s not found after creation", model.IndexName.ValueString()))
				return
			}

			// Compare existing state with desired state
			if !statesAreEqual(existingState, &model) {
				// Attempt to update the existing index
				err = r.marqoClient.UpdateIndex(model.IndexName.ValueString(), settings)
				if err != nil {
					resp.Diagnostics.AddError("Failed to Update Existing Index",
						fmt.Sprintf("Index %s exists but couldn't be updated to match the configuration: %s", model.IndexName.ValueString(), err.Error()))
					return
				}
				tflog.Info(ctx, fmt.Sprintf("Index %s updated to match configuration.", model.IndexName.ValueString()))
			} else {
				tflog.Info(ctx, fmt.Sprintf("Existing index %s matches configuration. No update needed.", model.IndexName.ValueString()))
			}

			// Set state to the (potentially updated) existing index
			diags = resp.State.Set(ctx, existingState)
			resp.Diagnostics.Append(diags...)
			resp.Diagnostics.AddWarning(fmt.Sprintf("Index %s already existed and has been imported into Terraform state.", model.IndexName.ValueString()),
				"Any differences between the existing index and your configuration have been resolved by updating the index.")
			return
		}

		resp.Diagnostics.AddError("Failed to Create Index", "Could not create index: "+err.Error())
		return
	}

	// Set the index name as the ID in the Terraform state
	//model.ID = model.IndexName
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
	//model.ID = model.IndexName
	diags = resp.State.Set(ctx, &model)
	resp.Diagnostics.Append(diags...)
}
