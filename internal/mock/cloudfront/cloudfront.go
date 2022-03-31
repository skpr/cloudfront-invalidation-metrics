package cloudfront

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

// MockCloudFrontDistributionClientInterface is a mock cloudfront interface.
type MockCloudFrontDistributionClientInterface interface {
	cloudfront.GetDistributionAPIClient
	cloudfront.GetInvalidationAPIClient
	cloudfront.ListDistributionsAPIClient
	cloudfront.GetInvalidationAPIClient
}

// MockCloudFrontClient is a mock cloudfront client.
type MockCloudFrontClient struct {
	MockCloudFrontDistributionClientInterface
	cloudfront.Client
}

// NewMockCloudFrontClient creates a new mock cloudfront client.
func NewMockCloudFrontClient() *MockCloudFrontClient {
	return &MockCloudFrontClient{}
}

func (c *MockCloudFrontClient) GetDistribution(ctx context.Context, params *cloudfront.GetDistributionInput, optFns ...func(*cloudfront.Options)) (*cloudfront.GetDistributionOutput, error) {
	return &cloudfront.GetDistributionOutput{}, nil
}
func (c *MockCloudFrontClient) GetInvalidation(ctx context.Context, params *cloudfront.GetInvalidationInput, optFns ...func(*cloudfront.Options)) (*cloudfront.GetInvalidationOutput, error) {
	return &cloudfront.GetInvalidationOutput{}, nil
}
func (c *MockCloudFrontClient) ListDistributions(ctx context.Context, params *cloudfront.ListDistributionsInput, optFns ...func(*cloudfront.Options)) (*cloudfront.ListDistributionsOutput, error) {
	return &cloudfront.ListDistributionsOutput{}, nil
}
func (c *MockCloudFrontClient) ListInvalidations(ctx context.Context, params *cloudfront.ListInvalidationsInput, optFns ...func(*cloudfront.Options)) (*cloudfront.ListInvalidationsOutput, error) {
	return &cloudfront.ListInvalidationsOutput{}, nil
}
