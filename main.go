package main

import (
	"fmt"
	"os"
	"github.com/cjhall1283/algo/config"
	"github.com/cjhall1283/algo/robinhood"
)

func main() {
	conf, err := config.Parse("./example_config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	robinhood.GetQuote(conf.Symbols)
}
