package main

import (
	//"fmt"
	//"github.com/iotaledger/giota"
)


func scrape_zmq(e *Exporter) {

	e.iota_zmq_seen_tx_count.Set(1)
	e.iota_zmq_txs_with_value_count.Set(1)
	e.iota_zmq_confirmed_tx_count.Set(1)
	e.iota_zmq_to_process.Set(1)
	e.iota_zmq_to_broadcast.Set(1)
	e.iota_zmq_to_request.Set(1)
	e.iota_zmq_to_reply.Set(1)
	e.iota_zmq_total_transactions.Set(1)

}