package messages

import (
	"encoding/json"
	"fmt"

	"github.com/hanchon/garnet/internal/logger"
	"github.com/hanchon/garnet/internal/txbuilder"
)

func connectMessage(ws *WebSocketContainer, registeredUsers map[string]User, p *[]byte) error {
	var connectMsg ConnectMessage
	err := json.Unmarshal(*p, &connectMsg)
	if err != nil {
		return err
	}

	v, ok := registeredUsers[connectMsg.User]
	if !ok {
		return fmt.Errorf("user not registered")
	}

	if v.Password != connectMsg.Password {
		return fmt.Errorf("incorrect password")
	}

	ws.User = connectMsg.User
	ws.Authenticated = true
	ws.WalletID = v.WalletID
	_, account, err := txbuilder.GetWallet(v.WalletID)
	if err != nil {
		logger.LogError(fmt.Sprintf("[backend] could not generate user wallet : %s", ws.User))
	}
	ws.WalletAddress = account.Address.Hex()

	logger.LogInfo(fmt.Sprintf("[backend] user connected: %s (%s)", ws.User, ws.WalletAddress))
	return nil
}

func getMatchStatus(p *[]byte) (string, error) {
	var msg GetMatchStatus
	err := json.Unmarshal(*p, &msg)
	if err != nil {
		return "", err
	}
	return msg.MatchID, nil
}
