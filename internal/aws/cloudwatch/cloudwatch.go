package cloudwatch

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/smithy-go/middleware"
)

// ClientInterface is a mock cloudwatch interface.
type ClientInterface interface {
	PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error)
}

// MockCloudWatchClient is a mock cloudwatch client.
type MockCloudWatchClient struct {
	MetricData []types.MetricDatum
}

// PutMetricData stores the metrics for later.
func (c *MockCloudWatchClient) PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error) {
	// Store the metrics for later.
	c.MetricData = append(c.MetricData, params.MetricData...)

	return &cloudwatch.PutMetricDataOutput{
		ResultMetadata: middleware.Metadata{},
	}, nil
}
