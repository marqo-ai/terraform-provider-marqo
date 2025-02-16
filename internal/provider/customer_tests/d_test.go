package customer_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// To run these tests manually:
// go test -v ./internal/provider/customer_tests -run TestSybill

func TestD(t *testing.T) {
	// Skip in normal test runs
	if testing.Short() {
		t.Skip("Skipping customer-specific test in short mode")
	}
	t.Parallel()

	// Test for deal index
	t.Run("deal_index", func(t *testing.T) {
		indexName := fmt.Sprintf("donotdelete_deal_%s", randomString(4))
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create initial index
				{
					Config: testDDealConfig(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "index_name", indexName),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.inference_type", "marqo.CPU"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "1"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_replicas", "0"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_shards", "1"),
						resource.TestCheckResourceAttr("marqo_index.test", "settings.storage_class", "marqo.balanced"),
					),
				},
				// Modify index
				{
					Config: testDDealConfigModified(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "settings.number_of_inferences", "2"),
					),
				},
				// Delete and recreate testing
				{
					Config: testAccEmptyConfig(),
					Check: resource.ComposeTestCheckFunc(
						func(s *terraform.State) error {
							fmt.Println("Starting Deal Index delete testing")
							return nil
						},
					),
				},
				// Recreate index
				{
					Config: testDDealConfig(indexName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("marqo_index.test", "index_name", indexName),
					),
				},
			},
		})
	})
}

func testDDealConfig(name string) string {
	return fmt.Sprintf(`
	resource "marqo_index" "test" {
		index_name = "%s"
		settings = {
			type = "unstructured"
			model = "open_clip/ViT-L-14/laion2b_s32b_b82k"
			inference_type = "marqo.CPU.large"
			number_of_inferences = 1
			number_of_replicas = 0
			number_of_shards = 1
			storage_class = "marqo.balanced"
		}
		timeouts = {
			create = "45m"
			update = "45m"
			delete = "20m"
		}
	}
	`, name)
}

func testDDealConfigModified(name string) string {
	return fmt.Sprintf(`
	resource "marqo_index" "test" {
		index_name = "%s"
		settings = {
			type = "unstructured"
			model = "open_clip/ViT-L-14/laion2b_s32b_b82k"
			inference_type = "marqo.CPU.large"
			number_of_inferences = 2
			number_of_replicas = 0
			number_of_shards = 1
			storage_class = "marqo.balanced"
		}
		timeouts = {
			create = "45m"
			update = "45m"
			delete = "20m"
		}
	}
	`, name)
}
