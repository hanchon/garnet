package mudhandlers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hanchon/garnet/internal/indexer/data/mudhelpers"
)

func ParseStoreSetRecord(log types.Log) (*mudhelpers.StorecoreStoreSetRecord, error) {
	event := new(mudhelpers.StorecoreStoreSetRecord)
	if err := UnpackLog(event, "StoreSetRecord", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func ParseStoreSetField(log types.Log) (*mudhelpers.StorecoreStoreSetField, error) {
	event := new(mudhelpers.StorecoreStoreSetField)
	if err := UnpackLog(event, "StoreSetField", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func ParseStoreDeleteRecord(log types.Log) (*mudhelpers.StorecoreStoreDeleteRecord, error) {
	event := new(mudhelpers.StorecoreStoreDeleteRecord)
	if err := UnpackLog(event, "StoreDeleteRecord", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func UnpackLog(out interface{}, eventName string, log types.Log) error {
	if log.Topics[0] != mudhelpers.StorecoreAbi.Events[eventName].ID {
		return fmt.Errorf("event signature mismatch")
	}
	if len(log.Data) > 0 {
		if err := mudhelpers.StorecoreAbi.UnpackIntoInterface(out, eventName, log.Data); err != nil {
			fmt.Println("unpack into interface")
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range mudhelpers.StorecoreAbi.Events[eventName].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return abi.ParseTopics(out, indexed, log.Topics[1:])
}
