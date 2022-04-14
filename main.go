package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: vkcs.Provider})
}
