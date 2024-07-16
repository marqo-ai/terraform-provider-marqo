package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestProviderUnit(t *testing.T) {
	p := New("test")()

	t.Run("schema", func(t *testing.T) {
		schemaResp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

		if schemaResp.Schema.Attributes == nil {
			t.Fatal("Schema attributes should not be nil")
		}

		if _, ok := schemaResp.Schema.Attributes["host"]; !ok {
			t.Fatal("Schema should have 'host' attribute")
		}

		if _, ok := schemaResp.Schema.Attributes["api_key"]; !ok {
			t.Fatal("Schema should have 'api_key' attribute")
		}
	})

	t.Run("resources", func(t *testing.T) {
		resourcesFunc := p.Resources(context.Background())
		if len(resourcesFunc) == 0 {
			t.Fatal("Provider should have at least one resource")
		}

		for _, rf := range resourcesFunc {
			r := rf()
			if _, ok := r.(resource.Resource); !ok {
				t.Fatalf("Resource does not implement resource.Resource")
			}
		}
	})

	t.Run("data_sources", func(t *testing.T) {
		dataSourcesFunc := p.DataSources(context.Background())
		if len(dataSourcesFunc) == 0 {
			t.Fatal("Provider should have at least one data source")
		}

		for _, dsf := range dataSourcesFunc {
			ds := dsf()
			if _, ok := ds.(datasource.DataSource); !ok {
				t.Fatalf("Data source does not implement datasource.DataSource")
			}
		}
	})
}
