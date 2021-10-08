package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayInstanceVolume(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "scaleway_instance_volume",
		Description:   "A volume is where you store your data inside your instance.",
		GetMatrixItem: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate: listInstanceVolumes,
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
			Hydrate:    getInstanceVolume,
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
				Name:        "state",
				Description: "The current state of the volume.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State").Transform(transform.ToString),
			},
			{
				Name:        "size",
				Description: "The size of the volume disk (in bytes).",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Size").Transform(transform.ToInt),
			},
			{
				Name:        "volume_type",
				Description: "The type of the volume.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VolumeType").Transform(transform.ToString),
			},
			{
				Name:        "creation_date",
				Description: "The time when the volume was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "export_uri",
				Description: "Specifies the volume NBD export URI.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExportURI"),
			},
			{
				Name:        "modification_date",
				Description: "The time when the volume was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "server",
				Description: "Specifies the server attached to the volume.",
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
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the volume resides.",
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

func listInstanceVolumes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_volume.listInstanceVolumes", "zone_parsing_error", err)
		return nil, err
	}

	quals := d.KeyColumnQuals
	if quals["zone"] != nil && quals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_volume.listInstanceVolumes", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	instanceApi := instance.NewAPI(client)

	req := &instance.ListVolumesRequest{
		Zone: parseZoneData,
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
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
	req.PerPage = scw.Uint32Ptr(uint32(maxResult))

	var count int

	for {
		resp, err := instanceApi.ListVolumes(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_instance_volume.listInstanceVolumes", "query_error", err)
			return nil, err
		}

		for _, volume := range resp.Volumes {
			d.StreamListItem(ctx, volume)

			// Increase the resource count by 1
			count++

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.TotalCount == uint32(count) {
			break
		} else if count == int(maxResult) {
			req.Page = scw.Int32Ptr(*req.Page + 1)
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getInstanceVolume(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_volume.getInstanceVolume", "zone_parsing_error", err)
		return nil, err
	}

	if d.KeyColumnQuals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_volume.getInstanceVolume", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	instanceApi := instance.NewAPI(client)

	id := d.KeyColumnQuals["id"].GetStringValue()
	instanceZone := d.KeyColumnQuals["zone"].GetStringValue()

	// No inputs
	if id == "" && instanceZone == "" {
		return nil, nil
	}

	data, err := instanceApi.GetVolume(&instance.GetVolumeRequest{
		VolumeID: id,
		Zone:     parseZoneData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_volume.getInstanceVolume", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data.Volume, nil
}
