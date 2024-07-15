package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccPreCheck(t *testing.T) {
	// Verify required environment variables are set
	if v := os.Getenv("MARQO_API_KEY"); v == "" {
		t.Fatal("MARQO_API_KEY must be set for acceptance tests")
	}
	// Add other required environment variables checks
}

func TestAccResourceIndex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIndexConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", "test_index"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.type", "unstructured"),
					// Add more checks based on resource schema
				),
			},
			// add more steps to test updates
		},
	})
}

const testAccResourceIndexConfig = `
resource "marqo_index" "test" {
  index_name = "test_index"
  settings = {
    type = "unstructured"
    vector_numeric_type = "float"
    treat_urls_and_pointers_as_images = true
    model = "open_clip/ViT-L-14/laion2b_s32b_b82k"
    normalize_embeddings = true
    inference_type = "marqo.CPU.small"
  }
}
`
