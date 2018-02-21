package main

import "fmt"
import "testing"
import "github.com/iotaledger/giota"


type neighborTest struct {
	tx_count int64
	result float64
}

func TestActiveNeighbor(t *testing.T) {

	tx := []neighborTest {
		{tx_count: 100, result: 1},
		{tx_count: 105, result: 1},
		{tx_count: 110, result: 1},
		{tx_count: 115, result: 1},
		{tx_count: 116, result: 1},
		{tx_count: 116, result: 1},
		{tx_count: 116, result: 1},
		{tx_count: 116, result: 1},
		{tx_count: 116, result: 0},
	}

	addr := "foo.com"

	a := float64(0)

	for i := range tx {
		n := giota.Neighbor{
				Address: giota.Address(addr), 
				NumberOfNewTransactions: tx[i].tx_count,
		}
		neighbormatrix.historyInc()
		neighbormatrix.register(n)

		a = GetActiveNeighbor(addr)
		if GetActiveNeighbor(addr) != tx[i].result {
			t.Errorf("Expected Neighbor Active (1), got %v", a)	
		}
	}

}


func TestActiveNeighbors(t *testing.T) {

	tx := []neighborTest {
		{tx_count: 100, result: 1},
		{tx_count: 105, result: 1},
		{tx_count: 110, result: 1},
		{tx_count: 115, result: 1},
		{tx_count: 116, result: 1},
		{tx_count: 116, result: 1},
		{tx_count: 116, result: 1},
		{tx_count: 116, result: 1},
		{tx_count: 116, result: 1},
	}

	addr := "foo.com"

	nl := []giota.Neighbor {
		{Address: giota.Address(addr)},
	}


	a := float64(0)

	for i := range tx {
		n := giota.Neighbor{
				Address: giota.Address(addr), 
				NumberOfNewTransactions: tx[i].tx_count,
		}

		fmt.Print(GetActiveNeighbors(nl))

		if i == 4 {
			a = GetActiveNeighbors(nl)
			if a != tx[i].result {
				t.Errorf("Expected One Active Neighbor, got %v", a)	
			}
		} else if i == 8 {
			a = GetActiveNeighbors(nl)
			if a != tx[i].result {
				t.Errorf("Expected No Active Neighbors, got %v", a)	
			}
		} else {
			neighbormatrix.historyInc()
			neighbormatrix.register(n)			
		}
	}
}
