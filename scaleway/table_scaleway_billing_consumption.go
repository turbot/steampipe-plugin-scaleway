package scaleway

import (
	"context"
	"fmt"

	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2alpha1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayBillingConsumption(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "scaleway_billing_consumption",
		Description: "A Scaleway Billing Consumption.",
		List: &plugin.ListConfig{
			Hydrate: listBillingConsumption,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "category",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "category",
				Description: "The ID of the project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "operation_path",
				Description: "The unique identifier of the product.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_id",
				Description: "The project ID of the consumption.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectID"),
			},
			{
				Name:        "description",
				Description: "Description of the consumption.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "Monetary value of the consumption.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

// // GET FUNCTION
func listBillingConsumption(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_project.listBillingConsumption", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Billing
	billingApi := billing.NewAPI(client)

	// Get organisationID from config to request IAM API
	organisationId := GetConfig(d.Connection).OrganizationID
	if organisationId == nil {
		err := fmt.Errorf("missing organization_id in scaleway.spc")
		plugin.Logger(ctx).Error("scaleway_project.listBillingConsumption", "query_error", err)
		return nil, err
	}

	req := &billing.GetConsumptionRequest{
		OrganizationID: *organisationId,
	}

	resp, err := billingApi.GetConsumption(req)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_project.listBillingConsumption", "query_error", err)
		return nil, err
	}

	for _, consumption := range resp.Consumptions {
		d.StreamListItem(ctx, consumption)
	}
	return nil, nil
}
