package provider

import (
    "context"
    "net/http"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure our provider satisfies the provider.Provider interface
var _ provider.Provider = &MarqoProvider{}


// MarqoProvider defines the provider implementation for Marqo.
type MarqoProvider struct {
	client *http.Client
	api_key string
}
	
// MarqoProviderData describes the provider data model.
type MarqoProviderData struct {
	APIKey types.String `tfsdk:"api_key"`
}

func (p *MarqoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "marqo"
	// Set the provider version here if you have one
}

func (p *MarqoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": {
				Type: types.StringType,
				Required: true,
				Description: "API Key for accessing Marqo API",
				Sensitive: true,
			},
		},
	}
}

func (p *MarqoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config MarqoProviderData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Correctly check if APIKey is set and then assign it
	if !config.APIKey.Null && !config.APIKey.Unknown {
		p.api_key = config.APIKey.Value
	} else {
		// Handle the case where APIKey is not set
		resp.Diagnostics.AddError(
			"Missing API Key",
			"API Key must be specified for accessing Marqo API",
		)
		return
	}

	p.client = &http.Client{
		// Optionally, set up HTTP client parameters (Timeout, Transport, etc.)
	}

	// Store the configured client in the provider data
	resp.ResourceData = p
}

func (p *MarqoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCreateIndexResource,
		NewDeleteIndexResource,
		NewReadIndicesStatusesResource,
	}
}
	
// If you have any data sources, implement them here
func (p *MarqoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Example: NewExampleDataSource, // Uncomment and replace with actual data source functions if you have any
	}
}

// New returns a new instance of the MarqoProvider.
func New() func() provider.Provider {
	return func() provider.Provider {
		return &MarqoProvider{}
	}
}