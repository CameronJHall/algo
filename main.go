package main

import (
	"fmt"
	"os"

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

	// Collect quotes
	quotes, err := robinhood.GetQuote(conf.Symbols)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the collected data
	for _, quote := range quotes.QuotesArray {
		fmt.Println(fmt.Sprintf("%s: %s", quote.Symbol, quote.Price))
	}
}
