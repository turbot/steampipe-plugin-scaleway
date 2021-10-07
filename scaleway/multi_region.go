package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/connection"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// Regions is the current known list of valid regions
func Regions() []scw.Region {
	return scw.AllRegions
}

func Zones() []scw.Zone {
	return scw.AllZones
}

var pluginQueryData *plugin.QueryData

func init() {
	pluginQueryData = &plugin.QueryData{
		ConnectionManager: connection.NewManager(),
	}
}

// BuildRegionList :: return a list of matrix items, one per region
func BuildRegionList(_ context.Context, connection *plugin.Connection) []map[string]interface{} {
	pluginQueryData.Connection = connection

	// cache matrix
	cacheKey := "RegionListMatrix"
	if cachedData, ok := pluginQueryData.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	allRegions := Regions()

	matrix := make([]map[string]interface{}, len(allRegions))
	for i, region := range allRegions {
		matrix[i] = map[string]interface{}{"region": region.String()}
	}

	// set cache
	pluginQueryData.ConnectionManager.Cache.Set(cacheKey, matrix)

	return matrix
}

// BuildZoneList :: return a list of matrix items, one per zone
func BuildZoneList(_ context.Context, connection *plugin.Connection) []map[string]interface{} {
	pluginQueryData.Connection = connection

	// cache matrix
	cacheKey := "ZoneListMatrix"
	if cachedData, ok := pluginQueryData.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	allZones := Zones()

	matrix := make([]map[string]interface{}, len(allZones))
	for i, zone := range allZones {
		matrix[i] = map[string]interface{}{"zone": zone.String()}
	}

	// set cache
	pluginQueryData.ConnectionManager.Cache.Set(cacheKey, matrix)

	return matrix
}
