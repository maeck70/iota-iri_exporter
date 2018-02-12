package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/iotaledger/giota"
	//"github.com/prometheus/common/log"
)

const NEIGHBOR_MAX = 32
const HISTORY_MAX = 5

var neighbormatrix = neighborMatrix{}

func getNeighborHash(n giota.Neighbor) string {
	s := fmt.Sprintf("%s %s %s %s", n.Address, n.NumberOfAllTransactions, n.NumberOfInvalidTransactions, n.NumberOfNewTransactions)
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

type neighborHistory struct {
	address string
	history [HISTORY_MAX]string //Neighbor
}

type neighborMatrix struct {
	ptr             int
	neighborhistory [NEIGHBOR_MAX]neighborHistory
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
	for p := 1; p < HISTORY_MAX; p++ {
		if compare != nh.history[p] {
			return true
		}
	}
	return false
}

func (nm *neighborMatrix) historyInc() {
	nm.ptr += 1
	if nm.ptr >= HISTORY_MAX {
		nm.ptr = 0
	}
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func btof(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func GetActiveNeighbor(addr string) float64 {
	return btof(neighbormatrix.isActive(addr))
}

func GetActiveNeighbors(neighborlist []giota.Neighbor) float64 {

	neighbormatrix.historyInc()
	active_count := 0
	for n := range neighborlist {
		neighbormatrix.register(neighborlist[n])
		active_count += btoi(neighbormatrix.isActive(string(neighborlist[n].Address)))
	}

	return float64(active_count)
}
