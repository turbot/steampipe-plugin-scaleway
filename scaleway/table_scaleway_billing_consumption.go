package scaleway

import (
	"context"
	"fmt"

	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayBillingConsumption(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "scaleway_billing_consumption",
		Description: "Scaleway Billing Consumption",
		List: &plugin.ListConfig{
			Hydrate: listBillingConsumption,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "category_name",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "category_name",
				Description: "The CategoryName: name of consumption category.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "product_name",
				Description: "The product name. For example, VPC Public Gateway S, VPC Public Gateway M for the VPC product.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_name",
				Description: "Identifies the reference based on the category.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_id",
				Description: "The project ID of the consumption.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectID"),
			},
			{
				Name:        "sku",
				Description: "The unique identifier of the product.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "unit",
				Description: "The unit of consumed quantity.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "Monetary value of the consumption.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "billed_quantity",
				Description: "Consumed quantity.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

// // LIST FUNCTION
func listBillingConsumption(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_billing_consumption.listBillingConsumption", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Billing
	billingApi := billing.NewAPI(client)

	// Get organisationID from config to request IAM API
	organisationId := GetConfig(d.Connection).OrganizationID
	if organisationId == nil {
		err := fmt.Errorf("missing organization_id in scaleway.spc")
		plugin.Logger(ctx).Error("scaleway_billing_consumption.listBillingConsumption", "query_error", err)
		return nil, err
	}

	quals := d.EqualsQuals

	req := &billing.ListConsumptionsRequest{
		Page:           scw.Int32Ptr(1),
		OrganizationID: organisationId,
	}

	// Additional filters
	if quals["category_name"] != nil {
		req.CategoryName = scw.StringPtr(quals["category_name"].GetStringValue())
	}

	// Retrieve the list of consumptions
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
		resp, err := billingApi.ListConsumptions(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_billing_consumption.listBillingConsumption", "query_error", err)
			return nil, err
		}

		for _, consumption := range resp.Consumptions {
			d.StreamListItem(ctx, consumption)
			// Increase the resource count by 1
			count++

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.TotalCount == uint64(count) {
			break
		}
		req.Page = scw.Int32Ptr(*req.Page + 1)
	}

	return nil, nil
}
