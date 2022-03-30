package push_metrics

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go/aws"
	"testing"
)

var queue = &Queue{
	Client:    cloudwatch.Client{},
	Namespace: "dev/null",
	QueueFull: false,
}

func TestAdd(t *testing.T) {

	if len(queue.Data) != 0 {
		t.FailNow()
	}

	for len(queue.Data) < AwsPayloadLimit {
		err := queue.Add(cwtypes.MetricDatum{
			MetricName: aws.String("TestResponse"),
			Value:      aws.Float64(1),
		})
		if err != nil {
			t.FailNow()
		}
	}
	if len(queue.Data) != AwsPayloadLimit {
		t.FailNow()
	}

	err := queue.Add(cwtypes.MetricDatum{
		MetricName: aws.String("TestResponse"),
		Value:      aws.Float64(1),
	})

	if err == nil {
		t.FailNow()
	}

	if len(queue.Data) > AwsPayloadLimit {
		t.FailNow()
	}

	if !queue.QueueFull {
		t.FailNow()
	}

}
func TestFlush(t *testing.T) {

	if len(queue.Data) != AwsPayloadLimit {
		t.FailNow()
	}

	// Expecting a failure here for now.
	data, err := queue.Flush()

	if err == nil {
		t.FailNow()
	}

	if len(data) != AwsPayloadLimit {
		t.FailNow()
	}

}
