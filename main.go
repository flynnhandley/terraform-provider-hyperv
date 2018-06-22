package main

import (
	"github.com/flynnhandley/hashicorp-plugins/hyperv"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: hyperv.Provider})
}
