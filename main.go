package main

import (
	"github.com/civo/terraform-provider-civo/civo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: civo.Provider,
	})
}
