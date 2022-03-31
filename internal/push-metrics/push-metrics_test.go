package push_metrics

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func populateTestData(queue *Queue) {
	for len(queue.Data) < AwsPayloadLimit {
		_ = queue.Add(types.MetricDatum{
			MetricName: aws.String("TestResponse"),
			Value:      aws.Float64(1),
		})
	}
}

func TestAdd(t *testing.T) {

	var queue = &Queue{
		Client:    cloudwatch.Client{},
		Namespace: "dev/null",
		QueueFull: false,
	}

	if len(queue.Data) != 0 {
		t.FailNow()
	}

	populateTestData(queue)
	if len(queue.Data) != AwsPayloadLimit {
		t.FailNow()
	}

	_ = queue.Add(types.MetricDatum{
		MetricName: aws.String("TestResponse"),
		Value:      aws.Float64(1),
	})

	if len(queue.Data) > AwsPayloadLimit {
		t.FailNow()
	}

	if !queue.QueueFull {
		t.FailNow()
	}

}
func TestFlush(t *testing.T) {

	var queue = &Queue{
		Client:    cloudwatch.Client{},
		Namespace: "dev/null",
		QueueFull: false,
	}

	populateTestData(queue)
	if len(queue.Data) != AwsPayloadLimit {
		t.FailNow()
	}

	if !queue.QueueFull {
		t.FailNow()
	}

	queue.Flush()

	if queue.QueueFull {
		t.FailNow()
	}

	if len(queue.Data) != 0 {
		t.FailNow()
	}

}
