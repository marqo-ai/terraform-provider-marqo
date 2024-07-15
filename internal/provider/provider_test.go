package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a new provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"marqo": providerserver.NewProtocol6WithError(New("test")()),
}

func TestAccPreCheck(t *testing.T) {
	// Remove the return statement
	// add code here to check if the environment is properly
	// configured for the test case. example:
	// if v := os.Getenv("MARQO_API_KEY"); v == "" {
	//     t.Fatal("MARQO_API_KEY must be set for acceptance tests")
	// }
}

func TestProvider(t *testing.T) {
	p := New("test")()
	var resp provider.MetadataResponse
	p.Metadata(context.Background(), provider.MetadataRequest{}, &resp)
	// Remove error checking as Metadata doesn't return an error
}
