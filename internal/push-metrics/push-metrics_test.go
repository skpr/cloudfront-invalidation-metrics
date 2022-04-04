package push_metrics

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	client "cloudfront-invalidation-metrics/internal/cloudwatch"
)

// testSetupQueue will return a testable and functional queue.
func testSetupQueue() *Queue {
	var queue = &Queue{
		Namespace: "dev/null",
		QueueFull: false,
	}
	for len(queue.Data) < AwsPayloadLimit {
		_ = queue.Add(types.MetricDatum{
			MetricName: aws.String("TestResponse"),
			Value:      aws.Float64(1),
		})
	}
	return queue
}

func TestAdd(t *testing.T) {

	queue := testSetupQueue()

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

	client := client.MockCloudWatchClient{}
	queue := testSetupQueue()
	if len(queue.Data) != AwsPayloadLimit {
		t.FailNow()
	}

	if !queue.QueueFull {
		t.FailNow()
	}

	t.Log(queue.QueueFull)
	queue.Flush(client)
	t.Log(queue.QueueFull)
	if queue.QueueFull {
		t.FailNow()
	}

	if len(queue.Data) != 0 {
		t.FailNow()
	}

}
