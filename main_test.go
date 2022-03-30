package main

import (
	"os"
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
	err := os.Setenv("METRICS_PUSH_DRYRUN", "TRUE")
	if err != nil {
		t.FailNow()
	}

	err = os.Setenv("METRICS_PUSH_DRYRUN", "")
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

	if outcome, _ := IsTimeRangeAcceptable(baselineFormat, sourceFormat.Add(time.Minute*-2)); !outcome {
		t.FailNow()
	}
	if outcome, _ := IsTimeRangeAcceptable(baselineFormat, sourceFormat.Add(time.Hour*-2)); outcome {
		t.FailNow()
	}
	if outcome, _ := IsTimeRangeAcceptable(baselineFormat, sourceFormat.Add(time.Hour*-24)); outcome {
		t.FailNow()
	}
}
