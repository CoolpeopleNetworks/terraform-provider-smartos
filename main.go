package main

import (
	"github.com/coolpeoplenetworks/terraform-provider-smartos/smartos"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return smartos.Provider()
		},
	})
}