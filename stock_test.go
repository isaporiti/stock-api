package server

import (
	"testing"
	"time"
)

func TestGetStockHistory(t *testing.T) {
	t.Run("provides stock history for the last N amount of days", func(t *testing.T) {
		t.Parallel()
		stockHistory := GetStockHistory("MSFT", 3)

		wantLength := 3
		gotLength := len(stockHistory)
		if wantLength != gotLength {
			t.Errorf("want stock history of length %d, got %d", wantLength, gotLength)
		}
		now := time.Now()
		wantNewestDate := now.Format("2006-01-02")
		gotNewestDate := stockHistory[0].Date
		if wantNewestDate != gotNewestDate {
			t.Errorf("want today's stock date to be %s, got %s", wantNewestDate, gotNewestDate)
		}
		wantOldestDate := now.AddDate(0, 0, -2).Format("2006-01-02")
		gotOldestDate := stockHistory[2].Date
		if wantOldestDate != gotOldestDate {
			t.Errorf("want oldest stock date to be %s, got %s", wantNewestDate, gotNewestDate)
		}
	})
}

func TestIsKnown(t *testing.T) {
	t.Run("returns false for unknown tickers", func(t *testing.T) {
		got := isKnown("MELI")
		want := false

		if want != got {
			t.Errorf("want MELI to be unknown, but it is")
		}
	})

	t.Run("returns true for known tickers", func(t *testing.T) {
		got := isKnown("AMZN")
		want := true

		if want != got {
			t.Errorf("want AMZN to be known, but it isn't")
		}
	})
}
