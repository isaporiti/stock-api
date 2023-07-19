package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	if user == "testA" {
		userStocks = []UserStock{
			{Ticker: "MSFT", Price: 10},
			{Ticker: "AAPL", Price: 20},
			{Ticker: "AMZN", Price: 4},
		}
	} else {
		userStocks = []UserStock{
			{Ticker: "FB", Price: 13},
			{Ticker: "NFLX", Price: 7},
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
