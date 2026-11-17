package scaleway

import (
	"context"
	"path"
	"slices"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

// Regions is the current known list of valid regions
func Regions() []scw.Region {
	return scw.AllRegions
}

func Zones() []scw.Zone {
	return scw.AllZones
}

// BuildRegionList :: return a list of matrix items, one per region
func BuildRegionList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// cache matrix
	cacheKey := "RegionListMatrix"
	if cachedData, ok := d.ConnectionCache.Get(ctx, cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	var allRegions []string

	// retrieve regions from connection config
	scalewayConfig := GetConfig(d.Connection)
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
		if err := d.ConnectionCache.Set(ctx, cacheKey, matrix); err != nil {
			plugin.Logger(ctx).Warn("scaleway.cache", "cache_key", cacheKey, "error", err)
		}

		return matrix
	}

	// Search for region configured using env, or use default region (i.e. fr-par)
	defaultScalewayRegion := GetDefaultScalewayRegion(ctx, d)
	matrix := []map[string]interface{}{
		{"region": defaultScalewayRegion},
	}

	// set cache
	if err := d.ConnectionCache.Set(ctx, cacheKey, matrix); err != nil {
		plugin.Logger(ctx).Warn("scaleway.cache", "cache_key", cacheKey, "error", err)
	}

	return matrix
}

// BuildZoneList :: return a list of matrix items, one per zone
func BuildZoneList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// cache matrix
	cacheKey := "ZoneListMatrix"
	if cachedData, ok := d.ConnectionCache.Get(ctx, cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	var allRegions []string

	// retrieve regions from connection config
	scalewayConfig := GetConfig(d.Connection)
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
	defaultRegion := GetDefaultScalewayRegion(ctx, d)
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
		if err := d.ConnectionCache.Set(ctx, cacheKey, matrix); err != nil {
			plugin.Logger(ctx).Warn("scaleway.cache", "cache_key", cacheKey, "error", err)
		}

		return matrix
	}

	// If nothing configured, fr-par-1 will be considered as default zone
	matrix := []map[string]interface{}{
		{"zone": "fr-par"},
	}

	// set cache
	if err := d.ConnectionCache.Set(ctx, cacheKey, matrix); err != nil {
		plugin.Logger(ctx).Warn("scaleway.cache", "cache_key", cacheKey, "error", err)
	}

	return matrix
}

// GetDefaultScalewayRegion returns the default region for Scaleway project
func GetDefaultScalewayRegion(ctx context.Context, d *plugin.QueryData) string {
	// have we already created and cached the service?
	serviceCacheKey := "GetDefaultScalewayRegion"
	if cachedData, ok := d.ConnectionCache.Get(ctx, serviceCacheKey); ok {
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
		if !slices.Contains(allRegions, region) {
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
	if !slices.Contains(allRegions, region) {
		region = "fr-par"
	}

	if err := d.ConnectionCache.Set(ctx, serviceCacheKey, region); err != nil {
		plugin.Logger(ctx).Warn("scaleway.cache", "cache_key", serviceCacheKey, "error", err)
	}
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
		if !slices.Contains(allRegions, region) {
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
