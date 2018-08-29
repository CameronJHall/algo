package robinhood

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var BASEURL = "https://api.robinhood.com/"

type Quote struct {
	Symbol string `json:"symbol"`
	Price  string `json:"last_trade_price"`
}

func GetQuote(symbols []string) {
	client := &http.Client{}
	for _, symbol := range symbols {
		quote := Quote{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%squotes/%s/", BASEURL, symbol), nil)
		if err != nil {
			fmt.Println(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		err = json.NewDecoder(resp.Body).Decode(&quote)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fmt.Sprintf("%s: %s", quote.Symbol, quote.Price))
	}
}
