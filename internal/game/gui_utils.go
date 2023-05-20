package game

import (
	"fmt"

	"github.com/hanchon/garnet/internal/gui"
	"github.com/jroimartin/gocui"
)

type ViewPosition struct {
	startX int
	startY int
	endX   int
	endY   int
}

func drawMana() string {
	return gui.ColorBlue("◆")
}

func drawHeart() string {
	return gui.ColorRed("♥")
}

func setBackgroundColor(viewName string, color gocui.Attribute, g *gocui.Gui) {
	view, err := g.View(viewName)
	if err == nil {
		view.BgColor = color
	}
}

func setBackgroundBoardPosition(x, y int64, color gocui.Attribute, g *gocui.Gui) {
	setBackgroundColor(fmt.Sprintf("%s%d%d", boardViewName, x, y), color, g)
}
