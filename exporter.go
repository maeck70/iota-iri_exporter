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
	"fmt"
	"net/http"
	"runtime"

	"github.com/iotaledger/giota"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Version is set during build to the git describe version
// (semantic version)-(commitish) form.
var Version = "0.3.1"

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9311").String()
	metricPath    = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	targetAddress = kingpin.Flag("web.iri-path", "URI of the IOTA IRI Node to scrape.").Default("http://localhost:14265").String()
	targetZmqAddress = kingpin.Flag("web.zmq-path", "URI of the IOTA IRI ZMQ Node to scrape.").Default("tcp://localhost:5556").String()
)

const (
	namespace = "iota-iri"
)

type Exporter struct {
	iriAddress string

	iota_node_info_totalScrapes               prometheus.Counter
	iota_node_info_duration                   prometheus.Gauge
	iota_node_info_available_processors       prometheus.Gauge
	iota_node_info_free_memory                prometheus.Gauge
	iota_node_info_max_memory                 prometheus.Gauge
	iota_node_info_total_memory               prometheus.Gauge
	iota_node_info_latest_milestone           prometheus.Gauge
	iota_node_info_latest_subtangle_milestone prometheus.Gauge
	iota_node_info_total_neighbors            prometheus.Gauge
	iota_node_info_total_tips                 prometheus.Gauge
	iota_node_info_total_transactions_queued  prometheus.Gauge
	iota_neighbors_info_total_neighbors       prometheus.Gauge
	iota_neighbors_info_active_neighbors      prometheus.Gauge
	iota_neighbors_new_transactions           *prometheus.GaugeVec
	iota_neighbors_random_transactions        *prometheus.GaugeVec
	iota_neighbors_all_transactions           *prometheus.GaugeVec
	iota_neighbors_invalid_transactions       *prometheus.GaugeVec
	iota_neighbors_sent_transactions          *prometheus.GaugeVec
	iota_neighbors_active                     *prometheus.GaugeVec
	iota_zmq_seen_tx_count                    prometheus.Gauge
	iota_zmq_txs_with_value_count             prometheus.Gauge
	iota_zmq_confirmed_tx_count               prometheus.Gauge
	iota_zmq_to_request                       prometheus.Gauge
	iota_zmq_to_process                       prometheus.Gauge
	iota_zmq_to_broadcast                     prometheus.Gauge
	iota_zmq_to_reply                         prometheus.Gauge
	iota_zmq_total_transactions               prometheus.Gauge
	iota_market_trade_price                   *prometheus.GaugeVec
	iota_market_trade_volume                  *prometheus.GaugeVec
	iota_market_high_price                    *prometheus.GaugeVec
	iota_market_low_price                     *prometheus.GaugeVec
}

func NewExporter(iriAddress string) *Exporter {
	e := &Exporter{
		iriAddress: iriAddress,

		iota_node_info_totalScrapes: prometheus.NewCounter(
			prometheus.CounterOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "scrapes_total",
				Name: "iota_node_info_scrapes_total",
				Help: "Total number of scrapes.",
			}),
	}

	metrics_nodeinfo(e)
	metrics_neighbors(e)
	metrics_zmq(e)
	metrics_bitfinex(e)

	return e
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

	ch <- e.iota_node_info_totalScrapes.Desc()

	describe_nodeinfo(e, ch)
	describe_neighbors(e, ch)
	describe_zmq(e, ch)
	describe_bitfinex(e, ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(ch)
	ch <- e.iota_node_info_totalScrapes

	collect_nodeinfo(e, ch)
	collect_neighbors(e, ch)
	collect_zmq(e, ch)
	collect_bitfinex(e, ch)
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {

	api := giota.NewAPI(e.iriAddress, nil)

	scrape_nodeinfo(e, api)
	scrape_neighbors(e, api)
	scrape_zmq(e)
	scrape_bitfinex(e)
}

func main() {
	kingpin.Version(fmt.Sprintf("iota-iri_exporter %s (built with %s)\n", Version, runtime.Version()))
	log.AddFlags(kingpin.CommandLine)
	kingpin.Parse()

	// landingPage contains the HTML served at '/'.
	// TODO: Make this nicer and more informative.
	var landingPage = []byte(`<html>
	<head><title>Iota-IRI Exporter</title></head>
	<body>
	<h1>Iota-IRI Node Exporter</h1>
	<p><a href='` + *metricPath + `'>Metrics</a></p>
	</body>
	</html>
	`)

	exporter := NewExporter(*targetAddress)
	prometheus.MustRegister(exporter)

	init_zmq(targetZmqAddress)

	http.Handle(*metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage) // nolint: errcheck
	})

	log.Infof("Starting %s_exporter Server on port %s monitoring %s", namespace, *listenAddress, *targetAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
