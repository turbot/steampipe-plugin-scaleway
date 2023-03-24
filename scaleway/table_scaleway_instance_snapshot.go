package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayInstanceSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "scaleway_instance_snapshot",
		Description:   "Snapshots contain the data of a specific volume at a particular point in time.",
		GetMatrixItemFunc: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate: listInstanceSnapshots,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "zone",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getInstanceSnapshot,
			KeyColumns: plugin.AllColumns([]string{"id", "zone"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier of the snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "state",
				Description: "The current state of the snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State").Transform(transform.ToString),
			},
			{
				Name:        "size",
				Description: "The size of the snapshot (in bytes).",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Size").Transform(transform.ToInt),
			},
			{
				Name:        "volume_type",
				Description: "Specifies the snapshot volume type .",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VolumeType").Transform(transform.ToString),
			},
			{
				Name:        "creation_date",
				Description: "The time when the snapshot was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "modification_date",
				Description: "The time when the snapshot was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "snapshot_base_volume",
				Description: "Specifies the volume on which the snapshot is based on.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "zone",
				Description: "Specifies the zone where the snapshot resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(transform.ToString),
			},
			{
				Name:        "project",
				Description: "The ID of the project where the snapshot resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the snapshot resides.",
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

//// LIST FUNCTION

func listInstanceSnapshots(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := d.EqualsQualString("zone")

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_snapshot.listInstanceSnapshots", "zone_parsing_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	if quals["zone"] != nil && quals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_snapshot.listInstanceSnapshots", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	instanceApi := instance.NewAPI(client)

	req := &instance.ListSnapshotsRequest{
		Zone: parseZoneData,
		Page: scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}

	// Retrieve the list of snapshots
	maxResult := int64(100)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < maxResult {
			maxResult = *limit
		}
	}
	req.PerPage = scw.Uint32Ptr(uint32(maxResult))

	var count int

	for {
		resp, err := instanceApi.ListSnapshots(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_instance_snapshot.listInstanceSnapshots", "query_error", err)
			return nil, err
		}

		for _, snapshot := range resp.Snapshots {
			d.StreamListItem(ctx, snapshot)

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

//// HYDRATE FUNCTIONS

func getInstanceSnapshot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zone := d.EqualsQualString("zone")

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_snapshot.getInstanceSnapshot", "zone_parsing_error", err)
		return nil, err
	}

	if d.EqualsQuals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_snapshot.getInstanceSnapshot", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	instanceApi := instance.NewAPI(client)

	id := d.EqualsQuals["id"].GetStringValue()
	snapshotZone := d.EqualsQuals["zone"].GetStringValue()

	// No inputs
	if id == "" && snapshotZone == "" {
		return nil, nil
	}

	data, err := instanceApi.GetSnapshot(&instance.GetSnapshotRequest{
		SnapshotID: id,
		Zone:       parseZoneData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_snapshot.getInstanceSnapshot", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data.Snapshot, nil
}
