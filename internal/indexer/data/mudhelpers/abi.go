package mudhelpers

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hanchon/garnet/internal/logger"
)

var StorecoreAbi abi.ABI
var events []string = []string{"StoreSetRecord", "StoreSetField", "StoreDeleteRecord"}

func init() {
	var err error
	StorecoreAbi, err = abi.JSON(strings.NewReader(string(StorecoreMetaData.ABI)))
	if err != nil {
		logger.LogError("failed to parse the store ABI")
		panic("")
	}
}

func GetStoreAbiEventID(eventName string) common.Hash {
	return StorecoreAbi.Events[eventName].ID
}

func (event *StorecoreStoreSetRecord) WorldAddress() string {
	return event.Raw.Address.Hex()
}

func (event *StorecoreStoreSetField) WorldAddress() string {
	return event.Raw.Address.Hex()
}

func (event *StorecoreStoreDeleteRecord) WorldAddress() string {
	return event.Raw.Address.Hex()
}
