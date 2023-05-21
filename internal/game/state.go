package game

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hanchon/garnet/internal/backend/messages"
	"github.com/jroimartin/gocui"
)

type SummonedUnit struct {
	ID            string
	Name          string
	Symbol        string
	PosX          int
	PosY          int
	MaxHealth     int
	CurrentHealth int
	Attack        int
	Movement      int
}

type MatchState struct {
	MatchId                string
	CurrentMana            int
	UnitsSummonedByPlayer1 []int
	UnitsSummonedByPlayer2 []int
	CurrentPlayerTurn      int
	Board                  []SummonedUnit
}

type GameState struct {
	MatchState           *MatchState
	ListOfAvailableGames []string
	Ws                   *websocket.Conn
	Username             string
	Password             string
	// UI
	ui                           *gocui.Gui
	done                         chan (struct{})
	keyPressed                   string
	currentScreen                string
	lastDbUpdate                 time.Time // Last time we got new available games
	lastRenderUpdate             time.Time // Last time we updated the displayed list of games
	listOfAvailableGamesToRender []string
	yOffset                      int
}

var testData = []string{
	"0x1000000000000000000000000000000000000000000000000000000000000034",
	"0x0200000000000000000000000000000000000000000000000000000000000038",
	"0x0010000000000000000000000000000000000000000000000000000000000088",
	"0x0210000000000000000000000000000000000000000000000000000000000024",
	"0x0002000000000000000000000000000000000000000000000000000000000044",
	"0x0200000000000000000000000000000000000000000000000000000000000094",
	"0x0000020000000000000000000000000000000000000000000000000000000014",
	"0x0000000001000000000000000000000000000000000000000000000000000024",
	"0x0000002000000000000000000000000000000000000000000000000000001034",
	"0x0000000200000000000000000000000000000000000000000000000000000034",
	"0x0000000301000000000000000000000000000000000000000000005555000035",
	"0x0000000000000000000000000000000000000000000000000000000000000038",
	"0x0000800000000000000000000000000000000000000000000000000000000088",
	"0x0000000000000000000000000000000000000000000000000001111000000024",
	"0x0000000000000000000000000000000000000000000000000000000000000044",
	"0x0000000000000000000010000000000000000000000000000000000000000094",
	"0x0000000000000000000000000000000000000000000000000000000000000014",
	"0x0000000000000000000010000000000000000000000000000000000000000024",
	"0x0000000000000000000000000000000000000000000000000000000000001034",
	"0x00000000000000000000000000000000000000j0000000000000000000000034",
	"0x00000000000000000000000000000000000000002h0000010000000000000035",
	"0x0000000000000000000000000000000000000jj000h000080000000000000038",
	"0x0000000000000000000000000000000020000000000h00002000000000000088",
	"0x00000000000000000000000000000000j0000000000h00008000000000000024",
	"0x0000000000000000000000000000000000000000000010000000000000000044",
	"0x00000000000000000000000000000000k0000000000000100000000000000094",
	"0x000000000000000000000000000000jj00000000000000010000000000000014",
	"0x0000000000000000000010000000100200000000000000000000000000000024",
	"0x0000000000000000000000000000000100000000000000000000000000001034",
	"0x0000000000000000000000000000100000000000000000000000000000000034",
	"0x0000000000000000000000000010000000000000000000000000000000000035",
	"0x0000000000000000000010h00100000000000000000000000000000000000034",
	"0x0000000000000000000000101000000000000000000000000000000000000034",
}

func NewGameState(ui *gocui.Gui, username string, password string) *GameState {
	return &GameState{
		MatchState:           nil,
		ListOfAvailableGames: testData,
		Ws:                   nil,
		Username:             username,
		Password:             password,
		// UI
		ui:                           ui,
		keyPressed:                   "",
		currentScreen:                "",
		lastDbUpdate:                 time.Unix(0, 1),
		listOfAvailableGamesToRender: []string{},
		yOffset:                      0,
		lastRenderUpdate:             time.Unix(0, 0),
	}
}

func InitWsConnection(gameState *GameState) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: ":6666", Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic("could not connect")
	}

	// Receive messages
	go func() {
		run := true
		for run {
			// var msg pubtypes.WebsocketEvent
			// err := c.ReadJSON(&msg)
			_, v, err := c.ReadMessage()
			if err != nil {
				//  TODO: remove connection from gameState
				run = false
				panic(err)
				// panic("could not decode message")
			}
			if strings.Contains(string(v), "connected") {
				msg := `{"msgtype":"creatematch"}`
				c.WriteMessage(websocket.TextMessage, []byte(msg))
			}
			if strings.Contains(string(v), "matchlist") {
				var msg messages.MatchList
				err := json.Unmarshal(v, &msg)
				if err != nil {
					// Invalid json
					panic(err)
				}
				gameState.ListOfAvailableGames = msg.Matches
			}
			// values.AddValue(msg)
			// textArea.TextLines = values.ToStringList()
		}
	}()

	return c
}
