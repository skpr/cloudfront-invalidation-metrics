package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	cloudfrontclient "cloudfront-invalidation-metrics/internal/aws/cloudfront"
	"cloudfront-invalidation-metrics/internal/logs"
	"cloudfront-invalidation-metrics/internal/metrics"
)

const (
	// CloudWatchNamespace is the CloudWatch Namespace to store metrics in.
	CloudWatchNamespace = "Skpr/CloudFront"

	// InvalidationLogGroupKey is the key used to store the log group name in the tags.
	InvalidationLogGroupKey = "invalidations.cloudfront.skpr.io/loggroup"
	// InvalidationLogStreamKey is the key used to store the log stream name in the tags.
	InvalidationLogStreamKey = "invalidations.cloudfront.skpr.io/stream"
)

// Start is an exported abstraction so that the application can be
// setup in a way that works for you, opposed to being a tightly
// coupled to provided and assumed Clients.
func Start(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get AWS client: %w", err)
	}

	dryRun := os.Getenv("CLOUDFRONT_INVALIDATION_METRICS_DRYRUN") != ""

	metricsClient, err := metrics.New(cloudwatch.NewFromConfig(cfg), CloudWatchNamespace, dryRun)
	if err != nil {
		return fmt.Errorf("failed to setup metricsClient: %w", err)
	}

	logsClient, err := logs.New(cloudwatchlogs.NewFromConfig(cfg), CloudWatchNamespace, dryRun)
	if err != nil {
		return fmt.Errorf("failed to setup logsClient: %w", err)
	}

	return Execute(ctx, cloudfront.NewFromConfig(cfg), metricsClient, logsClient)
}

// Execute will execute the given API calls against the input Clients.
func Execute(ctx context.Context, clientCloudFront cloudfrontclient.ClientInterface, metricsClient metrics.ClientInterface, logsClient logs.ClientInterface) error {
	distributions, err := clientCloudFront.ListDistributions(ctx, &cloudfront.ListDistributionsInput{})
	if err != nil {
		return fmt.Errorf("failed to get CloudFront distibution list: %w", err)
	}

	// FiveMinutesAgo is a variable storing time. It will be used to make a
	// time comparison between the time an invalidation was created and
	// five minutes ago, which is the fixed input time which this lambda
	// is intended to execute.
	fiveMinutesAgo := time.Now().Add(time.Minute * -5)

	for _, distribution := range distributions.DistributionList.Items {
		invalidations, err := clientCloudFront.ListInvalidations(ctx, &cloudfront.ListInvalidationsInput{
			DistributionId: distribution.Id,
		})
		if err != nil {
			return fmt.Errorf("failed to list invalidations: %w", err)
		}

		var (
			countInvalidations float64
			countPaths         float64
			invalidationPaths  []string
		)

		for _, invalidation := range invalidations.InvalidationList.Items {

			if !fiveMinutesAgo.Before(*invalidation.CreateTime) {
				break
			}

			// Include Invalidation in count as the timeframe is acceptable.
			countInvalidations++

			invalidationDetail, err := clientCloudFront.GetInvalidation(ctx, &cloudfront.GetInvalidationInput{
				DistributionId: distribution.Id,
				Id:             invalidation.Id,
			})
			if err != nil {
				return fmt.Errorf("failed to get invalidation detail: %w", err)
			}

			if invalidationDetail != nil {
				countPaths = countPaths + float64(*invalidationDetail.Invalidation.InvalidationBatch.Paths.Quantity)
				invalidationPaths = invalidationDetail.Invalidation.InvalidationBatch.Paths.Items
			}
		}

		// fetch log group and log name from distribution tags
		logGroupName, logStreamName, logExists, err := pullTags(ctx, clientCloudFront, distribution)
		if logExists {
			// send logs to cloudwatch
			err = logsClient.Flush(ctx, logGroupName, logStreamName, cloudwatchlogstypes.InputLogEvent{
				Message:   aws.String(fmt.Sprintf("{ \"InvalidationRequestID\": \"%g\", \"InvalidationPathCount\": %g, \"InvalidatedPaths\": [%s]}", countInvalidations, countPaths, strings.Join(invalidationPaths, ","))),
				Timestamp: aws.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
			})
			if err != nil {
				return fmt.Errorf("failed to push log: %w", err)
			}
		}

		err = metricsClient.Add(cloudwatchtypes.MetricDatum{
			MetricName: aws.String("InvalidationRequest"),
			Unit:       cloudwatchtypes.StandardUnitCount,
			Value:      aws.Float64(countInvalidations),
			Timestamp:  aws.Time(time.Now()),
			Dimensions: []cloudwatchtypes.Dimension{
				{
					Name:  aws.String("Distribution"),
					Value: aws.String(*distribution.Id),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to push metric: InvalidationRequest: %w", err)
		}

		err = metricsClient.Add(cloudwatchtypes.MetricDatum{
			MetricName: aws.String("InvalidationPathCounter"),
			Unit:       cloudwatchtypes.StandardUnitCount,
			Value:      aws.Float64(countPaths),
			Timestamp:  aws.Time(time.Now()),
			Dimensions: []cloudwatchtypes.Dimension{
				{
					Name:  aws.String("Distribution"),
					Value: aws.String(*distribution.Id),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to push metric: InvalidationPathCounter: %w", err)
		}
	}

	return metricsClient.Flush()
}

func main() {
	lambda.Start(Start)
}

// pullTags will pull the tags from the distribution and return the log group and log stream name.
func pullTags(ctx context.Context, clientCloudFront cloudfrontclient.ClientInterface, distribution types.DistributionSummary) (logGroupName string, logStreamName string, tagExists bool, err error) {
	tags, err := clientCloudFront.ListTagsForResource(ctx, &cloudfront.ListTagsForResourceInput{
		Resource: distribution.ARN,
	})
	if err != nil {
		return "", "", false, fmt.Errorf("failed to list tags: %w", err)
	}

	for _, tag := range tags.Tags.Items {
		if *tag.Key == InvalidationLogGroupKey {
			logGroupName = *tag.Value
		}

		if *tag.Key == InvalidationLogStreamKey {
			logStreamName = *tag.Value
		}
	}

	if logGroupName == "" || logStreamName == "" {
		return "", "", false, fmt.Errorf("failed to get log group or log stream name from tags")
	}

	return logGroupName, logStreamName, true, nil
}
