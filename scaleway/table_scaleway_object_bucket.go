package scaleway

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayObjectBucket(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "scaleway_object_bucket",
		Description:   "A Scaleway Object bucket is a public cloud storage resource available in Scaleway, an object storage offering.",
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listObjectBuckets,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user-defined name of the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "An unique identifier of the bucket.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "bucket_policy_is_public",
				Description: "The policy status for an Scaleway bucket, indicating whether the bucket is public.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Hydrate:     getBucketIsPublic,
				Transform:   transform.FromField("PolicyStatus.IsPublic"),
			},
			{
				Name:        "versioning_enabled",
				Description: "The versioning state of a bucket.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketVersioning,
				Transform:   transform.FromField("Status").Transform(handleNilString).Transform(transform.ToBool),
			},
			{
				Name:        "versioning_mfa_delete",
				Description: "The MFA Delete status of the versioning state.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBucketVersioning,
				Transform:   transform.FromField("MFADelete").Transform(handleNilString).Transform(transform.ToBool),
			},
			{
				Name:        "cors_rule",
				Description: "The cors configuration information set for the bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketCors,
				Transform:   transform.FromField("CORSRules"),
			},
			{
				Name:        "acl",
				Description: "The access control list (ACL) of a bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketACL,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "lifecycle_rules",
				Description: "The lifecycle configuration information of the bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketLifecycle,
				Transform:   transform.FromField("Rules"),
			},
			{
				Name:        "website",
				Description: "The website configuration for a bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketWebsite,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "policy",
				Description: "The resource-based policy access document for the bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketPolicy,
				Transform:   transform.FromField("Policy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "tags",
				Description: "A list of tags associated with the bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketTagging,
				Transform:   transform.FromValue(),
			},

			// Scaleway standard columns
			{
				Name:        "region",
				Description: "Specifies the region where the bucket resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project",
				Description: "The ID of the project where the bucket resides.",
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

type bucketInfo = struct {
	s3.Bucket
	Region  string
	Project string
}

//// LIST FUNCTION

func listObjectBuckets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create client
	client, err := getObjectSessionConfig(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.listObjectBuckets", "connection_error", err)
		return nil, err
	}

	resp, err := client.ListBucketsWithContext(ctx, &s3.ListBucketsInput{})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.listObjectBuckets", "query_error", err)
		return nil, err
	}

	bucketOwner := strings.Split(*resp.Owner.ID, ":")[1]
	for _, bucket := range resp.Buckets {
		d.StreamListItem(ctx, bucketInfo{*bucket, region, bucketOwner})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getBucketVersioning(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixRegion := plugin.GetMatrixItem(ctx)["region"].(string)
	bucket := h.Item.(bucketInfo)

	if matrixRegion != bucket.Region {
		return nil, nil
	}

	// Create client
	client, err := getObjectSessionConfig(ctx, d, bucket.Region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketVersioning", "connection_error", err)
		return nil, err
	}

	data, err := client.GetBucketVersioningWithContext(ctx, &s3.GetBucketVersioningInput{
		Bucket: bucket.Name,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketVersioning", "query_error", err)
		return nil, err
	}

	return data, nil
}

func getBucketIsPublic(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixRegion := plugin.GetMatrixItem(ctx)["region"].(string)
	bucket := h.Item.(bucketInfo)

	if matrixRegion != bucket.Region {
		return nil, nil
	}

	// Create client
	client, err := getObjectSessionConfig(ctx, d, bucket.Region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketIsPublic", "connection_error", err)
		return nil, err
	}

	data, err := client.GetBucketPolicyStatusWithContext(ctx, &s3.GetBucketPolicyStatusInput{
		Bucket: bucket.Name,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketIsPublic", "query_error", err)
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NoSuchBucketPolicy" {
				return &s3.GetBucketPolicyStatusOutput{}, nil
			}
		}
		return nil, err
	}

	return data, nil
}

func getBucketPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixRegion := plugin.GetMatrixItem(ctx)["region"].(string)
	bucket := h.Item.(bucketInfo)

	if matrixRegion != bucket.Region {
		return nil, nil
	}

	// Create client
	client, err := getObjectSessionConfig(ctx, d, bucket.Region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketIsPublic", "connection_error", err)
		return nil, err
	}

	bucketPolicy, err := client.GetBucketPolicyWithContext(ctx, &s3.GetBucketPolicyInput{
		Bucket: bucket.Name,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketPolicy", "query_error", err)
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NoSuchBucketPolicy" {
				return &s3.GetBucketPolicyOutput{}, nil
			}
		}
		return nil, err
	}

	return bucketPolicy, nil
}

func getBucketLifecycle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixRegion := plugin.GetMatrixItem(ctx)["region"].(string)
	bucket := h.Item.(bucketInfo)

	if matrixRegion != bucket.Region {
		return nil, nil
	}

	// Create client
	client, err := getObjectSessionConfig(ctx, d, bucket.Region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketIsPublic", "connection_error", err)
		return nil, err
	}

	data, err := client.GetBucketLifecycleConfigurationWithContext(ctx, &s3.GetBucketLifecycleConfigurationInput{
		Bucket: bucket.Name,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketLifecycle", "query_error", err)
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NoSuchLifecycleConfiguration" {
				return nil, nil
			}
		}
		return nil, err
	}

	return data, nil
}

func getBucketACL(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixRegion := plugin.GetMatrixItem(ctx)["region"].(string)
	bucket := h.Item.(bucketInfo)

	if matrixRegion != bucket.Region {
		return nil, nil
	}

	// Create client
	client, err := getObjectSessionConfig(ctx, d, bucket.Region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketIsPublic", "connection_error", err)
		return nil, err
	}

	data, err := client.GetBucketAclWithContext(ctx, &s3.GetBucketAclInput{
		Bucket: bucket.Name,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getBucketTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixRegion := plugin.GetMatrixItem(ctx)["region"].(string)
	bucket := h.Item.(bucketInfo)

	if matrixRegion != bucket.Region {
		return nil, nil
	}

	// Create client
	client, err := getObjectSessionConfig(ctx, d, bucket.Region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketIsPublic", "connection_error", err)
		return nil, err
	}

	bucketTags, _ := client.GetBucketTaggingWithContext(ctx, &s3.GetBucketTaggingInput{
		Bucket: bucket.Name,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketTagging", "query_error", err)
		return nil, err
	}

	return flattenObjectBucketTags(bucketTags.TagSet), nil
}

func getBucketCors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixRegion := plugin.GetMatrixItem(ctx)["region"].(string)
	bucket := h.Item.(bucketInfo)

	if matrixRegion != bucket.Region {
		return nil, nil
	}

	// Create client
	client, err := getObjectSessionConfig(ctx, d, bucket.Region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketIsPublic", "connection_error", err)
		return nil, err
	}

	data, _ := client.GetBucketCorsWithContext(ctx, &s3.GetBucketCorsInput{
		Bucket: bucket.Name,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketCors", "query_error", err)
		return nil, err
	}

	return data, nil
}

func getBucketWebsite(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixRegion := plugin.GetMatrixItem(ctx)["region"].(string)
	bucket := h.Item.(bucketInfo)

	if matrixRegion != bucket.Region {
		return nil, nil
	}

	// Create client
	client, err := getObjectSessionConfig(ctx, d, bucket.Region)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketIsPublic", "connection_error", err)
		return nil, err
	}

	data, _ := client.GetBucketWebsiteWithContext(ctx, &s3.GetBucketWebsiteInput{
		Bucket: bucket.Name,
	})
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_object_bucket.getBucketWebsite", "query_error", err)
		return nil, err
	}

	return data, nil
}
