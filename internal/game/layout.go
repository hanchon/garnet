package game

import "github.com/jroimartin/gocui"

const (
	// CardInfo
	topOffset      = 2
	leftOffset     = 2
	cardInfoWidth  = 24
	cardInfoHeight = 13
	// PlayerActions
	playerActionsTopOffset = cardInfoHeight + 1
	playerActionsWidth     = cardInfoWidth
	playerActionsHeight    = cardInfoHeight + 13
	// GameActions
	gameActionsTopOffset = playerActionsHeight + 1
	gameActionsWidth     = cardInfoWidth
	gameActionsHeight    = playerActionsHeight + 8
	// Board
	boardLeftOffset = leftOffset + cardInfoWidth
	boardTopOffset  = topOffset
	boardWidth      = boardLeftOffset + 82
	boardHeight     = topOffset + 32
)

func Layout(g *gocui.Gui) error {
	if err := cardInfo(
		ViewPosition{
			startX: leftOffset,
			startY: topOffset,
			endX:   cardInfoWidth,
			endY:   cardInfoHeight,
		},
		g,
	); err != nil {
		return err
	}

	if err := playerActions(
		ViewPosition{
			startX: leftOffset,
			startY: playerActionsTopOffset,
			endX:   playerActionsWidth,
			endY:   playerActionsHeight,
		},
		g,
	); err != nil {
		return err
	}

	if err := gameActions(
		ViewPosition{
			startX: leftOffset,
			startY: gameActionsTopOffset,
			endX:   gameActionsWidth,
			endY:   gameActionsHeight,
		},
		g,
	); err != nil {
		return err
	}

	if err := board(
		ViewPosition{
			startX: boardLeftOffset,
			startY: boardTopOffset,
			endX:   boardWidth,
			endY:   boardHeight,
		},
		g,
	); err != nil {
		return err
	}

	return nil
}
