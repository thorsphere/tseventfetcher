// Copyright (c) 2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public Licence v3.0
// that can be found in the LICENSE file.
package tseventfetcher_test

// Import standard library packages, testing, tserr, tsfio, tstrading and tseventfetcher.
import (
	"context" // context
	"testing" // testing
	"time"    // time

	"github.com/thorsphere/tserr"          // tserr
	"github.com/thorsphere/tseventfetcher" // tseventfetcher
	"github.com/thorsphere/tsfio"          // tsfio
	"github.com/thorsphere/tstrading"      // tstrading
)

var (
	// Define some sample events for testing purposes
	evNfp *tstrading.Event = &tstrading.Event{
		Name:     "Non-Farm Payrolls",
		Time:     time.Date(2024, 7, 5, 8, 30, 0, 0, time.UTC),
		Country:  "US",
		Actual:   new(200.0),
		Estimate: new(180.0),
		Previous: new(150.0),
		Unit:     "K",
		Impact:   tstrading.ImpactHigh,
		Source:   "Bureau of Labor Statistics",
	}
	evGdp24 *tstrading.Event = &tstrading.Event{
		Name:     "GDP Growth Rate",
		Time:     time.Date(2024, 7, 10, 8, 30, 0, 0, time.UTC),
		Country:  "US",
		Actual:   new(3.5),
		Estimate: new(3.0),
		Previous: new(2.8),
		Unit:     "%",
		Impact:   tstrading.ImpactMedium,
		Source:   "Bureau of Economic Analysis",
	}
	evGdp30 *tstrading.Event = &tstrading.Event{
		Name:     "GDP Growth Rate",
		Time:     time.Date(2030, 7, 10, 8, 30, 0, 0, time.UTC),
		Country:  "US",
		Actual:   nil,
		Estimate: nil,
		Previous: nil,
		Unit:     "%",
		Impact:   tstrading.ImpactLow,
		Source:   "Bureau of Economic Analysis",
	}
	// Define a slice of events for testing purposes
	evs []*tstrading.Event = []*tstrading.Event{
		evNfp,
		evGdp24,
		evGdp30,
	}
	// Define a sample period for testing purposes
	per *tstrading.Period = &tstrading.Period{
		From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2025, 7, 31, 23, 59, 59, 0, time.UTC),
	}
)

// TestFetchEvents tests the FetchEvents function by using the Mock provider
// to fetch events for a specified period and comparing the output to a golden file.
func TestFetchEvents(t *testing.T) {
	// Create a new instance of the Mock provider
	p := &Mock{}
	// Use the FetchEvents function to fetch events for the specified period using the Mock provider
	events, err := tseventfetcher.FetchEvents(context.Background(), p, per)
	// If there is an error fetching events, fail the test with the error message
	if err != nil {
		t.Fatal(err)
	}
	// Use the PrintEvents function to get a formatted string representation of the fetched events
	out := tseventfetcher.PrintEvents(events)
	// Compare the output to a golden file using the EvalGoldenFile function from the tsfio package,
	// and if there is an error, fail the test with the error message
	if e := tsfio.EvalGoldenFile(&tsfio.Testcase{Name: "fetchevents", Data: out}); e != nil {
		t.Fatal(e)
	}
}

// TestFetchEventsNilProvider tests the FetchEvents function with a nil provider and expects an error.
func TestFetchEventsNilProvider(t *testing.T) {
	// Use the FetchEvents function with a nil provider and a valid period, and expect an error
	_, err := tseventfetcher.FetchEvents(context.Background(), nil, per)
	// If there is no error, fail the test with a message indicating that a nil provider should have caused an error
	if err == nil {
		t.Fatal(tserr.NilExpected("FetchEvents for nil provider"))
	}
}

// TestFetchEventsNilPeriod tests the FetchEvents function with a nil period and expects an error.
func TestFetchEventsNilPeriod(t *testing.T) {
	// Create a new instance of the Mock provider
	p := &Mock{}
	// Use the FetchEvents function with a valid provider and a nil period, and expect an error
	_, err := tseventfetcher.FetchEvents(context.Background(), p, nil)
	if err == nil {
		t.Fatal(tserr.NilExpected("FetchEvents for nil period"))
	}
}

// TestFetchEventsProviderError tests the FetchEvents function with a provider
// that returns an error and expects the error to be propagated.
func TestFetchEventsProviderError(t *testing.T) {
	// Create a new instance of the MockErr provider, which returns an error when GetEvents is called
	p := &MockErr{}
	// Use the FetchEvents function with the MockErr provider and a valid period, and expect an error
	events, err := tseventfetcher.FetchEvents(context.Background(), p, per)
	// If there is no error, fail the test with a message indicating that an error was expected
	if err == nil {
		t.Fatal(tserr.NilFailed("FetchEvents"))
	}
	// If the error message is correct, the test passes, and we can also check that no events were returned
	l := int64(len(events))
	if l != 0 {
		t.Fatal(tserr.EqualInt(&tserr.EqualIntArgs{Var: "Length of events", Actual: l, Want: 0}))
	}
}
