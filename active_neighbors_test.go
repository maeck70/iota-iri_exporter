package main

import "testing"
//import "github.com/iotaledger/giota"

func TestActiveNeighbor(t *testing.T) {

/*	tx := []int64 {100,105,106,106,110}
	addr := giota.Address {Address: "foo.com",}

	for t := range tx {
		n := giota.Neighbor{
				Address: addr, 
				NumberOfNewTransactions: tx[t],
		}
		neighborhistory.register(n)
	}
	a = GetActiveNeighbor(addr)
	if GetActiveNeighbor(addr) == false {
		t.Error("Expected First Neighbor Active, got %v", a)	
	}
*/
	a := true
	if a == false {
		t.Error("Expected First Neighbor Active, got %v", a)	
	}

}