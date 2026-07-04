package tseventfetcher_test

import (
	"context"
	"testing"
	"time"

	"github.com/thorsphere/tseventfetcher"
	"github.com/thorsphere/tstrading"
)

func TestFMPProvider(t *testing.T) {

	period := &tstrading.Period{
		From: time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2026, 6, 10, 23, 59, 59, 0, time.UTC),
	}

	p := tseventfetcher.NewBLSProvider("669650da92d541279454149cb828be1d")
	events, err := p.GetEvents(context.Background(), period)
	if err != nil {
		t.Fatal(err)
	}
	tseventfetcher.PrintEvents(events)
}
