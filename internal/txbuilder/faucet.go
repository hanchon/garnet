package txbuilder

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hanchon/garnet/internal/logger"
)

func Faucet(addr string) (common.Hash, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return [32]byte{}, err
	}

	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		return [32]byte{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return [32]byte{}, err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return [32]byte{}, err
	}

	value := big.NewInt(9000000000000000000) // in wei (1 eth)
	gasLimit := uint64(100000)               // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return [32]byte{}, err
	}

	toAddress := common.HexToAddress(addr)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return [32]byte{}, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return [32]byte{}, err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return [32]byte{}, err
	}

	logger.LogDebug(fmt.Sprintf("[backend] faucet tx sent with hash: %s", signedTx.Hash().Hex()))

	return signedTx.Hash(), nil
}

func TransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	return client.TransactionReceipt(context.Background(), hash)
}

func WasTransactionSuccessful(hash common.Hash) (bool, error) {
	receipt, err := TransactionReceipt(hash)
	if err != nil {
		return false, err
	}
	return receipt.Status == types.ReceiptStatusSuccessful, nil
}
