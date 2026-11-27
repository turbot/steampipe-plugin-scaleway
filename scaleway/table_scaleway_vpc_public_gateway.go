package scaleway

import (
	"context"
	"slices"

	vpcgw "github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayVPCPublicGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:              "scaleway_vpc_public_gateway",
		Description:       "A Public Gateway sits at the border of Private Networks and provides them with internet access, through NAT and an SSH bastion.",
		GetMatrixItemFunc: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate: listVPCPublicGateways,
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
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getVPCPublicGateway,
			KeyColumns: plugin.AllColumns([]string{"id", "zone"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the Public Gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier of the Public Gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "status",
				Description: "The current status of the Public Gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status").Transform(transform.ToString),
			},
			{
				Name:        "type",
				Description: "The commercial offer type of the Public Gateway, such as VPC-GW-S.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bandwidth",
				Description: "The bandwidth available to the Public Gateway, in bits per second.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "ipv4",
				Description: "The public IPv4 address of the Public Gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("IPv4"),
			},
			{
				Name:        "gateway_networks",
				Description: "The Private Networks the Public Gateway is attached to, with their masquerade and default route settings.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "version",
				Description: "The version of the Public Gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "can_upgrade_to",
				Description: "The newer version the Public Gateway can be upgraded to, if any.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bastion_enabled",
				Description: "True if the SSH bastion is enabled on the Public Gateway.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "bastion_port",
				Description: "The port of the SSH bastion.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "bastion_allowed_ips",
				Description: "The IP ranges allowed to connect to the SSH bastion.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BastionAllowedIPs"),
			},
			{
				Name:        "smtp_enabled",
				Description: "True if SMTP traffic is allowed to pass through the Public Gateway.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("SMTPEnabled"),
			},
			{
				Name:        "is_legacy",
				Description: "True if the Public Gateway is a legacy one, not yet migrated to IPAM.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "created_at",
				Description: "The time when the Public Gateway was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the Public Gateway was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tags",
				Description: "A list of tags associated with the Public Gateway.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "zone",
				Description: "Specifies the zone where the Public Gateway resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(transform.ToString),
			},
			{
				Name:        "project",
				Description: "The ID of the project where the Public Gateway resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectID"),
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the Public Gateway resides.",
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

func listVPCPublicGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := d.EqualsQualString("zone")

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc_public_gateway.listVPCPublicGateways", "zone_parsing_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	if quals["zone"] != nil && quals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc_public_gateway.listVPCPublicGateways", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Public Gateways product
	vpcgwApi := vpcgw.NewAPI(client)

	// Public Gateways are not available in every zone of a region
	if !slices.Contains(vpcgwApi.Zones(), parseZoneData) {
		return nil, nil
	}

	req := &vpcgw.ListGatewaysRequest{
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

	// Retrieve the list of Public Gateways
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
		resp, err := vpcgwApi.ListGateways(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_vpc_public_gateway.listVPCPublicGateways", "query_error", err)
			return nil, err
		}

		for _, gateway := range resp.Gateways {
			d.StreamListItem(ctx, gateway)

			// Increase the resource count by 1
			count++

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Stop when the last page has been read, or when the API returns an empty
		// page, which would otherwise loop forever if resources went away mid-listing
		if len(resp.Gateways) == 0 || uint64(count) >= resp.TotalCount {
			break
		}
		req.Page = scw.Int32Ptr(*req.Page + 1)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getVPCPublicGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := d.EqualsQualString("zone")

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc_public_gateway.getVPCPublicGateway", "zone_parsing_error", err)
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
		plugin.Logger(ctx).Error("scaleway_vpc_public_gateway.getVPCPublicGateway", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Public Gateways product
	vpcgwApi := vpcgw.NewAPI(client)

	if !slices.Contains(vpcgwApi.Zones(), parseZoneData) {
		return nil, nil
	}

	data, err := vpcgwApi.GetGateway(&vpcgw.GetGatewayRequest{
		GatewayID: id,
		Zone:      parseZoneData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc_public_gateway.getVPCPublicGateway", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
