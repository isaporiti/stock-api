package server

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

type Stock struct {
	Ticker string  `json:"-"`
	Date   string  `json:"date"`
	Price  float64 `json:"price"`
}

type UserStock struct {
	Ticker string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type StockHistoryGetter func(ticker string, length int) []Stock

func GetStockHistory(ticker string, length int) []Stock {
	var stockHistory []Stock = make([]Stock, length)
	now := time.Now()
	for i := length; i > 1; i-- {
		iterationDate := now.AddDate(0, 0, -i+1)
		stock := Stock{Date: iterationDate.Format("2006-01-02"), Ticker: ticker}
		stock.Price = getPrice(stock.Ticker, stock.Date)
		stockHistory[i-1] = stock
	}
	stock := Stock{Date: now.Format("2006-01-02"), Ticker: ticker}
	stock.Price = getPrice(stock.Ticker, stock.Date)
	stockHistory[0] = stock
	return stockHistory
}

func getPrice(ticker string, date string) float64 {
	inputString := fmt.Sprintf("%v%s", date, ticker)
	hasher := sha256.New()
	hasher.Write([]byte(inputString))
	hashBytes := hasher.Sum(nil)
	numericHash := binary.BigEndian.Uint64(hashBytes[:8])
	lastFourDigits := numericHash % 10000
	return float64(lastFourDigits) / 100
}

var knownTickers = []string{
	"AAPL",
	"MSFT",
	"GOOG",
	"AMZN",
	"FB",
	"TSLA",
	"NVDA",
	"JPM",
	"BABA",
	"JNJ",
	"WMT",
	"PG",
	"PYPL",
	"DIS",
	"ADBE",
	"PFE",
	"V",
	"MA",
	"CRM",
	"NFLX",
}

func isKnown(ticker string) bool {
	for _, knownTicker := range knownTickers {
		if knownTicker == ticker {
			return true
		}
	}
	return false
}

func GetUserStock(user string) []UserStock {
	var userStock []UserStock
	date := time.Now().Format("2006-01-02")
	for i, ticker := range knownTickers {
		if i%2 == 0 {
			userStock = append(userStock, UserStock{Ticker: ticker, Price: getPrice(ticker, date)})
		}
	}
	return userStock
}
