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
	"crypto/sha1"
	"fmt"
	"github.com/iotaledger/giota"
	"github.com/prometheus/common/log"
)

const neighborMax = 32
const historyMax = 5

var neighbormatrix = neighborMatrix{}

func getNeighborHash(n giota.Neighbor) string {
	s := fmt.Sprintf("%s %d %d %d", n.Address, n.NumberOfAllTransactions, n.NumberOfInvalidTransactions, n.NumberOfNewTransactions)
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

type neighborHistory struct {
	address string
	history [historyMax]string //Neighbor
}

type neighborMatrix struct {
	ptr             int
	neighborhistory [neighborMax]neighborHistory
}

func (nm *neighborMatrix) findNeighbor(addr string) (*neighborHistory, int) {
	emptyslot := -1
	nh := neighborHistory{}
	n := 0
	for n = range nm.neighborhistory {
		if nm.neighborhistory[n].address == addr {
			nh = nm.neighborhistory[n]
			break
		} else if nm.neighborhistory[n].address == "" && emptyslot == -1 {
			nm.neighborhistory[n].address = addr
			emptyslot = n
			nh = nm.neighborhistory[n]
			break
		}
	}
	return &nh, n
}

func (nm *neighborMatrix) register(n giota.Neighbor) {
	_, slot := nm.findNeighbor(string(n.Address))
	nm.neighborhistory[slot].address = string(n.Address)
	nm.neighborhistory[slot].history[nm.ptr] = getNeighborHash(n)
}

func (nm *neighborMatrix) isActive(addr string) bool {
	nh, _ := nm.findNeighbor(string(addr))
	compare := nh.history[0]
	for p := 1; p < historyMax; p++ {
		if compare != nh.history[p] {
			return true
		}
	}
	return false
}

func (nm *neighborMatrix) historyInc() {
	nm.ptr++
	if nm.ptr >= historyMax {
		nm.ptr = 0
	}
}

func getActiveNeighbor(addr string) float64 {
	status := neighbormatrix.isActive(addr)
	log.Debugf("Neighbor with address %s active status is %v", addr, status)
	return btof(status)
}

func getActiveNeighbors(neighborlist []giota.Neighbor) float64 {

	neighbormatrix.historyInc()

	activeCount := 0
	for n := range neighborlist {
		neighbormatrix.register(neighborlist[n])
		activeCount += btoi(neighbormatrix.isActive(string(neighborlist[n].Address)))
	}
	log.Debugf("There are %v of %v active Neighbors.", activeCount, len(neighborlist))
	return float64(activeCount)
}
