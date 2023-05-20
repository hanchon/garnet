package game

import (
	"fmt"

	"github.com/hanchon/garnet/internal/gui"
	"github.com/jroimartin/gocui"
)

const (
	mulX = 8
	mulY = 3
)

const (
	boardViewName = "board"
	boardLimitX   = 9
	boardLimitY   = 9
)

func board(pos ViewPosition, g *gocui.Gui) error {
	if _, err := g.SetView(boardViewName, pos.startX, pos.startY, pos.endX, pos.endY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		offsetX := pos.startX + 1
		offsetY := pos.startY + 1
		endX := offsetX + mulX
		endY := offsetY + mulY
		for i := 0; i <= boardLimitX; i = i + 1 {
			for j := 0; j <= boardLimitY; j = j + 1 {
				if v, err := g.SetView(fmt.Sprintf("%s%d%d", boardViewName, i, j), offsetX, offsetY, endX, endY); err != nil {
					if err != gocui.ErrUnknownView {
						return err
					}
					if j == 0 && i == 0 {
						fmt.Fprintln(v, "10\u26A1")
						fmt.Fprintln(v, "     ♖")
					}
					if j == 0 && i == 1 {
						fmt.Fprintf(v, "%d %s\n", 10, gui.ColorRed("♥"))
						fmt.Fprintln(v, "     P")
					}

					// Player 1 base
					if j == 4 && i == 0 {
						fmt.Fprintf(v, "%s%s%s%s%s\n", "10", drawHeart(), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"))
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"))
					}
					if j == 4 && i == 1 {
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"))
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"))
					}
					if j == 5 && i == 0 {
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"))
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"))
					}
					if j == 5 && i == 1 {
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"))
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"), gui.ColorGreen("⛫"))
					}

					// Player 2 base
					if j == 4 && i == 8 {
						fmt.Fprintf(v, "%s%s%s%s%s\n", "10", drawHeart(), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"))
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"))
					}
					if j == 4 && i == 9 {
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"))
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"))
					}
					if j == 5 && i == 8 {
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"))
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"))
					}
					if j == 5 && i == 9 {
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"))
						fmt.Fprintf(v, "%s%s%s%s%s%s\n", gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"), gui.ColorRed("⛫"))
					}
				}
				offsetX = endX
				endX = offsetX + mulX
			}
			offsetX = pos.startX + 1
			endX = offsetX + mulX
			offsetY = endY
			endY = offsetY + mulY
		}
	}
	return nil
}
