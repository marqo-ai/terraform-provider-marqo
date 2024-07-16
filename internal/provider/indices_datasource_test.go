package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
					func(s *terraform.State) error {
						fmt.Println("Starting Create index")
						return nil
					},
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", "example_index_datasource_2"),
					testAccCheckIndexIsReady("example_index_datasource_2"),
					func(s *terraform.State) error {
						fmt.Println("Finished Create index")
						return nil
					},
				),
			},

			// Read indices
			{
				Config: testAccDataSourceIndicesConfig,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Read indices")
						return nil
					},
					resource.TestCheckResourceAttrSet("data.marqo_read_indices.test", "items.#"),
					resource.TestCheckResourceAttrSet("data.marqo_read_indices.test", "items.0.index_name"),
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "items.0.type", "unstructured"),
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "items.0.vector_numeric_type", "float"),
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "items.0.treat_urls_and_pointers_as_images", "true"),
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "items.0.model", "open_clip/ViT-L-14/laion2b_s32b_b82k"),
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "items.0.normalize_embeddings", "true"),
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "items.0.inference_type", "marqo.CPU.small"),
					//testAccCheckIndexIsReady("example_index_datasource_2"),
					func(s *terraform.State) error {
						fmt.Println("Finished Read indices")
						return nil
					},
				),
			},
		},
	})
}

const testAccDataSourceIndicesConfig = `
data "marqo_read_indices" "test" {
	id = "test_id_1"
}
`

const testAccIndexResourceConfig = `
resource "marqo_index" "test" {
  index_name = "example_index_datasource_2"
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
