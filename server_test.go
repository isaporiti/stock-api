package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestTickerHistory_unknown_ticker(t *testing.T) {
	t.Run("responds Not Found when requested an unknown ticker", func(t *testing.T) {
		t.Parallel()
		request := newStockHistoryRequest(t, "UNKNOWN")
		response := httptest.NewRecorder()
		server, err := NewStockServer(WithStockHistoryGetter(getStubStockHistory))
		if err != nil {
			t.Fatalf("could not create server: %s", err)
		}
		server.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusNotFound

		if want != got {
			t.Errorf("want status code %d, got %d", want, got)
		}
	})
}

func newStockHistoryRequest(t *testing.T, ticker string) *http.Request {
	urlPath := fmt.Sprintf("/tickers/%s/history", ticker)
	request, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		t.Fatal("could not create request", err)
	}
	request.SetBasicAuth("testA", "")
	return request
}

func TestTickerHistory_known_ticker(t *testing.T) {
	for _, ticker := range knownTickers {
		testName := fmt.Sprintf("responds OK when requested %s ticker", ticker)
		server, err := NewStockServer(WithStockHistoryGetter(getStubStockHistory))
		if err != nil {
			t.Fatalf("could not create server: %s", err)
		}

		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			request := newStockHistoryRequest(t, ticker)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			got := response.Result().StatusCode
			want := http.StatusOK

			if want != got {
				t.Errorf("want status code %d, got %d", want, got)
			}
		})
	}

	t.Run("responds with list of prices", func(t *testing.T) {
		t.Parallel()
		var err error
		request := newStockHistoryRequest(t, "MSFT")
		response := httptest.NewRecorder()
		server, err := NewStockServer(WithStockHistoryGetter(getStubStockHistory))
		if err != nil {
			t.Fatalf("could not create server: %s", err)
		}

		server.ServeHTTP(response, request)

		want := getStubStockHistory("MSFT", 2)
		var got []Stock
		json.NewDecoder(response.Body).Decode(&got)
		if len(got) != len(want) {
			t.Fatal("didn't get the expected amount of stock")
		}
		for i := range want {
			if !reflect.DeepEqual(want[i], got[i]) {
				t.Errorf("want %v stock, got %v", want[i], got[i])
			}
		}
	})
}

func getStubStockHistory(ticker string, length int) []Stock {
	return []Stock{
		{Date: "2023-07-19", Price: 6.1},
		{Date: "2023-07-19", Price: 10.5},
	}
}

func TestUserStocks(t *testing.T) {
	type testCase struct {
		user string
		want []string
	}
	date := time.Now().Format("2006-01-02")
	testCases := []testCase{
		{
			user: "testA",
			want: []string{"MSFT", "AAPL", "AMZN"},
		},
		{
			user: "testB",
			want: []string{"FB", "NFLX"},
		},
	}
	for _, test := range testCases {
		testName := fmt.Sprintf("returns stocks for user %s", test.user)
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			var err error
			request, err := http.NewRequest(http.MethodGet, "/tickers", nil)
			request.SetBasicAuth(test.user, "")
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			stockServer, err := NewStockServer()
			if err != nil {
				t.Fatal(err)
			}

			stockServer.ServeHTTP(response, request)

			var got []UserStock
			json.NewDecoder(response.Body).Decode(&got)
			if len(got) != len(test.want) {
				t.Fatal("didn't get the expected amount of user stock")
			}
			for i := range test.want {
				want := UserStock{Ticker: test.want[i], Price: getPrice(test.want[i], date)}
				if !reflect.DeepEqual(want, got[i]) {
					t.Errorf("want %v user stock, got %v", test.want[i], got[i])
				}
			}
		})
	}
}
