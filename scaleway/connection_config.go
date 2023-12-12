package scaleway

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type scalewayConfig struct {
	AccessKey      *string  `hcl:"access_key"`
	SecretKey      *string  `hcl:"secret_key"`
	OrganizationID *string  `hcl:"organization_id"`
	Regions        []string `hcl:"regions,optional"`
}

func ConfigInstance() interface{} {
	return &scalewayConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) scalewayConfig {
	if connection == nil || connection.Config == nil {
		return scalewayConfig{}
	}
	config, _ := connection.Config.(scalewayConfig)
	return config
}
