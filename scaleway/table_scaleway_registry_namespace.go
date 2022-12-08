package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayRegistryNamespace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:              "scaleway_registry_namespace",
		Description:       "Registry Namespaces allow you to manage your Container Registry in Scaleway.",
		GetMatrixItemFunc: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listRegistryNamespace,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "region",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getRegsitryNamespace,
			KeyColumns: plugin.AllColumns([]string{"id", "region"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description of the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "status",
				Description: "The current status of the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status").Transform(transform.ToString),
			},
			{
				Name:        "status_message",
				Description: "The current status of the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StatusMessage").Transform(transform.ToString),
			},
			{
				Name:        "endpoint",
				Description: "The endpoint reachable by docker of the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_public",
				Description: "The namespace visibility policy.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IsPublic").Transform(transform.ToBool),
			},
			{
				Name:        "size",
				Description: "Total size of the namespace, calculated as the sum of the size of all images in the namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Size").Transform(transform.ToString),
			},
			{
				Name:        "created_at",
				Description: "The time when the namespace was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the namespace was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "image_count",
				Description: "The number of images in the namespace.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ImageCount").Transform(transform.ToInt),
			},
			{
				Name:        "region",
				Description: "Specifies the region where the namespace is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(transform.ToString),
			},

			// Scaleway standard columns
			{
				Name:        "project",
				Description: "The ID of the project where the namespace resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the namespace resides.",
				Type:        proto.ColumnType_STRING,
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

func listRegistryNamespace(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_regsitry_namespace.listRegsitryNamespace", "region_parsing_error", err)
		return nil, err
	}

	quals := d.KeyColumnQuals
	if quals["region"] != nil && quals["region"].GetStringValue() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_regsitry_namespace.listRegsitryNamespace", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Registry product
	registryApi := registry.NewAPI(client)

	req := &registry.ListNamespacesRequest{
		Region: parseRegionData,
		Page:   scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}

	// Retrieve the list of namespaces
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
		resp, err := registryApi.ListNamespaces(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_regsitry_namespace.listRegsitryNamespace", "query_error", err)
			return nil, err
		}

		for _, namespace := range resp.Namespaces {
			d.StreamListItem(ctx, namespace)

			// Increase the resource count by 1
			count++

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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

func getRegsitryNamespace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_regsitry_namespace.getRegsitryNamespace", "region_parsing_error", err)
		return nil, err
	}

	if d.KeyColumnQuals["region"].GetStringValue() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_regsitry_namespace.getRegsitryNamespace", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Registry product
	registryApi := registry.NewAPI(client)

	id := d.KeyColumnQuals["id"].GetStringValue()
	RegistryRegion := d.KeyColumnQuals["region"].GetStringValue()

	// No inputs
	if id == "" && RegistryRegion == "" {
		return nil, nil
	}

	data, err := registryApi.GetNamespace(&registry.GetNamespaceRequest{
		NamespaceID: id,
		Region:      parseRegionData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_regsitry_namespace.getRegsitryNamespace", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
