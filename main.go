package main

import (
	"github.com/cjhall1283/algo/robinhood"
)

func main() {
	robinhood.GetQuote([]string{"MSFT", "TSLA"})
}
