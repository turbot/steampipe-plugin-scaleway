package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayKubernetesPool(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:              "scaleway_kubernetes_pool",
		Description:       "Kubernetes Pools allow you to manage your Container Kubernetes in Scaleway.",
		GetMatrixItemFunc: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate:       listKubernetesPools,
			ParentHydrate: listKubernetesClusters,
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
			Hydrate:    getKubernetesPool,
			KeyColumns: plugin.AllColumns([]string{"id", "region"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current status of the pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_id",
				Description: "The cluster ID of the pool.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterID"),
			},
			{
				Name:        "node_type",
				Description: "The type of the node.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NodeType").Transform(transform.ToString),
			},
			{
				Name:        "version",
				Description: "The Kubernetes version of the pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The time when the pool was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the pool was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "autoscaling",
				Description: "The enablement of the autoscaling feature for the pool.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Autoscaling").Transform(transform.ToBool),
			},
			{
				Name:        "size",
				Description: "The size (number of nodes) of the pool.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "min_size",
				Description: "The minimum size of the pool.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "max_size",
				Description: "The maximum size of the pool",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "autohealing",
				Description: "The enablement of the autohealing feature for the pool.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Autohealing").Transform(transform.ToBool),
			},
			{
				Name:        "placement_group_id",
				Description: "The placement group ID in which all the nodes of the pool will be created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PlacementGroupID").Transform(transform.ToString),
			},
			{
				Name:        "kubelet_args",
				Description: "The Kubelet arguments to be used by this pool.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KubeletArgs").Transform(transform.ToString),
			},
			{
				Name:        "upgrade_policy",
				Description: "The Pool upgrade policy.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "root_volume_type",
				Description: "The system volume disk type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "root_volume_size",
				Description: "The system volume disk size.",
				Type:        proto.ColumnType_JSON,
			},

			// Scaleway standard columns
			{
				Name:        "project",
				Description: "The ID of the project where the pool resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the pool resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "Specifies the region where the pool is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(transform.ToString),
			},
			{
				Name:        "zone",
				Description: "Specifies the zone where the pool is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(transform.ToString),
			},
			{
				Name:        "id",
				Description: "A unique identifier of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "tags",
				Description: "A list of tags associated with the pool.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
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

func listKubernetesPools(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_pool.listKubernetesPools", "region_parsing_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	if quals["region"] != nil && quals["region"].GetStringValue() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_pool.listKubernetesPools", "connection_error", err)
		return nil, err
	}

	// Get cluster details
	clusterData := h.Item.(*k8s.Cluster)

	// Create SDK objects for Scaleway Kubernetes product
	kubernetesApi := k8s.NewAPI(client)

	req := &k8s.ListPoolsRequest{
		Region:    parseRegionData,
		ClusterID: clusterData.ID,
		Page:      scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}

	// Retrieve the list of pools
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
		resp, err := kubernetesApi.ListPools(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_kubernetes_pool.listKubernetesPools", "query_error", err)
			return nil, err
		}

		for _, pool := range resp.Pools {
			d.StreamListItem(ctx, pool)

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

func getKubernetesPool(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_pool.getKubernetesPool", "region_parsing_error", err)
		return nil, err
	}

	if d.EqualsQuals["region"].GetStringValue() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_pool.getKubernetesPool", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Kubernetes product
	kubernetesApi := k8s.NewAPI(client)

	id := d.EqualsQuals["id"].GetStringValue()

	// No inputs
	if id == "" {
		return nil, nil
	}

	data, err := kubernetesApi.GetPool(&k8s.GetPoolRequest{
		PoolID: id,
		Region: parseRegionData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_pool.getKubernetesPool", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
