package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type StockServer struct {
	stockHistoryHandler StockHistoryHandler
	userStockHandler    UserStockHandle
	http.Handler
}

type Option func(*StockServer) error

func WithStockHistoryHandler(stockHistoryHandler StockHistoryHandler) Option {
	return func(s *StockServer) error {
		if stockHistoryHandler == nil {
			return fmt.Errorf("can't create a stock server without a stock history handler")
		}
		s.stockHistoryHandler = stockHistoryHandler
		return nil
	}
}

func WithUserStockHandler(userStockHandler UserStockHandle) Option {
	return func(s *StockServer) error {
		if userStockHandler == nil {
			return fmt.Errorf("can't create a stock server without a user stock handler")
		}
		s.userStockHandler = userStockHandler
		return nil
	}
}

func NewStockServer(options ...Option) (*StockServer, error) {
	stockServer := StockServer{}
	stockServer.stockHistoryHandler = HandleStockHistory
	stockServer.userStockHandler = HandleUserStock
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
	stockHistory, err := s.stockHistoryHandler(ticker)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	response.Header().Set("content-type", "application/json")
	json.NewEncoder(response).Encode(stockHistory)
}

func (s *StockServer) handleUserStocks(response http.ResponseWriter, request *http.Request) {
	user, _, _ := request.BasicAuth()
	userStocks := s.userStockHandler(user)
	json.NewEncoder(response).Encode(&userStocks)
}

func getTicker(request *http.Request) string {
	var ticker string
	ticker = strings.TrimPrefix(request.URL.Path, "/tickers/")
	ticker = strings.TrimSuffix(ticker, "/history")
	return ticker
}
