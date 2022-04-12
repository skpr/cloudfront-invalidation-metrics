package main

import (
	"testing"
	"time"
)

// TestIsTimeRangeAcceptable will test if a given input will return a positive
// or negative response from the IsTimeRangeAcceptable function. This regulates
// a specific timeframe around what metrics should be ingested.
func TestIsTimeRangeAcceptable(t *testing.T) {
	// Compare our values against Now() - X
	// X being the value of time taken from time.Now() to determine if the time
	// is within a specific time frame.
	sourceFormat := time.Now()

	// Time.Now() - An invalidation 2 minutes old must pass
	if _, err := IsTimeRangeAcceptable(sourceFormat.Add(time.Minute * -2)); err != nil {
		t.FailNow()
	}
	// Time.Now() - An invalidation 2 hours old must not pass
	if _, err := IsTimeRangeAcceptable(sourceFormat.Add(time.Hour * -2)); err == nil {
		t.FailNow()
	}
	// Time.Now() - An invalidation 2 days old must not pass
	if _, err := IsTimeRangeAcceptable(sourceFormat.Add((time.Hour * 24) * -2)); err == nil {
		t.FailNow()
	}
}
