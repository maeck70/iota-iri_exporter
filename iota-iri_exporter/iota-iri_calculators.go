package iota_iri_exporter

import (
	"github.com/iotaledger/giota"
)

func getActiveNeighbors(Neighbors []giota.Neighbor) float64 {

	return 0
}

/*import (
	"fmt"
	"github.com/iotaledger/giota"
)

const NEIGHBOR_LIMIT = 128
const HISTORY_LIMIT = 5

var (
	NeighborCols [NEIGHBOR_LIMIT]giota.Neighbor
	NeighborRows [HISTORY_LIMIT]NeighborCols
	rowPtr       = HISTORY_LIMIT
)

func getSum(n giota.Neighbor) int64 {
	// TODO: Update with all Neighbor Transaction types as soon as the Iota Go Lib has been updates to the current API
	return n.NumberOfAllTransactions + n.NumberOfInvalidTransactions + n.NumberOfNewTransactions
}

func getFullAddress(n giota.Neighbor) string {
	// TODO: Update with all Neighbor Transaction types as soon as the Iota Go Lib has been updates to the current API
	// Possibly insert this int the API
	//vfunc (n giota.Neighbor) FullAddress() string {
	return fmt.Sprintf("%s//%s", "tcp", n.Address)
}

func getActiveNeighbors(Neighbors []giota.Neighbor) float64 {

	// Rotate History
	rowPtr += 1
	if rowPtr > HISTORY_LIMIT {
		rowPtr = 1
	}

	// Insert current
	neighbor_cnt := len(Neighbors)
	if neighbor_cnt > NEIGHBOR_LIMIT {
		neighbor_cnt = NEIGHBOR_LIMIT
	}

	for n := 1; n < neighbor_cnt; n++ {
		NeighborRows[rowPtr].NeighborCols[n] = Neighbors[n]
	}

	return 0
}
*/
