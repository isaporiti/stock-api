package server

import (
	"fmt"
)

type StockHistoryHandler func(ticker string) (stockHistory []Stock, unknownTickerError error)

func HandleStockHistory(ticker string) (stockHistory []Stock, unknownTickerError error) {
	if !isKnown(ticker) {
		unknownTickerError = fmt.Errorf("ticker %s unknown", ticker)
	}
	stockHistory = GetStockHistory(ticker, 90)
	return
}
