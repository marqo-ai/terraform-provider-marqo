package customer_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// To run these tests manually:
// go test -v ./internal/provider/customer_tests -run TestB

func TestB(t *testing.T) {
	// Skip in normal test runs
	if testing.Short() {
		t.Skip("Skipping customer-specific test in short mode")
	}
	t.Parallel()
	// Test for production_secondary_photo index
	t.Run("mock_prod_photo", func(t *testing.T) {
		indexName := fmt.Sprintf("donotdelete_photo_%s", randomString(4))
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create initial index
				{
					Config: testBProductionSecondaryPhotoConfig(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "index_name", indexName),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.type", "unstructured"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.model", "open_clip/ViT-L-14/laion2b_s32b_b82k"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.GPU"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "6"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_replicas", "1"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_shards", "2"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.storage_class", "marqo.balanced"),
					),
				},
				// Modify index
				{
					Config: testBProductionSecondaryPhotoConfigModified(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "7"),
					),
				},
				// Delete and recreate testing
				{
					Config: testAccEmptyConfig(),
					Check: resource.ComposeTestCheckFunc(
						func(s *terraform.State) error {
							fmt.Println("Starting Secondary Photo delete testing")
							return nil
						},
					),
				},
				// Recreate index
				{
					Config: testBProductionSecondaryPhotoConfig(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "index_name", indexName),
					),
				},
			},
		})
	})

	// Test for prod_teal_primary_music_text index
	t.Run("mock_prod_music_text", func(t *testing.T) {
		indexName := fmt.Sprintf("test_mock_music_text_%s", randomString(6))
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create initial index
				{
					Config: testBProductionTenantConfig(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "index_name", indexName),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.type", "unstructured"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.model", "hf/e5-base-v2"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.GPU"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "3"),
					),
				},
				// Modify index
				{
					Config: testBProductionTenantConfigModified(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "4"),
					),
				},
				// Delete and recreate testing
				{
					Config: testAccEmptyConfig(),
					Check: resource.ComposeTestCheckFunc(
						func(s *terraform.State) error {
							fmt.Println("Starting Teal Music Text delete testing")
							return nil
						},
					),
				},
				// Recreate index
				{
					Config: testBProductionTenantConfig(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "index_name", indexName),
					),
				},
			},
		})
	})
}

func testBProductionSecondaryPhotoConfig(name string) string {
	return fmt.Sprintf(`
	resource "marqo_index" "test" {
		index_name = "%s"
		settings = {
			type = "unstructured"
			vector_numeric_type = "float"
			treat_urls_and_pointers_as_images = true
			treat_urls_and_pointers_as_media = false
			model = "open_clip/ViT-L-14/laion2b_s32b_b82k"
			normalize_embeddings = true
			inference_type = "marqo.GPU"
			number_of_inferences = 6
			number_of_replicas = 1
			number_of_shards = 2
			storage_class = "marqo.balanced"
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
			filter_string_max_length = 50
		}
	}
	`, name)
}

func testBProductionSecondaryPhotoConfigModified(name string) string {
	return fmt.Sprintf(`
	resource "marqo_index" "test" {
		index_name = "%s"
		settings = {
			type = "unstructured"
			vector_numeric_type = "float"
			treat_urls_and_pointers_as_images = true
			treat_urls_and_pointers_as_media = false
			model = "open_clip/ViT-L-14/laion2b_s32b_b82k"
			normalize_embeddings = true
			inference_type = "marqo.GPU"
			number_of_inferences = 1
			number_of_replicas = 1
			number_of_shards = 2
			storage_class = "marqo.balanced"
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
			filter_string_max_length = 50
		}
	}
	`, name)
}

func testBProductionTenantConfig(name string) string {
	return fmt.Sprintf(`
	resource "marqo_index" "test" {
		index_name = "%s"
		settings = {
			type = "unstructured"
			vector_numeric_type = "float"
			treat_urls_and_pointers_as_images = false
			treat_urls_and_pointers_as_media = false
			model = "hf/e5-base-v2"
			normalize_embeddings = true
			inference_type = "marqo.GPU"
			number_of_inferences = 1
			number_of_replicas = 1
			number_of_shards = 2
			storage_class = "marqo.balanced"
			text_preprocessing = {
				split_length = 20
				split_method = "sentence"
				split_overlap = 2
			}
			image_preprocessing = {}
			ann_parameters = {
				space_type = "prenormalized-angular"
				parameters = {
					ef_construction = 512
					m = 16
				}
			}
			filter_string_max_length = 50
		}
	}
	`, name)
}

func testBProductionTenantConfigModified(name string) string {
	return fmt.Sprintf(`
	resource "marqo_index" "test" {
		index_name = "%s"
		settings = {
			type = "unstructured"
			vector_numeric_type = "float"
			treat_urls_and_pointers_as_images = false
			treat_urls_and_pointers_as_media = false
			model = "hf/e5-base-v2"
			normalize_embeddings = true
			inference_type = "marqo.GPU"
			number_of_inferences = 2
			number_of_replicas = 1
			number_of_shards = 2
			storage_class = "marqo.balanced"
			text_preprocessing = {
				split_length = 20
				split_method = "sentence"
				split_overlap = 2
			}
			image_preprocessing = {}
			ann_parameters = {
				space_type = "prenormalized-angular"
				parameters = {
					ef_construction = 512
					m = 16
				}
			}
			filter_string_max_length = 50
		}
	}
	`, name)
}
