package game

import (
	"fmt"
	"strconv"
	"strings"

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
				if _, err := g.SetView(fmt.Sprintf("%s%d%d", boardViewName, j, i), offsetX, offsetY, endX, endY); err != nil {
					if err != gocui.ErrUnknownView {
						return err
					}
					// if j == 0 && i == 0 {
					// 	fmt.Fprintln(v, "10\u26A1")
					// 	fmt.Fprintln(v, "     ♖")
					// }
					// if j == 0 && i == 1 {
					// 	fmt.Fprintf(v, "%d %s\n", 10, gui.ColorRed("♥"))
					// 	fmt.Fprintln(v, "     P")
					// }
					// drawBase(10, 5, j, i, v)

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

func lineWithCastles(enemy bool) string {
	castle := gui.ColorGreen("⛫")
	if enemy == true {
		castle = gui.ColorRed("⛫")
	}
	return fmt.Sprintf("%s\n", strings.Repeat(castle, 6))
}

func drawBase(p1Health int, p2Health int, currentX int, currentY int, v *gocui.View, isPlayerOne bool) {
	// Player 1 base
	if currentX == 4 && currentY == 0 {
		fmt.Fprintf(v, "%d/10%s\n", p1Health, drawHeart())
		fmt.Fprintf(v, lineWithCastles(!isPlayerOne))
	}
	if (currentX == 4 && currentY == 1) || (currentX == 5 && currentY == 0) || (currentX == 5 && currentY == 1) {
		fmt.Fprintf(v, lineWithCastles(!isPlayerOne))
		fmt.Fprintf(v, lineWithCastles(!isPlayerOne))
	}

	// Player 2 base
	if currentX == 4 && currentY == 8 {
		fmt.Fprintf(v, "%d/10%s\n", p2Health, drawHeart())
		fmt.Fprintf(v, lineWithCastles(isPlayerOne))
	}
	if (currentX == 4 && currentY == 9) || (currentX == 5 && currentY == 8) || (currentX == 5 && currentY == 9) {
		fmt.Fprintf(v, lineWithCastles(isPlayerOne))
		fmt.Fprintf(v, lineWithCastles(isPlayerOne))
	}
}

func (gs *GameState) updateBoard() error {
	userCards := gs.GetUserCards()

	isPlayerOne := true
	if gs.Username == gs.BoardStatus.PlayerTwoUsermane {
		isPlayerOne = false
	}

	for i := 0; i <= boardLimitX; i = i + 1 {
		for j := 0; j <= boardLimitY; j = j + 1 {
			v, err := gs.ui.View(fmt.Sprintf("%s%d%d", boardViewName, i, j))
			if err != nil {
				return err
			}
			v.Clear()
			drawBase(10, 5, i, j, v, isPlayerOne)
			// TODO: move the loop outside the other 2 loops
			for _, card := range gs.BoardStatus.Cards {
				if card.Position.X == int64(i) && card.Position.Y == int64(j) {
					symbol := gui.ColorRed(fmt.Sprintf("(%s)", TypeOfCards[card.Type].Symbol))
					for _, userCard := range userCards {
						if userCard.ID == card.ID {
							symbol = gui.ColorGreen(fmt.Sprintf("(%s)", TypeOfCards[card.Type].Symbol))
						}
					}
					fmt.Fprintf(v, " %s/%s %s\n", strconv.FormatInt(card.CurrentHp, 10), strconv.FormatInt(card.MaxHp, 10), drawHeart())
					fmt.Fprintf(v, "  %s\n", symbol)
					break
				}
			}
		}
	}

	return nil
}
