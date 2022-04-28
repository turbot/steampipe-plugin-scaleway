package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayInstanceIP(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "scaleway_instance_ip",
		Description:   "A flexible IP address is an IP address which you hold independently of any server.",
		GetMatrixItem: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate: listInstanceIPs,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "zone",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getInstanceIP,
			KeyColumns: plugin.AllColumns([]string{"id", "zone"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "An unique identifier of the IP address.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "address",
				Description: "Specifies the IP address.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("Address").Transform(transform.ToString),
			},
			{
				Name:        "reverse",
				Description: "The reverse dns attached to this IP address.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "server",
				Description: "Specifies the server attached to the IP address.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags",
				Description: "A list of tags attached with the IP address.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "zone",
				Description: "Specifies the zone where the IP address resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(transform.ToString),
			},
			{
				Name:        "project",
				Description: "The ID of the project where the IP address resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the IP address resides.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Address").Transform(transform.ToString),
			},
		},
	}
}

//// LIST FUNCTION

func listInstanceIPs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_ip.listInstanceIPs", "zone_parsing_error", err)
		return nil, err
	}

	quals := d.KeyColumnQuals
	if quals["zone"] != nil && quals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_ip.listInstanceIPs", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	instanceApi := instance.NewAPI(client)

	req := &instance.ListIPsRequest{
		Zone: parseZoneData,
		Page: scw.Int32Ptr(1),
	}

	// Retrieve the list of IPs
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
		resp, err := instanceApi.ListIPs(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_instance_ip.listInstanceIPs", "query_error", err)
			return nil, err
		}

		for _, ip := range resp.IPs {
			d.StreamListItem(ctx, ip)

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

func getInstanceIP(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_ip.getInstanceIP", "zone_parsing_error", err)
		return nil, err
	}

	if d.KeyColumnQuals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_ip.getInstanceIP", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	instanceApi := instance.NewAPI(client)

	id := d.KeyColumnQuals["id"].GetStringValue()
	ipZone := d.KeyColumnQuals["zone"].GetStringValue()

	// No inputs
	if id == "" && ipZone == "" {
		return nil, nil
	}

	data, err := instanceApi.GetIP(&instance.GetIPRequest{
		IP:   id,
		Zone: parseZoneData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_ip.getInstanceIP", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data.IP, nil
}
