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
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	"cloudfront-invalidation-metrics/internal/push-metrics"
)

const (
	// CloudWatchNamespace is the CloudWatch Namespace to store metrics in.
	CloudWatchNamespace = "Skpr/CloudFront"
)

// Start is an exported abstraction so that the application can be
// setup in a way that works for you, opposed to being a tightly
// coupled to provided and assumed Clients.
func Start() error {

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("failed to get AWS client config: %s\n", err.Error())
	}

	clientCloudWatch := cloudwatch.NewFromConfig(cfg)
	clientCloudFront := cloudfront.NewFromConfig(cfg)

	return Execute(ctx, *clientCloudFront, *clientCloudWatch)
}

// Execute will execute the given API calls against the input Clients.
func Execute(ctx context.Context, clientCloudFront cloudfront.Client, clientCloudWatch cloudwatch.Client) error {
	distributions, err := clientCloudFront.ListDistributions(ctx, &cloudfront.ListDistributionsInput{})
	if err != nil {
		return fmt.Errorf("failed to get CloudFront distibution list: %w", err)
	}

	var data []cwtypes.MetricDatum
	dataQueue := push_metrics.Queue{
		Client:    clientCloudWatch,
		Namespace: CloudWatchNamespace,
	}

	for _, distribution := range distributions.DistributionList.Items {
		invalidations, err := clientCloudFront.ListInvalidations(ctx, &cloudfront.ListInvalidationsInput{
			DistributionId: distribution.Id,
		})
		if err != nil {
			return err
		}

		var (
			countInvalidations float64
			countPaths         float64
		)

		for _, invalidation := range invalidations.InvalidationList.Items {

			acceptable, err := IsTimeRangeAcceptable(time.Now(), *invalidation.CreateTime)
			if err != nil {
				continue
			}

			if acceptable {
				countInvalidations++
			}

			invalidationDetail, _ := clientCloudFront.GetInvalidation(ctx, &cloudfront.GetInvalidationInput{
				DistributionId: distribution.Id,
				Id:             invalidation.Id,
			})

			if invalidationDetail != nil {
				acceptable, err := IsTimeRangeAcceptable(time.Now(), *invalidationDetail.Invalidation.CreateTime)
				if err != nil {
					continue
				}

				if !acceptable {
					break
				}

				if acceptable {
					countPaths = countPaths + float64(*invalidationDetail.Invalidation.InvalidationBatch.Paths.Quantity)
				}

			}
		}

		data = append(data, cwtypes.MetricDatum{
			MetricName: aws.String("InvalidationRequest"),
			Unit:       cwtypes.StandardUnitCount,
			Value:      aws.Float64(countInvalidations),
			Timestamp:  aws.Time(time.Now()),
			Dimensions: []cwtypes.Dimension{
				{
					Name:  aws.String("Distribution"),
					Value: aws.String(*distribution.Id),
				},
			},
		})

		data = append(data, cwtypes.MetricDatum{
			MetricName: aws.String("InvalidationPathCounter"),
			Unit:       cwtypes.StandardUnitCount,
			Value:      aws.Float64(countPaths),
			Timestamp:  aws.Time(time.Now()),
			Dimensions: []cwtypes.Dimension{
				{
					Name:  aws.String("Distribution"),
					Value: aws.String(*distribution.Id),
				},
			},
		})
	}

	for _, queueItem := range data {
		if dataQueue.QueueFull {
			dataQueue.Flush()
		}
		if err = dataQueue.Add(queueItem); err != nil {
			return fmt.Errorf(err.Error())
		}
	}

	_ = dataQueue.Flush()

	return nil
}

// IsTimeRangeAcceptable will determine if an input time is within
// a given date range. It's intended here to be a frequency of every
// five minutes. timeline
func IsTimeRangeAcceptable(timeBaseline time.Time, timeSource time.Time) (bool, error) {
	// Format is the format of the invalidation entry timestamp from GetInvalidation.
	format := "2006-01-02 15:04:05 +0000 UTC"
	// Parse the time.Time which will be used to compare to timeBaseline.
	timestamp, err := time.Parse(format, fmt.Sprint(timeSource))
	if err != nil {
		return false, err
	}

	// Calculate what is the acceptable age of an invalidation to ingest.
	fiveMinutesAgo := timeBaseline.Add(time.Minute * -5)
	if timestamp.Before(fiveMinutesAgo) {
		return false, fmt.Errorf("input time exceeds constraints")
	}

	return true, nil
}

func main() {

	lambda.Start(Start)

}
