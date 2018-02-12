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
var Version = "0.0.2"

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9187").String()
	metricPath    = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	targetAddress = kingpin.Flag("web.iri-path", "URI of the IOTA IRI Node to scrape.").Default("http://localhost:14265").String()
)

const (
	namespace = "iota-iri"
)

type Exporter struct {
	iriAddress string

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
	iota_node_info_totalScrapes               prometheus.Counter
	iota_neighbors_info_total_neighbors       prometheus.Gauge
	iota_neighbors_info_active_neighbors      prometheus.Gauge
	iota_neighbors_new_transactions           *prometheus.GaugeVec
	iota_neighbors_random_transactions        *prometheus.GaugeVec
	iota_neighbors_all_transactions           *prometheus.GaugeVec
	iota_neighbors_invalid_transactions       *prometheus.GaugeVec
	iota_neighbors_sent_transactions          *prometheus.GaugeVec
	iota_neighbors_active                     *prometheus.GaugeVec
}

func NewExporter(iriAddress string) *Exporter {
	return &Exporter{
		iriAddress: iriAddress,

		iota_node_info_duration: prometheus.NewGauge(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "duration",
				Name: "iota_node_info_duration",
				Help: "Response time of getting Node Info.",
			}),

		iota_node_info_available_processors: prometheus.NewGauge(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "available_processors",
				Name: "iota_node_info_available_processors",
				Help: "Number of cores available in this Node.",
			}),

		iota_node_info_free_memory: prometheus.NewGauge(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "free_memory",
				Name: "iota_node_info_free_memory",
				Help: "Free Memory in this IRI instance.",
			}),

		iota_node_info_max_memory: prometheus.NewGauge(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "max_memory",
				Name: "iota_node_info_max_memory",
				Help: "Max Memory in this IRI instance.",
			}),

		iota_node_info_total_memory: prometheus.NewGauge(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "total_memory",
				Name: "iota_node_info_total_memory",
				Help: "Total Memory in this IRI instance.",
			}),

		iota_node_info_latest_milestone: prometheus.NewGauge(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "latest_milestone",
				Name: "iota_node_info_latest_milestone",
				Help: "Tangle milestone at the interval.",
			}),

		iota_node_info_latest_subtangle_milestone: prometheus.NewGauge(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "latest_subtangle_milestone",
				Name: "iota_node_info_latest_subtangle_milestone",
				Help: "Subtangle milestone at the interval.",
			}),

		iota_node_info_total_neighbors: prometheus.NewGauge(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "total_neighbors",
				Name: "iota_node_info_total_neighbors",
				Help: "Total neighbors at the interval.",
			}),

		iota_node_info_total_tips: prometheus.NewGauge(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "total_tips",
				Name: "iota_node_info_total_tips",
				Help: "Total tips at the interval.",
			}),

		iota_node_info_total_transactions_queued: prometheus.NewGauge(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "total_transactions_queued",
				Name: "iota_node_info_total_transactions_queued",
				Help: "Total open txs at the interval.",
			}),

		iota_node_info_totalScrapes: prometheus.NewCounter(
			prometheus.CounterOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "scrapes_total",
				Name: "iota_node_info_scrapes_total",
				Help: "Total number of scrapes.",
			}),

		iota_neighbors_info_total_neighbors: prometheus.NewGauge(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "total_neighbors_ws",
				Name: "iota_neighbors_info_total_neighbors",
				Help: "Total number of neighbors as received in the getNeighbors ws call.",
			}),

		iota_neighbors_info_active_neighbors: prometheus.NewGauge(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "total_neighbors_ws",
				Name: "iota_neighbors_info_active_neighbors",
				Help: "Total number of neighbors that are active.",
			}),

		iota_neighbors_new_transactions: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "neighbors_new_transactions",
				Name: "iota_neighbors_new_transactions",
				Help: "Number of New Transactions for a specific Neighbor.",
			},
			[]string{"id"},
		),

		iota_neighbors_random_transactions: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "neighbors_random_transactions",
				Name: "iota_neighbors_random_transactions",
				Help: "Number of Random Transactions for a specific Neighbor.",
			},
			[]string{"id"},
		),

		iota_neighbors_all_transactions: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "neighbors_all_transactions",
				Name: "iota_neighbors_all_transactions",
				Help: "Number of All Transaction Types for a specific Neighbor.",
			},
			[]string{"id"},
		),

		iota_neighbors_invalid_transactions: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "neighbors_invalid_transactions",
				Name: "iota_neighbors_invalid_transactions",
				Help: "Number of Invalid Transactions for a specific Neighbor.",
			},
			[]string{"id"},
		),

		iota_neighbors_sent_transactions: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "neighbors_sent_transactions",
				Name: "iota_neighbors_sent_transactions",
				Help: "Number of Invalid Transactions for a specific Neighbor.",
			},
			[]string{"id"},
		),

		iota_neighbors_active: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "neighbors_sent_transactions",
				Name: "iota_neighbors_active",
				Help: "Report if the Neighbor Active based on incoming transactions.",
			},
			[]string{"id"},
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

	ch <- e.iota_node_info_duration.Desc()
	ch <- e.iota_node_info_available_processors.Desc()
	ch <- e.iota_node_info_free_memory.Desc()
	ch <- e.iota_node_info_max_memory.Desc()
	ch <- e.iota_node_info_total_memory.Desc()
	ch <- e.iota_node_info_latest_milestone.Desc()
	ch <- e.iota_node_info_latest_subtangle_milestone.Desc()
	ch <- e.iota_node_info_total_neighbors.Desc()
	ch <- e.iota_node_info_total_tips.Desc()
	ch <- e.iota_node_info_total_transactions_queued.Desc()
	ch <- e.iota_node_info_totalScrapes.Desc()

	ch <- e.iota_neighbors_info_total_neighbors.Desc()
	ch <- e.iota_neighbors_info_active_neighbors.Desc()

	e.iota_neighbors_new_transactions.Describe(ch)
	e.iota_neighbors_random_transactions.Describe(ch)
	e.iota_neighbors_all_transactions.Describe(ch)
	e.iota_neighbors_invalid_transactions.Describe(ch)
	e.iota_neighbors_sent_transactions.Describe(ch)
	e.iota_neighbors_active.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(ch)
	ch <- e.iota_node_info_duration
	ch <- e.iota_node_info_available_processors
	ch <- e.iota_node_info_free_memory
	ch <- e.iota_node_info_max_memory
	ch <- e.iota_node_info_total_memory
	ch <- e.iota_node_info_latest_milestone
	ch <- e.iota_node_info_latest_subtangle_milestone
	ch <- e.iota_node_info_total_neighbors
	ch <- e.iota_node_info_total_tips
	ch <- e.iota_node_info_total_transactions_queued
	ch <- e.iota_node_info_totalScrapes

	ch <- e.iota_neighbors_info_total_neighbors
	ch <- e.iota_neighbors_info_active_neighbors

	e.iota_neighbors_new_transactions.Collect(ch)
	e.iota_neighbors_random_transactions.Collect(ch)
	e.iota_neighbors_all_transactions.Collect(ch)
	e.iota_neighbors_invalid_transactions.Collect(ch)
	e.iota_neighbors_sent_transactions.Collect(ch)
	e.iota_neighbors_active.Collect(ch)
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {

	// Get getNodeInfo metrics
	api := giota.NewAPI(e.iriAddress, nil)
	resp, err := api.GetNodeInfo()

	if err == nil {
		// Set response values into the predefined metrics
		e.iota_node_info_duration.Set(float64(resp.Duration))
		e.iota_node_info_available_processors.Set(float64(resp.JREAvailableProcessors))
		e.iota_node_info_free_memory.Set(float64(resp.JREFreeMemory))
		e.iota_node_info_max_memory.Set(float64(resp.JREMaxMemory))
		e.iota_node_info_total_memory.Set(float64(resp.JRETotalMemory))
		e.iota_node_info_latest_milestone.Set(float64(resp.LatestMilestoneIndex))
		e.iota_node_info_latest_subtangle_milestone.Set(float64(resp.LatestSolidSubtangleMilestoneIndex))
		e.iota_node_info_total_neighbors.Set(float64(resp.Neighbors))
		e.iota_node_info_total_tips.Set(float64(resp.Tips))
		e.iota_node_info_total_transactions_queued.Set(float64(resp.TransactionsToRequest))

		e.iota_node_info_totalScrapes.Inc()
	} else {
		log.Info(err)
	}

	// Get getNeighbors metrics
	resp2, err := api.GetNeighbors()

	/* --- Neighbors CURL response
	"address": "d5c52a6a.ftth.concepts.nl:15600",
	"connectionType": "tcp",
	"numberOfAllTransactions": 0,
	"numberOfInvalidTransactions": 0,
	"numberOfNewTransactions": 0,
	"numberOfRandomTransactionRequests": 0,
	"numberOfSentTransactions": 0
	*/

	if err == nil {
		neighbor_cnt := len(resp2.Neighbors)
		e.iota_neighbors_info_total_neighbors.Set(float64(neighbor_cnt))
		e.iota_neighbors_info_active_neighbors.Set(GetActiveNeighbors(resp2.Neighbors))
		for n := 1; n < neighbor_cnt; n++ {
			//log.Infof("Neighbor %s_is %s", string(resp2.Neighbors[n].Address), actify(GetActiveNeighbor(string(resp2.Neighbors[n].Address))))
			e.iota_neighbors_active.WithLabelValues(string(resp2.Neighbors[n].Address)).Set(float64(GetActiveNeighbor(string(resp2.Neighbors[n].Address))))
			// TODO: update to enable the two missing metrics from the getNeighbors api ass soon as this call has been updated.
			e.iota_neighbors_new_transactions.WithLabelValues(string(resp2.Neighbors[n].Address)).Set(float64(resp2.Neighbors[n].NumberOfNewTransactions))
			//e.iota_neighbors_random_transactions.WithLabelValues(string(resp2.Neighbors[n].Address)).Set(float64(resp2.Neighbors[n].NumberOfRandomTransactionRequests))
			e.iota_neighbors_all_transactions.WithLabelValues(string(resp2.Neighbors[n].Address)).Set(float64(resp2.Neighbors[n].NumberOfAllTransactions))
			e.iota_neighbors_invalid_transactions.WithLabelValues(string(resp2.Neighbors[n].Address)).Set(float64(resp2.Neighbors[n].NumberOfInvalidTransactions))
			//e.iota_neighbors_sent_transactions.WithLabelValues(string(resp2.Neighbors[n].Address)).Set(float64(resp2.Neighbors[n].NumberOfSentTransactions))
		}
	} else {
		log.Info(err)
	}
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

	http.Handle(*metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage) // nolint: errcheck
	})

	log.Infof("Starting %s_exporter Server on port %s monitoring %s", namespace, *listenAddress, *targetAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
