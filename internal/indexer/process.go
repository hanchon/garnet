package indexer

import (
	"context"
	"math/big"
	"time"

	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/indexer/eth"
	"github.com/hanchon/garnet/internal/logger"
)

func Process(database *data.Database, quit *bool) {
	logger.LogInfo("indexer is starting...")
	c := eth.GetEthereumClient("http://localhost:8545/")
	ctx := context.Background()
	chainId, err := c.ChainID(ctx)
	if err != nil {
		logger.LogError("could not get the latest height")
		// TODO: retry instead of panic
		panic("")
	}
	database.ChainID = chainId.String()

	height, err := c.BlockNumber(context.Background())
	if err != nil {
		logger.LogError("could not get the latest height")
		// TODO: retry instead of panic
		panic("")
	}

	eth.ProcessBlocks(c, database, nil, big.NewInt(int64(height)))

	for *quit == false {
		// TODO: handle control+c
		newHeight, err := c.BlockNumber(context.Background())
		if err != nil {
			logger.LogError("could not get the latest height")
			// TODO: retry instead of panic
			panic("")
		}

		if newHeight != height {
			eth.ProcessBlocks(c, database, big.NewInt(int64(height)), big.NewInt(int64(newHeight)))
			height = newHeight
		}

		database.LastHeight = newHeight

		time.Sleep(1 * time.Second)
	}
}
