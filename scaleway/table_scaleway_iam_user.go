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

func tableScalewayIamUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "scaleway_iam_user",
		Description: "Users allow you to connect to scaleway console in your organization.",
		List: &plugin.ListConfig{
			Hydrate: listIamUsers,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIamUser,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{
				Name:        "user_id",
				Description: "ID of user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "email",
				Description: "The email of user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deletable",
				Description: "The deletion status of user. Owner user cannot be deleted.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "created_at",
				Description: "The time when the key was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_login_at",
				Description: "The last login date.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the key was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "type",
				Description: "The type of the user.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "two_factor_enabled",
				Description: "The 2FA enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "status",
				Description: "The status of invitation for the user.",
				Type:        proto.ColumnType_STRING,
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
				Transform:   transform.FromField("ID"),
			},
		},
	}
}

//// LIST FUNCTION

func listIamUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_instance.listIamUsers", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	iamApi := iam.NewAPI(client)

	// Get organisationID from config to request IAM API
	organisationId := GetConfig(d.Connection).OrganizationID

	req := &iam.ListUsersRequest{
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
		resp, err := iamApi.ListUsers(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_instance.listIamUsers", "query_error", err)
		}

		for _, key := range resp.Users {
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

func getIamUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_iamt_user.getIamUser", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Instance product
	iamApi := iam.NewAPI(client)

	userId := d.EqualsQualString("id")

	// No inputs
	if userId == "" {
		return nil, nil
	}

	data, err := iamApi.GetUser(&iam.GetUserRequest{
		UserID: userId,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_iam_user.getIamUser", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
