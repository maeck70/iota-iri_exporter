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
	"github.com/pebbe/zmq4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"strings"
	"time"
)

// Size of the timeslice buffer, each array row represents one second
const TimesliceLimit = 300

type timeslicef struct {
	txAny            float64
	txValue          float64
	txConfirmed      float64
	txToProcess      float64
	txToBroadcast    float64
	txToReply        float64
	txNumberStoredTx float64
	txTxnToRequest   float64
}

type timeslice struct {
	txAny            int64
	txValue          int64
	txConfirmed      int64
	txToProcess      int64
	txToBroadcast    int64
	txToReply        int64
	txNumberStoredTx int64
	txTxnToRequest   int64
}

func (ts *timeslice) reset() {
	//ts.txAny = 0
	ts.txValue = 0
	ts.txConfirmed = 0
	//ts.txToProcess = 0
	//ts.txToBroadcast = 0
	//ts.txToReply = 0
}

type transaction struct {
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

var txTotal int64
var timeslicePtr int
var timeslicePtrPrev int
var timesliceSet [TimesliceLimit]timeslice

func metricsZmq(e *exporter) {

	e.iotaZmqSeenTxCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_seen_tx_count",
			Name: "iota_zmq_seen_tx_count",
			Help: "Count of transactions seen by zeroMQ.",
		})

	e.iotaZmqTxsWithValueCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_txs_with_value_count",
			Name: "iota_zmq_txs_with_value_count",
			Help: "Count of transactions seen by zeroMQ that have a non-zero value.",
		})

	e.iotaZmqConfirmedTxCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_confirmed_tx_count",
			Name: "iota_zmq_confirmed_tx_count",
			Help: "Count of transactions confirmed by zeroMQ.",
		})

	e.iotaZmqToProcess = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_to_process",
			Name: "iota_zmq_to_process",
			Help: "toProcess from RSTAT output of ZMQ.",
		})

	e.iotaZmqToBroadcast = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_to_broadcast",
			Name: "iota_zmq_to_broadcast",
			Help: "toBroadcast from RSTAT output of ZMQ.",
		})

	e.iotaZmqToRequest = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_to_request",
			Name: "iota_zmq_to_request",
			Help: "toRequest from RSTAT output of ZMQ.",
		})

	e.iotaZmqToReply = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_to_reply",
			Name: "iota_zmq_to_reply",
			Help: "toReply from RSTAT output of ZMQ.",
		})

	e.iotaZmqTotalTransactions = prometheus.NewGauge(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_total_transactions",
			Name: "iota_zmq_total_transactions",
			Help: "totalTransactions from RSTAT output of ZMQ.",
		})
}

func describeZmq(e *exporter, ch chan<- *prometheus.Desc) {
	ch <- e.iotaZmqSeenTxCount.Desc()
	ch <- e.iotaZmqTxsWithValueCount.Desc()
	ch <- e.iotaZmqConfirmedTxCount.Desc()
	ch <- e.iotaZmqToProcess.Desc()
	ch <- e.iotaZmqToBroadcast.Desc()
	ch <- e.iotaZmqToRequest.Desc()
	ch <- e.iotaZmqToReply.Desc()
	ch <- e.iotaZmqTotalTransactions.Desc()
}

func collectZmq(e *exporter, ch chan<- prometheus.Metric) {
	ch <- e.iotaZmqSeenTxCount
	ch <- e.iotaZmqTxsWithValueCount
	ch <- e.iotaZmqConfirmedTxCount
	ch <- e.iotaZmqToProcess
	ch <- e.iotaZmqToBroadcast
	ch <- e.iotaZmqToRequest
	ch <- e.iotaZmqToReply
	ch <- e.iotaZmqTotalTransactions
}

func scrapeZmq(e *exporter) {
	ts := getTimesliceAvg()
	e.iotaZmqSeenTxCount.Set(ts.txAny)
	e.iotaZmqTxsWithValueCount.Set(ts.txValue)
	e.iotaZmqConfirmedTxCount.Set(ts.txConfirmed)
	e.iotaZmqToProcess.Set(ts.txToProcess)
	e.iotaZmqToBroadcast.Set(ts.txToBroadcast)
	e.iotaZmqToRequest.Set(ts.txTxnToRequest)
	e.iotaZmqToReply.Set(ts.txToReply)
	e.iotaZmqTotalTransactions.Set(float64(txTotal))

	log.Debugf("total tx:          %v", txTotal)
	log.Debugf("txAny:            %v tps", ts.txAny)
	log.Debugf("txValue:          %v tx", ts.txValue)
	log.Debugf("txConfirmed:      %v tx", ts.txConfirmed)
	log.Debugf("txToProcess:      %v tx", int64(ts.txToProcess))
	log.Debugf("txToBroadcast:    %v tx", int64(ts.txToBroadcast))
	log.Debugf("txToReply:        %v tx", int64(ts.txToReply))
	log.Debugf("txNumberStoredTx: %v tx", int64(ts.txNumberStoredTx))
	log.Debugf("txTxnToRequest:   %v tx", int64(ts.txTxnToRequest))

}

func collectTimeslice(address *string) {

	socket, err := zmq4.NewSocket(zmq4.SUB)
	must(err)
	socket.SetSubscribe("") // TODO: Listen to only tx, sn, rstat

	err = socket.Connect(*address)
	must(err)

	log.Infof("Connected to IRI at address %s.", *address)

	for {

		msg, err := socket.Recv(0)
		must(err)

		parts := strings.Fields(msg)
		switch parts[0] {

		//transaction
		case "tx":
			txTotal++
			txn := transaction{
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
			timesliceSet[timeslicePtr].txAny++
			if txn.Value != 0 {
				timesliceSet[timeslicePtr].txValue++
			}

		// transaction confirmed
		case "sn":
			timesliceSet[timeslicePtr].txAny++
			/*				stat := sn{
								Index:       parts[1],
								Hash:        parts[2],
								AddressHash: parts[3],
								Trunk:       parts[4],
								Branch:      parts[5],
								Bundle:      parts[6],
							}
			*/timesliceSet[timeslicePtr].txConfirmed++

		case "rstat":
			stat := queue{
				ReceiveQueueSize:   stoi(parts[1]),
				BroadcastQueueSize: stoi(parts[2]),
				TxnToRequest:       stoi(parts[3]),
				ReplyQueueSize:     stoi(parts[4]),
				NumberOfStoredTxns: stoi(parts[5]),
			}
			// Note that these are total counts, no need to increment into the timeslice
			timesliceSet[timeslicePtr].txToProcess = stat.ReceiveQueueSize
			timesliceSet[timeslicePtr].txToBroadcast = stat.BroadcastQueueSize
			timesliceSet[timeslicePtr].txToReply = stat.ReplyQueueSize
			timesliceSet[timeslicePtr].txNumberStoredTx = stat.NumberOfStoredTxns
			timesliceSet[timeslicePtr].txTxnToRequest = stat.TxnToRequest
		}
	}
}

func advanceTimeslice() {
	// Advance to the next row in the timeslice array, rotate to 0 if final row is passed
	timeslicePtrNew := timeslicePtr
	timeslicePtrNew++
	if timeslicePtrNew == TimesliceLimit {
		timeslicePtrNew = 0
	}

	timesliceSet[timeslicePtrNew].reset()

	// Copy the rstat metricsto the next slice
	timesliceSet[timeslicePtrNew].txAny = timesliceSet[timeslicePtr].txAny
	timesliceSet[timeslicePtrNew].txToProcess = timesliceSet[timeslicePtr].txToProcess
	timesliceSet[timeslicePtrNew].txToBroadcast = timesliceSet[timeslicePtr].txToBroadcast
	timesliceSet[timeslicePtrNew].txToReply = timesliceSet[timeslicePtr].txToReply
	timesliceSet[timeslicePtrNew].txNumberStoredTx = timesliceSet[timeslicePtr].txNumberStoredTx
	timesliceSet[timeslicePtrNew].txTxnToRequest = timesliceSet[timeslicePtr].txTxnToRequest

	// Activate the new slice
	timeslicePtr = timeslicePtrNew
}

func manageTimeslice(address *string) {
	timeslicePtrPrev = 0
	timeslicePtr = 0

	// Start collector concurrently
	go collectTimeslice(address)

	// Rotate the timeslice array pointer every second
	for {
		time.Sleep(time.Duration(1) * time.Second)
		advanceTimeslice()
	}
}

func getTimesliceAvg() timeslicef {

	// Define the start of the time slice we are averaging on
	timeslicePtrAvg := timeslicePtr - 1
	if timeslicePtrAvg == -1 {

		timeslicePtrAvg = TimesliceLimit - 1
	}

	timesliceAvg := timeslicef{}
	timesliceCnt := float64(0)

	// Accumulate the transactions that happened in the timeslice
	for ts := timeslicePtrPrev; ts != timeslicePtrAvg; ts++ {
		if ts == TimesliceLimit {
			if timeslicePtrAvg == 0 {
				// Exception condition where we need to stop the loop
				break
			}
			ts = 0
		}
		timesliceAvg.txAny += float64(timesliceSet[ts].txAny)
		timesliceAvg.txValue += float64(timesliceSet[ts].txValue)
		timesliceAvg.txConfirmed += float64(timesliceSet[ts].txConfirmed)
		timesliceAvg.txToProcess = float64(timesliceSet[ts].txToProcess)
		timesliceAvg.txToBroadcast = float64(timesliceSet[ts].txToBroadcast)
		timesliceAvg.txToReply = float64(timesliceSet[ts].txToReply)
		timesliceAvg.txNumberStoredTx = float64(timesliceSet[ts].txNumberStoredTx)
		timesliceAvg.txTxnToRequest = float64(timesliceSet[ts].txTxnToRequest)
		timesliceCnt++
	}

	// Calculate the average of the timeslice
	if timesliceCnt > 1 {
		//timesliceAvg.txAny = timesliceAvg.txAny / timesliceCnt
		timesliceAvg.txValue = timesliceAvg.txValue / timesliceCnt
		timesliceAvg.txConfirmed = timesliceAvg.txConfirmed / timesliceCnt
		//timesliceAvg.txToProcess = timesliceAvg.txToProcess / timesliceCnt
		//timesliceAvg.txToBroadcast = timesliceAvg.txToBroadcast / timesliceCnt
		//timesliceAvg.txToReply = timesliceAvg.txToReply / timesliceCnt
		//timesliceAvg.txNumberStoredTx = timesliceAvg.txNumberStoredTx / timesliceCnt
		//timesliceAvg.txTxnToRequest = timesliceAvg.txTxnToRequest / timesliceCnt
	}

	// Maintain where we left off for the next call to this function
	timeslicePtrPrev = timeslicePtrAvg

	return timesliceAvg
}

func initZmq(address *string) {

	go manageTimeslice(address)

}
