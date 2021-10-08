package scaleway

import (
	"context"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func handleNilString(_ context.Context, d *transform.TransformData) (interface{}, error) {
	value := types.SafeString(d.Value)
	if value == "" {
		return "false", nil
	}
	return value, nil
}

func flattenObjectBucketTags(tagsSet []*s3.Tag) []string {
	var tags []string

	for _, tagSet := range tagsSet {
		if tagSet.Key != nil {
			tags = append(tags, *tagSet.Key)
		}
	}

	return tags
}
