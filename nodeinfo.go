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
	"github.com/iotaledger/giota"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func metricsNodeinfo(e *exporter) {

	e.iotaNodeInfoDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "duration",
			Name: "iota_node_info_duration",
			Help: "Response time of getting Node Info.",
		})

	e.iotaNodeInfoAvailableProcessors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "available_processors",
			Name: "iota_node_info_available_processors",
			Help: "Number of cores available in this Node.",
		})

	e.iotaNodeInfoFreeMemory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "free_memory",
			Name: "iota_node_info_free_memory",
			Help: "Free Memory in this IRI instance.",
		})

	e.iotaNodeInfoMaxMemory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "max_memory",
			Name: "iota_node_info_max_memory",
			Help: "Max Memory in this IRI instance.",
		})

	e.iotaNodeInfoTotalMemory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "total_memory",
			Name: "iota_node_info_total_memory",
			Help: "Total Memory in this IRI instance.",
		})

	e.iotaNodeInfoLatestMilestone = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "latest_milestone",
			Name: "iota_node_info_latest_milestone",
			Help: "Tangle milestone at the interval.",
		})

	e.iotaNodeInfoLatestSubtangleMilestone = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "latest_subtangle_milestone",
			Name: "iota_node_info_latest_subtangle_milestone",
			Help: "Subtangle milestone at the interval.",
		})

	e.iotaNodeInfoTotalNeighbors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "total_neighbors",
			Name: "iota_node_info_total_neighbors",
			Help: "Total neighbors at the interval.",
		})

	e.iotaNodeInfoTotalTips = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "total_tips",
			Name: "iota_node_info_total_tips",
			Help: "Total tips at the interval.",
		})

	e.iotaNodeInfoTotalTransactionsQueued = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "total_transactions_queued",
			Name: "iota_node_info_total_transactions_queued",
			Help: "Total open txs at the interval.",
		})
}

func describeNodeinfo(e *exporter, ch chan<- *prometheus.Desc) {
	ch <- e.iotaNodeInfoDuration.Desc()
	ch <- e.iotaNodeInfoAvailableProcessors.Desc()
	ch <- e.iotaNodeInfoFreeMemory.Desc()
	ch <- e.iotaNodeInfoMaxMemory.Desc()
	ch <- e.iotaNodeInfoTotalMemory.Desc()
	ch <- e.iotaNodeInfoLatestMilestone.Desc()
	ch <- e.iotaNodeInfoLatestSubtangleMilestone.Desc()
	ch <- e.iotaNodeInfoTotalNeighbors.Desc()
	ch <- e.iotaNodeInfoTotalTips.Desc()
	ch <- e.iotaNodeInfoTotalTransactionsQueued.Desc()
}

func collectNodeinfo(e *exporter, ch chan<- prometheus.Metric) {
	ch <- e.iotaNodeInfoDuration
	ch <- e.iotaNodeInfoAvailableProcessors
	ch <- e.iotaNodeInfoFreeMemory
	ch <- e.iotaNodeInfoMaxMemory
	ch <- e.iotaNodeInfoTotalMemory
	ch <- e.iotaNodeInfoLatestMilestone
	ch <- e.iotaNodeInfoLatestSubtangleMilestone
	ch <- e.iotaNodeInfoTotalNeighbors
	ch <- e.iotaNodeInfoTotalTips
	ch <- e.iotaNodeInfoTotalTransactionsQueued
}

func scrapeNodeinfo(e *exporter, api *giota.API) {
	resp, err := api.GetNodeInfo()

	if err == nil {
		// Set response values into the predefined metrics
		e.iotaNodeInfoDuration.Set(float64(resp.Duration))
		e.iotaNodeInfoAvailableProcessors.Set(float64(resp.JREAvailableProcessors))
		e.iotaNodeInfoFreeMemory.Set(float64(resp.JREFreeMemory))
		e.iotaNodeInfoMaxMemory.Set(float64(resp.JREMaxMemory))
		e.iotaNodeInfoTotalMemory.Set(float64(resp.JRETotalMemory))
		e.iotaNodeInfoLatestMilestone.Set(float64(resp.LatestMilestoneIndex))
		e.iotaNodeInfoLatestSubtangleMilestone.Set(float64(resp.LatestSolidSubtangleMilestoneIndex))
		e.iotaNodeInfoTotalNeighbors.Set(float64(resp.Neighbors))
		e.iotaNodeInfoTotalTips.Set(float64(resp.Tips))
		e.iotaNodeInfoTotalTransactionsQueued.Set(float64(resp.TransactionsToRequest))

		e.iotaNodeInfoTotalScrapes.Inc()
	} else {
		log.Info(err)
	}
}
