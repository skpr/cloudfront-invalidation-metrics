package logs

import (
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/stretchr/testify/assert"

	client "cloudfront-invalidation-metrics/internal/aws/cloudwatchlogs"
)

func TestFlush(t *testing.T) {
	cw := &client.MockCloudWatchLogsClient{}

	client, err := New(cw, "dev/null", false)
	assert.NoError(t, err)

	output, err := captureOutput(func() error {
		err := client.Flush(nil, "dev/test-group", "test", types.InputLogEvent{
			Message:   aws.String("test"),
			Timestamp: aws.Int64(0),
		})
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, "Creating log group: dev/test-group\nCreating log stream: test\nFlushed log events to dev/test-group/test\n", output)
}

func captureOutput(f func() error) (string, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out), err
}
