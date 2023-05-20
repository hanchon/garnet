package game

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func cardInfo(pos ViewPosition, g *gocui.Gui) error {
	if v, err := g.SetView("cardInfo", pos.startX, pos.startY, pos.endX, pos.endY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "      Card Info      ")
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintln(v, " Type -> Warrior")
		fmt.Fprintf(v, "  ◎ Health  : 6 %s\n", drawHeart())
		fmt.Fprintf(v, "  ◎ Attack  : 4 (2%s)\n", drawMana())
		fmt.Fprintf(v, "  ◎ Movement: 2 (2%s)\n", drawMana())
		fmt.Fprintln(v, " ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ")
		fmt.Fprintln(v, " Ability:")
		fmt.Fprintf(v, "  ◎ Drain Sword: (4%s)\n", drawMana())
	}
	return nil
}
