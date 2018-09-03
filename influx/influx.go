package influx

import (
	"fmt"

	"github.com/cjhall1283/algo/robinhood"
)

// GetQuoteILP will print quote data in ILP to stdout
func GetQuoteILP(quote robinhood.Quote) (ilpString string) {
	ilpString = fmt.Sprintf("quote,symbol=%s price=%s", quote.Symbol, quote.LastTradePrice)
	return
}
