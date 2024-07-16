package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceIndices(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create an index to ensure we have data to read
			{
				Config: testAccIndexResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", "example_index_3"),
				),
			},
			// Read indices
			{
				Config: testAccDataSourceIndicesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.marqo_read_indices.test", "items.#"),
					resource.TestCheckResourceAttrSet("data.marqo_read_indices.test", "items.0.index_name"),
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "items.0.type", "unstructured"),
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "items.0.vector_numeric_type", "float"),
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "items.0.treat_urls_and_pointers_as_images", "true"),
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "items.0.model", "open_clip/ViT-L-14/laion2b_s32b_b82k"),
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "items.0.normalize_embeddings", "true"),
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "items.0.inference_type", "marqo.CPU.small"),
				),
			},
		},
	})
}

const testAccDataSourceIndicesConfig = `
data "marqo_read_indices" "test" {}
`

const testAccIndexResourceConfig = `
resource "marqo_index" "test" {
  index_name = "example_index_3"
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
`
