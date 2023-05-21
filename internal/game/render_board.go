package game

import (
	"time"

	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/jroimartin/gocui"
)

func (gs *GameState) GetUserCards() []data.Card {
	// Check if the unit was summoned (spectator = player 1)
	owner := gs.BoardStatus.PlayerOne
	if gs.Username == gs.BoardStatus.PlayerTwoUsermane {
		owner = gs.BoardStatus.PlayerTwo
	}

	cards := []data.Card{}
	for _, v := range gs.BoardStatus.Cards {
		if v.Owner == owner {
			cards = append(cards, v)
		}
	}
	return cards
}

func (gs *GameState) UpdateBoard() {
	for {
		select {
		case <-gs.done:
			return
		case <-time.After(50 * time.Millisecond):
			if gs.BoardStatus != nil {
				rerender := false
				lastUpdate := gs.lastBoardStatusUpdate
				if gs.lastBoardRenderUpdate != lastUpdate {
					gs.lastBoardRenderUpdate = lastUpdate
					rerender = true
				}

				if rerender == true {
					gs.ui.Update(func(g *gocui.Gui) error {
						err := gs.updatePlayerActions()
						if err != nil {
							return err
						}
						err = gs.updateNotifications()
						if err != nil {
							return err
						}
						return nil
					})
				}
			}
		}
	}
}
