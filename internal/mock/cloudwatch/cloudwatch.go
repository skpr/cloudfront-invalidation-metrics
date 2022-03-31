package cloudwatch

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

// MockCloudWatchClientInterface is a mock cloudwatch interface.
type MockCloudWatchClientInterface interface {
	PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error)
}

// MockCloudWatchClient is a mock cloudwatch client.
type MockCloudWatchClient struct {
	MockCloudWatchClientInterface
	cloudwatch.Client
}

// NewMockCloudwatchClient creates a new mock cloudwatch client.
func NewMockCloudwatchClient() *MockCloudWatchClient {
	return &MockCloudWatchClient{}
}

func (c *MockCloudWatchClient) PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error) {
	return &cloudwatch.PutMetricDataOutput{}, nil
}
