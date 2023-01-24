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

func tableScalewayRegistryImage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:              "scaleway_registry_image",
		Description:       "Registry Images allow you to manage your docker images in Scaleway.",
		GetMatrixItemFunc: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate:       listRegistryImages,
			ParentHydrate: listRegistryNamespaces,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				}, 
				{
					Name:    "namespace_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getRegistryImage,
			KeyColumns: plugin.AllColumns([]string{"id"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The image name, unique in a namespace",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace_id",
				Description: "The unique ID of the namespace the image belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NamespaceID"),
			},
			{
				Name:        "id",
				Description: "A unique identifier of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "status",
				Description: "The current status of the Image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "visibility",
				Description: "A `public` image is pullable from internet without authentication, opposed to a `private` image. `inherit` will use the namespace `is_public` parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "Image size in bytes, calculated from the size of image layers. One layer used in two tags of the same image is counted once but one layer used in two images is counted twice.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "created_at",
				Description: "The time when the image was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the image was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "tags",
				Description: "A list of tags attached with the docker image.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "project",
				Description: "The ID of the project where the Image resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the Image resides.",
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

type imageInfo = struct {
	registry.Image
	InstanceID   string
	Region       scw.Region
	Project      string
	Organization string
}

//// LIST FUNCTION

func listRegistryImages(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_registry_image.listRegistryImages", "region_parsing_error", err)
		return nil, err
	}

	// Get Namespace details
	namespaceData := h.Item.(*registry.Namespace)

	quals := d.KeyColumnQuals

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_registry_image.listRegistryImages", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Registry product
	registryApi := registry.NewAPI(client)

	req := &registry.ListImagesRequest{
		Region:      parseRegionData,
		NamespaceID: &namespaceData.ID,
		Page:        scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}

	// Retrieve the list of Images
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
		resp, err := registryApi.ListImages(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_registry_image.listRegistryImages", "query_error", err)
			return nil, err
		}

		for _, image := range resp.Images {
			d.StreamListItem(ctx, imageInfo{*image, namespaceData.ID, namespaceData.Region, namespaceData.ProjectID, namespaceData.OrganizationID})

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

func getRegistryImage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_registry_image.getRegistryImage", "region_parsing_error", err)
		return nil, err
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_registry_image.getRegistryImage", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Registry product
	registryApi := registry.NewAPI(client)

	id := d.KeyColumnQuals["id"].GetStringValue()
	registryRegion := d.KeyColumnQuals["region"].GetStringValue()

	// No inputs
	if id == "" && registryRegion == "" {
		return nil, nil
	}

	data, err := registryApi.GetImage(&registry.GetImageRequest{
		ImageID: id,
		Region:  parseRegionData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_registry_image.getRegistryImage", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
