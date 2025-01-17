package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	cloudfrontclient "github.com/skpr/cloudfront-invalidation-metrics/internal/aws/cloudfront"
	"github.com/skpr/cloudfront-invalidation-metrics/internal/metrics"
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
		return fmt.Errorf("failed to get AWS client: %w", err)
	}

	dryRun := os.Getenv("CLOUDFRONT_INVALIDATION_METRICS_DRYRUN") != ""

	client, err := metrics.New(cloudwatch.NewFromConfig(cfg), CloudWatchNamespace, dryRun)
	if err != nil {
		return fmt.Errorf("failed to setup client: %w", err)
	}

	return Execute(ctx, cloudfront.NewFromConfig(cfg), client)
}

// Execute will execute the given API calls against the input Clients.
func Execute(ctx context.Context, clientCloudFront cloudfrontclient.ClientInterface, client metrics.ClientInterface) error {
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
			}
		}

		err = client.Add(types.MetricDatum{
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
		if err != nil {
			return fmt.Errorf("failed to push metric: InvalidationRequest: %w", err)
		}

		err = client.Add(types.MetricDatum{
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
		if err != nil {
			return fmt.Errorf("failed to push metric: InvalidationPathCounter: %w", err)
		}
	}

	return client.Flush()
}

func main() {
	lambda.Start(Start)
}
