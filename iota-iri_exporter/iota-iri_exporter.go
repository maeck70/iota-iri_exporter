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
var Version = "0.0.1"

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9187").String()
	metricPath    = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	targetAddress = kingpin.Flag("web.iri-path", "URI of the IOTA IRI Node to scrape.").Default("http://localhost:14265").String()
)

const (
	namespace = "iota-iri"
)

type Exporter struct {
	iriAddress                                string
	iota_node_info_latest_milestone           prometheus.Gauge
	iota_node_info_latest_subtangle_milestone prometheus.Gauge
	iota_node_info_totalScrapes               prometheus.Counter
}

func NewExporter(iriAddress string) *Exporter {
	return &Exporter{
		iriAddress: iriAddress,
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

		iota_node_info_totalScrapes: prometheus.NewCounter(
			prometheus.CounterOpts{
				//Namespace: namespace,
				//Subsystem: "exporter",
				//Name: "scrapes_total",
				Name: "iota_node_info_scrapes_total",
				Help: "Total number of scrapes.",
			}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.iota_node_info_latest_milestone.Desc()
	ch <- e.iota_node_info_latest_subtangle_milestone.Desc()
	ch <- e.iota_node_info_totalScrapes.Desc()
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(ch)
	ch <- e.iota_node_info_latest_milestone
	ch <- e.iota_node_info_latest_subtangle_milestone
	ch <- e.iota_node_info_totalScrapes
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {

	// Get the actual responses for the node
	api := giota.NewAPI(e.iriAddress, nil)
	resp, err := api.GetNodeInfo()

	if err != nil {
		log.Fatal(err)
	}

	e.iota_node_info_latest_milestone.Set(float64(resp.LatestMilestoneIndex))
	e.iota_node_info_latest_subtangle_milestone.Set(float64(resp.LatestSolidSubtangleMilestoneIndex))
	e.iota_node_info_totalScrapes.Inc()
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

	log.Infof("Starting Server: %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
