package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayKubernetesCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:              "scaleway_kubernetes_cluster",
		Description:       "Kubernetes Clusters allow you to manage your Container Kubernetes in Scaleway.",
		GetMatrixItemFunc: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listKubernetesClusters,
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
			Hydrate:    getKubernetesCluster,
			KeyColumns: plugin.AllColumns([]string{"id", "region"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current status of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status").Transform(transform.ToString),
			},
			{
				Name:        "type",
				Description: "The type of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Type").Transform(transform.ToString),
			},
			{
				Name:        "version",
				Description: "The Kubernetes version of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cni",
				Description: "The cluster visibility policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Cni"),
			},
			{
				Name:        "cluster_url",
				Description: "The Kubernetes API server URL of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterURL"),
			},
			{
				Name:        "dns_wildcard",
				Description: "The DNS wildcard resovling all the ready nodes of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DNSWildcard"),
			},
			{
				Name:        "created_at",
				Description: "The time when the cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when the cluster was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "autoscaler_config",
				Description: "The autoscaler config for the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AutoscalerConfig"),
			},
			{
				Name:        "dashboard_enabled",
				Description: "The enablement of the Kubernetes Dashboard in the cluster.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DashboardEnabled").Transform(transform.ToBool),
			},
			{
				Name:        "auto_upgrade",
				Description: "The auto upgrade configuration of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AutoUpgrade"),
			},
			{
				Name:        "upgrade_available",
				Description: "True if a new Kubernetes version is available.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "feature_gates",
				Description: "The list of enabled feature gates.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FeatureGates").Transform(transform.ToString),
			},
			{
				Name:        "admission_plugins",
				Description: "The list of enabled admission plugins.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AdmissionPlugins").Transform(transform.ToString),
			},
			{
				Name:        "open_id_connect_config",
				Description: "The OpenID Connect configuration of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("OpenIDConnectConfig"),
			},
			{
				Name:        "apiserver_cert_sans",
				Description: "The additional Subject Alternative Names for the Kubernetes API server certificate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApiserverCertSans").Transform(transform.ToString),
			},

			// Scaleway standard columns
			{
				Name:        "project",
				Description: "The ID of the project where the cluster resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the cluster resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "Specifies the region where the cluster is located.",
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
				Description: "A list of tags associated with the cluster.",
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

func listKubernetesClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString("region")

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_cluster.listKubernetesClusters", "region_parsing_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	if quals["region"] != nil && quals["region"].GetStringValue() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_cluster.listKubernetesClusters", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Kubernetes product
	kubernetesApi := k8s.NewAPI(client)

	req := &k8s.ListClustersRequest{
		Region: parseRegionData,
		Page:   scw.Int32Ptr(1),
	}
	// Additional filters
	if quals["name"] != nil {
		req.Name = scw.StringPtr(quals["name"].GetStringValue())
	}

	// Retrieve the list of clusters
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
		resp, err := kubernetesApi.ListClusters(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_kubernetes_cluster.listKubernetesClusters", "query_error", err)
			return nil, err
		}

		for _, cluster := range resp.Clusters {
			d.StreamListItem(ctx, cluster)

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

func getKubernetesCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString("region")

	parseRegionData, err := scw.ParseRegion(region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_cluster.getKubernetesCluster", "region_parsing_error", err)
		return nil, err
	}

	if d.EqualsQuals["region"].GetStringValue() != region {
		return nil, nil
	}

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_cluster.getKubernetesCluster", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Kubernetes product
	kubernetesApi := k8s.NewAPI(client)

	id := d.EqualsQuals["id"].GetStringValue()

	// No input id has been passed
	if id == "" {
		return nil, nil
	}

	data, err := kubernetesApi.GetCluster(&k8s.GetClusterRequest{
		ClusterID: id,
		Region:    parseRegionData,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_kubernetes_cluster.getKubernetesCluster", "query_error", err)
		if is404Error(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
