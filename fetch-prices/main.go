package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/kevinburke/portfolio/services/alphavantage"
)

type price struct {
	Date  string
	Price int64
}

func main() {
	flag.Parse()
	symbol := flag.Arg(0)
	if symbol == "" {
		log.Fatal("must provide a symbol")
	}
	client := alphavantage.NewClient(os.Getenv("ALPHAVANTAGE_KEY"))
	resp, err := client.GetRange(context.Background(), symbol, nil)
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	prices := make([]*price, len(resp.TimeSeries))
	for d, entry := range resp.TimeSeries {
		prices[i] = &price{Date: d, Price: entry.Close}
		i++
	}
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Date < prices[j].Date
	})
	for i := range prices {
		fmt.Printf("%s %d\n", prices[i].Date, prices[i].Price)
	}
}
