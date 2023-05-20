package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/indexer/eth"
	"github.com/hanchon/garnet/internal/logger"
)

func main() {
	fileName := "logFile.log"
	// open log file
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)

	// Index the database
	quit := false
	database := data.NewDatabase()
	go process(&database, &quit)
	ui := NewDebugUI()
	defer ui.ui.Close()

	go ui.ProcessIncomingData(&database)
	go ui.ProcessBlockchainInfo(&database)
	go ui.ProcessLatestEvents(&database)
	ui.Run()
	quit = true
}

func process(database *data.Database, quit *bool) {
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
