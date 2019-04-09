package main

import (
	"github.com/gig-tech/terraform-provider-ovc/ovc"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return ovc.Provider()
		},
	})
}
