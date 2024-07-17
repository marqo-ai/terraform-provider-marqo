package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceIndices(t *testing.T) {
	unstructured_index_name := fmt.Sprintf("unstructured_dsource_%s", randomString(9))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check if index exists and delete if it does
			{
				Config: testAccEmptyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExistsAndDelete(unstructured_index_name),
				),
			},
			// Create an index to ensure we have data to read
			{
				Config: testAccDataSourceIndexConfig(unstructured_index_name),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Create index")
						return nil
					},
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", unstructured_index_name),
					testAccCheckIndexIsReady(unstructured_index_name),
					func(s *terraform.State) error {
						fmt.Println("Finished Create index")
						return nil
					},
				),
			},

			// Read the index using the data source
			{
				Config: testAccDataSourceIndicesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexInDataSource("data.marqo_read_indices.test", unstructured_index_name),
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

func testAccEmptyConfig() string {
	return `
    # Empty config
    `
}

func testAccDataSourceIndexConfig(name string) string {
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

func testAccCheckIndexInDataSource(dataSourceName string, indexName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("Data source not found: %s", dataSourceName)
		}

		itemCount, err := strconv.Atoi(ds.Primary.Attributes["items.#"])
		if err != nil {
			return fmt.Errorf("Error parsing item count: %s", err)
		}

		for i := 0; i < itemCount; i++ {
			if ds.Primary.Attributes[fmt.Sprintf("items.%d.index_name", i)] == indexName {
				// Check attributes for this specific index
				attributesToCheck := map[string]string{
					"inference_type": "marqo.CPU.small",
					"storage_class":  "marqo.basic",
					// Add more attributes to check
				}

				for attr, expectedValue := range attributesToCheck {
					actualValue := ds.Primary.Attributes[fmt.Sprintf("items.%d.%s", i, attr)]
					if actualValue != expectedValue {
						return fmt.Errorf("Attribute %s does not match for index %s. Expected: %s, Got: %s",
							attr, indexName, expectedValue, actualValue)
					}
				}

				// If all checks pass, return nil (success)
				return nil
			}
		}

		return fmt.Errorf("Index %s not found in data source results", indexName)
	}
}
