package scaleway

import (
	"context"
	"fmt"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func getSessionConfig(ctx context.Context, d *plugin.QueryData) (*scw.Client, error) {
	// Load clientOptions from cache
	sessionCacheKey := "scaleway.clientoption"
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*scw.Client), nil
	}

	var accessKey, secretKey string

	opts := []scw.ClientOption{}

	// Load credentials from environment variables
	loadEnv := scw.LoadEnvProfile()

	if loadEnv.AccessKey != nil && loadEnv.SecretKey != nil {
		accessKey = *loadEnv.AccessKey
		secretKey = *loadEnv.SecretKey
	}

	// Get scaleway config
	scalewayConfig := GetConfig(d.Connection)

	if scalewayConfig.AccessKey != nil && scalewayConfig.SecretKey == nil {
		return nil, fmt.Errorf("partial credentials found in connection config, missing: secret_key")
	} else if scalewayConfig.SecretKey != nil && scalewayConfig.AccessKey == nil {
		return nil, fmt.Errorf("partial credentials found in connection config, missing: access_key")
	} else if scalewayConfig.AccessKey != nil && scalewayConfig.SecretKey != nil {
		accessKey = *scalewayConfig.AccessKey
		secretKey = *scalewayConfig.SecretKey
	}

	// No creds
	if accessKey == "" && secretKey == "" {
		return nil, fmt.Errorf("both access_key and secret_key must be configured")
	}

	opts = append(opts, scw.WithAuth(*scalewayConfig.AccessKey, *scalewayConfig.SecretKey))

	// Create clientOption using env variables
	opts = append(opts, scw.WithProfile(loadEnv))

	// Create client
	client, err := scw.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	// save clientOptions in cache
	d.ConnectionManager.Cache.Set(sessionCacheKey, client)

	return client, nil
}
