package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-cidr/cidr"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cidr.Provider})
}
