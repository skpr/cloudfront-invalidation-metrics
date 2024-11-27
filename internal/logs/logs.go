package logs

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	cloudwatchlogsclient "cloudfront-invalidation-metrics/internal/aws/cloudwatchlogs"
)

const (
	// AwsPayloadLimit is the maximum quality for a data-set to contain
	// before AWS will reject the payload.
	AwsPayloadLimit = 20
)

// ClientInterface for pushing metrics to CloudWatch.
type ClientInterface interface {
	Flush(ctx context.Context, logGroupName string, logStreamName string, logData types.InputLogEvent) error
}

// Client for pushing metrics to CloudWatch.
type Client struct {
	CloudWatchLogs cloudwatchlogsclient.ClientInterface
	Namespace      string
	LogData        []types.InputLogEvent
	DryRun         bool
}

//

// New client for pushing metrics to CloudWatch.
func New(cloudwatchlogs cloudwatchlogsclient.ClientInterface, namespace string, dryRun bool) (*Client, error) {
	return &Client{
		CloudWatchLogs: cloudwatchlogs,
		Namespace:      namespace,
		DryRun:         dryRun,
	}, nil
}

// Flush logs to CloudWatch.
func (c *Client) Flush(ctx context.Context, logGroupName string, logStreamName string, logData types.InputLogEvent) error {
	if c.DryRun {
		return nil
	}

	fmt.Printf("Creating log group: %s\n", logGroupName)
	_, err := c.CloudWatchLogs.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(logGroupName),
	})
	if err != nil {
		var awsErr *types.ResourceAlreadyExistsException
		if errors.As(err, &awsErr) {
			fmt.Printf("LogGroup already exists: %s\n", logGroupName)
		}
		return err
	}

	fmt.Printf("Creating log stream: %s\n", logStreamName)
	_, err = c.CloudWatchLogs.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	})
	if err != nil {
		var awsErr *types.ResourceAlreadyExistsException
		if errors.As(err, &awsErr) {
			fmt.Printf("LogStream already exists: %s\n", logGroupName)
		}
		return err
	}

	logEvent := []types.InputLogEvent{logData}

	_, err = c.CloudWatchLogs.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
		LogEvents:     logEvent,
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	})
	if err != nil {
		return fmt.Errorf("failed to put log events: %w", err)
	}

	fmt.Printf("Flushed log events to %s/%s\n", logGroupName, logStreamName)

	return err
}
