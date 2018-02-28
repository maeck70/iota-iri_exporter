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
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/pebbe/zmq4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"strings"
	"time"
)

type zmqAccumsf struct {
	txTotal          float64
	txAnyZero        float64
	txAnyNotZero     float64
	txValue          float64
	txConfirmed      float64
	txToProcess      float64
	txToBroadcast    float64
	txToReply        float64
	txNumberStoredTx float64
	txTxnToRequest   float64
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

type txRecord struct {
	Timestamp   int64
	TxIn        int64
	TxConfirmed int64
	TxAddress   string
	TxValue     int64
}

type zmqConfirmation struct {
	label    string
	duration float64
}

var zmqAccums zmqAccumsf
var zmqConfirmationSet []zmqConfirmation

func getTxLabel(c int64) string {
	label := "0"
	if c != 0 {
		label = "<> 0"
	}
	return label
}

func metricsZmq(e *exporter) {

	e.iotaZmqSeenTxCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_seen_tx_count",
			Name: "iota_zmq_seen_tx_count",
			Help: "Count of transactions seen by zeroMQ.",
		},
		[]string{"hasValue"},
	)

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

	e.iotaZmqConfirmationHisto = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			//Namespace: namespace,
			//Subsystem: "zmq",
			//Name: "zmq_total_transactions",
			Name:    "iota_zmq_tx_confirm_time",
			Help:    "Actual seconds it takes to confirm each tx.",
			Buckets: []float64{300, 600, 1200, 2400, 3600, 7200, 21600, 43200},
		},
		[]string{"hasValue"},
	)

}

func describeZmq(e *exporter, ch chan<- *prometheus.Desc) {

	e.iotaZmqSeenTxCount.Describe(ch)
	ch <- e.iotaZmqTxsWithValueCount.Desc()
	ch <- e.iotaZmqConfirmedTxCount.Desc()
	ch <- e.iotaZmqToProcess.Desc()
	ch <- e.iotaZmqToBroadcast.Desc()
	ch <- e.iotaZmqToRequest.Desc()
	ch <- e.iotaZmqToReply.Desc()
	ch <- e.iotaZmqTotalTransactions.Desc()
	e.iotaZmqConfirmationHisto.Describe(ch)
}

func collectZmq(e *exporter, ch chan<- prometheus.Metric) {

	e.iotaZmqSeenTxCount.Collect(ch)
	ch <- e.iotaZmqTxsWithValueCount
	ch <- e.iotaZmqConfirmedTxCount
	ch <- e.iotaZmqToProcess
	ch <- e.iotaZmqToBroadcast
	ch <- e.iotaZmqToRequest
	ch <- e.iotaZmqToReply
	ch <- e.iotaZmqTotalTransactions
	e.iotaZmqConfirmationHisto.Collect(ch)
}

func scrapeZmq(e *exporter) {

	e.iotaZmqSeenTxCount.WithLabelValues("<> 0").Set(zmqAccums.txAnyNotZero)
	e.iotaZmqSeenTxCount.WithLabelValues("0").Set(zmqAccums.txAnyZero)
	e.iotaZmqTxsWithValueCount.Set(zmqAccums.txValue)
	e.iotaZmqConfirmedTxCount.Set(zmqAccums.txConfirmed)
	e.iotaZmqToProcess.Set(zmqAccums.txToProcess)
	e.iotaZmqToBroadcast.Set(zmqAccums.txToBroadcast)
	e.iotaZmqToRequest.Set(zmqAccums.txTxnToRequest)
	e.iotaZmqToReply.Set(zmqAccums.txToReply)
	e.iotaZmqTotalTransactions.Set(zmqAccums.txTotal)

	for i := range zmqConfirmationSet {
		log.Infof("zmqConfirmationSet[%d] = [%s] %v", i, zmqConfirmationSet[i].label, zmqConfirmationSet[i].duration)
		e.iotaZmqConfirmationHisto.WithLabelValues(zmqConfirmationSet[i].label).Observe(zmqConfirmationSet[i].duration)
	}
	zmqConfirmationSet = nil

	log.Debugf("total tx:         %v tx", int64(zmqAccums.txTotal))
	log.Debugf("txAnyZero:        %v tx", int64(zmqAccums.txAnyZero))
	log.Debugf("txAnyNotZero:     %v tx", int64(zmqAccums.txAnyNotZero))
	log.Debugf("txValue:          %v tx", int64(zmqAccums.txValue))
	log.Debugf("txConfirmed:      %v tx", int64(zmqAccums.txConfirmed))
	log.Debugf("txToProcess:      %v tx", int64(zmqAccums.txToProcess))
	log.Debugf("txToBroadcast:    %v tx", int64(zmqAccums.txToBroadcast))
	log.Debugf("txToReply:        %v tx", int64(zmqAccums.txToReply))
	log.Debugf("txNumberStoredTx: %v tx", int64(zmqAccums.txNumberStoredTx))
	log.Debugf("txTxnToRequest:   %v tx", int64(zmqAccums.txTxnToRequest))

}

func collectZmqAccums(address *string) {

	for {

		socket, err := zmq4.NewSocket(zmq4.SUB)
		must(err)

		for _, topic := range []string{"tx", "sn", "rstat"} {
			err = socket.SetSubscribe(topic)
			must(err)
		}

		// Set ZMQ no received messages time-out
		err = socket.SetRcvtimeo(10 * time.Second)
		must(err)

		err = socket.Connect(*address)
		must(err)

		log.Infof("Connected to IRI at address %s.", *address)

		opts := badger.DefaultOptions
		opts.Dir = *databasePath
		opts.ValueDir = *databasePath
		db, err := badger.Open(opts)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		go badgerDBCleanup(db)

		for {

			msg, err := socket.Recv(0)
			if err == zmq4.ETIMEDOUT {
				log.Info("No ZMQ RStat msg received, reconnecting to zmq socket.")
				break
			} else if err != nil {
				panic(err)
			}

			parts := strings.Fields(msg)
			switch parts[0] {

			// Transaction
			case "tx":
				zmqAccums.txTotal++
				tx := transaction{
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

				processValueTx(db, &tx)
				if tx.Value != 0 {
					zmqAccums.txAnyNotZero++
					zmqAccums.txValue++
					log.Debug("ZMQ Tx with value msg received.")
				} else {
					zmqAccums.txAnyZero++
					//log.Debug("ZMQ Tx with zero value msg received.")
				}

			// Confirmed Transaction
			case "sn":
				sn := sn{
					Index:       parts[1],
					Hash:        parts[2],
					AddressHash: parts[3],
					Trunk:       parts[4],
					Branch:      parts[5],
					Bundle:      parts[6],
				}
				zmqAccums.txConfirmed++
				log.Debug("ZMQ Confirmed Tx msg received.")
				go processConfirmedTx(db, &sn)

			// RStat message (overall statistics)
			case "rstat":
				stat := queue{
					ReceiveQueueSize:   stoi(parts[1]),
					BroadcastQueueSize: stoi(parts[2]),
					TxnToRequest:       stoi(parts[3]),
					ReplyQueueSize:     stoi(parts[4]),
					NumberOfStoredTxns: stoi(parts[5]),
				}

				// Note that these are total counts, no need to increment into the timeslice
				zmqAccums.txToProcess = float64(stat.ReceiveQueueSize)
				zmqAccums.txToBroadcast = float64(stat.BroadcastQueueSize)
				zmqAccums.txToReply = float64(stat.ReplyQueueSize)
				zmqAccums.txNumberStoredTx = float64(stat.NumberOfStoredTxns)
				zmqAccums.txTxnToRequest = float64(stat.TxnToRequest)

				log.Debug("ZMQ RStat msg received.")
			}
		}
	}
}

func processValueTx(db *badger.DB, tx *transaction) {

	recttl := 15 * 24 * time.Hour // 15 Days
	err := db.Update(func(txn *badger.Txn) error {

		key := fmt.Sprintf("%s", tx.Hash)

		rec := txRecord{
			Timestamp:   stoi(tx.Timestamp),
			TxIn:        stoi(time.Now().UTC().Format("20060102150405")),
			TxConfirmed: 0,
			TxAddress:   tx.Address,
			TxValue:     tx.Value,
		}

		val, err := json.Marshal(rec)
		//log.Debugf("BadgerDB write: key(%s) value(%s)", key, val)
		err = txn.SetWithTTL([]byte(key), val, recttl)

		return err
	})
	if err != nil {
		log.Infof("BadgerDB error %v.", err)
	}
}

func processConfirmedTx(db *badger.DB, tx *sn) {

	recttl := 24 * time.Hour // 1 Day
	err := db.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("%s", tx.Hash)
		val, err := txn.Get([]byte(key))

		if err != badger.ErrKeyNotFound {

			v, _ := val.Value()
			log.Debugf("BadgerDB get: key(%s) value(%s)", key, v)

			rec := txRecord{}
			json.Unmarshal(v, &rec)

			log.Infof("rec: %v.", rec)

			rec.TxConfirmed = stoi(time.Now().UTC().Format("20060102150405"))
			v, _ = json.Marshal(rec)
			err = txn.SetWithTTL([]byte(key), v, recttl)
			c := zmqConfirmation{label: getTxLabel(rec.TxValue), duration: float64(rec.TxConfirmed - rec.TxIn)}
			zmqConfirmationSet = append(zmqConfirmationSet, c)

		} else {
			log.Debugf("BadgerDB get: Key(%s) not found", key)
		}
		return err
	})
	if err != nil {
		if err != badger.ErrKeyNotFound {
			log.Infof("BadgerDB error %v.", err)
		}
	}
}

func badgerDBCleanup(db *badger.DB) {

	// Cleanup every 15 minutes
	for {
		time.Sleep(15 * time.Minute)
		db.PurgeOlderVersions()
		db.RunValueLogGC(0.5)
		log.Info("BadgerDB purge.")
	}
}

func initZmq(address *string) {
	major, minor, patch := zmq4.Version()
	log.Infof("ZMQ version is %d.%d.%d", major, minor, patch)

	go collectZmqAccums(address)
}
