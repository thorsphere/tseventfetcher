// Copyright (c) 2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tseventfetcher

// Import standard library packages, tserr and tstrading
import (
	"context" // context for managing request-scoped values, cancellation signals, and deadlines
	"strings" // strings for building output

	"github.com/thorsphere/tserr"     // tserr for custom error handling
	"github.com/thorsphere/tstrading" // tstrading for trading data types
)

// Provider defines the interface for fetching economic events from a data source.
type Provider interface {
	// GetEvents returns a slice of Event for the specified date range.
	GetEvents(ctx context.Context, period *tstrading.Period) ([]*tstrading.Event, error)
}

// FetchEvents uses the provided Provider to fetch events for the specified period.
// If the provider or period is nil, it returns an error.
// Otherwise, it returns the fetched events and nil for the error.
func FetchEvents(ctx context.Context, p Provider, period *tstrading.Period) ([]*tstrading.Event, error) {
	// Check if the provider or period is nil, and return an error if so
	if (p == nil) || (period == nil) {
		return nil, tserr.NilPtr()
	}
	// Fetch events using the provider's GetEvents method
	events, err := p.GetEvents(ctx, period)
	// If there is an error fetching events, return the error
	if err != nil {
		return nil, err
	}
	// Filter out any nil events from the result
	filtered := events[:0]
	// Iterate over the events and append only non-nil ones
	for _, e := range events {
		if e != nil {
			// If the event is non-nil, append it to the filtered slice
			filtered = append(filtered, e)
		}
	}
	// Return the filtered events and nil for the error
	return filtered, nil
}

// PrintEvents prints the details of each event in the provided slice of events.
// If the slice is nil or empty, it returns an empty string.
// Otherwise, it builds a string representation of each event using the String method of the Event struct,
// and returns the combined string.
func PrintEvents(events []*tstrading.Event) string {
	// If there are no events, return an empty string
	if len(events) == 0 {
		return ""
	}
	// If there are events, iterate over each event and print its details using the String method of the Event struct.
	var out strings.Builder
	for _, event := range events {
		// Skip nil events
		if event == nil {
			continue
		}
		// Write the string representation of the event to the output
		out.WriteString(event.String())
		// Write a newline character to the output
		out.WriteString("\n")
	}
	// Return the string representation of the events
	return out.String()
}
