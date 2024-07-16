package provider

import (
	"context"
	"fmt"

	"terraform-provider-marqo/marqo"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	Created                      types.String           `tfsdk:"created"`
	IndexName                    types.String           `tfsdk:"index_name"`
	NumberOfShards               types.Int64            `tfsdk:"number_of_shards"`
	NumberOfReplicas             types.Int64            `tfsdk:"number_of_replicas"`
	IndexStatus                  types.String           `tfsdk:"index_status"`
	AllFields                    []AllFieldInput        `tfsdk:"all_fields"`
	TensorFields                 []string               `tfsdk:"tensor_fields"`
	NumberOfInferences           types.Int64            `tfsdk:"number_of_inferences"`
	StorageClass                 types.String           `tfsdk:"storage_class"`
	InferenceType                types.String           `tfsdk:"inference_type"`
	DocsCount                    types.String           `tfsdk:"docs_count"`
	StoreSize                    types.String           `tfsdk:"store_size"`
	DocsDeleted                  types.String           `tfsdk:"docs_deleted"`
	SearchQueryTotal             types.String           `tfsdk:"search_query_total"`
	TreatUrlsAndPointersAsImages types.Bool             `tfsdk:"treat_urls_and_pointers_as_images"`
	MarqoEndpoint                types.String           `tfsdk:"marqo_endpoint"`
	Type                         types.String           `tfsdk:"type"`
	VectorNumericType            types.String           `tfsdk:"vector_numeric_type"`
	Model                        types.String           `tfsdk:"model"`
	NormalizeEmbeddings          types.Bool             `tfsdk:"normalize_embeddings"`
	TextPreprocessing            TextPreprocessingModel `tfsdk:"text_preprocessing"` // Assuming no specific structure
	//ImagePreprocessing           types.Object           `tfsdk:"image_preprocessing"` // Assuming no specific structure
	AnnParameters         AnnParametersModel `tfsdk:"ann_parameters"` // Assuming no specific structure
	MarqoVersion          types.String       `tfsdk:"marqo_version"`
	FilterStringMaxLength types.Int64        `tfsdk:"filter_string_max_length"`
}

type TextPreprocessingModel struct {
	SplitLength  types.Int64  `tfsdk:"split_length"`
	SplitMethod  types.String `tfsdk:"split_method"`
	SplitOverlap types.Int64  `tfsdk:"split_overlap"`
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

	client, ok := req.ProviderData.(*marqo.Client)

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
	marqoClient *marqo.Client
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
				Computed:    true,
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

// ConvertMarqoAllFieldInputs converts a slice of marqo.AllFieldInput to a slice of AllFieldInput.
func ConvertMarqoAllFieldInputs(marqoFields []marqo.AllFieldInput) []AllFieldInput {
	allFieldsConverted := make([]AllFieldInput, len(marqoFields))
	for i, field := range marqoFields {
		featuresConverted := make([]types.String, len(field.Features))
		for j, feature := range field.Features {
			featuresConverted[j] = types.StringValue(feature)
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

// Read refreshes the Terraform state with the latest data.
func (d *indicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(context.TODO(), "Calling marqo client ListIndices")
	var model allIndicesResourceModel
	indices, err := d.marqoClient.ListIndices()
	if err != nil {
		resp.Diagnostics.AddError("Failed to List Indices", fmt.Sprintf("Could not list indices: %s", err.Error()))
		return
	}

	fmt.Println("Indices: ", indices)

	items := make([]indexModel, len(indices))
	for i, indexDetail := range indices {
		items[i] = indexModel{
			Created:                      types.StringValue(indexDetail.Created),
			IndexName:                    types.StringValue(indexDetail.IndexName),
			NumberOfShards:               types.Int64Value(indexDetail.NumberOfShards),
			NumberOfReplicas:             types.Int64Value(indexDetail.NumberOfReplicas),
			IndexStatus:                  types.StringValue(indexDetail.IndexStatus),
			AllFields:                    ConvertMarqoAllFieldInputs(indexDetail.AllFields),
			TensorFields:                 indexDetail.TensorFields,
			NumberOfInferences:           types.Int64Value(indexDetail.NumberOfInferences),
			StorageClass:                 types.StringValue(indexDetail.StorageClass),
			InferenceType:                types.StringValue(indexDetail.InferenceType),
			DocsCount:                    types.StringValue(indexDetail.DocsCount),
			StoreSize:                    types.StringValue(indexDetail.StoreSize),
			DocsDeleted:                  types.StringValue(indexDetail.DocsDeleted),
			SearchQueryTotal:             types.StringValue(indexDetail.SearchQueryTotal),
			TreatUrlsAndPointersAsImages: types.BoolValue(indexDetail.TreatUrlsAndPointersAsImages),
			MarqoEndpoint:                types.StringValue(indexDetail.MarqoEndpoint),
			Type:                         types.StringValue(indexDetail.Type),
			VectorNumericType:            types.StringValue(indexDetail.VectorNumericType),
			Model:                        types.StringValue(indexDetail.Model),
			NormalizeEmbeddings:          types.BoolValue(indexDetail.NormalizeEmbeddings),
			TextPreprocessing: TextPreprocessingModel{
				SplitLength:  types.Int64Value(indexDetail.TextPreprocessing.SplitLength),
				SplitMethod:  types.StringValue(indexDetail.TextPreprocessing.SplitMethod),
				SplitOverlap: types.Int64Value(indexDetail.TextPreprocessing.SplitOverlap),
			},
			//ImagePreprocessing: types.ObjectValue(map[string]interface{}, indexDetail.ImagePreprocessing),
			AnnParameters: AnnParametersModel{
				SpaceType: types.StringValue(indexDetail.AnnParameters.SpaceType),
				Parameters: parametersModel{
					EfConstruction: types.StringValue(fmt.Sprintf("%d", indexDetail.AnnParameters.Parameters.EfConstruction)),
					M:              types.StringValue(fmt.Sprintf("%d", indexDetail.AnnParameters.Parameters.M)),
				},
			},
			MarqoVersion:          types.StringValue(indexDetail.MarqoVersion),
			FilterStringMaxLength: types.Int64Value(indexDetail.FilterStringMaxLength),
		}
	}

	model.Items = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
