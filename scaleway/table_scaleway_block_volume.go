package scaleway

import (
	"context"
	"slices"

	block "github.com/scaleway/scaleway-sdk-go/api/block/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayBlockVolume(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:              "scaleway_block_volume",
		Description:       "A Block Storage volume is a network-attached volume that can be attached to an Instance.",
		GetMatrixItemFunc: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate: listBlockVolumes,
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
					Name:    "type",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getBlockVolume,
			KeyColumns: plugin.AllColumns([]string{"id", "zone"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the volume.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier of the volume.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "type",
				Description: "The type of the volume, such as sbs_5k or sbs_15k.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "The size of the volume (in bytes).",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Size").Transform(transform.ToInt),
			},
			{
				Name:        "status",
				Description: "The current status of the volume.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status").Transform(transform.ToString),
			},
			{
				Name:        "class",
				Description: "The storage class of the volume.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Specs.Class").Transform(transform.ToString),
			},
			{
				Name:        "perf_iops",
				Description: "The maximum IO/s of the volume, if it has a performance limit.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Specs.PerfIops"),
			},
			{
				Name:        "kms_key_id",
				Description: "The ID of the Key Manager key used to encrypt the volume, if it is encrypted with a customer-managed key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KmsKeyID"),
			},
			{
				Name:        "attached_resources",
				Description: "The resources the volume is attached to, such as an Instance server. Named after the API's references field, which is a reserved word in SQL.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("References"),
			},
			{
				Name:        "parent_snapshot_id",
				Description: "The ID of the snapshot the volume was created from, if any.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ParentSnapshotID"),
			},
			{
				Name:        "created_at",
				Description: "The time when the volume was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the volume was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_detached_at",
				Description: "The time when the volume was last detached from a resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tags",
				Description: "A list of tags associated with the volume.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "zone",
				Description: "Specifies the zone where the volume resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(transform.ToString),
			},
			{
				Name:        "project",
				Description: "The ID of the project where the volume resides.",
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

func listBlockVolumes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := d.EqualsQualString("zone")

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_block_volume.listBlockVolumes", "zone_parsing_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	if quals["zone"] != nil && quals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_block_volume.listBlockVolumes", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Block Storage product
	blockApi := block.NewAPI(client)

	// Block Storage is not available in every zone of a region
	if !slices.Contains(blockApi.Zones(), parseZoneData) {
		return nil, nil
	}

	req := &block.ListVolumesRequest{
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
	if quals["type"] != nil {
		req.VolumeType = scw.StringPtr(quals["type"].GetStringValue())
	}

	// Retrieve the list of volumes
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
		resp, err := blockApi.ListVolumes(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_block_volume.listBlockVolumes", "query_error", err)
			return nil, err
		}

		for _, volume := range resp.Volumes {
			d.StreamListItem(ctx, volume)

			// Increase the resource count by 1
			count++

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Stop when the last page has been read, or when the API returns an empty
		// page, which would otherwise loop forever if resources went away mid-listing
		if len(resp.Volumes) == 0 || uint64(count) >= resp.TotalCount {
			break
		}
		req.Page = scw.Int32Ptr(*req.Page + 1)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getBlockVolume(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := d.EqualsQualString("zone")

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_block_volume.getBlockVolume", "zone_parsing_error", err)
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
		plugin.Logger(ctx).Error("scaleway_block_volume.getBlockVolume", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Block Storage product
	blockApi := block.NewAPI(client)

	if !slices.Contains(blockApi.Zones(), parseZoneData) {
		return nil, nil
	}

	data, err := blockApi.GetVolume(&block.GetVolumeRequest{
		VolumeID: id,
		Zone:     parseZoneData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_block_volume.getBlockVolume", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
