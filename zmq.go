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
	"strings"
	"time"
	//"github.com/iotaledger/giota"
	"github.com/pebbe/zmq4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// Size of the timeslice buffer, each array row represents one second
const TIMESLICE_LIMIT = 300

type timeslicef struct {
	tx_any            float64
	tx_value          float64
	tx_confirmed      float64
	tx_toprocess      float64
	tx_tobroadcast    float64
	tx_toreply        float64
	tx_numberstoredtx float64
	tx_txntorequest   float64
}

type timeslice struct {
	tx_any            int64
	tx_value          int64
	tx_confirmed      int64
	tx_toprocess      int64
	tx_tobroadcast    int64
	tx_toreply        int64
	tx_numberstoredtx int64
	tx_txntorequest   int64
}

func (ts *timeslice) reset() {
	//ts.tx_any = 0
	ts.tx_value = 0
	ts.tx_confirmed = 0
	//ts.tx_toprocess = 0
	//ts.tx_tobroadcast = 0
	//ts.tx_toreply = 0
}

type Transaction struct {
	Hash         string
	Address      string
	Value        int64
	Tag          string
	Timestamp    string
	CurrentIndex string
	LastIndex    string
	Bundle       string
	Trunk        string
	Branch       string
	ArrivalDate  string
}

type sn struct {
	Index       string
	Hash        string
	AddressHash string
	Trunk       string
	Branch      string
	Bundle      string
}

type queue struct {
	ReceiveQueueSize   int64
	BroadcastQueueSize int64
	TxnToRequest       int64
	ReplyQueueSize     int64
	NumberOfStoredTxns int64
}

var tx_total int64
var timeslice_ptr int
var timeslice_ptr_prev int
var timeslice_set [TIMESLICE_LIMIT]timeslice
var address = "tcp://node21.heliumsushi.com:5556"

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
	ts := getTimesliceAvg()
	e.iota_zmq_seen_tx_count.Set(ts.tx_any)
	e.iota_zmq_txs_with_value_count.Set(ts.tx_value)
	e.iota_zmq_confirmed_tx_count.Set(ts.tx_confirmed)
	e.iota_zmq_to_process.Set(ts.tx_toprocess)
	e.iota_zmq_to_broadcast.Set(ts.tx_tobroadcast)
	e.iota_zmq_to_request.Set(ts.tx_txntorequest)
	e.iota_zmq_to_reply.Set(ts.tx_toreply)
	e.iota_zmq_total_transactions.Set(float64(tx_total))
}

func collectTimeslice() {

	socket, err := zmq4.NewSocket(zmq4.SUB)
	must(err)
	socket.SetSubscribe("") // TODO: Listen to only tx, sn, rstat

	err = socket.Connect(address)
	must(err)

	log.Infof("Connected to IRI at address %s.", address)

	for {

		msg, err := socket.Recv(0)
		must(err)

		parts := strings.Fields(msg)
		switch parts[0] {

		//Transaction
		case "tx":
			tx_total++
			txn := Transaction{
				Hash:         parts[1],
				Address:      parts[2],
				Value:        stoi(parts[3]),
				Tag:          parts[4],
				Timestamp:    parts[5],
				CurrentIndex: parts[6],
				LastIndex:    parts[7],
				Bundle:       parts[8],
				Trunk:        parts[9],
				Branch:       parts[10],
				ArrivalDate:  parts[11],
			}
			timeslice_set[timeslice_ptr].tx_any++
			if txn.Value != 0 {
				timeslice_set[timeslice_ptr].tx_value++
			}

		// Transaction confirmed
		case "sn":
			timeslice_set[timeslice_ptr].tx_any++
			/*				stat := sn{
								Index:       parts[1],
								Hash:        parts[2],
								AddressHash: parts[3],
								Trunk:       parts[4],
								Branch:      parts[5],
								Bundle:      parts[6],
							}
			*/timeslice_set[timeslice_ptr].tx_confirmed++

		case "rstat":
			stat := queue{
				ReceiveQueueSize:   stoi(parts[1]),
				BroadcastQueueSize: stoi(parts[2]),
				TxnToRequest:       stoi(parts[3]),
				ReplyQueueSize:     stoi(parts[4]),
				NumberOfStoredTxns: stoi(parts[5]),
			}
			// Note that these are total counts, no need to increment into the timeslice
			timeslice_set[timeslice_ptr].tx_toprocess = stat.ReceiveQueueSize
			timeslice_set[timeslice_ptr].tx_tobroadcast = stat.BroadcastQueueSize
			timeslice_set[timeslice_ptr].tx_toreply = stat.ReplyQueueSize
			timeslice_set[timeslice_ptr].tx_numberstoredtx = stat.NumberOfStoredTxns
			timeslice_set[timeslice_ptr].tx_txntorequest = stat.TxnToRequest
		}
	}
}

func advanceTimeslice() {
	// Advance to the next row in the timeslice array, rotate to 0 if final row is passed
	timeslice_ptr_new := timeslice_ptr
	timeslice_ptr_new++
	if timeslice_ptr_new == TIMESLICE_LIMIT {
		timeslice_ptr_new = 0
	}

	timeslice_set[timeslice_ptr_new].reset()

	// Copy the rstat metricsto the next slice
	timeslice_set[timeslice_ptr_new].tx_any = timeslice_set[timeslice_ptr].tx_any
	timeslice_set[timeslice_ptr_new].tx_toprocess = timeslice_set[timeslice_ptr].tx_toprocess
	timeslice_set[timeslice_ptr_new].tx_tobroadcast = timeslice_set[timeslice_ptr].tx_tobroadcast
	timeslice_set[timeslice_ptr_new].tx_toreply = timeslice_set[timeslice_ptr].tx_toreply
	timeslice_set[timeslice_ptr_new].tx_numberstoredtx = timeslice_set[timeslice_ptr].tx_numberstoredtx
	timeslice_set[timeslice_ptr_new].tx_txntorequest = timeslice_set[timeslice_ptr].tx_txntorequest

	// Activate the new slice
	timeslice_ptr = timeslice_ptr_new
}

func manageTimeslice() {
	timeslice_ptr_prev = 0
	timeslice_ptr = 0

	// Start collector concurrently
	go collectTimeslice()

	// Rotate the timeslice array pointer every second
	for {
		time.Sleep(time.Duration(1) * time.Second)
		advanceTimeslice()
	}
}

func getTimesliceAvg() timeslicef {

	// Define the start of the time slice we are averaging on
	timeslice_ptr_avg := timeslice_ptr - 1
	if timeslice_ptr_avg == -1 {

		timeslice_ptr_avg = TIMESLICE_LIMIT - 1
	}

	timeslice_avg := timeslicef{}
	timeslice_cnt := float64(0)

	// Accumulate the transactions that happened in the timeslice
	for ts := timeslice_ptr_prev; ts != timeslice_ptr_avg; ts++ {
		if ts == TIMESLICE_LIMIT {
			if timeslice_ptr_avg == 0 {
				// Exception condition where we need to stop the loop
				break
			}
			ts = 0
		}
		timeslice_avg.tx_any += float64(timeslice_set[ts].tx_any)
		timeslice_avg.tx_value += float64(timeslice_set[ts].tx_value)
		timeslice_avg.tx_confirmed += float64(timeslice_set[ts].tx_confirmed)
		timeslice_avg.tx_toprocess = float64(timeslice_set[ts].tx_toprocess)
		timeslice_avg.tx_tobroadcast = float64(timeslice_set[ts].tx_tobroadcast)
		timeslice_avg.tx_toreply = float64(timeslice_set[ts].tx_toreply)
		timeslice_avg.tx_numberstoredtx = float64(timeslice_set[ts].tx_numberstoredtx)
		timeslice_avg.tx_txntorequest = float64(timeslice_set[ts].tx_txntorequest)
		timeslice_cnt++
	}

	// Calculate the average of the timeslice
	if timeslice_cnt > 1 {
		//timeslice_avg.tx_any = timeslice_avg.tx_any / timeslice_cnt
		timeslice_avg.tx_value = timeslice_avg.tx_value / timeslice_cnt
		timeslice_avg.tx_confirmed = timeslice_avg.tx_confirmed / timeslice_cnt
		//timeslice_avg.tx_toprocess = timeslice_avg.tx_toprocess / timeslice_cnt
		//timeslice_avg.tx_tobroadcast = timeslice_avg.tx_tobroadcast / timeslice_cnt
		//timeslice_avg.tx_toreply = timeslice_avg.tx_toreply / timeslice_cnt
		//timeslice_avg.tx_numberstoredtx = timeslice_avg.tx_numberstoredtx / timeslice_cnt
		//timeslice_avg.tx_txntorequest = timeslice_avg.tx_txntorequest / timeslice_cnt
	}

	// Maintain where we left off for the next call to this function
	timeslice_ptr_prev = timeslice_ptr_avg

	return timeslice_avg
}

func init_zmq() {

	go manageTimeslice()

}
