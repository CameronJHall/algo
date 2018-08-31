package influx

import (
	"testing"

	"github.com/cjhall1283/algo/robinhood"
	"github.com/stretchr/testify/assert"
)

func TestGetQuoteILP(test *testing.T) {
	test.Run("Should return data in ILP", func(test *testing.T) {
		quote := robinhood.Quote{}
		quote.Symbol, quote.Price = "FB", "100.000000"
		ilpString := GetQuoteILP(quote)
		assert.Equal(test, ilpString, "quote,symbol=FB price=100.000000")
	})
}
