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
	symbol            string
	bid               float64
	bid_size          float64
	ask               float64
	ask_size          float64
	daily_change      float64
	daily_change_perc float64
	last_price        float64
	volume            float64
	high              float64
	low               float64
}

var tradingPairList = []string{
	"tIOTUSD",
	"tIOTEUR",
	"tIOTBTC",
	"tIOTETH",
	"tBTCUSD",
	"tBTCEUR",
	"tETHUSD"}

func metrics_bitfinex(e *Exporter) {
	e.iota_market_trade_price = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_sent_transactions",
			Name: "iota_market_trade_price",
			Help: "Latest price from Bitfinex.",
		},
		[]string{"pair"},
	)

	e.iota_market_trade_volume = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_sent_transactions",
			Name: "iota_market_trade_volume",
			Help: "Latest volume from Bitfinex.",
		},
		[]string{"pair"},
	)

	e.iota_market_high_price = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_sent_transactions",
			Name: "iota_market_high_price",
			Help: "Highest price from Bitfinex.",
		},
		[]string{"pair"},
	)

	e.iota_market_low_price = prometheus.NewGaugeVec(
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

func describe_bitfinex(e *Exporter, ch chan<- *prometheus.Desc) {
	e.iota_market_trade_price.Describe(ch)
	e.iota_market_trade_volume.Describe(ch)
	e.iota_market_high_price.Describe(ch)
	e.iota_market_low_price.Describe(ch)
}

func collect_bitfinex(e *Exporter, ch chan<- prometheus.Metric) {
	e.iota_market_trade_price.Collect(ch)
	e.iota_market_trade_volume.Collect(ch)
	e.iota_market_high_price.Collect(ch)
	e.iota_market_low_price.Collect(ch)
}

func scrape_bitfinex(e *Exporter) {
	// Get Bitfinex metrics

	url := "https://api.bitfinex.com/v2/tickers?symbols="
	for t := range tradingPairList {
		url += tradingPairList[t]
		if t < len(tradingPairList)-1 {
			url += ","
		}
	}
	req, err := http.NewRequest("GET", url, nil)
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
			tp.bid, _ = strconv.ParseFloat(s2[1], 64)
			tp.bid_size, _ = strconv.ParseFloat(s2[2], 64)
			tp.ask, _ = strconv.ParseFloat(s2[3], 64)
			tp.ask_size, _ = strconv.ParseFloat(s2[4], 64)
			tp.daily_change, _ = strconv.ParseFloat(s2[5], 64)
			tp.daily_change_perc, _ = strconv.ParseFloat(s2[6], 64)
			tp.last_price, _ = strconv.ParseFloat(s2[7], 64)
			tp.volume, _ = strconv.ParseFloat(s2[8], 64)
			tp.high, _ = strconv.ParseFloat(s2[9], 64)
			tp.low, _ = strconv.ParseFloat(s2[10], 64)

			e.iota_market_trade_price.WithLabelValues(tp.symbol).Set(tp.last_price)
			e.iota_market_trade_volume.WithLabelValues(tp.symbol).Set(tp.volume)
			e.iota_market_high_price.WithLabelValues(tp.symbol).Set(tp.high)
			e.iota_market_low_price.WithLabelValues(tp.symbol).Set(tp.low)
		}
	} else {
		log.Info(err)
	}
}
