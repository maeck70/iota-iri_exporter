package main

import (
	"log"
	"net/http"

	"github.com/iotaledger/giota"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	iota_node_http_request_counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "iota_node_http_request_counter",
			Help: "Number of requests to the exporter.",
		})

	iota_node_info_duration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "iota_node_info_duration",
			Help: "Response time of getting Node Info.",
		})

	iota_node_info_available_processors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "iota_node_info_available_processors",
			Help: "Number of cores available in this Node.",
		})

	iota_node_info_free_memory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "iota_node_info_free_memory",
			Help: "Free Memory in this IRI instance.",
		})

	iota_node_info_max_memory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "iota_node_info_max_memory",
			Help: "Max Memory in this IRI instance.",
		})

	iota_node_info_total_memory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "iota_node_info_total_memory",
			Help: "Total Memory in this IRI instance.",
		})

	iota_node_info_latest_milestone = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "iota_node_info_latest_milestone",
			Help: "Tangle milestone at the interval.",
		})

	iota_node_info_latest_subtangle_milestone = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "iota_node_info_latest_subtangle_milestone",
			Help: "Subtangle milestone at the interval.",
		})

	iota_node_info_total_neighbors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "iota_node_info_total_neighbors",
			Help: "Total neighbors at the interval.",
		})

	iota_node_info_total_tips = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "iota_node_info_total_tips",
			Help: "Total tips at the interval.",
		})

	iota_node_info_total_transactions_queued = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "iota_node_info_total_transactions_queued",
			Help: "Total open txs at the interval.",
		})
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(iota_node_http_request_counter)
	prometheus.MustRegister(iota_node_info_duration)
	prometheus.MustRegister(iota_node_info_available_processors)
	prometheus.MustRegister(iota_node_info_free_memory)
	prometheus.MustRegister(iota_node_info_max_memory)
	prometheus.MustRegister(iota_node_info_total_memory)
	prometheus.MustRegister(iota_node_info_latest_milestone)
	prometheus.MustRegister(iota_node_info_latest_subtangle_milestone)
	prometheus.MustRegister(iota_node_info_total_neighbors)
	prometheus.MustRegister(iota_node_info_total_tips)
	prometheus.MustRegister(iota_node_info_total_transactions_queued)

}

func handler(w http.ResponseWriter, r *http.Request) {
	filters := r.URL.Query()["collect[]"]
	//log.Debugln("collect query:", filters)

	api := giota.NewAPI("http://node01.heliumsushi.com:14265", nil)
	resp, err := api.GetNodeInfo()

	if err != nil {
		log.Fatal(err)
	}

	iota_node_http_request_counter.Inc()
	iota_node_info_duration.Set(float64(resp.Duration))
	iota_node_info_available_processors.Set(float64(resp.JREAvailableProcessors))
	iota_node_info_free_memory.Set(float64(resp.JREFreeMemory))
	iota_node_info_max_memory.Set(float64(resp.JREMaxMemory))
	iota_node_info_total_memory.Set(float64(resp.JRETotalMemory))
	iota_node_info_latest_milestone.Set(float64(resp.LatestMilestoneIndex))
	iota_node_info_latest_subtangle_milestone.Set(float64(resp.LatestSolidSubtangleMilestoneIndex))
	iota_node_info_total_neighbors.Set(float64(resp.Neighbors))
	iota_node_info_total_tips.Set(float64(resp.Tips))
	iota_node_info_total_transactions_queued.Set(float64(resp.TransactionsToRequest))

	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.HandlerFor(gatherers,
		promhttp.HandlerOpts{
			ErrorLog:      log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError,
		})
	h.ServeHTTP(w, r)
}

func main() {

	/*	log.Infoln("Starting node_exporter", version.Info())
		log.Infoln("Build context", version.BuildContext())
	*/
	// The Handler function provides a default handler to expose metrics
	// Via an HTTP server. "/metrics" is the usual endpoint for that.
	http.HandleFunc("/metrics", prometheus.InstrumentHandlerFunc("prometheus", handler))
	//http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9311", nil))

}
