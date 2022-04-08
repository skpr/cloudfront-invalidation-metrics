package cloudwatch

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/smithy-go/middleware"
)

// CloudWatchClientInterface is a mock cloudwatch interface.
type CloudWatchClientInterface interface {
	PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error)
}

// MockCloudWatchClient is a mock cloudwatch client.
type MockCloudWatchClient struct{}

func (c MockCloudWatchClient) PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error) {
	return &cloudwatch.PutMetricDataOutput{
		ResultMetadata: middleware.Metadata{},
	}, nil
}
