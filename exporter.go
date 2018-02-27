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
	"github.com/iotaledger/giota"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"runtime"
)

// Version is set during build to the git Describe version
// (semantic version)-(commitish) form.
var Version = "0.4.1"

var (
	listenAddress    = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9311").String()
	metricPath       = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	targetAddress    = kingpin.Flag("web.iri-path", "URI of the IOTA IRI Node to scrape.").Default("http://localhost:14265").String()
	targetZmqAddress = kingpin.Flag("web.zmq-path", "URI of the IOTA IRI ZMQ Node to scrape.").Default("tcp://localhost:5556").String()
)

const (
	namespace = "iota-iri"
)

type exporter struct {
	iriAddress string

	iotaNodeInfoTotalScrapes             prometheus.Counter
	iotaNodeInfoDuration                 prometheus.Gauge
	iotaNodeInfoAvailableProcessors      prometheus.Gauge
	iotaNodeInfoFreeMemory               prometheus.Gauge
	iotaNodeInfoMaxMemory                prometheus.Gauge
	iotaNodeInfoTotalMemory              prometheus.Gauge
	iotaNodeInfoLatestMilestone          prometheus.Gauge
	iotaNodeInfoLatestSubtangleMilestone prometheus.Gauge
	iotaNodeInfoTotalNeighbors           prometheus.Gauge
	iotaNodeInfoTotalTips                prometheus.Gauge
	iotaNodeInfoTotalTransactionsQueued  prometheus.Gauge
	iotaNeighborsInfoTotalNeighbors      prometheus.Gauge
	iotaNeighborsInfoActiveNeighbors     prometheus.Gauge
	iotaNeighborsNewTransactions         *prometheus.GaugeVec
	iotaNeighborsRandomTransactions      *prometheus.GaugeVec
	iotaNeighborsAllTransactions         *prometheus.GaugeVec
	iotaNeighborsInvalidTransactions     *prometheus.GaugeVec
	iotaNeighborsSentTransactions        *prometheus.GaugeVec
	iotaNeighborsActive                  *prometheus.GaugeVec
	iotaZmqSeenTxCount                   prometheus.Gauge
	iotaZmqTxsWithValueCount             prometheus.Gauge
	iotaZmqConfirmedTxCount              prometheus.Gauge
	iotaZmqToRequest                     prometheus.Gauge
	iotaZmqToProcess                     prometheus.Gauge
	iotaZmqToBroadcast                   prometheus.Gauge
	iotaZmqToReply                       prometheus.Gauge
	iotaZmqTotalTransactions             prometheus.Gauge
	iotaMarketTradePrice                 *prometheus.GaugeVec
	iotaMarketTradeVolume                *prometheus.GaugeVec
	iotaMarketHighPrice                  *prometheus.GaugeVec
	iotaMarketLowPrice                   *prometheus.GaugeVec
}

func newExporter(iriAddress string) *exporter {
	e := &exporter{
		iriAddress: iriAddress,

		iotaNodeInfoTotalScrapes: prometheus.NewCounter(
			prometheus.CounterOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "scrapes_total",
				Name: "iota_node_info_scrapes_total",
				Help: "Total number of scrapes.",
			}),
	}

	metricsNodeinfo(e)
	metricsNeighbors(e)
	metricsZmq(e)
	metricsBitfinex(e)

	return e
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {

	ch <- e.iotaNodeInfoTotalScrapes.Desc()

	describeNodeinfo(e, ch)
	describeNeighbors(e, ch)
	describeZmq(e, ch)
	describeBitfinex(e, ch)
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(ch)
	ch <- e.iotaNodeInfoTotalScrapes

	collectNodeinfo(e, ch)
	collectNeighbors(e, ch)
	collectZmq(e, ch)
	collectBitfinex(e, ch)
}

func (e *exporter) scrape(ch chan<- prometheus.Metric) {

	api := giota.NewAPI(e.iriAddress, nil)

	scrapeNodeinfo(e, api)
	scrapeNeighbors(e, api)
	scrapeZmq(e)
	scrapeBitfinex(e)
}

func main() {
	kingpin.Version(fmt.Sprintf("iota-iri_exporter %s (built with %s)\n", Version, runtime.Version()))
	log.AddFlags(kingpin.CommandLine)
	kingpin.Parse()

	// landingPage contains the HTML served at '/'.
	// TODO: Make this nicer and more informative.
	var landingPage = []byte(`<html>
	<head><title>Iota-IRI exporter</title></head>
	<body>
	<h1>Iota-IRI Node exporter</h1>
	<p><a href='` + *metricPath + `'>Metrics</a></p>
	</body>
	</html>
	`)

	exporter := newExporter(*targetAddress)
	prometheus.MustRegister(exporter)

	initZmq(targetZmqAddress)

	http.Handle(*metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage) // nolint: errcheck
	})

	log.Infof("Starting %s_exporter Server on port %s monitoring %s", namespace, *listenAddress, *targetAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
