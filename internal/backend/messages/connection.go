package messages

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hanchon/garnet/internal/backend/cors"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketContainer struct {
	Authenticated bool
	User          string
	WalletID      int
	Conn          *websocket.Conn
}

type User struct {
	Username string
	Password string
	WalletID int
}

type GlobalState struct {
	WalletIndex map[string]string
	WsSockets   map[string]*websocket.Conn
	// Simulate users database: map[user]password
	RegisteredUsers map[string]User
}

func NewGlobalState() GlobalState {
	return GlobalState{
		WalletIndex: make(map[string]string),
		WsSockets:   make(map[string]*websocket.Conn),
		RegisteredUsers: map[string]User{
			"user1": {Username: "user1", Password: "password1", WalletID: 0},
			"user2": {Username: "user2", Password: "password2", WalletID: 1},
		},
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
