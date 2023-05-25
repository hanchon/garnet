package messages

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hanchon/garnet/internal/backend/cors"
	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketContainer struct {
	Authenticated bool
	User          string
	WalletID      int
	WalletAddress string
	Conn          *websocket.Conn
}

type User struct {
	Username string
	Password string
	WalletID int
}

type GlobalState struct {
	done        chan (struct{})
	WalletIndex map[string]string
	WsSockets   map[string]*WebSocketContainer
	// Simulate users database: map[user]password
	RegisteredUsers   map[string]User
	Database          *data.Database
	LastBroadcastTime time.Time
}

func NewGlobalState(database *data.Database) GlobalState {
	return GlobalState{
		done:        make(chan struct{}),
		WalletIndex: make(map[string]string),
		WsSockets:   make(map[string]*WebSocketContainer),
		RegisteredUsers: map[string]User{
			"user1": {Username: "user1", Password: "password1", WalletID: 0},
			"user2": {Username: "user2", Password: "password2", WalletID: 1},
		},
		Database:          database,
		LastBroadcastTime: time.Now(),
	}
}

func (g *GlobalState) WebSocketConnectionHandler(response http.ResponseWriter, request *http.Request) {
	if cors.SetHandlerCorsForOptions(request, &response) {
		return
	}

	// TODO: Filter prod page or localhost for development
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		// Maybe log the error
		return
	}

	webSocket := WebSocketContainer{
		Authenticated: false,
		Conn:          ws,
	}

	g.WsHandler(&webSocket)
}

func (g *GlobalState) BroadcastUpdates() {
	for {
		select {
		case <-g.done:
			return
		case <-time.After(50 * time.Millisecond):
			if len(g.WsSockets) != 0 {
				timestamp := g.Database.LastUpdate
				if g.LastBroadcastTime != timestamp {
					logger.LogDebug(fmt.Sprintf("[backend] database was updated, broadcasting messages..."))

					w, ok := g.Database.Worlds[WorldID]
					if !ok {
						panic("world not found")
					}

					g.LastBroadcastTime = timestamp

					for _, v := range g.WsSockets {
						if v.Conn != nil {
							matchData := g.Database.GetBoardStatus(WorldID, v.WalletAddress)
							if matchData != nil {
								// TODO: save the user and wallet somewhere
								matchData.PlayerOneUsermane = "user1"
								matchData.PlayerTwoUsermane = "user2"
								msgToSend := BoardStatus{MsgType: "boardstatus", Status: *matchData}
								logger.LogDebug(fmt.Sprintf("[backend] sending match info %s to %s", matchData.MatchID, v.User))
								err := v.Conn.WriteJSON(msgToSend)
								if err != nil {
									logger.LogError(fmt.Sprintf("[backend] error sending transaction to client, unsubscribing: %s", err))
									// TODO: unsub from this connection, it requires a lock in the g.WsSockets variable to avoid breaking the loops
									v.Authenticated = false
									v.Conn = nil
								}
							} else {
								t := w.GetTableByName("Match")
								if t != nil {
									ret := []string{}
									for k := range *t.Rows {
										ret = append(ret, k)
									}
									msg := MatchList{MsgType: "matchlist", Matches: ret}
									err := v.Conn.WriteJSON(msg)
									if err != nil {
										logger.LogError(fmt.Sprintf("[backend] error sending transaction to client, unsubscribing: %s", err))
										// TODO: unsub from this connection, it requires a lock in the g.WsSockets variable to avoid breaking the loops
										v.Authenticated = false
										v.Conn = nil
									}
									logger.LogDebug(fmt.Sprintf("[backend] sending %d active matches", len(ret)))
								}
							}
						}
					}
				}
			}
		}
	}
}
