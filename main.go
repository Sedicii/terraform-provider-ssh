package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sedicii/terraform-provider-ssh/ssh"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ssh.Provider,
	})
}
