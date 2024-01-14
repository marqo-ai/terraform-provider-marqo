package provider

import (
    "net/http"
)

type ProviderConfiguration struct {
	APIClient *http.Client
	APIKey    string
}
