package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/spectrumwebco/terraform-provider-hashicorp-ovh/internal/provider"
)

var (
	version = "dev"
	commit  = ""
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debugMode,
		ProviderAddr: "registry.terraform.io/spectrumwebco/hashicorp-ovh",
		ProviderFunc: provider.New(version),
	}

	plugin.Serve(opts)
}
