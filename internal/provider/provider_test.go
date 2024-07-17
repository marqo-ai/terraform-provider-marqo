package provider

import (
	"fmt"
	"math/rand"
	"os"
	"terraform-provider-marqo/marqo"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"marqo": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("MARQO_HOST"); v == "" {
		t.Fatal("MARQO_HOST must be set for acceptance tests")
	}
	if v := os.Getenv("MARQO_API_KEY"); v == "" {
		t.Fatal("MARQO_API_KEY must be set for acceptance tests")
	}
}

func TestAccProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testProviderConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.marqo_read_indices.test", "id", "test_id_1"),
				),
			},
		},
	})
}

const testProviderConfig = `
provider "marqo" {}

data "marqo_read_indices" "test" {
	id = "test_id_1"
}
`

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func testAccCheckIndexExistsAndDelete(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Get environment variables
		host := os.Getenv("MARQO_HOST")
		apiKey := os.Getenv("MARQO_API_KEY")

		// Create a new Marqo client
		client, err := marqo.NewClient(&host, &apiKey)
		if err != nil {
			return fmt.Errorf("Error creating Marqo client: %s", err)
		}

		fmt.Printf("Checking if index %s exists...\n", name)

		indices, err := client.ListIndices()
		if err != nil {
			return fmt.Errorf("Error listing indices: %s", err)
		}

		var indexExists bool
		for _, index := range indices {
			if index.IndexName == name {
				indexExists = true
				break
			}
		}

		if !indexExists {
			fmt.Printf("Index %s does not exist. Proceeding with creation.\n", name)
			return nil
		}

		fmt.Printf("Index %s exists. Deleting...\n", name)
		err = client.DeleteIndex(name)
		if err != nil {
			return fmt.Errorf("Error deleting index %s: %s", name, err)
		}

		// Wait for the index to be deleted
		timeout := time.After(5 * time.Minute)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				return fmt.Errorf("Timed out waiting for index %s to be deleted", name)
			case <-ticker.C:
				indices, err := client.ListIndices()
				if err != nil {
					fmt.Printf("Error listing indices: %s\n", err)
					continue
				}
				indexExists := false
				for _, index := range indices {
					if index.IndexName == name {
						indexExists = true
						break
					}
				}
				if !indexExists {
					fmt.Printf("Index %s has been successfully deleted.\n", name)
					return nil
				}
				fmt.Printf("Waiting for index %s to be deleted...\n", name)
			}
		}
	}
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

		timeout := time.After(15 * time.Minute)
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		fmt.Printf("Waiting for index %s to be ready...\n", name)

		start := time.Now()
		for {
			select {
			case <-timeout:
				return fmt.Errorf("Index %s did not become ready within the 10-minute timeout period", name)
			case <-ticker.C:
				indices, err := client.ListIndices()
				if err != nil {
					fmt.Printf("Error listing indices: %s\n", err)
					continue
				}
				for _, index := range indices {
					if index.IndexName == name {
						fmt.Printf("Index %s status: %s (elapsed: %v)\n", name, index.IndexStatus, time.Since(start))
						if index.IndexStatus == "READY" {
							fmt.Printf("Index %s is now ready (total time: %v)\n", name, time.Since(start))
							return nil
						}
						break
					}
				}
				fmt.Printf("Index %s not ready yet, continuing to wait... (elapsed: %v)\n", name, time.Since(start))
			}
		}
	}
}
