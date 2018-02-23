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

import "testing"
import "github.com/maeck70/giota"

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
		{transactionCount: 120, result: 1},
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
			t.Errorf("Test %v: Expected Neighbor to be %v, got %v", i, tx[i].result, a)
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
		{transactionCount: 120, result: 1},
	}

	addr := "foo.com"

	nl := []giota.Neighbor{
		{Address: giota.Address(addr)},
	}

	for i := range tx {
		nl[0].NumberOfNewTransactions = tx[i].transactionCount
		a := getActiveNeighbors(nl)
		if a != tx[i].result {
			t.Errorf("Test %v: Expected %v Active Neighbor, got %v", i, tx[i].result, a)
		}
	}
}
