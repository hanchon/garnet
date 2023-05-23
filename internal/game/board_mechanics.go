package game

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (gs *GameState) isValidPlaceToMove(x int64, y int64) bool {
	if x == 4 && y == 0 || x == 5 && y == 0 || x == 4 && y == 1 || x == 5 && y == 1 || x == 4 && y == 8 || x == 5 && y == 8 || x == 4 && y == 9 || x == 5 && y == 9 {
		// Bases
		return false
	}

	for _, v := range gs.BoardStatus.Cards {
		// Placed cards
		if v.Position.X == x && v.Position.Y == y {
			return false
		}
	}

	return true
}

func (gs *GameState) setMovementPosition(x, y int64, color gocui.Attribute) {
	if gs.isValidPlaceToMove(x, y) {
		setBackgroundColor(fmt.Sprintf("%s%d%d", boardViewName, x, y), color, gs.ui)
	}
}

var attackBackgroundColor = gocui.ColorRed

func (gs *GameState) setAttackPosition(x, y int64, myUser string) {
	// Bases
	if myUser == gs.BoardStatus.PlayerOne {
		if (x == 4 || x == 5) && (y == 8 || y == 9) {
			setBackgroundColor(fmt.Sprintf("%s%d%d", boardViewName, x, y), attackBackgroundColor, gs.ui)
		}
	} else {
		if (x == 4 || x == 5) && (y == 0 || y == 1) {
			setBackgroundColor(fmt.Sprintf("%s%d%d", boardViewName, x, y), attackBackgroundColor, gs.ui)
		}

	}
	for _, v := range gs.BoardStatus.Cards {
		if v.Position.X == x && v.Position.Y == y && v.Owner != myUser {
			setBackgroundColor(fmt.Sprintf("%s%d%d", boardViewName, x, y), attackBackgroundColor, gs.ui)
			return
		}
	}

}

func (gs *GameState) GetUserWallet() string {
	myUser := gs.BoardStatus.PlayerTwo
	if gs.Username == gs.BoardStatus.PlayerOneUsermane {
		myUser = gs.BoardStatus.PlayerOne
	}
	return myUser
}

func (gs *GameState) drawAttackPlaces(x, y int64) {
	myUser := gs.GetUserWallet()

	gs.setAttackPosition(x, y+1, myUser)
	gs.setAttackPosition(x, y-1, myUser)
	gs.setAttackPosition(x+1, y, myUser)
	gs.setAttackPosition(x-1, y, myUser)
}

func (gs *GameState) drawMovementPlaces(x, y int64, speed int64) {

	if speed >= 1 {
		gs.setMovementPosition(x, y+1, gocui.ColorYellow)
		gs.setMovementPosition(x, y-1, gocui.ColorYellow)
		gs.setMovementPosition(x+1, y, gocui.ColorYellow)
		gs.setMovementPosition(x-1, y, gocui.ColorYellow)
	}
	if speed >= 2 {
		gs.setMovementPosition(x+1, y+1, gocui.ColorYellow)
		gs.setMovementPosition(x-1, y-1, gocui.ColorYellow)
		gs.setMovementPosition(x+1, y-1, gocui.ColorYellow)
		gs.setMovementPosition(x-1, y+1, gocui.ColorYellow)
		gs.setMovementPosition(x, y+2, gocui.ColorYellow)
		gs.setMovementPosition(x, y-2, gocui.ColorYellow)
		gs.setMovementPosition(x+2, y, gocui.ColorYellow)
		gs.setMovementPosition(x-2, y, gocui.ColorYellow)
	}
	if speed >= 3 {
		gs.setMovementPosition(x+1, y+2, gocui.ColorYellow)
		gs.setMovementPosition(x-1, y+2, gocui.ColorYellow)
		gs.setMovementPosition(x+1, y-2, gocui.ColorYellow)
		gs.setMovementPosition(x-1, y-2, gocui.ColorYellow)
		gs.setMovementPosition(x+2, y-1, gocui.ColorYellow)
		gs.setMovementPosition(x+2, y+1, gocui.ColorYellow)
		gs.setMovementPosition(x-2, y+1, gocui.ColorYellow)
		gs.setMovementPosition(x-2, y-1, gocui.ColorYellow)
		gs.setMovementPosition(x, y+3, gocui.ColorYellow)
		gs.setMovementPosition(x, y-3, gocui.ColorYellow)
		gs.setMovementPosition(x+3, y, gocui.ColorYellow)
		gs.setMovementPosition(x-3, y, gocui.ColorYellow)
	}
}
