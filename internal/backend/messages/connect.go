package messages

import (
	"encoding/json"
	"fmt"

	"github.com/hanchon/garnet/internal/database"
	"github.com/hanchon/garnet/internal/logger"
)

func connectMessage(ws *WebSocketContainer, usersDB *database.InMemoryDatabase, p *[]byte) error {
	var connectMsg ConnectMessage
	err := json.Unmarshal(*p, &connectMsg)
	if err != nil {
		return err
	}

	user, err := usersDB.Login(connectMsg.User, connectMsg.Password)
	if err != nil {
		return fmt.Errorf("incorrect credentials")
	}

	ws.User = connectMsg.User
	ws.Authenticated = true
	ws.WalletID = user.Index
	ws.WalletAddress = user.Address

	logger.LogInfo(fmt.Sprintf("[backend] user connected: %s (%s)", ws.User, ws.WalletAddress))
	return nil
}
