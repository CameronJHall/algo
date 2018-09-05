package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cjhall1283/algo/config"
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

	for _, symbol := range conf.Symbols {
		fmt.Println(robinhood.GetInstrumentID(symbol, client))
	}

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
			c <- quote
		}

		// Sleep so we don't overwhelm the API
		time.Sleep(time.Duration(conf.LiveFrequency) * time.Second)
	}
}

func histProducer(conf config.Config, client *http.Client, c chan robinhood.Quote) {
	// Collect quotes
	hist, err := robinhood.GetHistoricalQuotes(conf.Symbols, conf.HistFrequency, conf.HistRange, client)
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
	rollingAverages := make(map[string]float64)
	for {
		select {
		case quote := <-incomingQuotes:
			val, _ := strconv.ParseFloat(quote.LastTradePrice, 64)
			if _, exists := rollingAverages[quote.Symbol]; !exists {
				rollingAverages[quote.Symbol] = val
			} else {
				if val > (rollingAverages[quote.Symbol]*1.002) {
					outgoingDecisions <- map[string]int{quote.Symbol: 5}
				}
				if val < (rollingAverages[quote.Symbol]*0.998) {
					outgoingDecisions <- map[string]int{quote.Symbol: -5}
				}
				fmt.Println(fmt.Sprintf("[[%s]] Live: %f Rolling: %f Live/Rolling: %f", quote.Symbol, val, rollingAverages[quote.Symbol], val/rollingAverages[quote.Symbol]))
				rollingAverages[quote.Symbol] = (rollingAverages[quote.Symbol]*0.9 + val*0.1)
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
