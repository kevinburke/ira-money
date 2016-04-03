package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/GaryBoone/GoStats/stats"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Dataset struct {
	Data [][]interface{} `json:"data"`
}

type Response struct {
	Dataset Dataset `json:"dataset"`
}

// procedure taken from http://www.fool.com/knowledge-center/2015/09/12/how-to-calculate-annualized-volatility.aspx
func getStandardDeviation(dataArray [][]interface{}) float64 {
	yesterdaysPrice := dataArray[len(dataArray)-1][4].(float64)
	var avgs []float64
	for i := len(dataArray) - 2; i >= 0; i-- {
		data := dataArray[i]
		currentPrice := data[4].(float64)
		percentChange := 100 * (currentPrice - yesterdaysPrice) / (yesterdaysPrice)
		//fmt.Println(percentChange)
		//fmt.Println(data[0])
		//fmt.Println(data[4])
		//fmt.Printf("%s, price: %f, change: %f\n", data[0], currentPrice, percentChange)
		avgs = append(avgs, percentChange)
		yesterdaysPrice = currentPrice
	}
	return stats.StatsPopulationStandardDeviation(avgs)
}

func getCurrentPrice(stock string) (float64, error) {
	url := fmt.Sprintf("https://download.finance.yahoo.com/d/quotes.csv?s=%s&f=sl1", stock)
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	r := csv.NewReader(resp.Body)
	records, err := r.ReadAll()
	if err != nil {
		return -1, err
	}
	if len(records) != 1 {
		return -1, fmt.Errorf("Invalid response: %v", records)
	}
	return strconv.ParseFloat(records[0][1], 64)
}

// implementation taken from http://quant.stackexchange.com/a/25074/19960
func determinePrice(annualizedVolatility float64, days int, currentPrice float64, p float64) (float64, error) {
	if p < 0 || p > 1 {
		return -1, fmt.Errorf("invalid probability: %f", p)
	}
	years := float64(days) / 365
	return currentPrice * math.Pow((1+annualizedVolatility), -0.0314192*math.Sqrt(years)), nil
}

func main() {
	stock := flag.String("stock", "", "Which stock to evaluate")
	total := flag.Int("total", 0, "How many dollars you wish to spend")
	flag.Parse()
	if *stock == "" {
		log.Fatal(errors.New("Usage: main.go --stock=AAPL"))
	}
	f, err := os.Open(fmt.Sprintf("data/%s.json", strings.ToLower(*stock)))
	checkError(err)
	var r Response
	err = json.NewDecoder(f).Decode(&r)
	checkError(err)
	stddev := getStandardDeviation(r.Dataset.Data)
	fmt.Printf("the standard deviation is: %f\n", stddev)
	annualized := stddev * math.Sqrt(252)
	fmt.Printf("annualized: %f\n", annualized)
	currentPrice, err := getCurrentPrice(*stock)
	checkError(err)
	fmt.Println("current price:", currentPrice)
	limitPrice, err := determinePrice(annualized, 365, currentPrice, 0.98)
	checkError(err)
	shares := float64(*total) / limitPrice
	diff := shares - float64(*total)/currentPrice
	fmt.Printf(`Based on this stock's volatility, you should set a limit order for: $%.2f.

Compared with buying it at the current price, you'll be able to buy %.1f extra shares (a value of $%.2f)
`, limitPrice, diff, diff*limitPrice)

}
