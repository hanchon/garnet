package game

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const gameActionsViewName = "gameActions"

func gameActions(pos ViewPosition, g *gocui.Gui) error {
	if v, err := g.SetView(gameActionsViewName, pos.startX, pos.startY, pos.endX, pos.endY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "     GAME OPTIONS    ")
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintln(v, "   ✔ CREATE GAME ✔   ")
		fmt.Fprintln(v, "    ✔ JOIN GAME ✔    ")
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintln(v, "      ✔ QUIT ✔       ")
	}

	return nil
}
