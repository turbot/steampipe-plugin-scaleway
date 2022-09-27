package main

import (
	"github.com/turbot/steampipe-plugin-scaleway/scaleway"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: scaleway.Plugin})
}
