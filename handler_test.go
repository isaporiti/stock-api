package server

import "testing"

func TestHandleStockHistory(t *testing.T) {
	t.Run("returns error if unknown ticker provided", func(t *testing.T) {
		t.Parallel()
		_, err := HandleStockHistory("UNKNOWN")
		if err == nil {
			t.Error("want an error but got nil")
		}
	})

	t.Run("returns stock history for provided ticker", func(t *testing.T) {
		t.Parallel()
		stockHistory, err := HandleStockHistory("MSFT")
		if err != nil {
			t.Fatalf("want no error but got %s", err.Error())
		}
		if len(stockHistory) == 0 {
			t.Errorf("want non-empty stock history, got empty")
		}
	})
}

func TestHandleUserStock(t *testing.T) {
	t.Run("returns even tickers for testA", func(t *testing.T) {
		t.Parallel()
		got := HandleUserStock("testA")
		want := []string{
			"AAPL",
			"GOOG",
			"FB",
			"NVDA",
			"BABA",
			"WMT",
			"PYPL",
			"ADBE",
			"V",
			"CRM",
		}
		if len(want) != len(got) {
			t.Fatal("didn't get the expected amount of user stock")
		}
		for i := range got {
			if want[i] != got[i].Ticker {
				t.Errorf("want ticker %s, got %s", want[i], got[i].Ticker)
			}
		}
	})

	t.Run("returns odd tickers for testB", func(t *testing.T) {
		t.Parallel()
		got := HandleUserStock("testB")
		want := []string{
			"MSFT",
			"AMZN",
			"TSLA",
			"JPM",
			"JNJ",
			"PG",
			"DIS",
			"PFE",
			"MA",
			"NFLX",
		}
		if len(want) != len(got) {
			t.Fatal("didn't get the expected amount of user stock")
		}
		for i := range got {
			if want[i] != got[i].Ticker {
				t.Errorf("want ticker %s, got %s", want[i], got[i].Ticker)
			}
		}
	})
}
