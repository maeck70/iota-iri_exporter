/*
MIT License

Copyright (c) 2018 Marcel van Eck

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type tradingPair struct {
	symbol                string
	bid                   float64
	bidSize               float64
	ask                   float64
	askSize               float64
	dailyChange           float64
	dailyChangePercentage float64
	lastPrice             float64
	volume                float64
	high                  float64
	low                   float64
}

var tradingPairList = []string{
	"tIOTUSD",
	"tIOTEUR",
	"tIOTBTC",
	"tIOTETH",
	"tBTCUSD",
	"tBTCEUR",
	"tETHUSD"}

var bitfinexURL = "https://api.bitfinex.com/v2/tickers?symbols="

func metricsBitfinex(e *exporter) {
	e.iotaMarketTradePrice = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_sent_transactions",
			Name: "iota_market_trade_price",
			Help: "Latest price from Bitfinex.",
		},
		[]string{"pair"},
	)

	e.iotaMarketTradeVolume = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_sent_transactions",
			Name: "iota_market_trade_volume",
			Help: "Latest volume from Bitfinex.",
		},
		[]string{"pair"},
	)

	e.iotaMarketHighPrice = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_sent_transactions",
			Name: "iota_market_high_price",
			Help: "Highest price from Bitfinex.",
		},
		[]string{"pair"},
	)

	e.iotaMarketLowPrice = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_sent_transactions",
			Name: "iota_market_low_price",
			Help: "Lowest price from Bitfinex.",
		},
		[]string{"pair"},
	)
}

func describeBitfinex(e *exporter, ch chan<- *prometheus.Desc) {
	e.iotaMarketTradePrice.Describe(ch)
	e.iotaMarketTradeVolume.Describe(ch)
	e.iotaMarketHighPrice.Describe(ch)
	e.iotaMarketLowPrice.Describe(ch)
}

func collectBitfinex(e *exporter, ch chan<- prometheus.Metric) {
	e.iotaMarketTradePrice.Collect(ch)
	e.iotaMarketTradeVolume.Collect(ch)
	e.iotaMarketHighPrice.Collect(ch)
	e.iotaMarketLowPrice.Collect(ch)
}

func init() {
	// Expand the Bitfinex URL with the list of trading pairs
	for t := range tradingPairList {
		bitfinexURL += tradingPairList[t]
		if t < len(tradingPairList)-1 {
			bitfinexURL += ","
		}
	}
}

func scrapeBitfinex(e *exporter) {
	// Get Bitfinex metrics
	req, err := http.NewRequest("GET", bitfinexURL, nil)
	resp, _ := http.DefaultClient.Do(req)

	if err == nil {

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		// Convert the character stream to the struct
		s := strings.Split(string(body), "],[")
		s[0] = strings.TrimLeft(s[0], "[[")
		s[len(s)-1] = strings.TrimRight(s[len(s)-1], "]]")

		for n := range s {
			s2 := strings.Split(s[n], ",")

			tp := tradingPair{}
			tp.symbol = strings.Trim(s2[0], "\"")
			tp.symbol = strings.TrimLeft(tp.symbol, "t")

			/*		tp.bid, _ = strconv.ParseFloat(s2[1], 64)
					tp.bidSize, _ = strconv.ParseFloat(s2[2], 64)
					tp.ask, _ = strconv.ParseFloat(s2[3], 64)
					tp.askSize, _ = strconv.ParseFloat(s2[4], 64)
					tp.dailyChange, _ = strconv.ParseFloat(s2[5], 64)
					tp.dailyChangePercentage, _ = strconv.ParseFloat(s2[6], 64)*/
			tp.lastPrice, _ = strconv.ParseFloat(s2[7], 64)
			tp.volume, _ = strconv.ParseFloat(s2[8], 64)
			tp.high, _ = strconv.ParseFloat(s2[9], 64)
			tp.low, _ = strconv.ParseFloat(s2[10], 64)

			e.iotaMarketTradePrice.WithLabelValues(tp.symbol).Set(tp.lastPrice)
			e.iotaMarketTradeVolume.WithLabelValues(tp.symbol).Set(tp.volume)
			e.iotaMarketHighPrice.WithLabelValues(tp.symbol).Set(tp.high)
			e.iotaMarketLowPrice.WithLabelValues(tp.symbol).Set(tp.low)
		}
	} else {
		log.Info(err)
	}
}
