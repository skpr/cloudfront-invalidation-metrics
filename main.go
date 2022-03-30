package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"time"

	//"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

const (
	CloudWatchNamespace = "Skpr/CloudFront"
)

type clients struct {
	CloudFront *cloudfront.Client `json:"cloudfront"`
	CloudWatch *cloudwatch.Client `json:"cloudwatch"`
}

// Lambda is an exported abstraction so that the application
// can be used externally from Skpr or Lambda by writing your own
// main function which calls this.
func Lambda(ctx context.Context, clients clients) error {

	distributions, err := clients.CloudFront.ListDistributions(ctx, &cloudfront.ListDistributionsInput{})
	if err != nil {
		return fmt.Errorf("failed to get CloudFront distibution list: %w", err)
	}

	var data []cwtypes.MetricDatum

	for _, distribution := range distributions.DistributionList.Items {
		invalidations, err := clients.CloudFront.ListInvalidations(ctx, &cloudfront.ListInvalidationsInput{
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
			invalidationDetail, _ := clients.CloudFront.GetInvalidation(ctx, &cloudfront.GetInvalidationInput{
				DistributionId: distribution.Id,
				Id:             invalidation.Id,
			})

			acceptable, err := isTimeRangeAcceptable(invalidationDetail.Invalidation.CreateTime)
			if err != nil {
				return err
			}

			if !acceptable {
				break
			}

			countInvalidations++

			countPaths = countPaths + float64(*invalidationDetail.Invalidation.InvalidationBatch.Paths.Quantity)
		}

		// 20 item limit per payload, if the limit is met or exceeded, offload now.
		if len(data) >= 20 {
			if err = pushMetrics(ctx, clients.CloudWatch, &data); err != nil {
				fmt.Errorf(err.Error())
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

	if err = pushMetrics(ctx, clients.CloudWatch, &data); err != nil {
		fmt.Errorf(err.Error())
	}

	return nil
}

// pushMetrics will push the metrics found in the input and return the error from that.
// It will also empty out the data if completed successfully so that the variable can
// be repopulated as needed.
func pushMetrics(ctx context.Context, clientCloudWatch *cloudwatch.Client, data *[]cwtypes.MetricDatum) error {
	if dryrun := os.Getenv("LAMBDA_DRYRUN"); dryrun != "" {
		return nil
	}
	_, err := clientCloudWatch.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(CloudWatchNamespace),
		MetricData: *data,
	})
	// If no error was found, remove the data from memory.
	// This is so that more can be re-queued if necessary.
	if err == nil {
		data = &[]cwtypes.MetricDatum{}
	}
	return err
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

func main() {

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Errorf("failed to get AWS client config: %w", err)
	}

	clientCloudWatch := cloudwatch.NewFromConfig(cfg)
	clientCloudFront := cloudfront.NewFromConfig(cfg)

	lambda.Start(Lambda(ctx, clients{
		CloudFront: clientCloudFront,
		CloudWatch: clientCloudWatch,
	}))
}
