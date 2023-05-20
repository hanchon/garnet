package game

import "github.com/jroimartin/gocui"

func drawMovementPlaces(x, y int64, speed int64, g *gocui.Gui) {
	if speed >= 1 {
		setBackgroundBoardPosition(x, y+1, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x, y-1, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x+1, y, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x-1, y, gocui.ColorYellow, g)
	}
	if speed >= 2 {
		setBackgroundBoardPosition(x+1, y+1, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x-1, y-1, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x+1, y-1, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x-1, y+1, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x, y+2, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x, y-2, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x+2, y, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x-2, y, gocui.ColorYellow, g)
	}
	if speed >= 3 {
		setBackgroundBoardPosition(x+1, y+2, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x-1, y+2, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x+1, y-2, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x-1, y-2, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x+2, y-1, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x+2, y+1, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x-2, y+1, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x-2, y-1, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x, y+3, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x, y-3, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x+3, y, gocui.ColorYellow, g)
		setBackgroundBoardPosition(x-3, y, gocui.ColorYellow, g)
	}
}
