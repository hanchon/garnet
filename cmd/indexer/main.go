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
	process()
}

func process() {
	database := data.NewDatabase()
	logger.LogInfo("indexer is starting...")
	c := eth.GetEthereumClient("http://localhost:8545/")
	height, err := c.BlockNumber(context.Background())
	if err != nil {
		logger.LogError("could not get the latest height")
		// TODO: retry instead of panic
		panic("")
	}

	eth.ProcessBlocks(c, &database, nil, big.NewInt(int64(height)))

	for {
		newHeight, err := c.BlockNumber(context.Background())
		if err != nil {
			logger.LogError("could not get the latest height")
			// TODO: retry instead of panic
			panic("")
		}

		if newHeight != height {
			eth.ProcessBlocks(c, &database, big.NewInt(int64(height)), big.NewInt(int64(newHeight)))
			height = newHeight
		}

		time.Sleep(1 * time.Second)
	}
}
