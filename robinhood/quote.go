package robinhood

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// BASEURL - Setting this as a global for now
var BASEURL = "https://api.robinhood.com/"

// Quotes - Struct for holding an array of Quotes
type Quotes struct {
	QuotesArray []Quote `json:"results"`
}

// Quote - Struct for getting symbol and price from API response
type Quote struct {
	Symbol string `json:"symbol"`
	Price  string `json:"last_trade_price"`
}

// GetQuote - Get last trade price for a list of stocks
// No return, prints out stock prices
func GetQuote(symbols []string) (quotes Quotes, err error) {

	client := &http.Client{}
	quotes = Quotes{}

	if len(symbols) > 1630 {
		err = errors.New("too many symbols (max 1630)")
		return
	}

	// Create the new request
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%squotes/?symbols=%s", BASEURL, strings.Join(symbols, ",")),
		nil)
	if err != nil {
		return
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	// Decode the body using the initialized quotes struct
	err = json.NewDecoder(resp.Body).Decode(&quotes)
	if err != nil {
		return
	}

	return
}
