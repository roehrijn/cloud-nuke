package resources

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ram"
	"github.com/gruntwork-io/cloud-nuke/config"
	"github.com/gruntwork-io/go-commons/errors"
)

type ResourceShares struct {
	BaseAwsResource
	Client *ram.Client
	Arns   []string
}

func (s *ResourceShares) Init(cfg aws.Config) {
	s.Client = ram.NewFromConfig(cfg)
}

func (s *ResourceShares) ResourceName() string {
	return "resource-share"
}

func (s *ResourceShares) ResourceIdentifiers() []string {
	return s.Arns
}

func (s *ResourceShares) GetAndSetIdentifiers(ctx context.Context, configObj config.Config) ([]string, error) {
	identifiers, err := s.getAll(ctx, configObj)
	if err != nil {
		return nil, err
	}

	s.Arns = aws.ToStringSlice(identifiers)
	return s.Arns, nil
}

func (s *ResourceShares) Nuke(identifiers []string) error {
	if err := s.nukeAll(aws.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
