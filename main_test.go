package main

import (
	"os"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	err := Start()
	if err == nil {
		t.FailNow()
	}
}

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

//func TestExecute(t *testing.T) {}
func TestIsTimeRangeAcceptable(t *testing.T) {
	// todo fix time format differences.
	var err error
	format := "2006-01-02 15:04:05 +0000 UTC"

	testTimeOne, err := time.Parse(format, time.Now().Format(time.RFC3339))
	if err != nil {
	}
	testTimeTwo, err := time.Parse(format, time.Now().Add(time.Minute*-5).Format(time.RFC3339))
	if err != nil {
	}
	testTimeThree, err := time.Parse(format, time.Now().Add(time.Hour*-1).Format(time.RFC3339))
	if err != nil {
	}

	_, err = IsTimeRangeAcceptable(time.RFC3339, &testTimeOne)
	if err != nil {
	}
	_, err = IsTimeRangeAcceptable(time.RFC3339, &testTimeTwo)
	if err != nil {
	}
	_, err = IsTimeRangeAcceptable(time.RFC3339, &testTimeThree)
	if err == nil {
	}
}
