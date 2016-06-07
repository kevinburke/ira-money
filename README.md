# IRA Money

Code here is designed to implement: http://quant.stackexchange.com/a/25074/19960

## Usage:

```
go run erf.go main.go --stock={STOCK}
```

This will print out the limit order value you should set to have a ~98% chance
of having that order execute in the next 365 days.

## Getting Data

To measure annualized volatility, you'll need historical stock data. I was able
to find it using the EOD dataset from [Quandl](https://quandl.com). I only
needed a trial account.

Download the JSON dataset and save it to `data/<stock-name>.json`. This file
will get loaded when `main.go` runs.

## Running the file

Run `go run erf.go main.go --stock=AAPL --total=400`

```
the standard deviation is: 0.677742
annualized: 10.758814
current price: 72.77
Based on this stock's volatility, you should set a limit order for: $67.35.

Compared with buying it at the current price, you'll be able to buy 6.1 extra
shares (a value of $409.83)
```

## Definitions

To calculate annualized volatility, I took the standard
deviation of every daily closing price in the dataset, then
multiplied by the square root of 252, the number of trading
days in a year. The procedure is described in more detail here:
http://www.fool.com/knowledge-center/2015/09/12/how-to-calculate-annualized-vol
atility.aspx

The current price is taken from the Yahoo Finance API, an implementation can be
found in main.go.
