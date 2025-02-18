//go:build customer_a

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestA(t *testing.T) {
	// Skip in normal test runs
	if testing.Short() {
		t.Skip("Skipping mock prod tests in short mode")
	}
	t.Parallel()

	// Test for production-tenant-v2 index
	t.Run("mock_prod", func(t *testing.T) {
		indexName := fmt.Sprintf("donotdelete_prod_%s", randomString(4))
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccEmptyConfig(),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckIndexExistsAndDelete(indexName),
					),
				},
				// Create initial index
				{
					Config: testAProductionTenantConfig(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "index_name", indexName),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.type", "unstructured"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.model", "marqo-fashion-clip"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.GPU"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "2"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_replicas", "1"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_shards", "1"),
					),
				},
				// Modify index
				{
					Config: testAProductionTenantConfigModified(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "3"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_replicas", "1"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_shards", "1"),
					),
				},
				// Delete and recreate testing
				{
					Config: testAccEmptyConfig(),
					Check: resource.ComposeTestCheckFunc(
						func(s *terraform.State) error {
							fmt.Println("Starting Production Tenant delete testing")
							return nil
						},
					),
				},
				// Recreate index
				{
					Config: testAProductionTenantConfig(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "index_name", indexName),
					),
				},
			},
		})
	})
}

func testAProductionTenantConfig(name string) string {
	return fmt.Sprintf(`
	resource "marqo_index" "test" {
		index_name = "%s"
		timeouts = {
			create = "45m"
			update = "45m"
			delete = "20m"
		}
		settings = {
			type = "unstructured"
			vector_numeric_type = "float"
			model = "marqo-fashion-clip"
			model_properties = {
				dimensions = 512
				name = "ViT-B-16"
				type = "open_clip"
				url = "https://marqo-gcl-public.s3.us-west-2.amazonaws.com/marqo-fashionCLIP/marqo_fashionCLIP.pt"
			}
			normalize_embeddings = true
			inference_type = "marqo.GPU"
			number_of_inferences = 2
			number_of_replicas = 1
			number_of_shards = 1
			storage_class = "marqo.balanced"
			text_preprocessing = {
				split_length = 2
				split_method = "sentence"
				split_overlap = 0
			}
			treat_urls_and_pointers_as_images = true
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

func testAProductionTenantConfigModified(name string) string {
	return fmt.Sprintf(`
	resource "marqo_index" "test" {
		index_name = "%s"
		timeouts = {
			create = "45m"
			update = "45m"
			delete = "20m"
		}
		settings = {
			type = "unstructured"
			vector_numeric_type = "float"
			model = "marqo-fashion-clip"
			model_properties = {
				dimensions = 512
				name = "ViT-B-16"
				type = "open_clip"
				url = "https://marqo-gcl-public.s3.us-west-2.amazonaws.com/marqo-fashionCLIP/marqo_fashionCLIP.pt"
			}
			normalize_embeddings = true
			inference_type = "marqo.GPU"
			number_of_inferences = 3
			number_of_replicas = 1
			number_of_shards = 1
			storage_class = "marqo.balanced"
			text_preprocessing = {
				split_length = 2
				split_method = "sentence"
				split_overlap = 0
			}
			treat_urls_and_pointers_as_images = true
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
