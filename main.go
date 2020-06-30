package main

import (
	"bufio"
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

	"gonum.org/v1/gonum/stat"
)

type Dataset struct {
	Data [][]interface{} `json:"data"`
}

type Response struct {
	Dataset Dataset `json:"dataset"`
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getCSVStandardDeviation(prices []int64) float64 {
	avgs := make([]float64, 0)
	currentPrice := float64(prices[len(prices)-1])
	maxDiff := float64(-100)
	for i := len(prices) - 2; i >= 0; i-- {
		yesterdaysPrice := float64(prices[i])
		percentChange := math.Log(currentPrice / yesterdaysPrice)
		if math.Abs(percentChange) > maxDiff {
			maxDiff = math.Abs(percentChange)
		}
		avgs = append(avgs, percentChange)
		currentPrice = yesterdaysPrice
	}
	return stat.StdDev(avgs, nil)
}

// procedure taken from http://www.fool.com/knowledge-center/2015/09/12/how-to-calculate-annualized-volatility.aspx
func getStandardDeviation(dataArray [][]interface{}) float64 {
	if len(dataArray) == 0 {
		panic("no entries in data array")
	}
	yesterdaysPrice := dataArray[len(dataArray)-1][4].(float64)
	var avgs []float64
	for i := len(dataArray) - 2; i >= 0; i-- {
		data := dataArray[i]
		currentPrice := data[4].(float64)
		percentChange := math.Log(currentPrice / yesterdaysPrice)
		avgs = append(avgs, percentChange)
		yesterdaysPrice = currentPrice
	}
	return stat.StdDev(avgs, nil)
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

func getExponent(p float64) float64 {
	return -1 * math.Sqrt(math.Pi) * ErfInv(1-p)
}

// implementation taken from http://quant.stackexchange.com/a/25074/19960
func determinePrice(annualizedVolatility float64, days int, currentPrice float64, p float64) (float64, error) {
	if p < 0 || p > 1 {
		return -1, fmt.Errorf("invalid probability: %v", p)
	}
	years := float64(days) / 365
	exponent := getExponent(p)
	return currentPrice * math.Pow((1+annualizedVolatility), exponent*math.Sqrt(years)), nil
}

func main() {
	stock := flag.String("stock", "", "Which stock to evaluate")
	total := flag.Int("total", 0, "How many dollars you wish to spend")
	percent := flag.Float64("percent", 0.99, "Percent chance of executing the order (between 0 and 1)")
	useCSV := flag.Bool("csv", false, "CSV mode")
	flag.Parse()
	if *stock == "" {
		log.Fatal(errors.New("Usage: main.go --stock=AAPL"))
	}
	if *percent > 1 || *percent < 0 {
		log.Fatalf("Percentage must be between 0 and 1, was %f", *percent)
	}
	var stddev float64
	var currentPrice float64
	if *useCSV {
		f, err := os.Open(fmt.Sprintf("data/%s.csv", strings.ToLower(*stock)))
		checkError(err)
		bs := bufio.NewScanner(f)
		prices := make([]int64, 0)
		for bs.Scan() {
			parts := strings.Split(bs.Text(), " ")
			price, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			prices = append(prices, price)
		}
		stddev = getCSVStandardDeviation(prices)
		currentPrice = float64(prices[len(prices)-1]) / 10000
	} else {
		f, err := os.Open(fmt.Sprintf("data/%s.json", strings.ToLower(*stock)))
		checkError(err)
		var r Response
		err = json.NewDecoder(f).Decode(&r)
		checkError(err)
		stddev = getStandardDeviation(r.Dataset.Data)
	}
	fmt.Printf("the standard deviation is: %f\n", stddev)
	annualized := stddev * math.Sqrt(252)
	fmt.Printf("annualized: %f%%\n", annualized*100)
	fmt.Println("current price:", currentPrice)
	limitPrice, err := determinePrice(annualized, 365, currentPrice, *percent)
	checkError(err)
	shares := float64(*total) / limitPrice
	diff := shares - float64(*total)/currentPrice
	fmt.Printf(`Based on this stock's volatility, you should set a limit order for: $%.2f.

Compared with buying it at the current price, you'll be able to buy %.1f extra shares (a value of $%.2f)
`, limitPrice, diff, diff*limitPrice)

}
