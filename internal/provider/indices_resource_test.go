package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceCustomModelIndex(t *testing.T) {
	t.Parallel() // Enable parallel testing
	unstructured_custom_model_index_name := fmt.Sprintf("donotdelete_unstr_resrc_%s", randomString(6))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check if index exists and delete if it does
			{
				Config: testAccEmptyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExistsAndDelete(unstructured_custom_model_index_name),
				),
			},
			// Create Custom Model Index
			{
				Config: testAccResourceIndexConfigCustomModel(unstructured_custom_model_index_name),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Custom Model testing")
						return nil
					},
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", unstructured_custom_model_index_name),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.type", "unstructured"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.model", "custom-model"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.model_properties.url", "https://marqo-ecs-50-audio-test-dataset.s3.us-east-1.amazonaws.com/test-hf.zip"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.model_properties.dimensions", "384"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.model_properties.type", "hf"),
					testAccCheckIndexIsReady(unstructured_custom_model_index_name),
					func(s *terraform.State) error {
						fmt.Println("Custom Model testing completed")
						return nil
					},
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccResourceLangBindIndex(t *testing.T) {
	t.Parallel() // Enable parallel testing
	unstructured_langbind_index_name := fmt.Sprintf("donotdelete_unstr_resrc_%s", randomString(6))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check if index exists and delete if it does
			{
				Config: testAccEmptyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExistsAndDelete(unstructured_langbind_index_name),
				),
			},
			// Create and Read testing
			{
				Config: testAccResourceIndexConfig(unstructured_langbind_index_name),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Create and Read testing")
						return nil
					},
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", unstructured_langbind_index_name),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.type", "unstructured"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.vector_numeric_type", "float"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.treat_urls_and_pointers_as_images", "true"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.treat_urls_and_pointers_as_media", "true"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.model", "LanguageBind/Video_V1.5_FT_Audio_FT_Image"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.normalize_embeddings", "true"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.GPU"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "1"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_replicas", "0"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_shards", "1"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.storage_class", "marqo.basic"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.text_preprocessing.split_length", "2"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.text_preprocessing.split_method", "sentence"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.text_preprocessing.split_overlap", "0"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.video_preprocessing.split_length", "5"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.video_preprocessing.split_overlap", "1"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.audio_preprocessing.split_length", "5"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.audio_preprocessing.split_overlap", "1"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.ann_parameters.space_type", "prenormalized-angular"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.ann_parameters.parameters.ef_construction", "512"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.ann_parameters.parameters.m", "16"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.filter_string_max_length", "20"),
					testAccCheckIndexIsReady(unstructured_langbind_index_name),
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
				Config: testAccResourceIndexConfigUpdated(unstructured_langbind_index_name),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Update and Read testing")
						return nil
					},
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", unstructured_langbind_index_name),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.GPU"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "2"),
					testAccCheckIndexIsReady(unstructured_langbind_index_name),
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

func TestAccResourceStructuredIndex(t *testing.T) {
	t.Parallel() // Enable parallel testing
	structured_index_name := fmt.Sprintf("donotdelete_str_rsrc_%s", randomString(7))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check if index exists and delete if it does
			{
				Config: testAccEmptyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExistsAndDelete(structured_index_name),
				),
			},
			// Create and Read testing
			{
				Config: testAccResourceStructuredIndexConfig(structured_index_name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", structured_index_name),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.type", "structured"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.vector_numeric_type", "float"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.model", "open_clip/ViT-L-14/laion2b_s32b_b82k"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.normalize_embeddings", "true"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.CPU.large"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "1"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_replicas", "0"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_shards", "2"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.storage_class", "marqo.basic"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.text_preprocessing.split_length", "2"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.text_preprocessing.split_method", "sentence"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.text_preprocessing.split_overlap", "0"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.ann_parameters.space_type", "prenormalized-angular"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.ann_parameters.parameters.ef_construction", "512"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.ann_parameters.parameters.m", "16"),
					testAccCheckIndexIsReady(structured_index_name),
				),
			},
			// Check for no changes on re-apply
			{
				Config:             testAccResourceStructuredIndexConfig(structured_index_name),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccResourceMinimalIndex(t *testing.T) {
	t.Parallel() // Enable parallel testing
	minimal_index_name := fmt.Sprintf("donotdelete_min_%s", randomString(6))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Check if index exists and delete if it does
			{
				Config: testAccEmptyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexExistsAndDelete(minimal_index_name),
				),
			},
			// Create and Read testing
			{
				Config: testAccResourceMinimalIndexConfig(minimal_index_name),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Minimal Index testing")
						return nil
					},
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", minimal_index_name),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.type", "unstructured"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.model", "open_clip/ViT-L-14/laion2b_s32b_b82k"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.CPU.large"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "1"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_replicas", "0"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_shards", "1"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.storage_class", "marqo.basic"),
					resource.TestCheckResourceAttr("marqo_index.test", "timeouts.create", "45m"),
					resource.TestCheckResourceAttr("marqo_index.test", "timeouts.update", "45m"),
					resource.TestCheckResourceAttr("marqo_index.test", "timeouts.delete", "20m"),
					func(s *terraform.State) error {
						fmt.Println("Minimal Index testing completed")
						return nil
					},
				),
			},
			// Update testing
			{
				Config: testAccResourceMinimalIndexConfigUpdated(minimal_index_name),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Minimal Index update testing")
						return nil
					},
					resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "2"),
					resource.TestCheckResourceAttr("marqo_index.test", "timeouts.create", "45m"),
					resource.TestCheckResourceAttr("marqo_index.test", "timeouts.update", "45m"),
					resource.TestCheckResourceAttr("marqo_index.test", "timeouts.delete", "20m"),
					func(s *terraform.State) error {
						fmt.Println("Minimal Index update testing completed")
						return nil
					},
				),
			},
			// Delete and recreate testing
			{
				Config: testAccEmptyConfig(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Minimal Index delete testing")
						return nil
					},
				),
			},
			// Recreate testing
			{
				Config: testAccResourceMinimalIndexConfig(minimal_index_name),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Minimal Index recreation testing")
						return nil
					},
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", minimal_index_name),
					resource.TestCheckResourceAttr("marqo_index.test", "timeouts.create", "45m"),
					resource.TestCheckResourceAttr("marqo_index.test", "timeouts.update", "45m"),
					resource.TestCheckResourceAttr("marqo_index.test", "timeouts.delete", "20m"),
					func(s *terraform.State) error {
						fmt.Println("Minimal Index recreation testing completed")
						return nil
					},
				),
			},
			// Final deletion occurs automatically
		},
	})
}

func testAccResourceStructuredIndexConfig(name string) string {
	return fmt.Sprintf(`
    resource "marqo_index" "test" {
        index_name = "%s"
		timeouts = {
			create = "45m"
			update = "45m"
			delete = "20m"
		}
        settings = {
            type                = "structured"
            vector_numeric_type = "float"
            all_fields = [
                { "name" : "text_field_0", "type" : "text", "features" : ["lexical_search"] },
                { "name" : "text_field", "type" : "text", "features" : ["lexical_search"] },
                { "name" : "image_field", "type" : "image_pointer" },
                {
                    "name" : "multimodal_field",
                    "type" : "multimodal_combination",
                    "dependent_fields" : {
                        "imageField" : 0.8,
                        "textField" : 0.1
                    },
                },
            ]
            number_of_inferences = 1
            storage_class        = "marqo.basic"
            number_of_replicas   = 0
            number_of_shards     = 2
            tensor_fields        = ["multimodal_field"]
            model                = "open_clip/ViT-L-14/laion2b_s32b_b82k"
            normalize_embeddings = true
            inference_type       = "marqo.CPU.large"
            text_preprocessing = {
                split_length  = 2
                split_method  = "sentence"
                split_overlap = 0
            }
            image_preprocessing = {}
            ann_parameters = {
                space_type = "prenormalized-angular"
                parameters = {
                    ef_construction = 512
                    m               = 16
                }
            }
        }
    }
    `, name)
}

func testAccResourceIndexConfig(name string) string {
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
				treat_urls_and_pointers_as_images = true
				treat_urls_and_pointers_as_media = true
				model = "LanguageBind/Video_V1.5_FT_Audio_FT_Image"
				normalize_embeddings = true
				inference_type = "marqo.GPU"
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
				video_preprocessing = {
					split_length = 5
					split_overlap = 1
				}
				audio_preprocessing = {
					split_length = 5
					split_overlap = 1
				}
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
				treat_urls_and_pointers_as_images = true
				treat_urls_and_pointers_as_media = true
				model = "LanguageBind/Video_V1.5_FT_Audio_FT_Image"
				normalize_embeddings = true
				inference_type = "marqo.GPU"
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
				video_preprocessing = {
					split_length = 5
					split_overlap = 1
				}
				audio_preprocessing = {
					split_length = 5
					split_overlap = 1
				}
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

func testAccResourceIndexConfigCustomModel(name string) string {
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
            treat_urls_and_pointers_as_images = true
            treat_urls_and_pointers_as_media = true
            model = "custom-model"
            model_properties = {
                url = "https://marqo-ecs-50-audio-test-dataset.s3.us-east-1.amazonaws.com/test-hf.zip"
                dimensions = 384
                type = "hf"
            }
            normalize_embeddings = true
            inference_type = "marqo.CPU.large"
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

func testAccResourceMinimalIndexConfig(name string) string {
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
				model = "open_clip/ViT-L-14/laion2b_s32b_b82k"
				inference_type = "marqo.CPU.large"
				number_of_inferences = 1
				number_of_replicas = 0
				number_of_shards = 1
				storage_class = "marqo.basic"
			}
		}
	`, name)
}

func testAccResourceMinimalIndexConfigUpdated(name string) string {
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
				model = "open_clip/ViT-L-14/laion2b_s32b_b82k"
				inference_type = "marqo.CPU.large"
				number_of_inferences = 2
				number_of_replicas = 0
				number_of_shards = 1
				storage_class = "marqo.basic"
			}
		}
	`, name)
}
