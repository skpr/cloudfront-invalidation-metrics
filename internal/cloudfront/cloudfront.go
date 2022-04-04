package cloudfront

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/smithy-go/middleware"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

// CloudFrontClientInterface is a mock cloudfront client.
type CloudFrontClientInterface interface {
	GetDistribution(ctx context.Context, params *cloudfront.GetDistributionInput, optFns ...func(*cloudfront.Options)) (*cloudfront.GetDistributionOutput, error)
	GetInvalidation(ctx context.Context, params *cloudfront.GetInvalidationInput, optFns ...func(*cloudfront.Options)) (*cloudfront.GetInvalidationOutput, error)
	ListDistributions(ctx context.Context, params *cloudfront.ListDistributionsInput, optFns ...func(*cloudfront.Options)) (*cloudfront.ListDistributionsOutput, error)
	ListInvalidations(ctx context.Context, params *cloudfront.ListInvalidationsInput, optFns ...func(*cloudfront.Options)) (*cloudfront.ListInvalidationsOutput, error)
}

type MockCloudFrontClient struct{}

func (c MockCloudFrontClient) GetDistribution(ctx context.Context, params *cloudfront.GetDistributionInput, optFns ...func(*cloudfront.Options)) (*cloudfront.GetDistributionOutput, error) {
	return &cloudfront.GetDistributionOutput{
		Distribution: &types.Distribution{
			Id: aws.String("test-distribution-id"),
		},
	}, nil
}
func (c MockCloudFrontClient) GetInvalidation(ctx context.Context, params *cloudfront.GetInvalidationInput, optFns ...func(*cloudfront.Options)) (*cloudfront.GetInvalidationOutput, error) {
	return &cloudfront.GetInvalidationOutput{
		Invalidation: &types.Invalidation{
			CreateTime: aws.Time(time.Now()),
			Id:         aws.String("test-invalidation-id"),
			InvalidationBatch: &types.InvalidationBatch{
				Paths: &types.Paths{
					Quantity: aws.Int32(3),
					Items: []string{
						"/test-item-one",
						"/test-item-two",
						"/test-item-three",
					},
				},
			},
			Status: aws.String("Completed"),
		},
		ResultMetadata: middleware.Metadata{},
	}, nil
}
func (c MockCloudFrontClient) ListDistributions(ctx context.Context, params *cloudfront.ListDistributionsInput, optFns ...func(*cloudfront.Options)) (*cloudfront.ListDistributionsOutput, error) {
	return &cloudfront.ListDistributionsOutput{
		DistributionList: &types.DistributionList{
			Items: []types.DistributionSummary{
				{
					Id: aws.String("test-distribution-id"),
				},
			},
		},
		ResultMetadata: middleware.Metadata{},
	}, nil
}
func (c MockCloudFrontClient) ListInvalidations(ctx context.Context, params *cloudfront.ListInvalidationsInput, optFns ...func(*cloudfront.Options)) (*cloudfront.ListInvalidationsOutput, error) {
	return &cloudfront.ListInvalidationsOutput{
		InvalidationList: &types.InvalidationList{
			Items: []types.InvalidationSummary{
				{
					Id:         aws.String("test-invalidation-id"),
					Status:     aws.String("Completed"),
					CreateTime: aws.Time(time.Now()),
				},
			},
		},
		ResultMetadata: middleware.Metadata{},
	}, nil
}
