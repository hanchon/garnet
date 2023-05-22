package game

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hanchon/garnet/internal/backend/messages"
	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/logger"
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

const (
	EmptyAction  = ""
	SummonAction = "summon"
	MoveAction   = "move"
)

type GameState struct {
	MatchState           *MatchState
	ListOfAvailableGames []string
	Ws                   *websocket.Conn
	Connected            bool
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
	BoardStatus                  *data.MatchData
	lastBoardStatusUpdate        time.Time
	lastBoardRenderUpdate        time.Time
	UnitSelected                 string
	CurrentAction                string
}

func NewGameState(ui *gocui.Gui, username string, password string) *GameState {
	return &GameState{
		MatchState:           nil,
		ListOfAvailableGames: []string{},
		Ws:                   nil,
		Connected:            false,
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
		BoardStatus:                  nil,
		lastBoardStatusUpdate:        time.Unix(0, 1),
		lastBoardRenderUpdate:        time.Unix(0, 0),
		UnitSelected:                 "",
		CurrentAction:                EmptyAction,
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
				gameState.Connected = false
				panic(err)
				// panic("could not decode message")
			}
			if strings.Contains(string(v), "connected") {
				gameState.Connected = true
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

			if strings.Contains(string(v), "boardstatus") {
				var msg messages.BoardStatus
				err := json.Unmarshal(v, &msg)
				if err != nil {
					// Invalid json
					panic(err)
				}

				if gameState.BoardStatus == nil {
					gameState.ui.SetManagerFunc(GameLayout)
					if err := gameState.GameKeybindings(gameState.ui); err != nil {
						logger.LogError("[client] failed to load game keybindings")
					}
				}

				gameState.BoardStatus = &msg.Status
				gameState.lastBoardStatusUpdate = time.Now()
				logger.LogInfo(fmt.Sprintf("[client] processing board status message at:%s", gameState.lastBoardStatusUpdate))
			}
		}
	}()

	return c
}

func (gs *GameState) GetSelectedCard() (data.Card, error) {
	for _, card := range gs.BoardStatus.Cards {
		if card.ID == gs.UnitSelected {
			return card, nil
		}
	}
	return data.Card{}, fmt.Errorf("card not selected")
}
