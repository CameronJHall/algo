package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
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

	incomingQuotes := make(chan robinhood.Quote)
	outgoingDecisions := make(chan map[string]int)

	go liveProducer(conf, client, incomingQuotes)
	go algo(incomingQuotes, outgoingDecisions)
	go broker(outgoingDecisions)

	select {}
}

func liveProducer(conf config.Config, client *http.Client, c chan robinhood.Quote) {
	for {
		// Collect quotes
		quotes, err := robinhood.GetLiveQuotes(conf.Symbols, client)
		if err != nil {
			fmt.Println(err)
		}

		// Insert collected data into the chan
		for _, quote := range quotes.QuotesArray {
			fmt.Println(influx.GetQuoteILP(quote))
			c <- quote
		}

		// Sleep so we don't overwhelm the API
		time.Sleep(time.Duration(conf.QuoteFrequency) * time.Second)
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

func algo(incomingQuotes chan robinhood.Quote, outgoingDecisions chan map[string]int) {
	for {
		select {
		case quote := <-incomingQuotes:
			val, _ := strconv.ParseFloat(quote.LastTradePrice, 64)
			if val > 300 {
				outgoingDecisions <- map[string]int{quote.Symbol: 5}
			}
		}
	}
}

func broker(outgoingDecisions chan map[string]int) {
	for {
		select {
		case decision := <-outgoingDecisions:
			for k := range decision {
				fmt.Println(fmt.Sprintf("%s:%d", k, decision[k]))
			}
		}
	}
}
