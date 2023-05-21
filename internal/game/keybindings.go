package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hanchon/garnet/internal/backend/messages"
	"github.com/jroimartin/gocui"
)

func (gs *GameState) GameKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	// Board cells keybindings
	for i := 0; i <= 9; i = i + 1 {
		for j := 0; j <= 9; j = j + 1 {
			key := fmt.Sprintf("%s%d%d", boardViewName, i, j)
			if err := g.SetKeybinding(key, gocui.MouseLeft, gocui.ModNone, showMovementPlaces); err != nil {
				return err
			}
		}
	}

	// // Create game
	// if err := g.SetKeybinding(gameActionsViewName, gocui.MouseLeft, gocui.ModNone, gs.clickOnGameActions); err != nil {
	// 	return err
	// }

	// if err := g.SetKeybinding("msg", gocui.MouseLeft, gocui.ModNone, delMsg); err != nil {
	// 	return err
	// }
	return nil
}

func (gs *GameState) clickOnGameActions(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	if l, err := v.Line(cy); err == nil {
		if strings.Contains(l, "CREATE") {
			// CREATE GAME
			gs.Ws = InitWsConnection(gs)
			msg := messages.ConnectMessage{
				MsgType:  "connect",
				User:     "user1",
				Password: "password1",
			}
			gs.Ws.WriteJSON(msg)
		}

		if strings.Contains(l, "QUIT") {

			g.SetManagerFunc(GameLayout)
			if err := gs.GameKeybindings(g); err != nil {
				panic(err)
			}

			// if err := g.DeleteView(gameActionsViewName); err != nil {
			// 	return err
			// }
			// return nil

		}
	}
	return nil
}

func showMovementPlaces(g *gocui.Gui, v *gocui.View) error {
	xy := strings.Replace(v.Name(), "board", "", 1)
	x, err := strconv.ParseInt(string(xy[0]), 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse x")
	}
	y, err := strconv.ParseInt(string(xy[1]), 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse y")
	}

	drawMovementPlaces(x, y, 3, g)

	// maxX, maxY := g.Size()
	// if v2, err := g.SetView("msg", maxX/2-10, maxY/2, maxX/2+10, maxY/2+2); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	fmt.Fprintf(v2, "%s", v.Name())
	// }
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
