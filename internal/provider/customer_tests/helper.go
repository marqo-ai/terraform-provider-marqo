package customer_tests

import (
	"context"
	"marqo/go_marqo"
	marqoprovider "marqo/internal/provider"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Provider is the provider implementation.
type MarqoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	Host   types.String `tfsdk:"host"`
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *MarqoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "marqo"
	resp.Version = p.version
}

func (p *MarqoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *MarqoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("MARQO_HOST")
	apiKey := os.Getenv("MARQO_API_KEY")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddError(
			"Missing Marqo Host",
			"Cannot create client without host. Set the host in the provider configuration or use the MARQO_HOST environment variable.",
		)
		return
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing Marqo API Key",
			"Cannot create client without API key. Set the api_key in the provider configuration or use the MARQO_API_KEY environment variable.",
		)
		return
	}

	client, err := go_marqo.NewClient(&host, &apiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Marqo Client",
			"An unexpected error occurred when creating the Marqo client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Marqo Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *MarqoProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		marqoprovider.ManageIndicesResource,
	}
}

func (p *MarqoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MarqoProvider{
			version: version,
		}
	}
}

// Copy of provider test utilities
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"marqo": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("MARQO_HOST"); v == "" {
		t.Fatal("MARQO_HOST must be set for acceptance tests")
	}
	if v := os.Getenv("MARQO_API_KEY"); v == "" {
		t.Fatal("MARQO_API_KEY must be set for acceptance tests")
	}
}

func testAccEmptyConfig() string {
	return `
	provider "marqo" {}
	`
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
