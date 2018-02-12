package main

import (
	//"fmt"
	"github.com/iotaledger/giota"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func metrics_neighbors(e *Exporter) {
	e.iota_neighbors_info_total_neighbors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "total_neighbors_ws",
			Name: "iota_neighbors_info_total_neighbors",
			Help: "Total number of neighbors as received in the getNeighbors ws call.",
		})

	e.iota_neighbors_info_active_neighbors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "total_neighbors_ws",
			Name: "iota_neighbors_info_active_neighbors",
			Help: "Total number of neighbors that are active.",
		})

	e.iota_neighbors_new_transactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_new_transactions",
			Name: "iota_neighbors_new_transactions",
			Help: "Number of New Transactions for a specific Neighbor.",
		},
		[]string{"id"},
	)

	e.iota_neighbors_random_transactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_random_transactions",
			Name: "iota_neighbors_random_transactions",
			Help: "Number of Random Transactions for a specific Neighbor.",
		},
		[]string{"id"},
	)

	e.iota_neighbors_all_transactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_all_transactions",
			Name: "iota_neighbors_all_transactions",
			Help: "Number of All Transaction Types for a specific Neighbor.",
		},
		[]string{"id"},
	)

	e.iota_neighbors_invalid_transactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_invalid_transactions",
			Name: "iota_neighbors_invalid_transactions",
			Help: "Number of Invalid Transactions for a specific Neighbor.",
		},
		[]string{"id"},
	)

	e.iota_neighbors_sent_transactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "exporter",
			//Name: "neighbors_sent_transactions",
			Name: "iota_neighbors_sent_transactions",
			Help: "Number of Invalid Transactions for a specific Neighbor.",
		},
		[]string{"id"},
	)

	e.iota_neighbors_active = prometheus.NewGaugeVec(
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

func describe_neighbors(e *Exporter, ch chan<- *prometheus.Desc) {
	ch <- e.iota_neighbors_info_total_neighbors.Desc()
	ch <- e.iota_neighbors_info_active_neighbors.Desc()
	e.iota_neighbors_new_transactions.Describe(ch)
	e.iota_neighbors_random_transactions.Describe(ch)
	e.iota_neighbors_all_transactions.Describe(ch)
	e.iota_neighbors_invalid_transactions.Describe(ch)
	e.iota_neighbors_sent_transactions.Describe(ch)
	e.iota_neighbors_active.Describe(ch)
}

func collect_neighbors(e *Exporter, ch chan<- prometheus.Metric) {
	ch <- e.iota_neighbors_info_total_neighbors
	ch <- e.iota_neighbors_info_active_neighbors
	e.iota_neighbors_new_transactions.Collect(ch)
	e.iota_neighbors_random_transactions.Collect(ch)
	e.iota_neighbors_all_transactions.Collect(ch)
	e.iota_neighbors_invalid_transactions.Collect(ch)
	e.iota_neighbors_sent_transactions.Collect(ch)
	e.iota_neighbors_active.Collect(ch)
}

func scrape_neighbors(e *Exporter, api *giota.API) {
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
