package main

import (
	//"fmt"
	"github.com/iotaledger/giota"
	"github.com/prometheus/common/log"
)


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
