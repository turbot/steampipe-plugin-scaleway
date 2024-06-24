package scaleway

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

var getSessionConfigMemoized = plugin.HydrateFunc(getSession).Memoize(memoize.WithCacheKeyFunction(getSessionConfigCacheKey))

func getSessionConfig(ctx context.Context, d *plugin.QueryData) (*scw.Client, error) {
	h := &plugin.HydrateData{}
	cfg, err := getSessionConfigMemoized(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getSessionConfig", err)
		return nil, err
	}
	config := cfg.(*scw.Client)
	return config, nil
}

// Build a cache key for the call to getSessionConfigCacheKey.
func getSessionConfigCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := "getOrganizationId"
	return key, nil
}

func getSession(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getSessionConfigUncached(ctx, d)
}

// getSessionConfig :: returns Scaleway client to perform API requests
func getSessionConfigUncached(ctx context.Context, d *plugin.QueryData) (*scw.Client, error) {

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

	// Create client
	client, err := scw.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// getObjectSessionConfig :: returns S3 client to perform Object Storage API requests
func getObjectSessionConfig(ctx context.Context, d *plugin.QueryData, region string) (*s3.S3, error) {
	// Load clientOptions from cache
	sessionCacheKey := "scaleway.objectclient-" + region
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*s3.S3), nil
	}

	var accessKey, secretKey string

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

	// session default configuration
	sessionOptions := session.Options{
		Config: aws.Config{
			Region:      &region,
			Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
			Endpoint:    scw.StringPtr("https://s3." + region + ".scw.cloud"),
		},
	}

	s, err := session.NewSessionWithOptions(sessionOptions)
	if err != nil {
		return nil, err
	}
	client := s3.New(s)

	// save clientOptions in cache
	d.ConnectionManager.Cache.Set(sessionCacheKey, client)

	return client, nil
}
