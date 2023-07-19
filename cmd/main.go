package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	server "github.com/isaporiti/stock-api"
)

func main() {
	stockServer, err := server.NewStockServer()
	if err != nil {
		fmt.Printf("couldn't start server: %s", err)
		os.Exit(1)
	}
	fmt.Println("Stock Server listening on port 5001.")
	log.Fatal(http.ListenAndServe(":5001", stockServer))
}
