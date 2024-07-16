package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestProvider_Unit(t *testing.T) {
	p := New("test")()
	if p == nil {
		t.Fatal("Provider should not be nil")
	}
}

func TestProvider_Metadata_Unit(t *testing.T) {
	p := New("test")()
	var resp provider.MetadataResponse
	p.Metadata(context.Background(), provider.MetadataRequest{}, &resp)
	if resp.TypeName != "marqo" {
		t.Fatalf("Expected TypeName to be 'marqo', got '%s'", resp.TypeName)
	}
	if resp.Version != "test" {
		t.Fatalf("Expected Version to be 'test', got '%s'", resp.Version)
	}
}

func TestProvider_Schema_Unit(t *testing.T) {
	p := New("test")()
	var resp provider.SchemaResponse
	p.Schema(context.Background(), provider.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema attributes should not be nil")
	}

	if _, ok := resp.Schema.Attributes["host"]; !ok {
		t.Fatal("Schema should have 'host' attribute")
	}

	if _, ok := resp.Schema.Attributes["api_key"]; !ok {
		t.Fatal("Schema should have 'api_key' attribute")
	}
}
