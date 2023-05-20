package messages

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func writeMessage(ws *websocket.Conn, msg *string) error {
	return ws.WriteMessage(websocket.TextMessage, []byte(*msg))
}

func writeBytes(ws *websocket.Conn, msg []byte) error {
	return ws.WriteMessage(websocket.TextMessage, msg)
}

func removeConnection(ws *WebSocketContainer, g *GlobalState) {
	ws.Conn.Close()
	delete(g.WsSockets, ws.User)
}

func (g *GlobalState) WsHandler(ws *WebSocketContainer) {
	for {
		defer removeConnection(ws, g)
		// Read until error the client messages
		_, p, err := ws.Conn.ReadMessage()
		if err != nil {
			return
		}

		fmt.Println(string(p))

		var m BasicMessage
		err = json.Unmarshal(p, &m)
		if err != nil {
			return
		}

		switch m.MsgType {
		case "connect":
			if connectMessage(ws, g.RegisteredUsers, &p) != nil {
				return
			}

			// transactions.SendTransaction(0, "creatematch")

			// Send response
			msg := `{"msgtype":"connected", "status":true}`
			if writeMessage(ws.Conn, &msg) != nil {
				return
			}

			g.WsSockets[ws.User] = ws.Conn
			// index, ok := g.WalletIndex[ws.user]
			_, ok := g.WalletIndex[ws.User]
			if ok {
				// game := g.Games[index]
				// json, _ := game.GameStatus()
				// state := fmt.Sprintf(`{"msgtype":"boardstate","value":%s,"id": "%s"}`, string(json), game.Id)
				// if writeMessage(ws.conn, &state) != nil {
				// 	return
				// }
				fmt.Println("send all the games active")
			}
		case "getmatchstatus":
			{
				if ws.Authenticated == false {
					return
				}
				// id, err := getMatchStatus(&p)
				// if err != nil {
				// 	// TODO: kill the connectiong and remove this ws from the list
				// 	return
				// }
				// if DummyData.MatchID == id {
				// 	value, err := json.Marshal(RespMatchStatus{MsgType: "respmatchstatus", MatchData: DummyData})
				// 	if err != nil {
				// 		// TODO: kill the connectiong and remove this ws from the list
				// 		return
				// 	}
				// 	if writeBytes(ws.conn, value) != nil {
				// 		// TODO: kill the connectiong and remove this ws from the list
				// 		return
				// 	}
				//
				// 	return
				// }
			}
		}

		// print out that message for clarity
		fmt.Println(string(p))
	}
}
