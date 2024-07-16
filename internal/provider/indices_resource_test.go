package provider

import (
	"fmt"
	"os"
	"terraform-provider-marqo/marqo"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceIndex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceIndexConfig("example_index_6"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", "example_index_6"),
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
					testAccCheckIndexIsReady("example_index_6"),
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
				Config: testAccResourceIndexConfig("example_index_6_updated"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						fmt.Println("Starting Update and Read testing")
						return nil
					},
					resource.TestCheckResourceAttr("marqo_index.test", "index_name", "example_index_6_updated"),
					resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.CPU.large"),
					testAccCheckIndexIsReady("example_index_6_updated"),
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

func testAccCheckIndexIsReady(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Get environment variables
		host := os.Getenv("MARQO_HOST")
		apiKey := os.Getenv("MARQO_API_KEY")

		// Create a new Marqo client
		client, err := marqo.NewClient(&host, &apiKey)
		if err != nil {
			return fmt.Errorf("Error creating Marqo client: %s", err)
		}

		timeout := time.After(10 * time.Minute)
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		fmt.Printf("Waiting for index %s to be ready...\n", name)

		for {
			select {
			case <-timeout:
				return fmt.Errorf("Index %s did not become ready within the timeout period", name)
			case <-ticker.C:
				indices, err := client.ListIndices()
				if err != nil {
					fmt.Printf("Error listing indices: %s\n", err)
					continue
				}
				for _, index := range indices {
					if index.IndexName == name {
						fmt.Printf("Index %s status: %s\n", name, index.IndexStatus)
						if index.IndexStatus == "READY" {
							fmt.Printf("Index %s is now ready\n", name)
							return nil
						}
						break
					}
				}
				fmt.Printf("Index %s not ready yet, continuing to wait...\n", name)
			}
		}
	}
}
