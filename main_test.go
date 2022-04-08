package main

import (
	"context"
	"testing"
	"time"

	"cloudfront-invalidation-metrics/internal/cloudfront"
	"cloudfront-invalidation-metrics/internal/cloudwatch"
)

// TestStart will test the Start function for a nil value completion.
func TestStart(t *testing.T) {
	err := Start(context.TODO())
	if err == nil {
		t.FailNow()
	}
}

// TestExecute tests the guts of the Lambda.
func TestExecute(t *testing.T) {
	ctx := context.Background()
	clientCloudFront := cloudfront.MockCloudFrontClient{}
	clientCloudWatch := cloudwatch.MockCloudWatchClient{}

	err := Execute(ctx, clientCloudFront, clientCloudWatch)
	if err != nil {
		t.FailNow()
	}
}

// TestIsTimeRangeAcceptable will test if a given input will return a positive
// or negative response from the IsTimeRangeAcceptable function. This regulates
// a specific timeframe around what metrics should be ingested.
func TestIsTimeRangeAcceptable(t *testing.T) {
	// Compare our values against Now() - X
	// X being the value of time taken from time.Now() to determine if the time
	// is within a specific time frame.
	sourceFormat := time.Now()

	// Time.Now() - 2 minutes should pass
	if _, err := IsTimeRangeAcceptable(sourceFormat.Add(time.Minute * -2)); err != nil {
		t.FailNow()
	}
	// Time.Now() - 2 hours should not pass
	if _, err := IsTimeRangeAcceptable(sourceFormat.Add(time.Hour * -2)); err == nil {
		t.FailNow()
	}
	// Time.Now() - 24 days should not pass
	if _, err := IsTimeRangeAcceptable(sourceFormat.Add((time.Hour * 24) * -2)); err == nil {
		t.FailNow()
	}
}
