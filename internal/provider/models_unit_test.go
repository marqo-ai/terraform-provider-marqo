package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestModelPropertiesModel_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		model    *ModelPropertiesModel
		expected bool
	}{
		{
			name:     "nil model",
			model:    nil,
			expected: true,
		},
		{
			name: "empty model",
			model: &ModelPropertiesModel{
				Name:             types.StringNull(),
				Dimensions:       types.StringValue(""),
				Type:             types.StringNull(),
				Tokens:           types.StringValue(""),
				Url:              types.StringNull(),
				TrustRemoteCode:  types.BoolValue(false),
				IsMarqtunedModel: types.BoolValue(false),
				ModelLocation:    nil,
			},
			expected: true,
		},
		{
			name: "model with values",
			model: &ModelPropertiesModel{
				Name:             types.StringValue("test"),
				Dimensions:       types.StringValue("384"),
				Type:             types.StringValue("hf"),
				Tokens:           types.StringValue("0"),
				Url:              types.StringValue("https://example.com"),
				TrustRemoteCode:  types.BoolValue(false),
				IsMarqtunedModel: types.BoolValue(false),
				ModelLocation:    nil,
			},
			expected: false,
		},
		{
			name: "model with only model location",
			model: &ModelPropertiesModel{
				Name:             types.StringNull(),
				Dimensions:       types.StringValue(""),
				Type:             types.StringNull(),
				Tokens:           types.StringValue(""),
				Url:              types.StringNull(),
				TrustRemoteCode:  types.BoolValue(false),
				IsMarqtunedModel: types.BoolValue(false),
				ModelLocation: &ModelLocationModel{
					AuthRequired: types.BoolValue(true),
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.model.IsEmpty()
			if got != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestModelLocationModel_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		location *ModelLocationModel
		expected bool
	}{
		{
			name:     "nil location",
			location: nil,
			expected: true,
		},
		{
			name: "empty location",
			location: &ModelLocationModel{
				AuthRequired: types.BoolValue(false),
				S3:           nil,
				Hf:           nil,
			},
			expected: true,
		},
		{
			name: "location with S3",
			location: &ModelLocationModel{
				AuthRequired: types.BoolValue(false),
				S3: &S3LocationModel{
					Bucket: types.StringValue("test-bucket"),
					Key:    types.StringValue("test-key"),
				},
				Hf: nil,
			},
			expected: false,
		},
		{
			name: "location with HF",
			location: &ModelLocationModel{
				AuthRequired: types.BoolValue(false),
				S3:           nil,
				Hf: &HfLocationModel{
					RepoId:   types.StringValue("test-repo"),
					Filename: types.StringValue("test-file"),
				},
			},
			expected: false,
		},
		{
			name: "location with auth required only",
			location: &ModelLocationModel{
				AuthRequired: types.BoolValue(true),
				S3:           nil,
				Hf:           nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.location.IsEmpty()
			if got != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}
