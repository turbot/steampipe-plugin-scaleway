package scaleway

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-scaleway"

// Plugin creates this (scaleway) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel().Transform(transform.NullIfZeroValue),
		DefaultGetConfig: &plugin.GetConfig{},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"scaleway_account_ssh_key":         tableScalewayAccountSSHKey(ctx),
			"scaleway_instance_ip":             tableScalewayInstanceIP(ctx),
			"scaleway_instance_security_group": tableScalewayInstanceSecurityGroup(ctx),
			"scaleway_instance_server":         tableScalewayInstanceServer(ctx),
			"scaleway_instance_volume":         tableScalewayInstanceVolume(ctx),
			"scaleway_rdb_database":            tableScalewayRDBDatabase(ctx),
			"scaleway_rdb_instance":            tableScalewayRDBInstance(ctx),
			"scaleway_vpc_private_network":     tableScalewayVPCPrivateNetwork(ctx),
		},
	}

	return p
}
