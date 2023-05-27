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

	// Set the log output to a file (stdin, stdout, stderror used by GUI)
	fileName := "client.log"
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(game.WelcomeScreenLayout)

	state := game.NewGameState(g, os.Args[1], os.Args[2])

	if err := state.WelcomeScreenKeybindings(g); err != nil {
		log.Panicln(err)
	}

	state.Ws = game.InitWsConnection(state)
	msg := messages.ConnectMessage{
		MsgType:  "connect",
		User:     state.Username,
		Password: state.Password,
	}
	if err := state.Ws.WriteJSON(msg); err != nil {
		// Could not send the connect message to the server
		log.Panicln(err)
	}

	go state.UpdateMatches()
	go state.UpdateBoard()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	// TODO: close the state channel
}
