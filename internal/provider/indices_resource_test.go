package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceIndex(t *testing.T) {
	unstructured_index_name := fmt.Sprintf("unstructured_resource_%s", randomString(8))
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
			// Create and Read testing
			{
				Config: testAccResourceIndexConfig(unstructured_index_name),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Create and Read testing")
						return nil
					},
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", unstructured_index_name),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.type", "unstructured"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.vector_numeric_type", "float"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.treat_urls_and_pointers_as_images", "true"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.model", "hf/e5-small-v2"),
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
					testAccCheckIndexIsReady(unstructured_index_name),
					func(s *terraform.State) error {
						fmt.Println("Create and Read testing completed")
						return nil
					},
				),
			},
			// ImportState testing
			/*
				{
					ResourceName:      "marqo_index.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
			*/
			// Update and Read testing
			{
				Config: testAccResourceIndexConfigUpdated(unstructured_index_name),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Update and Read testing")
						return nil
					},
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", unstructured_index_name),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.CPU.large"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "2"),
					testAccCheckIndexIsReady(unstructured_index_name),
					func(s *terraform.State) error {
						fmt.Println("Update and Read testing completed")
						return nil
					},
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccResourceIndexConfig(name string) string {
	return fmt.Sprintf(`
		resource "marqo-terraform_index" "test" {
		index_name = "%s"
		settings = {
				type = "unstructured"
				vector_numeric_type = "float"
				treat_urls_and_pointers_as_images = true
				model = "hf/e5-small-v2"
				normalize_embeddings = true
				inference_type = "marqo.CPU.small"
				number_of_inferences = 1
				number_of_replicas = 0
				number_of_shards = 1
				storage_class = "marqo.basic"
				all_fields = []
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

func testAccResourceIndexConfigUpdated(name string) string {
	return fmt.Sprintf(`
		resource "marqo-terraform_index" "test" {
		index_name = "%s"
		settings = {
			type = "unstructured"
			vector_numeric_type = "float"
			treat_urls_and_pointers_as_images = true
			model = "hf/e5-small-v2"
			normalize_embeddings = true
			inference_type = "marqo.CPU.large"
			number_of_inferences = 2
			number_of_replicas = 0
			number_of_shards = 1
			storage_class = "marqo.basic"
			all_fields = []
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
