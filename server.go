package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type StockServer struct {
	getStockHistory StockHistoryGetter
	http.Handler
}

type Option func(*StockServer) error

func WithStockHistoryGetter(stockHistoryGetter StockHistoryGetter) Option {
	return func(s *StockServer) error {
		if stockHistoryGetter == nil {
			return fmt.Errorf("can't create a stock server without a stock history getter")
		}
		s.getStockHistory = stockHistoryGetter
		return nil
	}
}

func NewStockServer(options ...Option) (*StockServer, error) {
	stockServer := StockServer{}
	stockServer.getStockHistory = GetStockHistory
	router := http.NewServeMux()
	router.HandleFunc("/tickers", basicAuthMiddleware(stockServer.handleUserStocks))
	router.HandleFunc("/tickers/", basicAuthMiddleware(stockServer.handleStockHistory))
	stockServer.Handler = router

	var err error
	for _, applyOption := range options {
		err = applyOption(&stockServer)
		if err != nil {
			return nil, err
		}
	}
	return &stockServer, nil
}

func (s *StockServer) handleStockHistory(response http.ResponseWriter, request *http.Request) {
	ticker := getTicker(request)
	if isKnown(ticker) {
		history := s.getStockHistory(ticker, 90)
		response.Header().Set("content-type", "application/json")
		json.NewEncoder(response).Encode(history)
		return
	}
	response.WriteHeader(http.StatusNotFound)
}

func (s *StockServer) handleUserStocks(response http.ResponseWriter, request *http.Request) {
	user, _, _ := request.BasicAuth()
	var userStocks []UserStock
	date := time.Now().Format("2006-01-02")
	if user == "testA" {
		userStocks = []UserStock{
			{Ticker: "MSFT", Price: getPrice("MSFT", date)},
			{Ticker: "AAPL", Price: getPrice("AAPL", date)},
			{Ticker: "AMZN", Price: getPrice("AMZN", date)},
		}
	} else {
		userStocks = []UserStock{
			{Ticker: "FB", Price: getPrice("FB", date)},
			{Ticker: "NFLX", Price: getPrice("NFLX", date)},
		}
	}
	json.NewEncoder(response).Encode(&userStocks)
}

func getTicker(request *http.Request) string {
	var ticker string
	ticker = strings.TrimPrefix(request.URL.Path, "/tickers/")
	ticker = strings.TrimSuffix(ticker, "/history")
	return ticker
}
