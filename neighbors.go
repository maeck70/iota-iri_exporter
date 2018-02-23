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
	//"fmt"
	"github.com/maeck70/giota"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func metricsNeighbors(e *exporter) {
	e.iotaNeighborsInfoTotalNeighbors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "total_neighbors_ws",
			Name: "iota_neighbors_info_total_neighbors",
			Help: "Total number of neighbors as received in the getNeighbors ws call.",
		})

	e.iotaNeighborsInfoActiveNeighbors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "total_neighbors_ws",
			//Name: "iotaNeighborsInfoActiveNeighbors", // This is the naming in the Grafana dashboard
			Name: "iota_neighbors_active_neighbors",
			Help: "Total number of neighbors that are active.",
		})

	e.iotaNeighborsNewTransactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_new_transactions",
			Name: "iota_neighbors_new_transactions",
			Help: "Number of New Transactions for a specific Neighbor.",
		},
		[]string{"id"},
	)

	e.iotaNeighborsRandomTransactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_random_transactions",
			Name: "iota_neighbors_random_transactions",
			Help: "Number of Random Transactions for a specific Neighbor.",
		},
		[]string{"id"},
	)

	e.iotaNeighborsAllTransactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_all_transactions",
			Name: "iota_neighbors_all_transactions",
			Help: "Number of All transaction Types for a specific Neighbor.",
		},
		[]string{"id"},
	)

	e.iotaNeighborsInvalidTransactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_invalid_transactions",
			Name: "iota_neighbors_invalid_transactions",
			Help: "Number of Invalid Transactions for a specific Neighbor.",
		},
		[]string{"id"},
	)

	e.iotaNeighborsSentTransactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_sent_transactions",
			Name: "iota_neighbors_sent_transactions",
			Help: "Number of Invalid Transactions for a specific Neighbor.",
		},
		[]string{"id"},
	)

	e.iotaNeighborsActive = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_sent_transactions",
			Name: "iota_neighbors_active",
			Help: "Report if the Neighbor Active based on incoming transactions.",
		},
		[]string{"id"},
	)
}

func describeNeighbors(e *exporter, ch chan<- *prometheus.Desc) {
	ch <- e.iotaNeighborsInfoTotalNeighbors.Desc()
	ch <- e.iotaNeighborsInfoActiveNeighbors.Desc()
	e.iotaNeighborsNewTransactions.Describe(ch)
	e.iotaNeighborsRandomTransactions.Describe(ch)
	e.iotaNeighborsAllTransactions.Describe(ch)
	e.iotaNeighborsInvalidTransactions.Describe(ch)
	e.iotaNeighborsSentTransactions.Describe(ch)
	e.iotaNeighborsActive.Describe(ch)
}

func collectNeighbors(e *exporter, ch chan<- prometheus.Metric) {
	ch <- e.iotaNeighborsInfoTotalNeighbors
	ch <- e.iotaNeighborsInfoActiveNeighbors
	e.iotaNeighborsNewTransactions.Collect(ch)
	e.iotaNeighborsRandomTransactions.Collect(ch)
	e.iotaNeighborsAllTransactions.Collect(ch)
	e.iotaNeighborsInvalidTransactions.Collect(ch)
	e.iotaNeighborsSentTransactions.Collect(ch)
	e.iotaNeighborsActive.Collect(ch)
}

func scrapeNeighbors(e *exporter, api *giota.API) {
	// Get getNeighbors metrics
	resp2, err := api.GetNeighbors()

	if err == nil {
		neighborCount := len(resp2.Neighbors)
		e.iotaNeighborsInfoTotalNeighbors.Set(float64(neighborCount))
		e.iotaNeighborsInfoActiveNeighbors.Set(getActiveNeighbors(resp2.Neighbors))
		for n := 1; n < neighborCount; n++ {
			address := string(resp2.Neighbors[n].Address)
			e.iotaNeighborsActive.WithLabelValues(address).Set(
				float64(getActiveNeighbor(address)))
			e.iotaNeighborsNewTransactions.WithLabelValues(address).Set(
				float64(resp2.Neighbors[n].NumberOfNewTransactions))
			e.iotaNeighborsRandomTransactions.WithLabelValues(address).Set(
				float64(resp2.Neighbors[n].NumberOfRandomTransactionRequests))
			e.iotaNeighborsAllTransactions.WithLabelValues(address).Set(
				float64(resp2.Neighbors[n].NumberOfAllTransactions))
			e.iotaNeighborsInvalidTransactions.WithLabelValues(address).Set(
				float64(resp2.Neighbors[n].NumberOfInvalidTransactions))
			e.iotaNeighborsSentTransactions.WithLabelValues(address).Set(
				float64(resp2.Neighbors[n].NumberOfSentTransactions))
		}
	} else {
		log.Info(err)
	}
}
