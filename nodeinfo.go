package main

import (
	//"fmt"
	"github.com/iotaledger/giota"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func metrics_nodeinfo(e *Exporter) {

	e.iota_node_info_duration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "duration",
			Name: "iota_node_info_duration",
			Help: "Response time of getting Node Info.",
		})

	e.iota_node_info_available_processors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "available_processors",
			Name: "iota_node_info_available_processors",
			Help: "Number of cores available in this Node.",
		})

	e.iota_node_info_free_memory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "free_memory",
			Name: "iota_node_info_free_memory",
			Help: "Free Memory in this IRI instance.",
		})

	e.iota_node_info_max_memory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "max_memory",
			Name: "iota_node_info_max_memory",
			Help: "Max Memory in this IRI instance.",
		})

	e.iota_node_info_total_memory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "total_memory",
			Name: "iota_node_info_total_memory",
			Help: "Total Memory in this IRI instance.",
		})

	e.iota_node_info_latest_milestone = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "latest_milestone",
			Name: "iota_node_info_latest_milestone",
			Help: "Tangle milestone at the interval.",
		})

	e.iota_node_info_latest_subtangle_milestone = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "latest_subtangle_milestone",
			Name: "iota_node_info_latest_subtangle_milestone",
			Help: "Subtangle milestone at the interval.",
		})

	e.iota_node_info_total_neighbors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "total_neighbors",
			Name: "iota_node_info_total_neighbors",
			Help: "Total neighbors at the interval.",
		})

	e.iota_node_info_total_tips = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "total_tips",
			Name: "iota_node_info_total_tips",
			Help: "Total tips at the interval.",
		})

	e.iota_node_info_total_transactions_queued = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "total_transactions_queued",
			Name: "iota_node_info_total_transactions_queued",
			Help: "Total open txs at the interval.",
		})
}

func describe_nodeinfo(e *Exporter, ch chan<- *prometheus.Desc) {
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
}

func collect_nodeinfo(e *Exporter, ch chan<- prometheus.Metric) {
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
}

func scrape_nodeinfo(e *Exporter, api *giota.API) {
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
}
