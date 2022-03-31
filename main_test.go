package main

import (
	"cloudfront-invalidation-metrics/internal/mock/cloudfront"
	"cloudfront-invalidation-metrics/internal/mock/cloudwatch"
	"context"
	"testing"
	"time"
)

// TestStart will test the Start function for a nil value completion.
func TestStart(t *testing.T) {
	err := Start()
	if err == nil {
		t.FailNow()
	}
}

// TestExecute tests the guts of the Lambda.
func TestExecute(t *testing.T) {
	ctx := context.Background()
	cloudFrontClient := cloudfront.NewMockCloudFrontClient()
	cloudWatchClient := cloudwatch.NewMockCloudwatchClient()

	// todo.... how?
	err := Execute(ctx, *cloudFrontClient, *cloudWatchClient)
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
