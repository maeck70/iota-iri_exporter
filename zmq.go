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
	"math/rand"
	"time"
	//"github.com/iotaledger/giota"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/k0kubun/pp"
	"github.com/pebbe/zmq4"
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

	e.iota_zmq_seen_tx_count.Set(data[0])
	e.iota_zmq_txs_with_value_count.Set(data[1])
	e.iota_zmq_confirmed_tx_count.Set(data[2])
	e.iota_zmq_to_process.Set(data[3])
	e.iota_zmq_to_broadcast.Set(data[4])
	e.iota_zmq_to_request.Set(data[0]+data[1])
	e.iota_zmq_to_reply.Set(data[1]+data[2])
	e.iota_zmq_total_transactions.Set(data[3]+data[4])

}



var data [5]float64
var zmq_url = "localhost:5556"

func zmq_collector() {
	for {
		i := rand.Intn(5)
		data[i]++

		r := rand.Intn(100000)
		time.Sleep(time.Duration(r) * time.Microsecond)	
	}
}


func init_zmq() {

	//go zmq_collector()
	
	log.Infof("ZMQ Initialized on listener %s.", zmq_url)

}


