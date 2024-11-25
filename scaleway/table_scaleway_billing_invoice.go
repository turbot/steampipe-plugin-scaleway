package scaleway

import (
	"context"
	"fmt"
	"math"

	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2beta1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableScalewayBillingInvoice(ctx context.Context) *plugin.Table {
	plugin.Logger(ctx).Debug("Initializing Scaleway invoices table")
	return &plugin.Table{
		Name:        "scaleway_billing_invoice",
		Description: "Scaleway Billing Invoice",
		List: &plugin.ListConfig{
			Hydrate: listScalewayInvoices,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization_id", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
				{Name: "billing_period", Require: plugin.Optional, Operators: []string{">=", "<="}},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getScalewayInvoice,
			KeyColumns: plugin.SingleColumn("id"),
			// When an incorrect Invoice ID is provided, no "not found" error is returned.
			// Instead, a timeout error is encountered, as shown below.
			// Error: rpc error: code = DeadlineExceeded desc = scaleway-sdk-go: error executing request: Get "https://api.scaleway.com/account/v3/projects/5575456ffhgfh": dial tcp: lookup api.scaleway.com: i/o timeout.
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier of the invoices.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "organization_id",
				Type:        proto.ColumnType_STRING,
				Description: "The organization ID associated with the invoices.",
				Transform:   transform.FromField("OrganizationID"),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of the invoices.",
			},
			{
				Name:        "state",
				Type:        proto.ColumnType_STRING,
				Description: "The current state of the invoices.",
			},
			{
				Name:        "number",
				Type:        proto.ColumnType_INT,
				Description: "The invoices number.",
			},
			{
				Name:        "seller_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the seller.",
			},
			{
				Name:        "organization_name",
				Type:        proto.ColumnType_STRING,
				Description: "The organization name associated with the invoices.",
			},
			{
				Name:        "start_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The start date of the billing period.",
			},
			{
				Name:        "stop_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The end date of the billing period.",
			},
			{
				Name:        "billing_period",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The billing period for the invoices.",
			},
			{
				Name:        "issued_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date when the invoices was issued.",
			},
			{
				Name:        "due_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The due date for the invoices payment.",
			},
			{
				Name:        "total_untaxed_amount",
				Type:        proto.ColumnType_DOUBLE,
				Description: "The total untaxed amount of the invoices.",
				Transform:   transform.FromField("TotalUntaxed").Transform(extractAmount),
			},
			{
				Name:        "total_taxed_amount",
				Type:        proto.ColumnType_DOUBLE,
				Description: "The total taxed amount of the invoices.",
				Transform:   transform.FromField("TotalTaxed").Transform(extractAmount),
			},
			{
				Name:        "total_discount_amount",
				Type:        proto.ColumnType_DOUBLE,
				Description: "The total discount amount of the invoice (always positive).",
				Transform:   transform.FromField("TotalDiscount").Transform(extractAmount),
			},
			{
				Name:        "total_undiscount_amount",
				Type:        proto.ColumnType_DOUBLE,
				Description: "The total undiscounted amount of the invoices.",
				Transform:   transform.FromField("TotalUndiscount").Transform(extractAmount),
			},
			{
				Name:        "currency",
				Type:        proto.ColumnType_STRING,
				Description: "The currency used for all monetary values in the invoices.",
				Transform:   transform.FromField("TotalTaxed").Transform(extractCurrency),
			},

			// Scaleway standard columns
			{
				Name:        "organization",
				Description: "The ID of the organization where the server resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OrganizationID"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Number"),
			},
		},
	}
}

//// LIST FUNCTION

func listScalewayInvoices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client configuration
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_billing_invoice.listScalewayInvoices", "connection_error", err)
		return nil, err
	}

	billingAPI := billing.NewAPI(client)

	// Prepare the request
	req := &billing.ListInvoicesRequest{}

	// Get the organization_id from the config
	scalewayConfig := GetConfig(d.Connection)
	var organizationID string
	if scalewayConfig.OrganizationID != nil {
		organizationID = *scalewayConfig.OrganizationID
	}

	// Check if organization_id is specified in the query parameter
	if d.EqualsQualString("organization_id") != "" {
		organizationID = d.EqualsQualString("organization_id")
	}

	// Set the organization_id in the request if it's available
	if organizationID != "" {
		req.OrganizationID = &organizationID
	}

	if d.EqualsQualString("type") != "" {
		req.InvoiceType = billing.InvoiceType(d.EqualsQualString("type"))
	}

	quals := d.Quals

	if quals["billing_period"] != nil {
		for _, q := range quals["billing_period"].Quals {
			billingPeriod := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">=":
				req.BillingPeriodStartAfter = &billingPeriod
			case "<=":
				req.BillingPeriodStartBefore = &billingPeriod
			}
		}
	}

	var count int

	for {
		// Make the API request to list invoices
		resp, err := billingAPI.ListInvoices(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_billing_invoice.listScalewayInvoices", "api_error", err)
			return nil, err
		}

		for _, invoice := range resp.Invoices {

			d.StreamListItem(ctx, invoice)

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

//// GYDRATE FUNCTION

func getScalewayInvoice(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client configuration
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_billing_invoice.getScalewayInvoice", "connection_error", err)
		return nil, err
	}

	// Empty check
	if d.EqualsQualString("id") == "" {
		return nil, nil
	}

	billingAPI := billing.NewAPI(client)

	// Prepare the request
	req := &billing.GetInvoiceRequest{
		InvoiceID: d.EqualsQualString("id"),
	}

	resp, err := billingAPI.GetInvoice(req)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_billing_invoice.getScalewayInvoice", "api_error", err)
		return nil, err
	}

	if resp != nil {
		return resp, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func extractAmount(_ context.Context, d *transform.TransformData) (interface{}, error) {

	if d.Value == nil {
		return nil, nil
	}

	money, ok := d.Value.(*scw.Money)
	if !ok || money == nil {
		return nil, nil
	}

	amount := float64(money.Units) + float64(money.Nanos)/1e9
	if d.ColumnName == "total_discount_amount" {
		return math.Abs(amount), nil
	}
	return amount, nil
}

func extractCurrency(ctx context.Context, d *transform.TransformData) (interface{}, error) {

	if d.Value == nil {
		return nil, nil
	}

	switch v := d.Value.(type) {
	case *scw.Money:
		if v == nil {
			return nil, nil
		}
		return v.CurrencyCode, nil
	default:
		plugin.Logger(ctx).Warn("extractCurrency: unexpected type", "type", fmt.Sprintf("%T", d.Value))
		return nil, nil
	}
}
