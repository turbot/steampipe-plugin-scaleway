package scaleway

import (
	"context"
	"path"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v4/connection"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
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

	var allRegions []string

	// retrieve regions from connection config
	scalewayConfig := GetConfig(connection)
	if scalewayConfig.Regions != nil {
		regions := Regions()
		for _, pattern := range scalewayConfig.Regions {
			for _, validRegion := range regions {
				if ok, _ := path.Match(pattern, validRegion.String()); ok {
					allRegions = append(allRegions, validRegion.String())
				}
			}
		}
	}

	// Build regions matrix using config regions
	if len(allRegions) > 0 {
		uniqueRegions := unique(allRegions)

		if len(getInvalidRegions(uniqueRegions)) > 0 {
			panic("\n\nConnection config have invalid regions: " + strings.Join(getInvalidRegions(uniqueRegions), ","))
		}

		// validate regions list
		matrix := make([]map[string]interface{}, len(uniqueRegions))
		for i, region := range uniqueRegions {
			matrix[i] = map[string]interface{}{"region": region}
		}

		// set cache
		pluginQueryData.ConnectionManager.Cache.Set(cacheKey, matrix)

		return matrix
	}

	// Search for region configured using env, or use default region (i.e. fr-par)
	defaultScalewayRegion := GetDefaultScalewayRegion(pluginQueryData)
	matrix := []map[string]interface{}{
		{"region": defaultScalewayRegion},
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

	var allRegions []string

	// retrieve regions from connection config
	scalewayConfig := GetConfig(connection)
	if scalewayConfig.Regions != nil {
		regions := Regions()
		for _, pattern := range scalewayConfig.Regions {
			for _, validRegion := range regions {
				if ok, _ := path.Match(pattern, validRegion.String()); ok {
					allRegions = append(allRegions, validRegion.String())
				}
			}
		}
	}

	// Get default region
	defaultRegion := GetDefaultScalewayRegion(pluginQueryData)
	allRegions = append(allRegions, defaultRegion)

	// Build regions matrix using config regions
	if len(allRegions) > 0 {
		uniqueRegions := unique(allRegions)

		if len(getInvalidRegions(uniqueRegions)) > 0 {
			panic("\n\nConnection config have invalid regions: " + strings.Join(getInvalidRegions(uniqueRegions), ","))
		}

		var allZones []scw.Zone
		for _, region := range uniqueRegions {
			zones := parseRegion(region).GetZones()
			allZones = append(allZones, zones...)
		}

		// validate regions list
		matrix := make([]map[string]interface{}, len(allZones))
		for i, zone := range allZones {
			matrix[i] = map[string]interface{}{"zone": zone.String()}
		}

		// set cache
		pluginQueryData.ConnectionManager.Cache.Set(cacheKey, matrix)

		return matrix
	}

	// If nothing configured, fr-par-1 will be considered as default zone
	matrix := []map[string]interface{}{
		{"zone": "fr-par"},
	}

	// set cache
	pluginQueryData.ConnectionManager.Cache.Set(cacheKey, matrix)

	return matrix
}

// GetDefaultScalewayRegion returns the default region for Scaleway project
func GetDefaultScalewayRegion(d *plugin.QueryData) string {
	// have we already created and cached the service?
	serviceCacheKey := "GetDefaultScalewayRegion"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(string)
	}
	var allRegions []string
	for _, i := range Regions() {
		allRegions = append(allRegions, i.String())
	}

	// get config info
	scalewayConfig := GetConfig(d.Connection)

	var regions []string
	var region string

	if scalewayConfig.Regions != nil {
		regions = scalewayConfig.Regions
		region = regions[0]
	} else {
		// Load credentials from environment variables
		loadEnv := scw.LoadEnvProfile()

		if loadEnv.DefaultRegion != nil {
			region = *loadEnv.DefaultRegion
		}

		if region != "" {
			regions = []string{region}
		}

		// https://registry.terraform.io/providers/scaleway/scaleway/latest/docs#arguments-reference
		if !helpers.StringSliceContains(allRegions, region) {
			regions = []string{"fr-par"}
		}
	}

	validPatterns := []string{}
	invalidPatterns := []string{}
	for _, namePattern := range regions {
		validRegions := []string{}
		for _, validRegion := range allRegions {
			if ok, _ := path.Match(namePattern, validRegion); ok {
				validRegions = append(validRegions, validRegion)
			}
		}
		if len(validRegions) == 0 {
			invalidPatterns = append(invalidPatterns, namePattern)
		} else {
			validPatterns = append(validPatterns, namePattern)
		}
	}

	if len(validPatterns) == 0 {
		panic("\nconnection config have invalid \"regions\": " + strings.Join(invalidPatterns, ", ") + ". Edit your connection configuration file and then restart Steampipe")
	}

	// https://registry.terraform.io/providers/scaleway/scaleway/latest/docs#arguments-reference
	if !helpers.StringSliceContains(allRegions, region) {
		region = "fr-par"
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, region)
	return region
}

// Return invalid regions from a region list
func getInvalidRegions(regions []string) []string {
	var allRegions []string
	for _, i := range Regions() {
		allRegions = append(allRegions, i.String())
	}

	invalidRegions := []string{}
	for _, region := range regions {
		if !helpers.StringSliceContains(allRegions, region) {
			invalidRegions = append(invalidRegions, region)
		}
	}
	return invalidRegions
}

// Returns a list of unique items
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func parseRegion(region string) scw.Region {
	parsedRegion, err := scw.ParseRegion(region)
	if err != nil {
		return ""
	}
	return parsedRegion
}
