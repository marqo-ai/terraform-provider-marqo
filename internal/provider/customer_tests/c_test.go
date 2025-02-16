package customer_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// To run these tests manually:
// go test -v ./internal/provider/customer_tests -run TestRedbubble

func TestC(t *testing.T) {
	// Skip in normal test runs
	if testing.Short() {
		t.Skip("Skipping customer-specific test in short mode")
	}
	t.Parallel()

	// Test for sr-marqo-prod-2-3 index
	t.Run("mock_prod_sr", func(t *testing.T) {
		indexName := fmt.Sprintf("donotdelete_sr_%s", randomString(4))
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create initial index
				{
					Config: testCProductionTenantConfig(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "index_name", indexName),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.type", "structured"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.model", "ViT-L-14-finetune-2-2"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.GPU"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "1"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_replicas", "1"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_shards", "2"),
					),
				},
				// Modify index
				{
					Config: testCProductionTenantConfigModified(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "2"),
					),
				},
				// Delete and recreate testing
				{
					Config: testAccEmptyConfig(),
					Check: resource.ComposeTestCheckFunc(
						func(s *terraform.State) error {
							fmt.Println("Starting SR Marqo Prod delete testing")
							return nil
						},
					),
				},
				// Recreate index
				{
					Config: testCProductionTenantConfig(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "index_name", indexName),
					),
				},
			},
		})
	})
}

func testCProductionTenantConfig(name string) string {
	return fmt.Sprintf(`
	resource "marqo_index" "test" {
		index_name = "%s"
		settings = {
			type = "structured"
			vector_numeric_type = "float"
			model = "ViT-L-14-finetune-2-2"
			model_properties = {
				dimensions = 768
				name = "ViT-L-14"
				type = "open_clip"
				url = "https://7b4d1a66-507d-43f1-b99f-7368b655de46.s3.amazonaws.com/e5a7d9c7-0736-4301-a037-b1307f43a314/24d04007-a0df-49ba-a712-1cd02b44d2d9.pt"
			}
			normalize_embeddings = true
			inference_type = "marqo.GPU"
			number_of_inferences = 1
			number_of_replicas = 1
			number_of_shards = 2
			storage_class = "marqo.balanced"
			all_fields = [
				{
					name = "work_id"
					type = "int"
					features = ["filter"]
				},
				{
					name = "image_url"
					type = "image_pointer"
				},
				{
					name = "prefixed_tags"
					type = "text"
				},
				{
					name = "prefixed_title"
					type = "text"
				},
				{
					name = "tags"
					type = "text"
					features = ["lexical_search"]
				},
				{
					name = "title"
					type = "text"
					features = ["lexical_search"]
				},
				{
					name = "image_with_meta_tensor"
					type = "multimodal_combination"
					dependent_fields = {
						image_url = 0.8
						prefixed_tags = 0.05
						prefixed_title = 0.1
					}
				}
			]
			tensor_fields = ["image_with_meta_tensor"]
			ann_parameters = {
				space_type = "prenormalized-angular"
				parameters = {
					ef_construction = 512
					m = 16
				}
			}
		}
	}
	`, name)
}

func testCProductionTenantConfigModified(name string) string {
	return fmt.Sprintf(`
	resource "marqo_index" "test" {
		index_name = "%s"
		settings = {
			type = "structured"
			vector_numeric_type = "float"
			model = "ViT-L-14-finetune-2-2"
			model_properties = {
				dimensions = 768
				name = "ViT-L-14"
				type = "open_clip"
				url = "https://7b4d1a66-507d-43f1-b99f-7368b655de46.s3.amazonaws.com/e5a7d9c7-0736-4301-a037-b1307f43a314/24d04007-a0df-49ba-a712-1cd02b44d2d9.pt"
			}
			normalize_embeddings = true
			inference_type = "marqo.GPU"
			number_of_inferences = 2
			number_of_replicas = 1
			number_of_shards = 2
			storage_class = "marqo.balanced"
			all_fields = [
				{
					name = "work_id"
					type = "int"
					features = ["filter"]
				},
				{
					name = "image_url"
					type = "image_pointer"
				},
				{
					name = "prefixed_tags"
					type = "text"
				},
				{
					name = "prefixed_title"
					type = "text"
				},
				{
					name = "tags"
					type = "text"
					features = ["lexical_search"]
				},
				{
					name = "title"
					type = "text"
					features = ["lexical_search"]
				},
				{
					name = "image_with_meta_tensor"
					type = "multimodal_combination"
					dependent_fields = {
						image_url = 0.8
						prefixed_tags = 0.05
						prefixed_title = 0.1
					}
				}
			]
			tensor_fields = ["image_with_meta_tensor"]
			ann_parameters = {
				space_type = "prenormalized-angular"
				parameters = {
					ef_construction = 512
					m = 16
				}
			}
		}
	}
	`, name)
}
