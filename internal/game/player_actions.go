package game

import (
	"fmt"

	"github.com/hanchon/garnet/internal/gui"
	"github.com/jroimartin/gocui"
)

func playerActions(pos ViewPosition, g *gocui.Gui) error {
	if v, err := g.SetView("playerActions", pos.startX, pos.startY, pos.endX, pos.endY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintf(v, " Current Mana: 10%s\n", drawMana())
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintf(v, " Summon (3%s): ◯◉◉\n", drawMana())
		fmt.Fprintf(v, "  ◯  %s Vaan Strife\n", gui.ColorGreen("(♚)"))
		fmt.Fprintf(v, "  ◉  %s Felguard\n", gui.ColorGreen("(♛)"))
		fmt.Fprintf(v, "  ◉  %s Makimachi\n", gui.ColorGreen("(♜)"))
		fmt.Fprintf(v, "  ◉  %s Freya\n", gui.ColorGreen("(♝)"))
		fmt.Fprintf(v, "  ◉  %s Madmartigan\n", gui.ColorGreen("(♞)"))
		fmt.Fprintf(v, "  ◉  %s Jaina\n", gui.ColorGreen("(♟)"))
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintln(v, "    ✔ END TURN ✔     ")
		fmt.Fprintln(v, "─────────────────────")
	}
	return nil
}
