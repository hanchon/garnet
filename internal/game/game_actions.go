package game

import (
	"fmt"
	"strings"

	"github.com/hanchon/garnet/internal/gui"
	"github.com/jroimartin/gocui"
)

const gameActionsViewName = "gameActions"

func (gs *GameState) generateTurn() string {
	if gs.BoardStatus != nil {
		if gs.GetUserWallet() == gs.BoardStatus.CurrentPlayer {
			return gui.ColorGreen(" \u2603 YOUR TURN")
		}
	}
	return gui.ColorRed(" \u2620 ENEMY TURN")
}

func (gs *GameState) generateUsername() string {
	username := ""
	if gs.BoardStatus != nil {
		username = strings.ToUpper(gs.Username)
	}
	return fmt.Sprintf(" \u26c7 NAME: %s", username)
}

func gameActions(pos ViewPosition, g *gocui.Gui) error {
	if v, err := g.SetView(gameActionsViewName, pos.startX, pos.startY, pos.endX, pos.endY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "   \u2663  GAME INFO  \u2663   ")
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintf(v, " \u23f1 CURRENT TURN: %d      \n", 0)
		fmt.Fprintln(v, gui.ColorGreen(" \u2603 YOUR TURN"))
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintln(v, " \u26c7 NAME: USER1")
	}

	return nil
}

func (gs *GameState) updateGameActions() error {
	v, err := gs.ui.View(gameActionsViewName)
	if err != nil {
		return err
	}
	if gs.BoardStatus != nil {
		v.Clear()
		fmt.Fprintln(v, "   \u2663  GAME INFO  \u2663   ")
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintf(v, " \u23f1 CURRENT TURN: %d\n", gs.BoardStatus.CurrentTurn)
		fmt.Fprintln(v, gs.generateTurn())
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintln(v, gs.generateUsername())
	}
	return nil
}
