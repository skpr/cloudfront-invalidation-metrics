package push_metrics

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	cloudwatchclient "cloudfront-invalidation-metrics/internal/cloudwatch"
)

const (
	// AwsPayloadLimit is the maximum quality for a data-set to contain
	// before AWS will reject the payload.
	AwsPayloadLimit = 20
)

type Queue struct {
	QueueInterface
	Namespace string              `json:"namespace"`
	QueueFull bool                `json:"full"`
	Data      []types.MetricDatum `json:"data"`
}

type QueueInterface interface {
	Add(datum types.MetricDatum) error
	Flush() error
}

func (Queue *Queue) Add(data types.MetricDatum) error {
	if len(Queue.Data) < AwsPayloadLimit {
		Queue.Data = append(Queue.Data, data)
		if len(Queue.Data) == AwsPayloadLimit {
			Queue.QueueFull = true
		}
	} else {
		Queue.QueueFull = true
		return fmt.Errorf("error adding to queue: queue size is full")
	}

	return nil
}

func (Queue *Queue) Flush(clientCLoudWatch cloudwatchclient.CloudWatchClientInterface) error {
	var err error

	if Queue.Data == nil {
		Queue.Data = []types.MetricDatum{}
	}

	if dryrun := os.Getenv("METRICS_PUSH_DRYRUN"); dryrun != "" {
		return nil
	}

	_, err = clientCLoudWatch.PutMetricData(context.Background(), &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(Queue.Namespace),
		MetricData: Queue.Data,
	})

	if err == nil {
		Queue.Data = []types.MetricDatum{}
	}
	
	Queue.QueueFull = len(Queue.Data) == AwsPayloadLimit

	return err
}
