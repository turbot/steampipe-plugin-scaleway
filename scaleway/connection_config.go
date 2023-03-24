package scaleway

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type scalewayConfig struct {
	AccessKey      *string  `cty:"access_key"`
	SecretKey      *string  `cty:"secret_key"`
	OrganizationID *string  `cty:"organization_id"`
	Regions        []string `cty:"regions"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"access_key": {
		Type: schema.TypeString,
	},
	"secret_key": {
		Type: schema.TypeString,
	},
	"organization_id": {
		Type: schema.TypeString,
	},
	"regions": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
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
