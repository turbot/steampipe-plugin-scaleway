package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayInstanceSecurityGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "scaleway_instance_security_group",
		Description:   "A security group is a set of firewall rules on a set of instances.",
		GetMatrixItemFunc: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate: listInstanceSecurityGroups,
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
			Hydrate:    getInstanceSecurityGroup,
			KeyColumns: plugin.AllColumns([]string{"id", "zone"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the security group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier of the security group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "description",
				Description: "The security group's description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_default",
				Description: "Indicates whether it is default security group for this project ID, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "organization_default",
				Description: "Indicates whether it is default security group for this organization ID, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "creation_date",
				Description: "The time when the security group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "enable_default_security",
				Description: "Indicates whether SMTP is blocked on IPv4 and IPv6, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "inbound_default_policy",
				Description: "Specifies the default inbound policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InboundDefaultPolicy").Transform(transform.ToString),
			},
			{
				Name:        "outbound_default_policy",
				Description: "Specifies the default outbound policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OutboundDefaultPolicy").Transform(transform.ToString),
			},
			{
				Name:        "modification_date",
				Description: "The time when the security group was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "stateful",
				Description: "Indicates whether the security group is stateful, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "servers",
				Description: "A list of tags associated with the security group.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "zone",
				Description: "Specifies the zone where the security group resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(transform.ToString),
			},
			{
				Name:        "project",
				Description: "The ID of the project where the security group resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the security group resides.",
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

func listInstanceSecurityGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_security_group.listInstanceSecurityGroups", "zone_parsing_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	if quals["zone"] != nil && quals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_security_group.listInstanceSecurityGroups", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	instanceApi := instance.NewAPI(client)

	req := &instance.ListSecurityGroupsRequest{
		Zone: parseZoneData,
		Page: scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}

	// Retrieve the list of security groups
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
		resp, err := instanceApi.ListSecurityGroups(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_instance_security_group.listInstanceSecurityGroups", "query_error", err)
			return nil, err
		}

		for _, securityGroup := range resp.SecurityGroups {
			d.StreamListItem(ctx, securityGroup)

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

func getInstanceSecurityGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zone := plugin.GetMatrixItem(ctx)["zone"].(string)

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_security_group.getInstanceSecurityGroup", "zone_parsing_error", err)
		return nil, err
	}

	if d.EqualsQuals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_security_group.getInstanceSecurityGroup", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	instanceApi := instance.NewAPI(client)

	id := d.EqualsQuals["id"].GetStringValue()
	sgZone := d.EqualsQuals["zone"].GetStringValue()

	// No inputs
	if id == "" && sgZone == "" {
		return nil, nil
	}

	data, err := instanceApi.GetSecurityGroup(&instance.GetSecurityGroupRequest{
		SecurityGroupID: id,
		Zone:            parseZoneData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance_security_group.getInstanceSecurityGroup", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data.SecurityGroup, nil
}
