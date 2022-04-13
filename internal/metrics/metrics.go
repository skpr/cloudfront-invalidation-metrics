package metrics

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	cloudwatchclient "cloudfront-invalidation-metrics/internal/aws/cloudwatch"
)

const (
	// AwsPayloadLimit is the maximum quality for a data-set to contain
	// before AWS will reject the payload.
	AwsPayloadLimit = 20
)

// ClientInterface for pushing metrics to CloudWatch.
type ClientInterface interface {
	Add(datum types.MetricDatum) error
	Flush() error
}

// Client for pushing metrics to CloudWatch.
type Client struct {
	CloudWatch cloudwatchclient.CloudWatchClientInterface
	Namespace  string
	Data       []types.MetricDatum
	DryRun     bool
}

// New client for pushing metrics to CloudWatch.
func New(cloudwatch cloudwatchclient.CloudWatchClientInterface, namespace string, dryRun bool) (*Client, error) {
	return &Client{
		CloudWatch: cloudwatch,
		Namespace:  namespace,
		DryRun:     dryRun,
	}, nil
}

// Add metrics to Client.
func (c *Client) Add(data types.MetricDatum) error {
	if len(c.Data) == AwsPayloadLimit {
		err := c.Flush()
		if err != nil {
			return err
		}
	}

	c.Data = append(c.Data, data)

	return nil
}

// Flush metrics to CloudWatch.
func (c *Client) Flush() error {
	if c.DryRun {
		return nil
	}

	if len(c.Data) == 0 {
		return nil
	}

	_, err := c.CloudWatch.PutMetricData(context.Background(), &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(c.Namespace),
		MetricData: c.Data,
	})
	if err != nil {
		return err
	}

	c.Data = []types.MetricDatum{}

	return err
}
