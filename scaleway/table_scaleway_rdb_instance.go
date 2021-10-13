package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayRDBInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "scaleway_rdb_instance",
		Description:   "A Database Instance is composed of one or more Nodes, depending of the is_ha_cluster setting.",
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listRDBInstances,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "region",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getRDBInstance,
			KeyColumns: plugin.AllColumns([]string{"id", "region"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "status",
				Description: "The current state of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status").Transform(transform.ToString),
			},
			{
				Name:        "engine",
				Description: "The database engine of the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The time when the instance was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "is_has_cluster",
				Description: "Indicates whether High-Availability is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "node_type",
				Description: "The node type of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_schedule",
				Description: "Describes the backup schedule of the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "endpoint",
				Description: "Describes the endpoint of the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "init_settings",
				Description: "A list of engine settings to be set at database initialization.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "read_replicas",
				Description: "A list of read replicas of the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "settings",
				Description: "A list of advanced settings of the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "volume",
				Description: "Describes the volume attached with the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags",
				Description: "A list of tags associated with the instance.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "region",
				Description: "Specifies the region where the instance resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(transform.ToString),
			},
			{
				Name:        "project",
				Description: "The ID of the project where the instance resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectID"),
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the instance resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OrganizationID"),
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

func listRDBInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_rdb_instance.listRDBInstances", "region_parsing_error", err)
		return nil, err
	}

	quals := d.KeyColumnQuals
	if quals["region"] != nil && quals["region"].GetStringValue() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_rdb_instance.listRDBInstances", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway RDB product
	rdbApi := rdb.NewAPI(client)

	req := &rdb.ListInstancesRequest{
		Region: parseRegionData,
		Page:   scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}

	// Retrieve the list of instances
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
		resp, err := rdbApi.ListInstances(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_rdb_instance.listRDBInstances", "query_error", err)
			return nil, err
		}

		for _, instance := range resp.Instances {
			d.StreamListItem(ctx, instance)

			// Increase the resource count by 1
			count++

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.TotalCount == uint32(count) {
			break
		}
		req.Page = scw.Int32Ptr(*req.Page + 1)

	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRDBInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_rdb_instance.getRDBInstance", "region_parsing_error", err)
		return nil, err
	}

	if d.KeyColumnQuals["region"].GetStringValue() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_rdb_instance.getRDBInstance", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway RDB product
	rdbApi := rdb.NewAPI(client)

	id := d.KeyColumnQuals["id"].GetStringValue()
	instanceZone := d.KeyColumnQuals["zone"].GetStringValue()

	// No inputs
	if id == "" && instanceZone == "" {
		return nil, nil
	}

	data, err := rdbApi.GetInstance(&rdb.GetInstanceRequest{
		InstanceID: id,
		Region:     parseRegionData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_rdb_instance.getRDBInstance", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
