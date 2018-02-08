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
	//"github.com/maeck70/iota-iri_exporter/mock_getnodeinfo"
)

// Version is set during build to the git describe version
// (semantic version)-(commitish) form.
var Version = "0.0.1"

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9187").String()
	metricPath    = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	targetAddress = kingpin.Flag("web.iri-path", "URI of the IOTA IRI Node to scrape.").Default("http://localhost:14265").String()
	//targetAddress = kingpin.Flag("web.iri-path", "URI of the IOTA IRI Node to scrape.").Default("http://node21.heliumsushi.com:14265").String()
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

/*func NewExporter() *Exporter {
	return &Exporter{

		iota_node_info_latest_milestone: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "iota_node_info_latest_milestone",
				Help: "Tangle milestone at the interval.",
				Subsystem: exporter,
			}),

		iota_node_info_latest_subtangle_milestone: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "iota_node_info_latest_subtangle_milestone",
				Help: "Subtangle milestone at the interval.",
				Subsystem: exporter,
			}),
	}
}
*/

func iriGetNodeInfo() {

	// Get the actual responses for the node
	//api := giota.NewAPI("http://node21.heliumsushi.com:14265", nil)
	api := giota.NetAPI(targetAddress, nil)
	resp, err := api.GetNodeInfo()

	if err != nil {
		log.Fatal(err)
	}

	//// Mock the response from the node
	//resp := mock_getnodeinfo.GetNodeInfo{
	//	Duration: 100,
	//	JREAvailableProcessors : 4,
	//	JREFreeMemory : 5000000,
	//	JREMaxMemory : 8500000,
	//	JRETotalMemory : 11000000,
	//	LatestMilestoneIndex : 454678,
	//	LatestSolidSubtangleMilestoneIndex : 454677,
	//	Neighbors : 7,
	//	Tips : 6555,
	//	TransactionsToRequest : 12,
	//}

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

	/*	exporter := NewExporter()
		prometheus.MustRegister(exporter)
	*/
	http.Handle(*metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage) // nolint: errcheck
	})

	log.Infof("Starting Server: %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
