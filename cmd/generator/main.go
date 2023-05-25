package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hanchon/garnet/internal/logger"
	"github.com/hanchon/garnet/internal/txbuilder"
)

func main() {
	init := 100
	end := 115

	logger.LogInfo("Sending coins to all the wallets")

	for i := init; i < end; i++ {
		_, account, err := txbuilder.GetWallet(i)
		if err != nil {
			logger.LogError(fmt.Sprintf("Error generating the wallet %d\n", i))
		}
		err = txbuilder.Faucet(account.Address.String())
		if err != nil {
			log.Panic(err)
		}
	}
	time.Sleep(time.Second)

	for i := init; i < end; i++ {
		err := txbuilder.SendTransaction(i, "creatematch")
		if err != nil {
			log.Panic(err)
		}
		time.Sleep(1 * time.Second)
	}
}
