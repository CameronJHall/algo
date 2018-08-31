package influx

import (
	"fmt"

	"github.com/cjhall1283/algo/robinhood"
)

// PrintQuoteILP will print quote data in ILP to stdout
func PrintQuoteILP(quote robinhood.Quote) {
	fmt.Println(fmt.Sprintf("quote,symbol=%s price=%s", quote.Symbol, quote.Price))
}
