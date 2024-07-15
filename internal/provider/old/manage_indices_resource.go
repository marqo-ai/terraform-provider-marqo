package provider

import (
	"context"
	"fmt"

	"terraform-provider-marqo/marqo"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &indicesResource{}
	_ resource.ResourceWithConfigure = &indicesResource{}
)

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
	SplitLength  types.String `tfsdk:"split_length"`
	SplitMethod  types.String `tfsdk:"split_method"`
	SplitOverlap types.String `tfsdk:"split_overlap"`
}

type AnnParametersModel struct {
	SpaceType  types.String    `tfsdk:"space_type"`
	Parameters parametersModel `tfsdk:"parameters"`
}

type parametersModel struct {
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

// ManageIndicesResource is a helper function to simplify the provider implementation.
func ManageIndicesResource() resource.Resource {
	return &indicesResource{}
}

// orderResource is the resource implementation.
type indicesResource struct {
	marqoClient *marqo.Client
}

// Metadata returns the resource type name.
func (r *indicesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_indices"
}

// Schema defines the schema for the resource.
func (r *indicesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *indicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model allIndicesResourceModel
	req.Plan.Get(ctx, &model)

	for _, item := range model.Items {
		indexName := item.IndexName.ValueString()
		settings := make(map[string]interface{})

		// Reforming index settings based on common fields in schema and models
		settings["indexSettings"] = map[string]interface{}{
			"annParameters": map[string]interface{}{
				"parameters": map[string]interface{}{
					"efConstruction": item.AnnParameters.Parameters.EfConstruction,
					"m":              item.AnnParameters.Parameters.M,
				},
				"spaceType": item.AnnParameters.SpaceType,
			},
			"filterStringMaxLength": item.FilterStringMaxLength,
			//"imagePreprocessing":    item. // Adjusted based on common fields
			"model":               item.Model,
			"normalizeEmbeddings": item.NormalizeEmbeddings,
			"textPreprocessing": map[string]interface{}{
				"splitLength":  item.TextPreprocessing.SplitLength,
				"splitMethod":  item.TextPreprocessing.SplitMethod,
				"splitOverlap": item.TextPreprocessing.SplitOverlap,
			},
			"treatUrlsAndPointersAsImages": item.TreatUrlsAndPointersAsImages,
			"type":                         item.Type,
			"vectorNumericType":            item.VectorNumericType,
			"marqoVersion":                 item.MarqoVersion,  // Added based on common fields
			"marqoEndpoint":                item.MarqoEndpoint, // Added based on common fields
		}

		err := r.marqoClient.CreateIndex(indexName, settings)
		if err != nil {
			resp.Diagnostics.AddError("Failed to Create Index", fmt.Sprintf("Could not create index '%s': %s", indexName, err.Error()))
			return
		}

		// Set ID and other state attributes as needed
		model.ID = types.StringValue(indexName)
		resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *indicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(context.TODO(), "Calling marqo client ListIndices")
	var model allIndicesResourceModel
	indices, err := r.marqoClient.ListIndices()
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
				SplitLength:  types.StringValue(fmt.Sprintf("%d", indexDetail.TextPreprocessing.SplitLength)),
				SplitMethod:  types.StringValue(indexDetail.TextPreprocessing.SplitMethod),
				SplitOverlap: types.StringValue(fmt.Sprintf("%d", indexDetail.TextPreprocessing.SplitOverlap)),
			},
			//ImagePreprocessing: types.ObjectValue(map[string]interface{}, indexDetail.ImagePreprocessing),
			AnnParameters: AnnParametersModel{
				SpaceType: types.StringValue(indexDetail.AnnParameters.SpaceType),
				Parameters: parametersModel{
					EfConstruction: types.Int64Value(indexDetail.AnnParameters.Parameters.EfConstruction),
					M:              types.Int64Value(indexDetail.AnnParameters.Parameters.M),
				},
			},
			MarqoVersion:          types.StringValue(indexDetail.MarqoVersion),
			FilterStringMaxLength: types.Int64Value(indexDetail.FilterStringMaxLength),
		}
	}

	model.Items = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *indicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *indicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model allIndicesResourceModel
	req.State.Get(ctx, &model)

	err := r.marqoClient.DeleteIndex(model.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to Delete Index", fmt.Sprintf("Could not delete index '%s': %s", model.ID.ValueString(), err.Error()))
		return
	}

	// Remove the resource from state by setting it to nil
	resp.State.RemoveResource(ctx)
}
