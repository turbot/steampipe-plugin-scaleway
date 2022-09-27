package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/vpc/v1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayVPCPrivateNetwork(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "scaleway_vpc_private_network",
		Description:   "A VPC private network allows interconnecting your instances in an isolated and private network.",
		GetMatrixItemFunc: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate: listVPCPrivateNetworks,
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
			Hydrate:    getVPCPrivateNetwork,
			KeyColumns: plugin.AllColumns([]string{"id", "zone"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the private network.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier of the private network.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "created_at",
				Description: "The time when the private network was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the private network was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tags",
				Description: "A list of tags associated with the private network.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "zone",
				Description: "Specifies the zone where the private network resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(transform.ToString),
			},
			{
				Name:        "project",
				Description: "The ID of the project where the private network resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectID"),
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the private network resides.",
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

func listVPCPrivateNetworks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc_private_network.listVPCPrivateNetworks", "zone_parsing_error", err)
		return nil, err
	}

	quals := d.KeyColumnQuals
	if quals["zone"] != nil && quals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc_private_network.listVPCPrivateNetworks", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway VPC product
	vpcApi := vpc.NewAPI(client)

	req := &vpc.ListPrivateNetworksRequest{
		Zone: parseZoneData,
		Page: scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}

	// Retrieve the list of private networks
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
		resp, err := vpcApi.ListPrivateNetworks(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_vpc_private_network.listVPCPrivateNetworks", "query_error", err)
			return nil, err
		}

		for _, network := range resp.PrivateNetworks {
			d.StreamListItem(ctx, network)

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

func getVPCPrivateNetwork(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc_private_network.getVPCPrivateNetwork", "zone_parsing_error", err)
		return nil, err
	}

	if d.KeyColumnQuals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc_private_network.listVPCPrivateNetworks", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway VPC product
	vpcApi := vpc.NewAPI(client)

	id := d.KeyColumnQuals["id"].GetStringValue()
	vpcZone := d.KeyColumnQuals["zone"].GetStringValue()

	// No inputs
	if id == "" && vpcZone == "" {
		return nil, nil
	}

	data, err := vpcApi.GetPrivateNetwork(&vpc.GetPrivateNetworkRequest{
		PrivateNetworkID: id,
		Zone:             parseZoneData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc_private_network.getVPCPrivateNetwork", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
