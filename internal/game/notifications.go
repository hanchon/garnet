package game

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	notificationsViewName = "notifications"
)

func (gs *GameState) updateNotifications() error {
	v, err := gs.ui.View(notificationsViewName)
	if err != nil {
		return err
	}

	currentPlayer := gs.BoardStatus.PlayerTwoUsermane
	if gs.BoardStatus.CurrentPlayer == gs.BoardStatus.PlayerOne {
		currentPlayer = gs.BoardStatus.PlayerOneUsermane
	}
	username := gs.BoardStatus.PlayerTwoUsermane
	if gs.Username == gs.BoardStatus.PlayerOneUsermane {
		username = gs.BoardStatus.PlayerOneUsermane
	}
	v.Clear()

	fmt.Fprintf(
		v,
		"%s%s                   %s%d            %s%s",
		"Player: ",
		username,
		"Current turn: ",
		gs.BoardStatus.CurrentTurn,
		"Current player: ",
		currentPlayer,
	)
	return nil
}

func notifications(pos ViewPosition, g *gocui.Gui) error {
	if v, err := g.SetView(notificationsViewName, pos.startX, pos.startY, pos.endX, pos.endY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintf(v, "%s                   %s            %s", "Player: user1", "Current turn: 0", "Current player: user1")
	}
	return nil
}
