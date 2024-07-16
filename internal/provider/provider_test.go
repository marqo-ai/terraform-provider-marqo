package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

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
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig,
				Check: resource.ComposeTestCheckFunc(
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
