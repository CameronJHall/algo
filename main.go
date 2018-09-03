package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cjhall1283/algo/config"
	"github.com/cjhall1283/algo/influx"
	"github.com/cjhall1283/algo/robinhood"
)

func main() {
	// Get config
	conf, err := config.Parse("./example_config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	c := make(chan robinhood.Quote)
	go histProducer(conf, client, c)
	go algo(c)

	select {}
}

func algo(c chan robinhood.Quote) {
	for {
		select {
		case quote := <-c:
			fmt.Println(influx.GetQuoteILP(quote))
		}
	}
}

func liveProducer(conf config.Config, client *http.Client, c chan robinhood.Quote) {
	for {
		// Collect quotes
		quotes, err := robinhood.GetLiveQuotes(conf.Symbols, client)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Insert collected data into the chan
		for _, quote := range quotes.QuotesArray {
			c <- quote
		}

		// Sleep so we don't overwhelm the API
		time.Sleep(10 * time.Second)
	}
}

func histProducer(conf config.Config, client *http.Client, c chan robinhood.Quote) {
	// Collect quotes
	hist, err := robinhood.GetHistoricalQuotes(conf.Symbols, "day", "year", client)
	if err != nil {
		fmt.Println(err)
	}

	for _, res := range hist.Results {
		for _, histQuote := range res.Timeseries {
			quote := robinhood.Quote{Symbol: res.Symbol, LastTradePrice: histQuote.OpenPrice}
			c <- quote
		}
	}
}
