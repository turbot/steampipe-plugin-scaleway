package scaleway

import (
	"context"
	"slices"

	lb "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayLB(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:              "scaleway_lb",
		Description:       "A Load Balancer distributes incoming traffic across a set of backend servers.",
		GetMatrixItemFunc: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate: listLBs,
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
			Hydrate:    getLB,
			KeyColumns: plugin.AllColumns([]string{"id", "zone"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the Load Balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier of the Load Balancer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "description",
				Description: "The description of the Load Balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current status of the Load Balancer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status").Transform(transform.ToString),
			},
			{
				Name:        "type",
				Description: "The commercial offer type of the Load Balancer, such as LB-S.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip",
				Description: "The flexible IP addresses attached to the Load Balancer.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("IP"),
			},
			{
				Name:        "instances",
				Description: "The underlying instances of the Load Balancer, and their status.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "frontend_count",
				Description: "The number of frontends of the Load Balancer.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "backend_count",
				Description: "The number of backends of the Load Balancer.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "private_network_count",
				Description: "The number of Private Networks attached to the Load Balancer.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "route_count",
				Description: "The number of routes configured on the Load Balancer.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "ssl_compatibility_level",
				Description: "The SSL compatibility level of the Load Balancer, which sets the TLS versions and ciphers its frontends accept.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SslCompatibilityLevel").Transform(transform.ToString),
			},
			{
				Name:        "subscriber",
				Description: "The subscriber the Load Balancer sends its alerts to, if any.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "created_at",
				Description: "The time when the Load Balancer was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the Load Balancer was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tags",
				Description: "A list of tags associated with the Load Balancer.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "zone",
				Description: "Specifies the zone where the Load Balancer resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(transform.ToString),
			},
			{
				Name:        "project",
				Description: "The ID of the project where the Load Balancer resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectID"),
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the Load Balancer resides.",
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

func listLBs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := d.EqualsQualString("zone")

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_lb.listLBs", "zone_parsing_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	if quals["zone"] != nil && quals["zone"].GetStringValue() != zone {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_lb.listLBs", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Load Balancer product
	lbApi := lb.NewZonedAPI(client)

	// Load Balancers are not available in every zone of a region
	if !slices.Contains(lbApi.Zones(), parseZoneData) {
		return nil, nil
	}

	req := &lb.ZonedAPIListLBsRequest{
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

	// Retrieve the list of Load Balancers
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
		resp, err := lbApi.ListLBs(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_lb.listLBs", "query_error", err)
			return nil, err
		}

		for _, item := range resp.LBs {
			d.StreamListItem(ctx, item)

			// Increase the resource count by 1
			count++

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Stop when the last page has been read, or when the API returns an empty
		// page, which would otherwise loop forever if resources went away mid-listing
		if len(resp.LBs) == 0 || uint32(count) >= resp.TotalCount {
			break
		}
		req.Page = scw.Int32Ptr(*req.Page + 1)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getLB(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	zone := d.EqualsQualString("zone")

	parseZoneData, err := scw.ParseZone(zone)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_lb.getLB", "zone_parsing_error", err)
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
		plugin.Logger(ctx).Error("scaleway_lb.getLB", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Load Balancer product
	lbApi := lb.NewZonedAPI(client)

	if !slices.Contains(lbApi.Zones(), parseZoneData) {
		return nil, nil
	}

	data, err := lbApi.GetLB(&lb.ZonedAPIGetLBRequest{
		LBID: id,
		Zone: parseZoneData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_lb.getLB", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
