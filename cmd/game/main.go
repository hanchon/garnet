package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hanchon/garnet/internal/backend/messages"
	"github.com/hanchon/garnet/internal/game"
	"github.com/jroimartin/gocui"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("ERROR: username and password missing \nTo run the game execute: make run user1 password1")
		return
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(game.WelcomeScreenLayout)

	state := game.NewGameState(g, os.Args[1], os.Args[2])

	state.Ws = game.InitWsConnection(state)
	msg := messages.ConnectMessage{
		MsgType:  "connect",
		User:     state.Username,
		Password: state.Password,
	}
	state.Ws.WriteJSON(msg)

	if err := state.WelcomeScreenKeybindings(g); err != nil {
		log.Panicln(err)
	}

	go state.UpdateMatches()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	// TODO: close the state channel
}
