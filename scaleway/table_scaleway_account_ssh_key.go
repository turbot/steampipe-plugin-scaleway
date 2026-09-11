package scaleway

import (
	"context"

	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayAccountSSHKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "scaleway_account_ssh_key",
		Description: "SSH keys to access servers provisioned on Scaleway.",
		List: &plugin.ListConfig{
			Hydrate: listAccountSSHKeys,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getAccountSSHKey,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the SSH key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID of the SSH key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "created_at",
				Description: "The time when the key was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "fingerprint",
				Description: "Specifies the key fingerprint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_key",
				Description: "The SSH public key string",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The time when the key was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "disabled",
				Description: "True if the SSH key is disabled.",
				Type:        proto.ColumnType_BOOL,
			},

			// Scaleway standard columns
			{
				Name:        "project",
				Description: "The ID of the project where the server resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectID"),
			},
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
				Transform:   transform.FromField("Name"),
			},
		},
	}
}

//// LIST FUNCTION

func listAccountSSHKeys(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_account_ssh_key.listAccountSSHKeys", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway IAM product
	iamApi := iam.NewAPI(client)

	// Get organisationID from config to request IAM API
	organisationId := GetConfig(d.Connection).OrganizationID

	req := &iam.ListSSHKeysRequest{
		Page:           scw.Int32Ptr(1),
		OrganizationID: organisationId,
	}

	// Additional filter
	if d.EqualsQuals["name"] != nil {
		req.Name = scw.StringPtr(d.EqualsQuals["name"].GetStringValue())
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
		resp, err := iamApi.ListSSHKeys(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_account_ssh_key.listAccountSSHKeys", "query_error", err)
			return nil, err
		}

		for _, key := range resp.SSHKeys {
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

func getAccountSSHKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_account_ssh_key.getAccountSSHKey", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway IAM product
	iamApi := iam.NewAPI(client)

	id := d.EqualsQuals["id"].GetStringValue()

	// No inputs
	if id == "" {
		return nil, nil
	}

	data, err := iamApi.GetSSHKey(&iam.GetSSHKeyRequest{
		SSHKeyID: id,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_account_ssh_key.getAccountSSHKey", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
