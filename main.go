package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/moofish32/terraform-provider-cidr/cidr"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cidr.Provider})
}
