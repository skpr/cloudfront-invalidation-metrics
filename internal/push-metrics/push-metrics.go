package push_metrics

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

const (
	//AwsPayloadLimit is the maximum quality for a data-set to contain
	// before AWS will reject the payload.
	AwsPayloadLimit = 20
)

type Queue struct {
	QueueInterface
	Client    cloudwatch.Client     `json:"client"`
	Namespace string                `json:"namespace"`
	QueueFull bool                  `json:"full"`
	Data      []cwtypes.MetricDatum `json:"data"`
}

type QueueInterface interface {
	Add(datum cwtypes.MetricDatum) error
	Flush() error
}

func (Queue *Queue) Add(data cwtypes.MetricDatum) error {
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

func (Queue *Queue) Flush() error {
	var err error
	if dryrun := os.Getenv("METRICS_PUSH_DRYRUN"); dryrun != "" {
		return nil
	}
	_, err = Queue.Client.PutMetricData(context.Background(), &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(Queue.Namespace),
		MetricData: Queue.Data,
	})

	Queue.Data = []cwtypes.MetricDatum{}
	Queue.QueueFull = len(Queue.Data) == AwsPayloadLimit

	return err
}
