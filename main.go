package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	cloudfrontclient "cloudfront-invalidation-metrics/internal/cloudfront"
	cloudwatchclient "cloudfront-invalidation-metrics/internal/cloudwatch"
	push_metrics "cloudfront-invalidation-metrics/internal/push-metrics"
)

const (
	// CloudWatchNamespace is the CloudWatch Namespace to store metrics in.
	CloudWatchNamespace = "Skpr/CloudFront"
)

// Start is an exported abstraction so that the application can be
// setup in a way that works for you, opposed to being a tightly
// coupled to provided and assumed Clients.
func Start(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("failed to get AWS client config: %s\n", err.Error())
	}

	clientCloudWatch := cloudwatch.NewFromConfig(cfg)
	clientCloudFront := cloudfront.NewFromConfig(cfg)

	return Execute(ctx, clientCloudFront, clientCloudWatch)
}

// Execute will execute the given API calls against the input Clients.
func Execute(ctx context.Context, clientCloudFront cloudfrontclient.CloudFrontClientInterface, clientCloudWatch cloudwatchclient.CloudWatchClientInterface) error {
	distributions, err := clientCloudFront.ListDistributions(ctx, &cloudfront.ListDistributionsInput{})
	if err != nil {
		return fmt.Errorf("failed to get CloudFront distibution list: %w", err)
	}

	var data []types.MetricDatum
	dataQueue := push_metrics.Queue{
		Namespace: CloudWatchNamespace,
	}

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
		)

		for _, invalidation := range invalidations.InvalidationList.Items {
			acceptable, err := IsTimeRangeAcceptable(*invalidation.CreateTime)
			if err != nil {
				return fmt.Errorf("failed to determine if invalidation is in range: %w", err)
			}

			if acceptable {
				countInvalidations++
			}

			invalidationDetail, err := clientCloudFront.GetInvalidation(ctx, &cloudfront.GetInvalidationInput{
				DistributionId: distribution.Id,
				Id:             invalidation.Id,
			})
			if err != nil {
				return fmt.Errorf("failed to get invalidation detail: %w", err)
			}

			if invalidationDetail != nil {
				if acceptable {
					countPaths = countPaths + float64(*invalidationDetail.Invalidation.InvalidationBatch.Paths.Quantity)
				}
			}
		}

		data = append(data, types.MetricDatum{
			MetricName: aws.String("InvalidationRequest"),
			Unit:       types.StandardUnitCount,
			Value:      aws.Float64(countInvalidations),
			Timestamp:  aws.Time(time.Now()),
			Dimensions: []types.Dimension{
				{
					Name:  aws.String("Distribution"),
					Value: aws.String(*distribution.Id),
				},
			},
		})

		data = append(data, types.MetricDatum{
			MetricName: aws.String("InvalidationPathCounter"),
			Unit:       types.StandardUnitCount,
			Value:      aws.Float64(countPaths),
			Timestamp:  aws.Time(time.Now()),
			Dimensions: []types.Dimension{
				{
					Name:  aws.String("Distribution"),
					Value: aws.String(*distribution.Id),
				},
			},
		})
	}

	for _, queueItem := range data {
		if dataQueue.QueueFull {
			dataQueue.Flush(clientCloudWatch)
		}

		if err = dataQueue.Add(queueItem); err != nil {
			return fmt.Errorf(err.Error())
		}
	}

	err = dataQueue.Flush(clientCloudWatch)

	return err
}

// IsTimeRangeAcceptable will determine if an input time is within
// a given date range. It's intended here to be a frequency of every
// five minutes.
func IsTimeRangeAcceptable(input time.Time) (bool, error) {
	// Calculate what is the acceptable age of an invalidation to ingest.
	fiveMinutesAgo := time.Now().Add(time.Minute * -5)
	if input.Before(fiveMinutesAgo) {
		return false, fmt.Errorf("input time is not in a reportable time frame")
	}
	return true, nil
}

func main() {
	lambda.Start(Start)
}
