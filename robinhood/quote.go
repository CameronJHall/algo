package robinhood

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// BASEURL - Setting this as a global for now
var BASEURL = "https://api.robinhood.com/"

type HistoricalQuotes struct {
	Results []Result `json:"results"`
}

type Result struct {
	Symbol     string            `json:"symbol"`
	Timeseries []HistoricalQuote `json:"historicals"`
}

type HistoricalQuote struct {
	OpenPrice  string `json:"open_price"`
	ClosePrice string `json:"close_price"`
	HighPrice  string `json:"high_price"`
	LowPrice   string `json:"low_price"`
}

// Quotes - Struct for holding an array of Quotes
type Quotes struct {
	QuotesArray []Quote `json:"results"`
}

// Quote - Struct for getting symbol and price from API response
type Quote struct {
	AskPrice              string `json:"ask_price"`
	AskSize               int    `json:"ask_size"`
	BidPrice              string `json:"bid_price"`
	BidSize               int    `json:"bid_size"`
	LastTradePrice        string `json:"last_trade_price"`
	PreviousClose         string `json:"previous_close"`
	AdjustedPreviousClose string `json:"adjusted_previous_close"`
	PreviousCloseDate     string `json:"previous_close_date"`
	Symbol                string `json:"symbol"`
}

// GetLiveQuote - Get last trade price for a list of stocks
// No return, prints out stock prices
func GetLiveQuotes(symbols []string, client *http.Client) (quotes Quotes, err error) {
	quotes = Quotes{}

	if len(symbols) > 1630 {
		err = errors.New("too many symbols (max 1630)")
		return
	}

	// Create the new request
	req, err := http.NewRequest("GET", BASEURL+"quotes/", nil)
	if err != nil {
		return
	}

	query := req.URL.Query()
	query.Add("symbols", strings.Join(symbols, ","))
	req.URL.RawQuery = query.Encode()

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	// Decode the body using the initialized quotes struct
	err = json.NewDecoder(resp.Body).Decode(&quotes)
	if err != nil {
		return
	}

	return
}

// GetHistoricalQuotes returns a HistoricalQuotes object
func GetHistoricalQuotes(symbols []string, interval, span string, client *http.Client) (history HistoricalQuotes, err error) {
	history = HistoricalQuotes{}

	// Create the new request
	req, err := http.NewRequest("GET", BASEURL+"quotes/historicals/", nil)
	if err != nil {
		return
	}

	query := req.URL.Query()
	query.Add("symbols", strings.Join(symbols, ","))
	query.Add("interval", interval)
	query.Add("span", span)
	req.URL.RawQuery = query.Encode()

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	// Decode the body using the initialized quotes struct
	err = json.NewDecoder(resp.Body).Decode(&history)
	if err != nil {
		return
	}

	return
}
