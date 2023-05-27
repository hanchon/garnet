package txbuilder

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/hanchon/garnet/internal/logger"
)

var (
	//go:embed IWorld.abi.json
	iworldAbiJSON []byte

	IWorldABI abi.ABI
)

type AbiStruct []struct {
	Inputs []struct {
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"inputs"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Anonymous bool   `json:"anonymous,omitempty"`
	Outputs   []struct {
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"outputs,omitempty"`
	StateMutability string `json:"stateMutability,omitempty"`
}

func init() {
	// We need to remove everything that is type error because it breaks the abi decoder
	var raw AbiStruct
	err := json.Unmarshal(iworldAbiJSON, &raw)
	if err != nil {
		logger.LogError(fmt.Sprintf("failed to unmarshal abi json: %s", err))
		panic("could not unmarshal json")
	}

	withoutErrors := make(AbiStruct, 0)
	for _, v := range raw {
		if v.Type != "error" {
			withoutErrors = append(withoutErrors, v)
		}
	}

	fixedToAbi, err := json.Marshal(withoutErrors)
	if err != nil {
		logger.LogError(fmt.Sprintf("failed to marshal the fixed data: %s", err))
		panic("failed to marshal the fixed data")
	}

	abiDecoded, err := abi.JSON(strings.NewReader(string(fixedToAbi)))
	if err != nil {
		logger.LogError(fmt.Sprintf("error decoding IWorld abi: %s", err))
		panic("error decoding IWorld abi")
	}
	IWorldABI = abiDecoded
}
