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
	baselineFormat, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", "2006-01-02 15:04:05 +0000 UTC")
	sourceFormat, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", "2006-01-02 15:04:05 +0000 UTC")

	// Time.Now() - 2 minutes should return true (not false)
	if outcome, _ := IsTimeRangeAcceptable(baselineFormat, sourceFormat.Add(time.Minute*-2)); !outcome {
		t.FailNow()
	}
	// Time.Now() - 2 hours should return false (not true)
	if outcome, _ := IsTimeRangeAcceptable(baselineFormat, sourceFormat.Add(time.Hour*-2)); outcome {
		t.FailNow()
	}
	// Time.Now() - 24 hours should return false (not true)
	if outcome, _ := IsTimeRangeAcceptable(baselineFormat, sourceFormat.Add(time.Hour*-24)); outcome {
		t.FailNow()
	}
}
