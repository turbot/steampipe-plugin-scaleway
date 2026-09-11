package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayVPC(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:              "scaleway_vpc",
		Description:       "A VPC is a regional network in which private networks are created.",
		GetMatrixItemFunc: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listVPCs,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "region",
					Require: plugin.Optional,
				},
				{
					Name:      "is_default",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getVPC,
			KeyColumns: plugin.AllColumns([]string{"id", "region"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the VPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier of the VPC.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "is_default",
				Description: "True if the VPC is the default one of the project.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IsDefault"),
			},
			{
				Name:        "private_network_count",
				Description: "The number of private networks within the VPC.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("PrivateNetworkCount"),
			},
			{
				Name:        "created_at",
				Description: "The time when the VPC was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the VPC was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tags",
				Description: "A list of tags associated with the VPC.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "region",
				Description: "Specifies the region where the VPC resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(transform.ToString),
			},
			{
				Name:        "project",
				Description: "The ID of the project where the VPC resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectID"),
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the VPC resides.",
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

func listVPCs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString("region")

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc.listVPCs", "region_parsing_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	if quals["region"] != nil && quals["region"].GetStringValue() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc.listVPCs", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway VPC product
	vpcApi := vpc.NewAPI(client)

	req := &vpc.ListVPCsRequest{
		Region: parseRegionData,
		Page:   scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}
	if quals["is_default"] != nil {
		req.IsDefault = scw.BoolPtr(quals["is_default"].GetBoolValue())
	}

	// Non-Equals Qual Map handling
	if d.Quals["is_default"] != nil {
		for _, q := range d.Quals["is_default"].Quals {
			value := q.Value.GetBoolValue()
			if q.Operator == "<>" {
				req.IsDefault = scw.BoolPtr(!value)
			}
		}
	}

	// Retrieve the list of VPCs
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
		resp, err := vpcApi.ListVPCs(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_vpc.listVPCs", "query_error", err)
			return nil, err
		}

		for _, v := range resp.Vpcs {
			d.StreamListItem(ctx, v)

			// Increase the resource count by 1
			count++

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Stop when the last page has been read, or when the API returns an empty
		// page, which would otherwise loop forever if resources went away mid-listing
		if len(resp.Vpcs) == 0 || uint32(count) >= resp.TotalCount {
			break
		}
		req.Page = scw.Int32Ptr(*req.Page + 1)

	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getVPC(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString("region")

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc.getVPC", "region_parsing_error", err)
		return nil, err
	}

	if d.EqualsQuals["region"].GetStringValue() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc.getVPC", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway VPC product
	vpcApi := vpc.NewAPI(client)

	id := d.EqualsQuals["id"].GetStringValue()

	// No inputs
	if id == "" || region == "" {
		return nil, nil
	}

	data, err := vpcApi.GetVPC(&vpc.GetVPCRequest{
		VpcID:  id,
		Region: parseRegionData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_vpc.getVPC", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
