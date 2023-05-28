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
			if connectMessage(ws, g.UsersDatabase, &p) != nil {
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
					err := ws.Conn.WriteJSON(msg)
					if err != nil {
						// TODO: close the connection
						return
					}
					logger.LogDebug(fmt.Sprintf("[backend] sending %d active matches", len(ret)))
				}
			}

		case "placecard":
			if !ws.Authenticated {
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
				logger.LogDebug("[backend] error creating transaction to place card: invalid length")
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
			if !ws.Authenticated {
				return
			}
			err = txbuilder.SendTransaction(ws.WalletID, "creatematch")
			if err != nil {
				// TODO: send response saying that the game could not be created
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to creatematch: %s", err))
			}

		case "joinmatch":
			if !ws.Authenticated {
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
				logger.LogDebug("[backend] error creating transaction to join match: invalid length")
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

		case "endturn":
			if !ws.Authenticated {
				return
			}
			logger.LogDebug("[backend] processing endturn request")

			var msg EndTurn
			err := json.Unmarshal(p, &msg)
			if err != nil {
				logger.LogError(fmt.Sprintf("[backend] error decoding endturn message: %s", err))
				return
			}

			logger.LogDebug(fmt.Sprintf("[backend] creating endturn tx: %s", msg.MatchID))

			id, err := hexutil.Decode(msg.MatchID)
			if err != nil {
				// TODO: send response saying that the game could not be created
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to endturn: %s", err))
				return
			}

			if len(id) != 32 {
				logger.LogDebug("[backend] error creating transaction to endturn: invalid length")
				return
			}

			// It must be array instead of slice
			var idArray [32]byte
			copy(idArray[:], id)

			err = txbuilder.SendTransaction(ws.WalletID, "endturn", idArray)
			if err != nil {
				// TODO: send response saying that the game could not be created
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to endturn: %s", err))
				return
			}

		case "movecard":
			if !ws.Authenticated {
				return
			}

			logger.LogDebug("[backend] processing move card request")

			var msg MoveCard
			err := json.Unmarshal(p, &msg)
			if err != nil {
				logger.LogError(fmt.Sprintf("[backend] error decoding move card message: %s", err))
				return
			}

			id, err := hexutil.Decode(msg.CardID)
			if err != nil {
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to move card: %s", err))
				return
			}

			if len(id) != 32 {
				logger.LogDebug("[backend] error creating transaction to move card: invalid length")
				return
			}

			// It must be array instead of slice
			var idArray [32]byte
			copy(idArray[:], id)

			err = txbuilder.SendTransaction(ws.WalletID, "movecard", idArray, uint32(msg.X), uint32(msg.Y))
			if err != nil {
				// TODO: send response saying that the game could not be created
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to move card: %s", err))
			}

		case "attack":
			if !ws.Authenticated {
				return
			}

			logger.LogDebug("[backend] processing attack request")

			var msg Attack
			err := json.Unmarshal(p, &msg)
			if err != nil {
				logger.LogError(fmt.Sprintf("[backend] error decoding attack message: %s", err))
				return
			}

			id, err := hexutil.Decode(msg.CardID)
			if err != nil {
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to attack: %s", err))
				return
			}

			if len(id) != 32 {
				logger.LogDebug("[backend] error creating transaction to attack: invalid length")
				return
			}

			// It must be array instead of slice
			var idArray [32]byte
			copy(idArray[:], id)

			err = txbuilder.SendTransaction(ws.WalletID, "attack", idArray, uint32(msg.X), uint32(msg.Y))
			if err != nil {
				// TODO: send response saying that the game could not be created
				logger.LogDebug(fmt.Sprintf("[backend] error creating transaction to attack: %s", err))
			}
		}
	}
}
