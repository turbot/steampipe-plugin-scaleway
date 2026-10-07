package scaleway

import (
	"context"
	"fmt"

	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayIamOrganization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "scaleway_iam_organization",
		Description: "An Organization is the top level of the Scaleway hierarchy, owning every Project, resource, user and invoice.",
		// The Organization API exposes no list endpoint: an API key belongs to exactly one
		// Organization, so the list is that single Organization, read from the connection config.
		List: &plugin.ListConfig{
			Hydrate: listIamOrganizations,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIamOrganization,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "ID of the Organization.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Description: "Name of the Organization.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "alias",
				Description: "Alias of the Organization, as used in its console login URL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "login_password_enabled",
				Description: "Whether login with a password is enabled for the Organization.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("LoginPasswordEnabled"),
			},
			{
				Name:        "login_magic_code_enabled",
				Description: "Whether login with an authentication code is enabled for the Organization.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("LoginMagicCodeEnabled"),
			},
			{
				Name:        "login_oauth2_enabled",
				Description: "Whether login through OAuth2 is enabled for the Organization.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("LoginOauth2Enabled"),
			},
			{
				Name:        "login_saml_enabled",
				Description: "Whether login through SAML is enabled for the Organization.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("LoginSamlEnabled"),
			},

			// Scaleway standard columns
			{
				Name:        "organization",
				Description: "The ID of the Organization.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
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

func listIamOrganizations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_iam_organization.listIamOrganizations", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway IAM product
	iamApi := iam.NewAPI(client)

	// Get organizationID from config to request IAM API. getSessionConfig sets no default
	// Organization on the client, so an unset organization_id would reach the API as an empty
	// path segment rather than defaulting.
	organizationId := GetConfig(d.Connection).OrganizationID
	if organizationId == nil {
		err := fmt.Errorf("missing organization_id in scaleway.spc")
		plugin.Logger(ctx).Error("scaleway_iam_organization.listIamOrganizations", "query_error", err)
		return nil, err
	}

	organization, err := iamApi.GetOrganization(&iam.GetOrganizationRequest{
		OrganizationID: *organizationId,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_iam_organization.listIamOrganizations", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	d.StreamListItem(ctx, organization)

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamOrganization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_iam_organization.getIamOrganization", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway IAM product
	iamApi := iam.NewAPI(client)

	organizationId := d.EqualsQualString("id")

	// No inputs
	if organizationId == "" {
		return nil, nil
	}

	// An API key belongs to a single Organization, so skip the API call for any other Organization
	configOrganizationId := GetConfig(d.Connection).OrganizationID
	if configOrganizationId != nil && *configOrganizationId != organizationId {
		return nil, nil
	}

	data, err := iamApi.GetOrganization(&iam.GetOrganizationRequest{
		OrganizationID: organizationId,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_iam_organization.getIamOrganization", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
