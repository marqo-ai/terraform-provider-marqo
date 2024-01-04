package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/integralist/terraform-provider-mock/mock"
	"github.com/marqo-ai/marqo_terraform_provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: *schema.Provider {
			return marqo.Provider()
		}
	})
}
