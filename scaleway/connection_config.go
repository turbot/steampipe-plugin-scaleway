package scaleway

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type scalewayConfig struct {
	AccessKey *string `cty:"access_key"`
	SecretKey *string `cty:"secret_key"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"access_key": {
		Type: schema.TypeString,
	},
	"secret_key": {
		Type: schema.TypeString,
	},
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
