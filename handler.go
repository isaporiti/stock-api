package server

import (
	"fmt"
	"time"
)

type StockHistoryHandler func(ticker string) (stockHistory []Stock, unknownTickerError error)

func HandleStockHistory(ticker string) (stockHistory []Stock, unknownTickerError error) {
	if !isKnown(ticker) {
		unknownTickerError = fmt.Errorf("ticker %s unknown", ticker)
	}
	stockHistory = GetStockHistory(ticker, 90)
	return
}

type UserStockHandle func(user string) (userStocks []UserStock)

func HandleUserStock(user string) (userStocks []UserStock) {
	date := time.Now().Format("2006-01-02")
	for i, ticker := range knownTickers {
		if user == "testA" && i%2 == 0 {
			userStocks = append(userStocks, UserStock{Ticker: ticker, Price: getPrice(ticker, date)})
			continue
		}
		if user == "testB" && i%2 != 0 {
			userStocks = append(userStocks, UserStock{Ticker: ticker, Price: getPrice(ticker, date)})
		}
	}
	return
}
