package messages

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("connected")
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
