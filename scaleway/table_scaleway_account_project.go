package scaleway

import (
	"context"
	"fmt"

	"github.com/scaleway/scaleway-sdk-go/api/account/v3"
	"github.com/scaleway/scaleway-sdk-go/scw"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayAccountProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "scaleway_account_project",
		Description: "A Scaleway Account Project.",
		List: &plugin.ListConfig{
			Hydrate: listAccountProjects,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getAccountProject,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Description: "Name of the project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The time when the project was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the project was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "Description of the project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization_id",
				Description: "Organization ID of the project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OrganizationID"),
			},
		},
	}
}

//// LIST FUNCTION

func listAccountProjects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_project.listProjects", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Account product
	accountApi := account.NewProjectAPI(client)

	// Get organisationID from config to request IAM API
	organisationId := GetConfig(d.Connection).OrganizationID
	if organisationId == nil {
		err := fmt.Errorf("missing organization_id in scaleway.spc")
		plugin.Logger(ctx).Error("scaleway_project.listProjects", "query_error", err)
		return nil, err
	}

	quals := d.EqualsQuals

	req := &account.ProjectAPIListProjectsRequest{
		OrganizationID: *organisationId,
		Page:           scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
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
		resp, err := accountApi.ListProjects(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_project.listProjects", "query_error", err)
			return nil, err
		}

		for _, project := range resp.Projects {
			d.StreamListItem(ctx, project)

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

//// HYDRATE FUNCTIONS

func getAccountProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_account_project.getAccountProject", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway IAM product
	accountApi := account.NewProjectAPI(client)

	projectId := d.EqualsQuals["id"].GetStringValue()

	// No inputs
	if projectId == "" {
		return nil, nil
	}

	data, err := accountApi.GetProject(&account.ProjectAPIGetProjectRequest{
		ProjectID: projectId,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_account_project.getAccountProject", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
