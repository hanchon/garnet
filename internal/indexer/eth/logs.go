package eth

import (
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hanchon/garnet/internal/indexer/data/mudhelpers"
)

func OrderLogs(logs []types.Log) []types.Log {
	// TODO: remove this, there is no reorgs on evmos.
	// Filter removed logs due to chain reorgs.
	filteredLogs := []types.Log{}
	for _, log := range logs {
		if !log.Removed {
			filteredLogs = append(filteredLogs, log)
		}
	}

	// Order logs.
	sort.SliceStable(filteredLogs, func(i, j int) bool {
		first := filteredLogs[i]
		second := filteredLogs[j]
		if first.BlockNumber < second.BlockNumber {
			return true
		} else if second.BlockNumber < first.BlockNumber {
			return false
		} else {
			if first.TxIndex < second.TxIndex {
				return true
			} else if second.TxIndex < first.TxIndex {
				return false
			} else {
				return first.Index < second.Index
			}
		}
	})

	return filteredLogs
}

func QueryForStoreLogs(initBlockHeight *big.Int, endBlockHeight *big.Int) ethereum.FilterQuery {
	if initBlockHeight == nil {
		initBlockHeight = big.NewInt(1)
	}

	// TODO: we should query the blockchain to get the latest block
	if endBlockHeight == nil {
		endBlockHeight = big.NewInt(999999999)
	}

	return ethereum.FilterQuery{
		FromBlock: initBlockHeight,
		ToBlock:   endBlockHeight,
		// Topics:    [][]common.Hash{{}},
		Topics: [][]common.Hash{{
			mudhelpers.GetStoreAbiEventID("StoreSetRecord"),
			mudhelpers.GetStoreAbiEventID("StoreSetField"),
			mudhelpers.GetStoreAbiEventID("StoreDeleteRecord"),
		}},
		Addresses: []common.Address{},
	}
}
