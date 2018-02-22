package main

import "testing"
import "github.com/iotaledger/giota"

type neighborTest struct {
	transactionCount int64
	result           float64
}

func TestActiveNeighbor(t *testing.T) {

	tx := []neighborTest{
		{transactionCount: 100, result: 1},
		{transactionCount: 105, result: 1},
		{transactionCount: 110, result: 1},
		{transactionCount: 115, result: 1},
		{transactionCount: 116, result: 1},
		{transactionCount: 116, result: 1},
		{transactionCount: 116, result: 1},
		{transactionCount: 116, result: 1},
		{transactionCount: 116, result: 0},
	}

	addr := "foo.com"

	a := float64(0)

	for i := range tx {
		n := giota.Neighbor{
			Address:                 giota.Address(addr),
			NumberOfNewTransactions: tx[i].transactionCount,
		}
		neighbormatrix.historyInc()
		neighbormatrix.register(n)

		a = getActiveNeighbor(addr)
		if getActiveNeighbor(addr) != tx[i].result {
			t.Errorf("Expected Neighbor Active (1), got %v", a)
		}
	}

}

func TestActiveNeighbors(t *testing.T) {

	tx := []neighborTest{
		{transactionCount: 100, result: 1},
		{transactionCount: 105, result: 1},
		{transactionCount: 110, result: 1},
		{transactionCount: 115, result: 1},
		{transactionCount: 116, result: 1},
		{transactionCount: 116, result: 1},
		{transactionCount: 116, result: 1},
		{transactionCount: 116, result: 1},
		{transactionCount: 116, result: 0},
	}

	addr := "foo.com"

	nl := []giota.Neighbor{
		{Address: giota.Address(addr)},
	}

	//a := float64(0)

	for i := range tx {
		nl[0].NumberOfNewTransactions = tx[i].transactionCount
		a := getActiveNeighbors(nl)
		if a != tx[i].result {
			t.Errorf("Test %v: Expected %v Active Neighbor, got %v", i, tx[i].result, a)
		}
	}
}
