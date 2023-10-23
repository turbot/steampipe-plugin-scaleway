package scaleway

import (
	"context"

	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayIamAPIKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "scaleway_iam_api_key",
		Description: "API keys allow you to securely connect to scaleway console in your organization.",
		List: &plugin.ListConfig{
			Hydrate: listIamAPIKeys,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "access_key",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIamAPIKey,
			KeyColumns: plugin.SingleColumn("access_key"),
		},
		Columns: []*plugin.Column{
			{
				Name:        "access_key",
				Description: "The access key of API key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "secret_key",
				Description: "The secret key of API Key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SecretKey").Transform(transform.ToString),
			},
			{
				Name:        "application_id",
				Description: "ID of application bearer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationID").Transform(transform.ToString),
			},
			{
				Name:        "created_at",
				Description: "Creation date and time of API key.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "user_id",
				Description: "ID of user bearer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserID").Transform(transform.ToString),
			},
			{
				Name:        "updated_at",
				Description: "Last update date and time of API key.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "expires_at",
				Description: "The expiration date and time of API key.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "default_project_id",
				Description: "The default project ID specified for this API key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefaultProjectID").Transform(transform.ToString),
			},
			{
				Name:        "editable",
				Description: "Whether or not the API key is editable.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "creation_ip",
				Description: "The IP Address of the device which created the API key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CreationIP").Transform(transform.ToString),
			},
			{
				Name:        "description",
				Description: "Description of API key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description").Transform(transform.ToString),
			},
		},
	}
}

//// LIST FUNCTION

func listIamAPIKeys(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_iam_api_key.listIamAPIKeys", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway IAM product
	iamApi := iam.NewAPI(client)

	// Get organisationID from config to request IAM API
	organisationId := GetConfig(d.Connection).OrganizationID

	req := &iam.ListAPIKeysRequest{
		Page:           scw.Int32Ptr(1),
		OrganizationID: organisationId,
	}

	// Retrieve the list of servers
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
		resp, err := iamApi.ListAPIKeys(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_iam_api_key.listIamAPIKeys", "query_error", err)
		}

		for _, key := range resp.APIKeys {
			d.StreamListItem(ctx, key)

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

func getIamAPIKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_iam_api_key.getIamAPIKey", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway IAM product
	iamApi := iam.NewAPI(client)

	accessKey := d.EqualsQuals["access_key"].GetStringValue()

	// No inputs
	if accessKey == "" {
		return nil, nil
	}

	data, err := iamApi.GetAPIKey(&iam.GetAPIKeyRequest{
		AccessKey: accessKey,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_iam_api_key.getIamAPIKey", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
