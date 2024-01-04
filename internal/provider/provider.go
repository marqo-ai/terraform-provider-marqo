package provider

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API Key for accessing your API",
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"index": NewIndexResource(), // Use NewIndexResource here
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := data.Get("api_key").(string)

	var diags diag.Diagnostics

	client := &http.Client{
		// Optionally, set up HTTP client parameters (Timeout, Transport, etc.)
	}

	return &ProviderConfiguration{
		APIClient: client,
		APIKey:    apiKey,
	}, diags
}
