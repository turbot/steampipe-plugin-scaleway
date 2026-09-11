package main

import (
	"github.com/turbot/steampipe-plugin-scaleway/v2/scaleway"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: scaleway.Plugin})
}
