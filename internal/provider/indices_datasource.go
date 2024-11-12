package provider

import (
	"context"
	"fmt"
	"marqo/go_marqo"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &indicesDataSource{}
	_ datasource.DataSourceWithConfigure = &indicesDataSource{}
)

// ManageIndicesResource is a helper function to simplify the provider implementation.
func ReadIndicesDataSource() datasource.DataSource {
	return &indicesDataSource{}
}

// orderResourceModel maps the resource schema data.
type allIndicesResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Items       []indexModel `tfsdk:"items"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// indexModel maps index detail data.
type indexModel struct {
	Created                      types.String             `tfsdk:"created"`
	IndexName                    types.String             `tfsdk:"index_name"`
	NumberOfShards               types.String             `tfsdk:"number_of_shards"`
	NumberOfReplicas             types.String             `tfsdk:"number_of_replicas"`
	IndexStatus                  types.String             `tfsdk:"index_status"`
	AllFields                    []AllFieldInput          `tfsdk:"all_fields"`
	TensorFields                 []string                 `tfsdk:"tensor_fields"`
	NumberOfInferences           types.String             `tfsdk:"number_of_inferences"`
	StorageClass                 types.String             `tfsdk:"storage_class"`
	InferenceType                types.String             `tfsdk:"inference_type"`
	DocsCount                    types.String             `tfsdk:"docs_count"`
	StoreSize                    types.String             `tfsdk:"store_size"`
	DocsDeleted                  types.String             `tfsdk:"docs_deleted"`
	SearchQueryTotal             types.String             `tfsdk:"search_query_total"`
	TreatUrlsAndPointersAsImages types.Bool               `tfsdk:"treat_urls_and_pointers_as_images"`
	TreatUrlsAndPointersAsMedia  types.Bool               `tfsdk:"treat_urls_and_pointers_as_media"`
	MarqoEndpoint                types.String             `tfsdk:"marqo_endpoint"`
	Type                         types.String             `tfsdk:"type"`
	VectorNumericType            types.String             `tfsdk:"vector_numeric_type"`
	Model                        types.String             `tfsdk:"model"`
	ModelProperties              *ModelPropertiesModel    `tfsdk:"model_properties"`
	NormalizeEmbeddings          types.Bool               `tfsdk:"normalize_embeddings"`
	TextPreprocessing            *TextPreprocessingModel  `tfsdk:"text_preprocessing"`
	ImagePreprocessing           *ImagePreprocessingModel `tfsdk:"image_preprocessing"`
	VideoPreprocessing           *VideoPreprocessingModel `tfsdk:"video_preprocessing"`
	AudioPreprocessing           *AudioPreprocessingModel `tfsdk:"audio_preprocessing"`
	AnnParameters                *AnnParametersModel      `tfsdk:"ann_parameters"`
	MarqoVersion                 types.String             `tfsdk:"marqo_version"`
	FilterStringMaxLength        types.String             `tfsdk:"filter_string_max_length"`
}

type ModelPropertiesModel struct {
	Name             types.String        `tfsdk:"name"`
	Dimensions       types.Int64         `tfsdk:"dimensions"`
	Type             types.String        `tfsdk:"type"`
	Tokens           types.Int64         `tfsdk:"tokens"`
	ModelLocation    *ModelLocationModel `tfsdk:"model_location"`
	Url              types.String        `tfsdk:"url"`
	TrustRemoteCode  types.Bool          `tfsdk:"trust_remote_code"`
	IsMarqtunedModel types.Bool          `tfsdk:"is_marqtuned_model"`
}

type ModelLocationModel struct {
	S3           *S3LocationModel `tfsdk:"s3"`
	Hf           *HfLocationModel `tfsdk:"hf"`
	AuthRequired types.Bool       `tfsdk:"auth_required"`
}

type S3LocationModel struct {
	Bucket types.String `tfsdk:"bucket"`
	Key    types.String `tfsdk:"key"`
}

type HfLocationModel struct {
	RepoId   types.String `tfsdk:"repo_id"`
	Filename types.String `tfsdk:"filename"`
}

type TextPreprocessingModel struct {
	SplitLength  types.String `tfsdk:"split_length"`
	SplitMethod  types.String `tfsdk:"split_method"`
	SplitOverlap types.String `tfsdk:"split_overlap"`
}

type VideoPreprocessingModel struct {
	SplitLength  types.String `tfsdk:"split_length"`
	SplitOverlap types.String `tfsdk:"split_overlap"`
}

type AudioPreprocessingModel struct {
	SplitLength  types.String `tfsdk:"split_length"`
	SplitOverlap types.String `tfsdk:"split_overlap"`
}

type AnnParametersModel struct {
	SpaceType  types.String    `tfsdk:"space_type"`
	Parameters parametersModel `tfsdk:"parameters"`
}

type parametersModel struct {
	EfConstruction types.String `tfsdk:"ef_construction"`
	M              types.String `tfsdk:"m"`
}

// Configure adds the provider configured client to the resource.
func (d *indicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	fmt.Println("Configure called")

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

	d.marqoClient = client
}

// orderResource is the resource implementation.
type indicesDataSource struct {
	marqoClient *go_marqo.Client
}

// Metadata returns the resource type name.
func (d *indicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_read_indices"
}

// Schema defines the schema for the resource.
func (d *indicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the resource.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "The last time the resource was updated.",
			},
			"items": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"index_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the index",
						},
						"marqo_version": schema.StringAttribute{
							Computed:    true,
							Description: "The version of Marqo used by the index",
						},
						"marqo_endpoint": schema.StringAttribute{
							Computed:    true,
							Description: "The Marqo endpoint used by the index",
						},
						"filter_string_max_length": schema.StringAttribute{
							Computed:    true,
							Description: "The filter string max length",
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
						"created": schema.StringAttribute{
							Computed:    true,
							Description: "The creation date of the index",
						},
						"number_of_shards": schema.StringAttribute{
							Computed:    true,
							Description: "The number of shards for the index",
						},
						"number_of_replicas": schema.StringAttribute{
							Computed:    true,
							Description: "The number of replicas for the index",
						},
						"index_status": schema.StringAttribute{
							Computed:    true,
							Description: "The status of the index",
						},
						"number_of_inferences": schema.StringAttribute{
							Computed:    true,
							Description: "The number of inferences made by the index",
						},
						"storage_class": schema.StringAttribute{
							Computed:    true,
							Description: "The storage class of the index",
						},
						"inference_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of inference used by the index",
						},
						"docs_count": schema.StringAttribute{
							Computed:    true,
							Description: "The number of documents in the index",
						},
						"store_size": schema.StringAttribute{
							Computed:    true,
							Description: "The size of the index storage",
						},
						"docs_deleted": schema.StringAttribute{
							Computed:    true,
							Description: "The number of documents deleted from the index",
						},
						"search_query_total": schema.StringAttribute{
							Computed:    true,
							Description: "The total number of search queries made on the index",
						},
						"treat_urls_and_pointers_as_images": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates if URLs and pointers should be treated as images",
						},
						"treat_urls_and_pointers_as_media": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates if URLs and pointers should be treated as media",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the index",
						},
						"vector_numeric_type": schema.StringAttribute{
							Computed:    true,
							Description: "The numeric type of the vector",
						},
						"model": schema.StringAttribute{
							Computed:    true,
							Description: "The model used by the index",
						},
						"model_properties": schema.SingleNestedAttribute{
							Computed: true,
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"name":       schema.StringAttribute{Computed: true},
								"dimensions": schema.StringAttribute{Computed: true},
								"type":       schema.StringAttribute{Computed: true},
								"tokens":     schema.StringAttribute{Computed: true},
								"model_location": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"s3": schema.SingleNestedAttribute{
											Computed: true,
											Optional: true,
											Attributes: map[string]schema.Attribute{
												"bucket": schema.StringAttribute{Computed: true},
												"key":    schema.StringAttribute{Computed: true},
											},
										},
										"hf": schema.SingleNestedAttribute{
											Computed: true,
											Optional: true,
											Attributes: map[string]schema.Attribute{
												"repo_id":  schema.StringAttribute{Computed: true},
												"filename": schema.StringAttribute{Computed: true},
											},
										},
										"auth_required": schema.BoolAttribute{Computed: true},
									},
								},
								"url":                schema.StringAttribute{Computed: true},
								"trust_remote_code":  schema.StringAttribute{Computed: true},
								"is_marqtuned_model": schema.BoolAttribute{Computed: true},
							},
						},
						"normalize_embeddings": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates if embeddings should be normalized",
						},
						"text_preprocessing": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"split_length": schema.StringAttribute{
									Computed:    true,
									Description: "The split length for text preprocessing",
								},
								"split_method": schema.StringAttribute{
									Computed:    true,
									Description: "The split method for text preprocessing",
								},
								"split_overlap": schema.StringAttribute{
									Computed:    true,
									Description: "The split overlap for text preprocessing",
								},
							},
						},
						"image_preprocessing": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"patch_method": schema.StringAttribute{
									Computed:    true,
									Description: "The patch method for image preprocessing",
								},
							},
						},
						"video_preprocessing": schema.SingleNestedAttribute{
							Computed: true,
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"split_length":  schema.StringAttribute{Computed: true, Optional: true},
								"split_overlap": schema.StringAttribute{Computed: true, Optional: true},
							},
						},
						"audio_preprocessing": schema.SingleNestedAttribute{
							Computed: true,
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"split_length":  schema.StringAttribute{Computed: true, Optional: true},
								"split_overlap": schema.StringAttribute{Computed: true, Optional: true},
							},
						},
						"ann_parameters": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"space_type": schema.StringAttribute{
									Computed:    true,
									Description: "The space type for ANN parameters",
								},
								"parameters": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"ef_construction": schema.StringAttribute{
											Computed:    true,
											Description: "The efConstruction parameter for ANN",
										},
										"m": schema.StringAttribute{
											Computed:    true,
											Description: "The m parameter for ANN",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// ConvertMarqoAllFieldInputs converts a slice of go_marqo.AllFieldInput to a slice of AllFieldInput.
func ConvertMarqoAllFieldInputs(marqoFields []go_marqo.AllFieldInput) []AllFieldInput {
	allFieldsConverted := make([]AllFieldInput, len(marqoFields))
	for i, field := range marqoFields {
		featuresConverted := make([]types.String, 0)
		for _, feature := range field.Features {
			featuresConverted = append(featuresConverted, types.StringValue(feature))
		}
		dependentFieldsConverted := make(map[string]types.Float64)
		for key, value := range field.DependentFields {
			dependentFieldsConverted[key] = types.Float64Value(value)
		}
		allFieldsConverted[i] = AllFieldInput{
			Name:            types.StringValue(field.Name),
			Type:            types.StringValue(field.Type),
			Features:        featuresConverted,
			DependentFields: dependentFieldsConverted,
		}
	}
	return allFieldsConverted
}

func convertModelLocation(location go_marqo.ModelLocation) *ModelLocationModel {
	modelLocation := &ModelLocationModel{
		AuthRequired: types.BoolValue(location.AuthRequired),
	}

	if location.S3 != nil {
		modelLocation.S3 = &S3LocationModel{
			Bucket: types.StringValue(location.S3.Bucket),
			Key:    types.StringValue(location.S3.Key),
		}
	}

	if location.Hf != nil {
		modelLocation.Hf = &HfLocationModel{
			RepoId:   types.StringValue(location.Hf.RepoId),
			Filename: types.StringValue(location.Hf.Filename),
		}
	}

	return modelLocation
}

// Read refreshes the Terraform state with the latest data.
func (d *indicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(context.TODO(), "Calling marqo client ListIndices")
	var model allIndicesResourceModel

	// Retrieve the id from the Terraform configuration
	diags := req.Config.GetAttribute(ctx, path.Root("id"), &model.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	indices, err := d.marqoClient.ListIndices()
	if err != nil {
		resp.Diagnostics.AddError("Failed to List Indices", fmt.Sprintf("Could not list indices: %s", err.Error()))
		return
	}

	// Use the id from the configuration
	model.ID = types.StringValue(model.ID.ValueString())

	inferenceTypeMap := map[string]string{
		"CPU.SMALL": "marqo.CPU.small",
		"CPU.LARGE": "marqo.CPU.large",
		"GPU":       "marqo.GPU",
	}

	storaceClassMap := map[string]string{
		"BASIC":       "marqo.basic",
		"BALANCED":    "marqo.balanced",
		"PERFORMANCE": "marqo.performance",
	}

	items := make([]indexModel, len(indices))
	for i, indexDetail := range indices {
		inferenceType := indexDetail.InferenceType
		if mappedValue, exists := inferenceTypeMap[inferenceType]; exists {
			inferenceType = mappedValue
		}

		storageClass := indexDetail.StorageClass
		if mappedValue, exists := storaceClassMap[storageClass]; exists {
			storageClass = mappedValue
		}

		// Handle image_preprocessing.patch_method
		items[i] = indexModel{
			Created:                      types.StringValue(indexDetail.Created),
			IndexName:                    types.StringValue(indexDetail.IndexName),
			NumberOfShards:               types.StringValue(fmt.Sprintf("%d", indexDetail.NumberOfShards)),
			NumberOfReplicas:             types.StringValue(fmt.Sprintf("%d", indexDetail.NumberOfReplicas)),
			IndexStatus:                  types.StringValue(indexDetail.IndexStatus),
			AllFields:                    ConvertMarqoAllFieldInputs(indexDetail.AllFields),
			TensorFields:                 indexDetail.TensorFields,
			NumberOfInferences:           types.StringValue(fmt.Sprintf("%d", indexDetail.NumberOfInferences)),
			StorageClass:                 types.StringValue(storageClass),
			InferenceType:                types.StringValue(inferenceType),
			DocsCount:                    types.StringValue(indexDetail.DocsCount),
			StoreSize:                    types.StringValue(indexDetail.StoreSize),
			DocsDeleted:                  types.StringValue(indexDetail.DocsDeleted),
			SearchQueryTotal:             types.StringValue(indexDetail.SearchQueryTotal),
			TreatUrlsAndPointersAsImages: types.BoolValue(indexDetail.TreatUrlsAndPointersAsImages),
			TreatUrlsAndPointersAsMedia:  types.BoolValue(indexDetail.TreatUrlsAndPointersAsMedia),
			MarqoEndpoint:                types.StringValue(indexDetail.MarqoEndpoint),
			Type:                         types.StringValue(indexDetail.Type),
			VectorNumericType:            types.StringValue(indexDetail.VectorNumericType),
			Model:                        types.StringValue(indexDetail.Model),
			ModelProperties: &ModelPropertiesModel{
				Name:             types.StringValue(indexDetail.ModelProperties.Name),
				Dimensions:       types.Int64Value(indexDetail.ModelProperties.Dimensions),
				Type:             types.StringValue(indexDetail.ModelProperties.Type),
				Tokens:           types.Int64Value(indexDetail.ModelProperties.Tokens),
				ModelLocation:    convertModelLocation(indexDetail.ModelProperties.ModelLocation),
				Url:              types.StringValue(indexDetail.ModelProperties.Url),
				TrustRemoteCode:  types.BoolValue(indexDetail.ModelProperties.TrustRemoteCode),
				IsMarqtunedModel: types.BoolValue(indexDetail.ModelProperties.IsMarqtunedModel),
			},
			NormalizeEmbeddings: types.BoolValue(indexDetail.NormalizeEmbeddings),
			TextPreprocessing: &TextPreprocessingModel{
				SplitLength:  types.StringValue(fmt.Sprintf("%d", indexDetail.TextPreprocessing.SplitLength)),
				SplitMethod:  types.StringValue(indexDetail.TextPreprocessing.SplitMethod),
				SplitOverlap: types.StringValue(fmt.Sprintf("%d", indexDetail.TextPreprocessing.SplitOverlap)),
			},
			// ImagePreprocessing
			AnnParameters: &AnnParametersModel{
				SpaceType: types.StringValue(indexDetail.AnnParameters.SpaceType),
				Parameters: parametersModel{
					EfConstruction: types.StringValue(fmt.Sprintf("%d", indexDetail.AnnParameters.Parameters.EfConstruction)),
					M:              types.StringValue(fmt.Sprintf("%d", indexDetail.AnnParameters.Parameters.M)),
				},
			},
			MarqoVersion:          types.StringValue(indexDetail.MarqoVersion),
			FilterStringMaxLength: types.StringValue(fmt.Sprintf("%d", indexDetail.FilterStringMaxLength)),
		}

		// Handle model properties
		// In the resource
		if items[i].ModelProperties.IsEmpty() {
			items[i].ModelProperties = nil
		}

		// Remove null fields
		if items[i].InferenceType.IsNull() {
			items[i].InferenceType = types.StringNull()
		}
	}

	// Set the last_updated field
	model.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	model.Items = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
