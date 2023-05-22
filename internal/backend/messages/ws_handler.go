package messages

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
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

			g.WsSockets[ws.User] = ws

			w, ok := g.Database.Worlds[WorldID]
			if !ok {
				panic("world not found")
			}

			matchData := g.Database.GetBoardStatus(WorldID, ws.WalletAddress)
			if matchData != nil {
				// TODO: save the user and wallet somewhere
				matchData.PlayerOneUsermane = "user1"
				matchData.PlayerTwoUsermane = "user2"
				msgToSend := BoardStatus{MsgType: "boardstatus", Status: *matchData}
				logger.LogDebug(fmt.Sprintf("[backend] sending match info %s to %s", matchData.MatchID, ws.User))
				err := ws.Conn.WriteJSON(msgToSend)
				if err != nil {
					panic("could not send the board status")
				}
			} else {
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

		case "placecard":
			if ws.Authenticated == false {
				return
			}

			logger.LogDebug("[backend] processing place card request")

			var msg PlaceCard
			err := json.Unmarshal(p, &msg)
			if err != nil {
				logger.LogError(fmt.Sprintf("[backend] error decoding place card message: %s", err))
				return
			}

			id, err := hexutil.Decode(msg.CardID)
			if err != nil {
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to place card: %s", err))
				return
			}

			if len(id) != 32 {
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to place card: invalid length"))
				return
			}

			// It must be array instead of slice
			var idArray [32]byte
			copy(idArray[:], id)

			err = txbuilder.SendTransaction(ws.WalletID, "placecard", idArray, uint32(msg.X), uint32(msg.Y))
			if err != nil {
				// TODO: send response saying that the game could not be created
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to place card: %s", err))
			}

		case "creatematch":
			if ws.Authenticated == false {
				return
			}
			err = txbuilder.SendTransaction(ws.WalletID, "creatematch")
			if err != nil {
				// TODO: send response saying that the game could not be created
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to creatematch: %s", err))
			}

		case "joinmatch":
			if ws.Authenticated == false {
				return
			}
			logger.LogDebug("[backend] processing join match request")

			var msg JoinMatch
			err := json.Unmarshal(p, &msg)
			if err != nil {
				logger.LogError(fmt.Sprintf("[backend] error decoding join match message: %s", err))
				return
			}

			logger.LogDebug(fmt.Sprintf("[backend] creating join match tx: %s", msg.MatchID))

			id, err := hexutil.Decode(msg.MatchID)
			if err != nil {
				// TODO: send response saying that the game could not be created
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to join match: %s", err))
				return
			}

			if len(id) != 32 {
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to join match: invalid length"))
				return
			}

			// It must be array instead of slice
			var idArray [32]byte
			copy(idArray[:], id)

			err = txbuilder.SendTransaction(ws.WalletID, "joinmatch", idArray)
			if err != nil {
				// TODO: send response saying that the game could not be created
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to join match: %s", err))
				return
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
