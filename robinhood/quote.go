package robinhood

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// BASEURL - Setting this as a global for now
var BASEURL = "https://api.robinhood.com/"

// Quote - Struct for getting symbol and price from API response
type Quote struct {
	Symbol string `json:"symbol"`
	Price  string `json:"last_trade_price"`
}

// GetQuote - Get last trade price for a list of stocks
// No return, prints out stock prices
func GetQuote(symbols []string) {

	client := &http.Client{}
	// Iterate through the provided symbols
	for _, symbol := range symbols {
		quote := Quote{}

		// Create the new request (creating request first in the case
		// that other functions may require headers and such we can keep structure)
		req, err := http.NewRequest("GET", fmt.Sprintf("%squotes/%s/", BASEURL, symbol), nil)
		if err != nil {
			fmt.Println(err)
		}

		// Execute request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		// Decode the body using hte initialized quote object
		err = json.NewDecoder(resp.Body).Decode(&quote)
		if err != nil {
			fmt.Println(err)
		}

		// Print the collected data
		fmt.Println(fmt.Sprintf("%s: %s", quote.Symbol, quote.Price))
	}
}
