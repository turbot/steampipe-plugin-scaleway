package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayRDBDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "scaleway_rdb_database",
		Description:   "A RDB database is a logical databases on your instance.",
		GetMatrixItemFunc: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate:       listRDBDatabases,
			ParentHydrate: listRDBInstances,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "region",
					Require: plugin.Optional,
				},
				{
					Name:    "owner",
					Require: plugin.Optional,
				},
				{
					Name:      "managed",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_id",
				Description: "An unique identifier of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceID"),
			},
			{
				Name:        "managed",
				Description: "Indicates whether database is managed, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "size",
				Description: "The size of the database.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Size").Transform(transform.ToString),
			},
			{
				Name:        "owner",
				Description: "Specifies the name of the owner of the database.",
				Type:        proto.ColumnType_STRING,
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
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the instance resides.",
				Type:        proto.ColumnType_STRING,
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

type databaseInfo = struct {
	rdb.Database
	InstanceID   string
	Region       scw.Region
	Project      string
	Organization string
}

//// LIST FUNCTION

func listRDBDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_rdb_instance.listRDBDatabases", "region_parsing_error", err)
		return nil, err
	}

	// Get Instance details
	instanceData := h.Item.(*rdb.Instance)

	quals := d.EqualsQuals
	if instanceData.Region.String() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_rdb_instance.listRDBDatabases", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway RDB product
	rdbApi := rdb.NewAPI(client)

	req := &rdb.ListDatabasesRequest{
		Region:     parseRegionData,
		InstanceID: instanceData.ID,
		Page:       scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}
	if quals["owner"] != nil {
		req.Owner = scw.StringPtr(quals["owner"].GetStringValue())
	}

	if d.EqualsQuals["managed"] != nil {
		req.Managed = scw.BoolPtr(d.EqualsQuals["managed"].GetBoolValue())
	}

	// Non-Equals Qual Map handling
	if d.Quals["managed"] != nil {
		for _, q := range d.Quals["managed"].Quals {
			value := q.Value.GetBoolValue()
			if q.Operator == "<>" {
				req.Managed = scw.BoolPtr(false)
				if !value {
					req.Managed = scw.BoolPtr(true)
				}
			}
		}
	}

	// Retrieve the list of databases
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
		resp, err := rdbApi.ListDatabases(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_rdb_instance.listRDBDatabases", "query_error", err)
			return nil, err
		}

		for _, database := range resp.Databases {
			d.StreamListItem(ctx, databaseInfo{*database, instanceData.ID, instanceData.Region, instanceData.ProjectID, instanceData.OrganizationID})

			// Increase the resource count by 1
			count++

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
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
