package provider

import (
	"context"
	"os"

	"terraform-provider-marqo/marqo"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &marqoProvider{}
)

// marqoProviderModel maps provider schema data to a Go type.
type marqoProviderModel struct {
	Host   types.String `tfsdk:"host"`
	APIKey types.String `tfsdk:"api_key"`
}

// marqoProvider is the provider implementation.
type marqoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &marqoProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *marqoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "marqo"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *marqoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "The Marqo API host. Can be set with MARQO_HOST environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The Marqo API key. Can be set with MARQO_API_KEY environment variable.",
			},
		},
	}
}

// Configure prepares a Marqo API client for data sources and resources.
func (p *marqoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Marqo client")

	var config marqoProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Marqo API Host",
			"The provider cannot create the Marqo API client as there is an unknown configuration value for the Marqo API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MARQO_HOST environment variable.",
		)
	}

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Marqo API Key",
			"The provider cannot create the Marqo API client as there is an unknown configuration value for the Marqo API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MARQO_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("MARQO_HOST")
	apiKey := os.Getenv("MARQO_API_KEY")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Marqo API Host",
			"The provider cannot create the Marqo API client as there is a missing or empty value for the Marqo API host. "+
				"Set the host value in the configuration or use the MARQO_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Marqo API Key",
			"The provider cannot create the Marqo API client as there is a missing or empty value for the Marqo API key. "+
				"Set the api_key value in the configuration or use the MARQO_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "marqo_host", host)
	ctx = tflog.SetField(ctx, "marqo_api_key", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "marqo_api_key")

	tflog.Debug(ctx, "Creating Marqo client")

	client, err := marqo.NewClient(&host, &apiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Marqo API Client",
			"An unexpected error occurred when creating the Marqo API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Marqo Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Marqo client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *marqoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		ReadIndicesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *marqoProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		ManageIndicesResource,
	}
}
