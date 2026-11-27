package scaleway

import (
	"context"
	"slices"

	redis "github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayRedisCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:              "scaleway_redis_cluster",
		Description:       "A Managed Database for Redis cluster is a fully managed Redis deployment of one or more nodes.",
		GetMatrixItemFunc: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate: listRedisClusters,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "zone",
					Require: plugin.Optional,
				},
				{
					Name:    "project",
					Require: plugin.Optional,
				},
				{
					Name:    "version",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getRedisCluster,
			KeyColumns: plugin.AllColumns([]string{"id", "zone"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "status",
				Description: "The current status of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status").Transform(transform.ToString),
			},
			{
				Name:        "version",
				Description: "The Redis engine version of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_type",
				Description: "The node type of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_size",
				Description: "The number of nodes in the cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "tls_enabled",
				Description: "True if TLS is enabled on the cluster.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("TLSEnabled"),
			},
			{
				Name:        "endpoints",
				Description: "The endpoints of the cluster, each on a Private Network or on the public network, with its IPs and port.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "acl_rules",
				Description: "The ACL rules restricting which IP ranges can reach the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ACLRules"),
			},
			{
				Name:        "cluster_settings",
				Description: "The advanced settings of the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "user_name",
				Description: "The name of the user created along with the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "upgradable_versions",
				Description: "The Redis engine versions the cluster can be upgraded to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "created_at",
				Description: "The time when the cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the cluster was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tags",
				Description: "A list of tags associated with the cluster.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "zone",
				Description: "Specifies the zone where the cluster resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(transform.ToString),
			},
			{
				Name:        "project",
				Description: "The ID of the project where the cluster resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectID"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		},
	}
}

//// LIST FUNCTION

func listRedisClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := d.EqualsQualString("zone")

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_redis_cluster.listRedisClusters", "zone_parsing_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	if quals["zone"] != nil && quals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_redis_cluster.listRedisClusters", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Managed Database for Redis product
	redisApi := redis.NewAPI(client)

	// Managed Database for Redis is not available in every zone of a region
	if !slices.Contains(redisApi.Zones(), parseZoneData) {
		return nil, nil
	}

	req := &redis.ListClustersRequest{
		Zone: parseZoneData,
		Page: scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}
	if quals["project"] != nil {
		req.ProjectID = scw.StringPtr(quals["project"].GetStringValue())
	}
	if quals["version"] != nil {
		req.Version = scw.StringPtr(quals["version"].GetStringValue())
	}

	// Retrieve the list of clusters
	maxResult := int64(100)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < maxResult {
			maxResult = *limit
		}
	}
	req.PageSize = scw.Uint32Ptr(uint32(maxResult))

	var count int

	for {
		resp, err := redisApi.ListClusters(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_redis_cluster.listRedisClusters", "query_error", err)
			return nil, err
		}

		for _, cluster := range resp.Clusters {
			d.StreamListItem(ctx, cluster)

			// Increase the resource count by 1
			count++

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Stop when the last page has been read, or when the API returns an empty
		// page, which would otherwise loop forever if resources went away mid-listing
		if len(resp.Clusters) == 0 || uint32(count) >= resp.TotalCount {
			break
		}
		req.Page = scw.Int32Ptr(*req.Page + 1)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRedisCluster(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := d.EqualsQualString("zone")

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_redis_cluster.getRedisCluster", "zone_parsing_error", err)
		return nil, err
	}

	if d.EqualsQuals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	id := d.EqualsQuals["id"].GetStringValue()

	// No inputs
	if id == "" || zone == "" {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_redis_cluster.getRedisCluster", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Managed Database for Redis product
	redisApi := redis.NewAPI(client)

	if !slices.Contains(redisApi.Zones(), parseZoneData) {
		return nil, nil
	}

	data, err := redisApi.GetCluster(&redis.GetClusterRequest{
		ClusterID: id,
		Zone:      parseZoneData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_redis_cluster.getRedisCluster", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
