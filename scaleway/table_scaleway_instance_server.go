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

func tableScalewayInstanceServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "scaleway_instance_server",
		Description:   "Scaleway Instance Server",
		GetMatrixItem: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate: listInstanceServers,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "commercial_type",
					Require: plugin.Optional,
				},
				{
					Name:    "zone",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getInstanceServer,
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
				Description: "An unique identifier of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "state",
				Description: "The current state of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State").Transform(transform.ToString),
			},
			{
				Name:        "creation_date",
				Description: "The time when the server was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "arch",
				Description: "The server state detail.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Arch").Transform(transform.ToString),
			},
			{
				Name:        "boot_type",
				Description: "The server boot type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BootType").Transform(transform.ToString),
			},
			{
				Name:        "commercial_type",
				Description: "The server commercial type (eg. GP1-M).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dynamic_ip_required",
				Description: "Indicates whether a dynamic IP is required, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DynamicIPRequired"),
			},
			{
				Name:        "enable_ipv6",
				Description: "Indicates whether IPv6 is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("EnableIPv6"),
			},
			{
				Name:        "host_name",
				Description: "The hostname of the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Hostname"),
			},
			{
				Name:        "protected",
				Description: "Indicates whether the server protection option is activated, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "private_ip",
				Description: "The private IP address of the server.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "modification_date",
				Description: "The time when the server was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "state_detail",
				Description: "The server state detail.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allowed_actions",
				Description: "A list of allowed actions on the server.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "image",
				Description: "Describes the information on the server image.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "public_ip",
				Description: "Specifies the information about the public IP.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PublicIP"),
			},
			{
				Name:        "location",
				Description: "Specifies the server location.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ipv6",
				Description: "The server IPv6 address.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("IPv6"),
			},
			{
				Name:        "boot_script",
				Description: "The server bootscript",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Bootscript"),
			},
			{
				Name:        "volumes",
				Description: "A list of server volumes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_group",
				Description: "Describes the server security group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "maintenances",
				Description: "The server planned maintenances.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "placement_group",
				Description: "The server placement group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_nics",
				Description: "The server private NICs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags",
				Description: "A list of tags associated with the server.",
				Type:        proto.ColumnType_JSON,
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

func listInstanceServers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_server.listInstanceServers", "zone_parsing_error", err)
		return nil, err
	}

	quals := d.KeyColumnQuals
	if quals["zone"] != nil && quals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_server.listInstanceServers", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	instanceApi := instance.NewAPI(client)

	req := &instance.ListServersRequest{
		Zone: parseZoneData,
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}
	if quals["commercial_type"] != nil {
		req.CommercialType = scw.StringPtr(quals["commercial_type"].GetStringValue())
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
	req.PerPage = scw.Uint32Ptr(uint32(maxResult))

	var count int

	for {
		resp, err := instanceApi.ListServers(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_instance_server.listInstanceServers", "query_error", err)
			return nil, err
		}

		for _, instance := range resp.Servers {
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
		} else if count == int(maxResult) {
			req.Page = scw.Int32Ptr(*req.Page + 1)
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getInstanceServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_server.listInstanceServers", "zone_parsing_error", err)
		return nil, err
	}

	if d.KeyColumnQuals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_server.getInstanceServer", "connection_error", err)
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

	data, err := instanceApi.GetServer(&instance.GetServerRequest{
		ServerID: id,
		Zone:     parseZoneData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_server.getInstanceServer", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data.Server, nil
}
