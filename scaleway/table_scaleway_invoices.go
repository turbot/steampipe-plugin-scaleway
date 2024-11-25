package scaleway

import (
	"context"
	"fmt"
	"math"
	"time"

	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2beta1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableScalewayInvoices(ctx context.Context) *plugin.Table {
	plugin.Logger(ctx).Debug("Initializing Scaleway invoices table")
	return &plugin.Table{
		Name:        "scaleway_invoices",
		Description: "invoices in your Scaleway account.",
		List: &plugin.ListConfig{
			Hydrate: listScalewayInvoices,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization_id", Require: plugin.Optional},
			},
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
		},
	}
}

func extractAmount(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Debug("extractAmount", "field", d.ColumnName, "raw_value", d.Value)

	if d.Value == nil {
		plugin.Logger(ctx).Warn("extractAmount: nil value", "field", d.ColumnName)
		return nil, nil
	}

	money, ok := d.Value.(*scw.Money)
	if !ok || money == nil {
		plugin.Logger(ctx).Warn("extractAmount: unexpected type or nil", "type", fmt.Sprintf("%T", d.Value))
		return nil, nil
	}

	amount := float64(money.Units) + float64(money.Nanos)/1e9
	if d.ColumnName == "total_discount_amount" {
		return math.Abs(amount), nil
	}
	return amount, nil
}

func extractCurrency(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Debug("extractCurrency", "field", d.ColumnName, "raw_value", d.Value)

	if d.Value == nil {
		plugin.Logger(ctx).Warn("extractCurrency: nil value", "field", d.ColumnName)
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

type invoicesItem struct {
	ID               string              `json:"id"`
	OrganizationID   string              `json:"organization_id"`
	OrganizationName string              `json:"organization_name"`
	StartDate        *time.Time          `json:"start_date"`
	StopDate         *time.Time          `json:"stop_date"`
	IssuedDate       *time.Time          `json:"issued_date"`
	DueDate          *time.Time          `json:"due_date"`
	TotalUntaxed     *scw.Money          `json:"total_untaxed"`
	TotalTaxed       *scw.Money          `json:"total_taxed"`
	TotalDiscount    *scw.Money          `json:"total_discount"`
	TotalUndiscount  *scw.Money          `json:"total_undiscount"`
	TotalTax         *scw.Money          `json:"total_tax"`
	Type             billing.InvoiceType `json:"type"`
	Number           int32               `json:"number"`
	State            string              `json:"state"`
	BillingPeriod    *time.Time          `json:"billing_period"`
	SellerName       string              `json:"seller_name"`
}

func listScalewayInvoices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client configuration
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_invoices.listScalewayInvoices", "connection_error", err)
		return nil, err
	}

	// Check if the client is properly configured
	if client == nil {
		return nil, fmt.Errorf("scaleway client is not properly configured")
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

	// Check if organization_id is specified in the query
	if d.EqualsQualString("organization_id") != "" {
		organizationID = d.EqualsQualString("organization_id")
	}

	// Log a warning if organization_id is not provided by either config or query
	if organizationID == "" {
		plugin.Logger(ctx).Warn("scaleway_invoices.listScalewayInvoices", "warning", "No organization_id provided in config or query")
	}

	// Set the organization_id in the request if it's available
	if organizationID != "" {
		req.OrganizationID = &organizationID
	}

	// Make the API request to list invoices
	resp, err := billingAPI.ListInvoices(req)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_invoices.listScalewayInvoices", "api_error", err)
		return nil, err
	}

	for _, invoices := range resp.Invoices {
		plugin.Logger(ctx).Debug("raw invoices data", "invoices", fmt.Sprintf("%+v", invoices))

		item := invoicesItem{
			ID:               invoices.ID,
			OrganizationID:   invoices.OrganizationID,
			OrganizationName: invoices.OrganizationName,
			StartDate:        invoices.StartDate,
			StopDate:         invoices.StopDate,
			IssuedDate:       invoices.IssuedDate,
			DueDate:          invoices.DueDate,
			TotalUntaxed:     invoices.TotalUntaxed,
			TotalTaxed:       invoices.TotalTaxed,
			TotalDiscount:    invoices.TotalDiscount,
			TotalUndiscount:  invoices.TotalUndiscount,
			TotalTax:         invoices.TotalTax,
			Type:             invoices.Type,
			Number:           invoices.Number,
			State:            invoices.State,
			BillingPeriod:    invoices.BillingPeriod,
			SellerName:       invoices.SellerName,
		}
		d.StreamListItem(ctx, item)
	}

	return nil, nil
}
