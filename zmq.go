package main

import (
	//"fmt"
	//"github.com/iotaledger/giota"
	"github.com/prometheus/client_golang/prometheus"
)

func metrics_zmq(e *Exporter) {

	e.iota_zmq_seen_tx_count = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_seen_tx_count",
			Name: "iota_zmq_seen_tx_count",
			Help: "Count of transactions seen by zeroMQ.",
		})

	e.iota_zmq_txs_with_value_count = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_txs_with_value_count",
			Name: "iota_zmq_txs_with_value_count",
			Help: "Count of transactions seen by zeroMQ that have a non-zero value.",
		})

	e.iota_zmq_confirmed_tx_count = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_confirmed_tx_count",
			Name: "iota_zmq_confirmed_tx_count",
			Help: "Count of transactions confirmed by zeroMQ.",
		})

	e.iota_zmq_to_process = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_to_process",
			Name: "iota_zmq_to_process",
			Help: "toProcess from RSTAT output of ZMQ.",
		})

	e.iota_zmq_to_broadcast = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_to_broadcast",
			Name: "iota_zmq_to_broadcast",
			Help: "toBroadcast from RSTAT output of ZMQ.",
		})

	e.iota_zmq_to_request = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_to_request",
			Name: "iota_zmq_to_request",
			Help: "toRequest from RSTAT output of ZMQ.",
		})

	e.iota_zmq_to_reply = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_to_reply",
			Name: "iota_zmq_to_reply",
			Help: "toReply from RSTAT output of ZMQ.",
		})

	e.iota_zmq_total_transactions = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_total_transactions",
			Name: "iota_zmq_total_transactions",
			Help: "totalTransactions from RSTAT output of ZMQ.",
		})
}

func describe_zmq(e *Exporter, ch chan<- *prometheus.Desc) {
	ch <- e.iota_zmq_seen_tx_count.Desc()
	ch <- e.iota_zmq_txs_with_value_count.Desc()
	ch <- e.iota_zmq_confirmed_tx_count.Desc()
	ch <- e.iota_zmq_to_process.Desc()
	ch <- e.iota_zmq_to_broadcast.Desc()
	ch <- e.iota_zmq_to_request.Desc()
	ch <- e.iota_zmq_to_reply.Desc()
	ch <- e.iota_zmq_total_transactions.Desc()
}

func collect_zmq(e *Exporter, ch chan<- prometheus.Metric) {
	ch <- e.iota_zmq_seen_tx_count
	ch <- e.iota_zmq_txs_with_value_count
	ch <- e.iota_zmq_confirmed_tx_count
	ch <- e.iota_zmq_to_process
	ch <- e.iota_zmq_to_broadcast
	ch <- e.iota_zmq_to_request
	ch <- e.iota_zmq_to_reply
	ch <- e.iota_zmq_total_transactions
}

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
