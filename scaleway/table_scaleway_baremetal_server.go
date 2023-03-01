package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayBaremetalServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:              "scaleway_baremetal_server",
		Description:       "A Compute Instance bare metal is a physical server in Scaleway.",
		GetMatrixItemFunc: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate: listBaremetalServers,
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
			Hydrate:    getBaremetalServer,
			KeyColumns: plugin.AllColumns([]string{"id", "zone"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "A unique identifier of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "description",
				Description: "The cdescription of the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current state of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status").Transform(transform.ToString),
			},
			{
				Name:        "updated_at",
				Description: "The time when the server was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_at",
				Description: "The time when the server was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "boot_type",
				Description: "The server boot type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BootType").Transform(transform.ToString),
			},
			{
				Name:        "offer_name",
				Description: "The offer name of the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "offer_id",
				Description: "The offer ID of the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags",
				Description: "A list of tags associated with the server.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ips",
				Description: "A list of IPs associated with the server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("IPs"),
			},
			{
				Name:        "domain",
				Description: "The domain of the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "install",
				Description: "The install of the server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Install"),
			},
			{
				Name:        "ping_status",
				Description: "The status of the ping.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PingStatus").Transform(transform.ToString),
			},
			{
				Name:        "options",
				Description: "The options of the server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Options"),
			},
			{
				Name:        "rescue_server",
				Description: "The rescue boot of the server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RescueServer"),
			},

			// Scaleway standard columns
			{
				Name:        "zone",
				Description: "Specifies the zone where the server resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(transform.ToString),
			},
			{
				Name:        "project",
				Description: "The ID of the project where the server resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the server resides.",
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

func listBaremetalServers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_baremetal_server.listBaremetalServers", "zone_parsing_error", err)
		return nil, err
	}

	quals := d.KeyColumnQuals
	if quals["zone"] != nil && quals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_baremetal_server.listBaremetalServers", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	baremetalApi := baremetal.NewAPI(client)

	req := &baremetal.ListServersRequest{
		Zone: parseZoneData,
		Page: scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}
	// Retrieve the list of servers
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
		resp, err := baremetalApi.ListServers(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_baremetal_server.listBaremetalServers", "query_error", err)
			return nil, err
		}

		for _, baremetal := range resp.Servers {
			d.StreamListItem(ctx, baremetal)

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

func getBaremetalServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_baremetal_server.listBaremetalServers", "zone_parsing_error", err)
		return nil, err
	}

	if d.KeyColumnQuals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_baremetal_server.getBaremetalServer", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Baremetal product
	baremetalApi := baremetal.NewAPI(client)

	id := d.KeyColumnQuals["id"].GetStringValue()
	baremetalZone := d.KeyColumnQuals["zone"].GetStringValue()

	// No inputs
	if id == "" && baremetalZone == "" {
		return nil, nil
	}

	data, err := baremetalApi.GetServer(&baremetal.GetServerRequest{
		ServerID: id,
		Zone:     parseZoneData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_baremetal_server.getBaremetalServer", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
