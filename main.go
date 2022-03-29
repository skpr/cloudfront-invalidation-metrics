package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

const (
	CloudWatchNamespace = "Skpr/CloudFront"
)

// StartLambda is an exported abstraction so that the application
// can be used externally from Skpr or Lambda by writing your own
// main function which calls this.
func StartLambda(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get AWS client config: %w", err)
	}

	clientCloudWatch := cloudwatch.NewFromConfig(cfg)
	clientCloudFront := cloudfront.NewFromConfig(cfg)

	// @todo, Run a function that passes context and these clients (interfaces ftw).

	distributions, err := clientCloudFront.ListDistributions(context.TODO(), &cloudfront.ListDistributionsInput{})
	if err != nil {
		return fmt.Errorf("failed to get CloudFront distibution list: %w", err)
	}

	var data []cwtypes.MetricDatum

	for _, distribution := range distributions.DistributionList.Items {
		invalidations, err := clientCloudFront.ListInvalidations(ctx, &cloudfront.ListInvalidationsInput{
			DistributionId: distribution.Id,
			// @todo, Descending order + break once past 5 minutes????
		})
		if err != nil {
			return err
		}

		var (
			countInvalidations = 0
			countPaths         = 0
		)

		for _, invalidation := range invalidations.InvalidationList.Items {
			invalidationDetail, _ := clientCloudFront.GetInvalidation(ctx, &cloudfront.GetInvalidationInput{
				DistributionId: distribution.Id,
				Id:             invalidation.Id,
			})

			acceptable, err := isTimeRangeAcceptable(invalidationDetail.Invalidation.CreateTime)
			if err != nil {
				return err
			}

			if !acceptable {
				// @todo, Consider break?
				continue
			}

			countInvalidations++

			countPaths = countPaths + *invalidationDetail.Invalidation.InvalidationBatch.Paths.Quantity
		}

		data = append(data, cwtypes.MetricDatum{
			MetricName: aws.String("InvalidationRequest"),
			Unit:       cwtypes.StandardUnitCount,
			Value:      aws.Float64(1),
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
			Value:      aws.Float64(1),
			Timestamp:  aws.Time(time.Now()),
			Dimensions: []cwtypes.Dimension{
				{
					Name:  aws.String("Distribution"),
					Value: aws.String(*distribution.Id),
				},
			},
		})
	}

	// @todo, Determine if we need to account for limits.

	_, err = clientCloudWatch.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(CloudWatchNamespace),
		MetricData: data,
	})

	return nil
}

func main() {
	lambda.Start(StartLambda)
}

// isTimeRangeAcceptable will determine if an input time is within
// a given date range. It's intended here to be a frequency of every
// five minutes.
func isTimeRangeAcceptable(timeSource *time.Time) (bool, error) {
	format := "2006-01-02 15:04:05 +0000 UTC"
	timestamp, err := time.Parse(format, fmt.Sprint(timeSource))
	if err != nil {
		return false, err
	}

	fiveMinutesAgo := time.Now().Add(time.Minute * -5)
	if timestamp.Before(fiveMinutesAgo) {
		return false, errors.New("input time exceeds constraints")
	}

	return true, nil
}
