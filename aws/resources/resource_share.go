package resources

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ram"
	"github.com/aws/aws-sdk-go-v2/service/ram/types"
	"github.com/gruntwork-io/cloud-nuke/config"
	"github.com/gruntwork-io/go-commons/errors"
	"github.com/hashicorp/go-multierror"
)

func (s *ResourceShares) getAll(ctx context.Context, configObj config.Config) ([]*string, error) {
	var allResourceShares []*string
	paginator := ram.NewGetResourceSharesPaginator(s.Client, &ram.GetResourceSharesInput{ResourceOwner: "SELF"})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, resourceShare := range page.ResourceShares {
			if s.shouldInclude(&resourceShare, configObj) {
				allResourceShares = append(allResourceShares, resourceShare.ResourceShareArn)
			}
		}
	}

	return allResourceShares, nil
}

func (s *ResourceShares) nukeAll(arns []*string) error {
	if len(arns) == 0 {
		return nil
	}

	var allErrs *multierror.Error
	for _, arn := range arns {
		_, err := s.Client.DeleteResourceShare(s.Context, &ram.DeleteResourceShareInput{
			ResourceShareArn: arn,
		})
		if err != nil {
			allErrs = multierror.Append(allErrs, err)
		}
	}

	finalErr := allErrs.ErrorOrNil()
	if finalErr != nil {
		return errors.WithStackTrace(finalErr)
	}

	return nil
}

func (s *ResourceShares) shouldInclude(share *types.ResourceShare, configObj config.Config) bool {
	if share == nil {
		return false
	}

	return configObj.ResourceShare.ShouldInclude(config.ResourceValue{
		Name: share.Name,
		Time: share.CreationTime,
		Tags: ConvertRAMTagsToMap(share.Tags),
	})
}

func ConvertRAMTagsToMap(tags []types.Tag) map[string]string {
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[*tag.Key] = *tag.Value
	}

	return tagMap
}
