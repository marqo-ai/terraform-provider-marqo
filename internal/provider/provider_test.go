package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
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

func TestAccProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testProviderConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "id", "marqo_read_indices"),
				),
			},
		},
	})
}

const testProviderConfig = `
provider "marqo" {}

data "marqo_read_indices" "test" {}
`
