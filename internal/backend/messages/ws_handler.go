package messages

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/hanchon/garnet/internal/logger"
	"github.com/hanchon/garnet/internal/txbuilder"
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

const WorldID = "0x5FbDB2315678afecb367f032d93F642f64180aa3"

func (g *GlobalState) WsHandler(ws *WebSocketContainer) {
	for {
		defer removeConnection(ws, g)
		// Read until error the client messages
		_, p, err := ws.Conn.ReadMessage()
		if err != nil {
			return
		}

		// TODO: log ip address
		logger.LogDebug(fmt.Sprintf("[backend] incoming message: %s", string(p)))

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

			// Send response
			msg := `{"msgtype":"connected", "status":true}`
			if writeMessage(ws.Conn, &msg) != nil {
				return
			}

			logger.LogDebug(fmt.Sprintf("[backend] senging message: %s", msg))

			g.WsSockets[ws.User] = ws.Conn

			w, ok := g.Database.Worlds[WorldID]
			if !ok {
				panic("world not found")
			}

			t := w.GetTableByName("Match")
			if t != nil {
				ret := []string{}
				for k := range *t.Rows {
					ret = append(ret, k)
				}
				msg := MatchList{MsgType: "matchlist", Matches: ret}
				ws.Conn.WriteJSON(msg)
				logger.LogDebug(fmt.Sprintf("[backend] sending %d active matches", len(ret)))
			}
			// index, ok := g.WalletIndex[ws.user]
			// _, ok := g.WalletIndex[ws.User]
			// if ok {
			// 	// game := g.Games[index]
			// 	// json, _ := game.GameStatus()
			// 	// state := fmt.Sprintf(`{"msgtype":"boardstate","value":%s,"id": "%s"}`, string(json), game.Id)
			// 	// if writeMessage(ws.conn, &state) != nil {
			// 	// 	return
			// 	// }
			// 	// fmt.Println("send all the games active")
			// }

		case "creatematch":
			if ws.Authenticated == false {
				return
			}
			err = txbuilder.SendTransaction(ws.WalletID, "creatematch")
			if err != nil {
				// TODO: send response saying that the game could not be created
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to creatematch: %s", err))
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
		// fmt.Println(string(p))
	}
}
