package provider

import (
	"fmt"
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
			// Create and Read testing
			{
				Config: testAccResourceIndexConfig("example_index_3"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", "example_index_3"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.type", "unstructured"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.vector_numeric_type", "float"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.treat_urls_and_pointers_as_images", "true"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.model", "open_clip/ViT-L-14/laion2b_s32b_b82k"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.normalize_embeddings", "true"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.CPU.small"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "1"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_replicas", "0"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_shards", "1"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.storage_class", "marqo.basic"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.text_preprocessing.split_length", "2"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.text_preprocessing.split_method", "sentence"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.text_preprocessing.split_overlap", "0"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.ann_parameters.space_type", "prenormalized-angular"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.ann_parameters.parameters.ef_construction", "512"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.ann_parameters.parameters.m", "16"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.filter_string_max_length", "20"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "marqo_index.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccResourceIndexConfig("example_index_3_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", "example_index_3_updated"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.CPU.large"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccResourceIndexConfig(name string) string {
	return fmt.Sprintf(`
resource "marqo_index" "test" {
  index_name = "%s"
  settings = {
    type = "unstructured"
    vector_numeric_type = "float"
    treat_urls_and_pointers_as_images = true
    model = "open_clip/ViT-L-14/laion2b_s32b_b82k"
    normalize_embeddings = true
    inference_type = "marqo.CPU.small"
    all_fields = []
    number_of_inferences = 1
    number_of_replicas = 0
    number_of_shards = 1
    storage_class = "marqo.basic"
    text_preprocessing = {
      split_length = 2
      split_method = "sentence"
      split_overlap = 0
    }
    image_preprocessing = {}
    ann_parameters = {
      space_type = "prenormalized-angular"
      parameters = {
        ef_construction = 512
        m = 16
      }
    }
    filter_string_max_length = 20
  }
}
`, name)
}
