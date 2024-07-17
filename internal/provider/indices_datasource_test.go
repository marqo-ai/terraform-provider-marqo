package provider

import (
	"fmt"
	"strconv"
	"terraform-provider-marqo/marqo"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceIndices(t *testing.T) {
	unstructured_index_name := fmt.Sprintf("unstructured_dsource_%s", randomString(9))
	var foundIndex *marqo.IndexDetail

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

			// Read indices
			{
				Config: testAccDataSourceIndicesConfig,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Read indices")
						return nil
					},
					func(s *terraform.State) error {
						var err error
						foundIndex, err = findIndexByName(unstructured_index_name)(s)
						if err != nil {
							t.Logf("Error finding index: %v", err)
							return err
						}
						if foundIndex == nil {
							t.Log("Index was not found but no error was returned")
							return fmt.Errorf("Index was not found but no error was returned")
						}
						t.Logf("Found index: %+v", *foundIndex)
						return nil
					},
					func(s *terraform.State) error {
						if foundIndex == nil {
							t.Log("foundIndex is nil, cannot compare with data source")
							return fmt.Errorf("foundIndex is nil, cannot compare with data source")
						}
						err := compareFoundIndexWithDataSource("data.marqo_read_indices.test", foundIndex)(s)
						if err != nil {
							t.Logf("Error comparing found index with data source: %v", err)
						}
						return err
					},
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

func compareFoundIndexWithDataSource(dataSourceName string, foundIndex *marqo.IndexDetail) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("Data source not found: %s", dataSourceName)
		}

		itemCount, err := strconv.Atoi(ds.Primary.Attributes["items.#"])
		if err != nil {
			return fmt.Errorf("Error parsing item count: %s", err)
		}

		inferenceTypeMap := map[string]string{
			"CPU.SMALL": "marqo.CPU.small",
			"CPU.LARGE": "marqo.CPU.large",
			"GPU":       "marqo.GPU",
		}
		storageClassMap := map[string]string{
			"BASIC":       "marqo.basic",
			"BALANCED":    "marqo.balanced",
			"PERFORMANCE": "marqo.performance",
		}

		for i := 0; i < itemCount; i++ {
			if ds.Primary.Attributes[fmt.Sprintf("items.%d.index_name", i)] == foundIndex.IndexName {
				attributesToCheck := map[string]string{
					"type":                              foundIndex.Type,
					"vector_numeric_type":               foundIndex.VectorNumericType,
					"treat_urls_and_pointers_as_images": strconv.FormatBool(foundIndex.TreatUrlsAndPointersAsImages),
					"model":                             foundIndex.Model,
					"normalize_embeddings":              strconv.FormatBool(foundIndex.NormalizeEmbeddings),
					"inference_type":                    inferenceTypeMap[foundIndex.InferenceType],
					"storage_class":                     storageClassMap[foundIndex.StorageClass],
				}
				for attr, expectedValue := range attributesToCheck {
					dsValue := ds.Primary.Attributes[fmt.Sprintf("items.%d.%s", i, attr)]
					if dsValue != expectedValue {
						return fmt.Errorf("Attribute %s does not match. Data source: %s, Expected (mapped from API): %s", attr, dsValue, expectedValue)
					}
				}
				return nil
			}
		}

		return fmt.Errorf("Index %s not found in data source results", foundIndex.IndexName)
	}
}
