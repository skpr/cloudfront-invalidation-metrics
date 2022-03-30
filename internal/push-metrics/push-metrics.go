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
	Client    cloudwatch.Client     `json:"client"`
	Namespace string                `json:"namespace"`
	QueueFull bool                  `json:"full"`
	Data      []cwtypes.MetricDatum `json:"data"`
}

func (Queue *Queue) Add(data cwtypes.MetricDatum) error {
	if len(Queue.Data) < AwsPayloadLimit {
		Queue.QueueFull = false
		Queue.Data = append(Queue.Data, data)
	} else {
		Queue.QueueFull = true
		return fmt.Errorf("error adding to queue: queue size is full")
	}

	return nil
}

func (Queue *Queue) Flush() ([]cwtypes.MetricDatum, error) {
	if dryrun := os.Getenv("METRICS_PUSH_DRYRUN"); dryrun != "" {
		return Queue.Data, nil
	}

	_, err := Queue.Client.PutMetricData(context.Background(), &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(Queue.Namespace),
		MetricData: Queue.Data,
	})

	if err != nil {
		return Queue.Data, fmt.Errorf(err.Error())
	}

	return []cwtypes.MetricDatum{}, nil
}
