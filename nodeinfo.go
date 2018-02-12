package main

import (
	//"fmt"
	"github.com/iotaledger/giota"
	"github.com/prometheus/common/log"
)


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
