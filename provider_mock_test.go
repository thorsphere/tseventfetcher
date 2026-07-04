// Copyright (c) 2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tseventfetcher_test

// Import standard library packages, tserr and tstrading
import (
	"context" // context for managing request-scoped values and cancellation

	"github.com/thorsphere/tserr" 		// tserr for custom error handling
	"github.com/thorsphere/tstrading"	// tstrading for trading data types
)

// Mock is a mock implementation of the Provider interface for testing purposes.
type Mock struct{}

// GetEvents returns a slice of Event for the specified date range, filtering the events based on the provided period.
func (p *Mock) GetEvents(ctx context.Context, period *tstrading.Period) ([]tstrading.Event, error) {
	// Check if the provider or period is nil, and return an error if so
	if (p == nil) || (period == nil) {
		return nil, tserr.NilPtr()
	}
	// Filter events based on the provided period and return the matching events
	evlist := []tstrading.Event{}
	for _, event := range evs {
		// Check if the event's time is within the specified period, and if so, add it to the list of events to return
		if event.Time.After(period.From) && event.Time.Before(period.To) {
			// If the event is within the period, append it to the list of events to return
			evlist = append(evlist, *event)
		}
	}
	// Return the list of events that match the specified period and nil for the error
	return evlist, nil
}

// MockErr is a mock implementation of the Provider interface that returns an error when GetEvents is called.
type MockErr struct{}

// GetEvents returns an error indicating that the operation is forbidden,
// simulating a failure scenario for testing purposes.
func (p *MockErr) GetEvents(ctx context.Context, period *tstrading.Period) ([]tstrading.Event, error) {
	return nil, tserr.Forbidden("MockErr")
}
