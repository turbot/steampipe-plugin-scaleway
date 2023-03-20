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

func tableScalewayInstanceImage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "scaleway_instance_image",
		Description:   "Images are backups of your instances.",
		GetMatrixItemFunc: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate: listInstanceImages,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:      "public",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
				{
					Name:    "zone",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getInstanceImage,
			KeyColumns: plugin.AllColumns([]string{"id", "zone"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier of the image.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "state",
				Description: "The current state of the image.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State").Transform(transform.ToString),
			},
			{
				Name:        "public",
				Description: "Indicates whether the image is public, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "creation_date",
				Description: "The time when the image was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "arch",
				Description: "The architecture the image is compatible with. Possible values are 'x86_64' and 'arm'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Arch").Transform(transform.ToString),
			},
			{
				Name:        "from_server",
				Description: "The ID of the server the image if based from.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_boot_script",
				Description: "Describes the default bootscript for this image.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DefaultBootscript"),
			},
			{
				Name:        "extra_volumes",
				Description: "Describes the extra volumes for this image.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DefaultBootscript"),
			},
			{
				Name:        "root_volume",
				Description: "Describes the root volume in this image.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "zone",
				Description: "Specifies the zone where the image resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(transform.ToString),
			},
			{
				Name:        "project",
				Description: "The ID of the project where the image resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the image resides.",
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

func listInstanceImages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_image.listInstanceImages", "zone_parsing_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	if quals["zone"] != nil && quals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_image.listInstanceImages", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	instanceApi := instance.NewAPI(client)

	req := &instance.ListImagesRequest{
		Zone: parseZoneData,
		Page: scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}

	if d.EqualsQuals["public"] != nil {
		req.Public = scw.BoolPtr(d.EqualsQuals["public"].GetBoolValue())
	}

	// Non-Equals Qual Map handling
	if d.Quals["public"] != nil {
		for _, q := range d.Quals["public"].Quals {
			value := q.Value.GetBoolValue()
			if q.Operator == "<>" {
				req.Public = scw.BoolPtr(false)
				if !value {
					req.Public = scw.BoolPtr(true)
				}
			}
		}
	}

	// Retrieve the list of images
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
		resp, err := instanceApi.ListImages(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_instance_image.listInstanceImages", "query_error", err)
			return nil, err
		}

		for _, image := range resp.Images {
			d.StreamListItem(ctx, image)

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

func getInstanceImage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_image.getInstanceImage", "zone_parsing_error", err)
		return nil, err
	}

	if d.EqualsQuals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_image.getInstanceImage", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	instanceApi := instance.NewAPI(client)

	id := d.EqualsQuals["id"].GetStringValue()
	imageZone := d.EqualsQuals["zone"].GetStringValue()

	// No inputs
	if id == "" && imageZone == "" {
		return nil, nil
	}

	data, err := instanceApi.GetImage(&instance.GetImageRequest{
		ImageID: id,
		Zone:    parseZoneData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_image.getInstanceImage", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data.Image, nil
}
