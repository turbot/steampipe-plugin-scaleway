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

func tableScalewayKubernetesNode(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:              "scaleway_kubernetes_node",
		Description:       "Kubernetes Nodes allow you to manage your Container Kubernetes in Scaleway.",
		GetMatrixItemFunc: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate:       listKubernetesNodes,
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
			Hydrate:    getKubernetesNode,
			KeyColumns: plugin.AllColumns([]string{"id", "region"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current status of the node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_id",
				Description: "The cluster ID of the node.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterID"),
			},
			{
				Name:        "pool_id",
				Description: "The pool ID of the node.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PoolID").Transform(transform.ToString),
			},
			{
				Name:        "provider_id",
				Description: "It is prefixed by instance type and location information.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProviderID").Transform(transform.ToString),
			},
			{
				Name:        "public_ip_v4",
				Description: "The public IPv4 address of the node.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPV4").Transform(transform.ToString),
			},
			{
				Name:        "public_ip_v6",
				Description: "The public IPv6 address of the node.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PublicIPV6").Transform(transform.ToString),
			},
			{
				Name:        "conditions",
				Description: "These conditions contain the Node Problem Detector condition.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Conditions"),
			},
			{
				Name:        "error_message",
				Description: "The details of the error, if any occured when managing the node.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ErrorMessage").Transform(transform.ToString),
			},
			{
				Name:        "created_at",
				Description: "The time when the node was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the node was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Scaleway standard columns
			{
				Name:        "project",
				Description: "The ID of the project where the node resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the node resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "Specifies the region where the node is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(transform.ToString),
			},
			{
				Name:        "zone",
				Description: "Specifies the zone where the node is located.",
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
				Description: "A list of tags associated with the node.",
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

func listKubernetesNodes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_node.listKubernetesNodes", "region_parsing_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	if quals["region"] != nil && quals["region"].GetStringValue() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_node.listKubernetesNodes", "connection_error", err)
		return nil, err
	}

	// Get cluster details
	clusterData := h.Item.(*k8s.Cluster)

	// Create SDK objects for Scaleway Kubernetes product
	kubernetesApi := k8s.NewAPI(client)

	req := &k8s.ListNodesRequest{
		Region:    parseRegionData,
		ClusterID: clusterData.ID,
		Page:      scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}

	// Retrieve the list of nodes
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
		resp, err := kubernetesApi.ListNodes(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_kubernetes_node.listKubernetesNodes", "query_error", err)
			return nil, err
		}

		for _, node := range resp.Nodes {
			d.StreamListItem(ctx, node)

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

func getKubernetesNode(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_node.getKubernetesNode", "region_parsing_error", err)
		return nil, err
	}

	if d.EqualsQuals["region"].GetStringValue() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_node.getKubernetesNode", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Kubernetes product
	kubernetesApi := k8s.NewAPI(client)

	id := d.EqualsQuals["id"].GetStringValue()

	// No inputs
	if id == "" {
		return nil, nil
	}

	data, err := kubernetesApi.GetNode(&k8s.GetNodeRequest{
		NodeID: id,
		Region: parseRegionData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_node.getKubernetesNode", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
