package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

// Params is the configuration for the app/lambda.
type Params struct {
	NameSpace string
	DryRun    bool
}

// ErrorData is a simple and controlled way to standardize output from the application.
type ErrorData struct {
	Action       string    `json:"action"`
	Distribution string    `json:"distribution"`
	Message      string    `json:"message"`
	Status       string    `json:"status"`
	Value        string    `json:"value"`
	Time         time.Time `json:"time"`
}

// LambdaApp is the Lambda construct with any configuration.
type LambdaApp struct {
	Params        Params
	Distributions *cloudfront.ListDistributionsOutput
}

// printMessage prints a message in a consistent format for AWS CloudWatch.
func printMessage(action, message, status, distribution, value string) {
	data := ErrorData{
		Action:       action,
		Distribution: distribution,
		Status:       status,
		Message:      message,
		Value:        value,
		Time:         time.Now(),
	}
	output, _ := json.Marshal(data)
	fmt.Println(string(output))
}

// getDistributions will run all ListDistributions actions and return the response.
func (app *LambdaApp) getDistributions() *cloudfront.ListDistributionsOutput {
	printMessage("ListDistributions", "", "Pending", "", "")

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	svc := cloudfront.NewFromConfig(cfg)

	result, err := svc.ListDistributions(context.TODO(), &cloudfront.ListDistributionsInput{})
	if err != nil {
		printMessage("ListDistributions", "", "Error", "", "")
		return &cloudfront.ListDistributionsOutput{}
	}
	printMessage("ListDistributions", "", "Success", "", "")
	return result
}

// getInvalidations will run all ListInvalidations actions and return the response.
func (app *LambdaApp) getInvalidations(distribution string) (*cloudfront.ListInvalidationsOutput, error) {

	printMessage("ListInvalidations", "", "Pending", distribution, "")

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	svc := cloudfront.NewFromConfig(cfg)

	input := &cloudfront.ListInvalidationsInput{
		DistributionId: aws.String(distribution),
	}

	result, err := svc.ListInvalidations(context.TODO(), input)
	if err != nil {
		printMessage("ListInvalidations", "", "Error", "", "")
		return &cloudfront.ListInvalidationsOutput{}, err
	}

	printMessage("ListInvalidations", "", "Success", "", "")

	return result, nil
}

// prepareToPush will structure the data ready to send to CloudWatch.
func (app *LambdaApp) prepareToPush(distribution *cftypes.DistributionSummary, input *cloudfront.ListInvalidationsOutput) error {

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	svc := cloudwatch.NewFromConfig(cfg)

	for _, item := range input.InvalidationList.Items {

		var dataSet []cwtypes.MetricDatum

		if err := isTimeRangeAcceptable(item.CreateTime); err != nil {
			return nil
		}

		// Extract any fields we need
		var projectEnvironment string
		var projectName string
		for _, v := range distribution.Aliases.Items {
			if strings.HasSuffix(v, "skpr.io") || strings.HasSuffix(v, "skpr.dev") || strings.HasSuffix(v, "skpr.live") {
				projectEnvironment = strings.Split(distribution.Aliases.Items[0], ".")[0]
				projectName = strings.Split(distribution.Aliases.Items[0], ".")[1]
			}
		}

		if projectName == "" || projectEnvironment == "" {
			printMessage("", "Cannot parse project information", "Warning", *distribution.Id, "")
			continue
		}

		dataSet = append(dataSet, cwtypes.MetricDatum{
			MetricName: aws.String("InvalidationRequest"),
			Unit:       cwtypes.StandardUnitCount,
			Value:      aws.Float64(1),
			Timestamp:  aws.Time(time.Now()),
			Dimensions: []cwtypes.Dimension{
				{
					Name:  aws.String("Distribution"),
					Value: aws.String(*distribution.Id),
				},
				{
					Name:  aws.String("Project"),
					Value: aws.String(projectName),
				},
				{
					Name:  aws.String("Environment"),
					Value: aws.String(projectEnvironment),
				},
			},
		})

		if app.Params.DryRun {
			printMessage("", "Ignoring InvalidationRequest data-set from push cycle - dry-run is enabled", "Warning", *distribution.Id, "")
			continue
		}

		if err := app.sendData(svc, &dataSet); err == nil {
			printMessage("PutMetricData", "InvalidationRequest", "Success", *distribution.Id, "1")
		} else {
			printMessage("PutMetricData", "InvalidationRequest", "Error", *distribution.Id, "1")
			return err
		}
	}

	return nil
}

func (app *LambdaApp) sendInvalidationCount(distribution *cftypes.DistributionSummary) error {

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	svc := cloudwatch.NewFromConfig(cfg)
	svct := cloudfront.NewFromConfig(cfg)
	pathCount := float64(0)
	invalidations, err := svct.ListInvalidations(context.TODO(), &cloudfront.ListInvalidationsInput{
		DistributionId: distribution.Id,
	})
	if err != nil {
		return err
	}

	for _, invalidation := range invalidations.InvalidationList.Items {

		invalidationDetail, _ := svct.GetInvalidation(context.TODO(), &cloudfront.GetInvalidationInput{
			DistributionId: distribution.Id,
			Id:             invalidation.Id,
		})

		if err := isTimeRangeAcceptable(invalidationDetail.Invalidation.CreateTime); err == nil {
			pathCount++
		} else {
			continue
		}
	}

	var dataSet []cwtypes.MetricDatum

	dataSet = append(dataSet, cwtypes.MetricDatum{
		MetricName: aws.String("InvalidationPathCounter"),
		Unit:       cwtypes.StandardUnitCount,
		Value:      aws.Float64(pathCount),
		Timestamp:  aws.Time(time.Now()),
		Dimensions: []cwtypes.Dimension{
			{
				Name:  aws.String("Distribution"),
				Value: aws.String(*distribution.Id),
			},
		},
	})

	if app.Params.DryRun {
		printMessage("", "Ignoring InvalidationPathCounter data-set from push cycle - dry-run is enabled", "Warning", *distribution.Id, "")
		return nil
	}

	if err := app.sendData(svc, &dataSet); err == nil {
		printMessage("PutMetricData", "InvalidationPathCounter", "Success", *distribution.Id, fmt.Sprint(pathCount))
	} else {
		printMessage("PutMetricData", "InvalidationPathCounter", "Error", *distribution.Id, fmt.Sprint(pathCount))
		return err
	}
	return nil
}

// sendData will push the data to CloudWatch
func (app *LambdaApp) sendData(svc *cloudwatch.Client, data *[]cwtypes.MetricDatum) error {
	_, err := svc.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(app.Params.NameSpace),
		MetricData: *data,
	})

	return err
}

// isTimeRangeAcceptable will determine if an input time is within
// a given date range. It's intended here to be a frequency of every
// five minutes.
func isTimeRangeAcceptable(timeSource *time.Time) error {
	format := "2006-01-02 15:04:05 +0000 UTC"
	timestamp, err := time.Parse(format, fmt.Sprint(timeSource))
	if err != nil {
		return err
	}

	fiveMinutesAgo := time.Now().Add(time.Minute * -5)
	if timestamp.Before(fiveMinutesAgo) {
		return errors.New("input time exceeds constraints")
	}

	return nil
}

// StartLambda is an exported abstraction so that the application
// can be used externally from Skpr or Lambda by writing your own
// main function which calls this.
func StartLambda(app LambdaApp) error {
	app.Distributions = app.getDistributions()
	for _, distro := range app.Distributions.DistributionList.Items {
		_ = app.sendInvalidationCount(&distro)
		if invalidations, err := app.getInvalidations(*distro.Id); err == nil {
			if len(invalidations.InvalidationList.Items) > 0 {
				if err := app.prepareToPush(&distro, invalidations); err != nil {
					return err
				}
			} else {
				printMessage("PutMetricData", "", "Warning", "", "")
			}
		} else {
			return err
		}
	}
	printMessage("", "Application completed", "Success", "", "")
	return nil
}

func main() {

	app := &LambdaApp{
		Params: Params{
			// AWS CloudWatch NameSpace to store metrics in
			NameSpace: "Skpr/CloudFront",
			// For testing of authentication reasons. we'll not allow pushing of metrics.
			DryRun: false,
		},
	}

	lambda.Start(StartLambda(*app))
}
