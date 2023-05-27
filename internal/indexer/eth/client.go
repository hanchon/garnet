package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/indexer/data/mudhelpers"
	"github.com/hanchon/garnet/internal/indexer/eth/mudhandlers"
	"github.com/hanchon/garnet/internal/logger"
)

func GetEthereumClient(wsURL string) *ethclient.Client {
	var client *ethclient.Client
	var err error
	client, err = ethclient.Dial(wsURL)
	if err != nil {
		// TODO: add retry in case of failure instead of panic
		logger.LogError("could not connect to the ethereum client")
		panic("")
	}
	return client
}

func ProcessBlocks(c *ethclient.Client, db *data.Database, initBlockHeight *big.Int, endBlockHeight *big.Int) {
	logs, err := c.FilterLogs(context.Background(), QueryForStoreLogs(initBlockHeight, endBlockHeight))
	if err != nil {
		// TODO: add retry in case of failure instead of panic
		logger.LogError("error filtering blocks")
		panic("")
	}
	logs = OrderLogs(logs)
	logger.LogInfo(fmt.Sprintf("processing logs up to %d", endBlockHeight))

	for _, v := range logs {
		if v.Topics[0].Hex() == mudhelpers.GetStoreAbiEventID("StoreSetRecord").Hex() {
			event, err := mudhandlers.ParseStoreSetRecord(v)
			if err != nil {
				logger.LogError(fmt.Sprintf("error decoding message:%s", err))
				// TODO: what should we do here?
				break
			}
			switch mudhelpers.PaddedTableId(event.TableId) {
			case mudhelpers.SchemaTableId():
				logger.LogInfo("processing and creating schema table")
				mudhandlers.HandleSchemaTableEvent(event, db)
			case mudhelpers.MetadataTableId():
				logger.LogInfo("processing and updating a schema with metadata")
				mudhandlers.HandleMetadataTableEvent(event, db)
			default:
				logger.LogInfo("processing a generic table event like adding a row")
				mudhandlers.HandleGenericTableEvent(event, db)
			}
		}

		if v.Topics[0].Hex() == mudhelpers.GetStoreAbiEventID("StoreSetField").Hex() {
			event, err := mudhandlers.ParseStoreSetField(v)
			logger.LogInfo("processing store set field message")
			if err != nil {
				logger.LogError(fmt.Sprintf("error decoding message for store set field:%s\n", err))
			} else {
				mudhandlers.HandleSetFieldEvent(event, db)
			}
		}
		if v.Topics[0].Hex() == mudhelpers.GetStoreAbiEventID("StoreDeleteRecord").Hex() {
			logger.LogInfo("processing store delete record message")
			event, err := mudhandlers.ParseStoreDeleteRecord(v)
			if err != nil {
				logger.LogError(fmt.Sprintf("error decoding message for store delete record:%s\n", err))
			} else {
				mudhandlers.HandleDeleteRecordEvent(event, db)
			}
		}
	}
}
