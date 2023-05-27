package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hanchon/garnet/internal/backend/messages"
	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/logger"
	"github.com/jroimartin/gocui"
)

func (gs *GameState) GameKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	// Board cells keybindings
	for i := 0; i <= 9; i++ {
		for j := 0; j <= 9; j++ {
			key := fmt.Sprintf("%s%d%d", boardViewName, i, j)
			if err := g.SetKeybinding(key, gocui.MouseLeft, gocui.ModNone, gs.boardMouseActionsHandler); err != nil {
				return err
			}
		}
	}

	if err := g.SetKeybinding(playerActionsViewName, gocui.MouseLeft, gocui.ModNone, gs.selectCardFromPlayerActions); err != nil {
		return err
	}

	// // Create game
	// if err := g.SetKeybinding(gameActionsViewName, gocui.MouseLeft, gocui.ModNone, gs.clickOnGameActions); err != nil {
	// 	return err
	// }

	return g.SetKeybinding("msg", gocui.MouseLeft, gocui.ModNone, delMsg)
}

func (gs *GameState) selectCardFromPlayerActions(_ *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()

	// End turn
	if cy == 10 {
		logger.LogInfo(fmt.Sprintf("[client] sending end turn transaction for match %s", gs.BoardStatus.MatchID))
		msg := messages.EndTurn{MsgType: "endturn", MatchID: gs.BoardStatus.MatchID}
		if err := gs.Ws.WriteJSON(msg); err != nil {
			// TODO: close the connection (?)
			return nil
		}
		gs.notificationMessages = append(gs.notificationMessages, "sending end turn transaction")
		if err := gs.updateNotifications(); err != nil {
			// TODO: move all the updates to a game cycle
			return err
		}
		gs.CurrentAction = EndTurn
		if err := gs.updatePlayerActions(); err != nil {
			// TODO: move all the updates to a game cycle
			return err
		}
		return nil
	}
	// Select card checks
	if cy < 3 || cy > 8 {
		return nil
	}
	cardType := cy - 3
	userCards := gs.GetUserCards()
	currentCard := data.Card{ID: ""}
	totalSummons := 0
	for _, card := range userCards {
		if card.Type == int64(cardType) {
			gs.UnitSelected = card.ID
			currentCard = card
		}
		if card.Placed {
			totalSummons++
		}
	}

	if currentCard.ID == "" {
		// The card should always exist
		return nil
	}

	gs.ui.Update(func(g *gocui.Gui) error {
		err := gs.updateCardInfo()
		if err != nil {
			return err
		}

		err = gs.updatePlayerActions()
		if err != nil {
			return err
		}

		err = gs.updateBoard()
		if err != nil {
			return err
		}

		// Make sure that the card was not already summoned
		if currentCard.Placed {
			gs.selectCard(currentCard.Position.X, currentCard.Position.Y)
			logger.LogDebug("[client] the card was already summoned")
			gs.notificationMessages = append(gs.notificationMessages, "the selected card was already summoned")
			if err := gs.updateNotifications(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
			return nil
		}

		// Make sure we have at least 3 mana
		if gs.BoardStatus.CurrentMana < 3 {
			logger.LogInfo("[client] not enough mana to summon")
			gs.notificationMessages = append(gs.notificationMessages, "there is not enough mana tu summon a new card")
			if err := gs.updateNotifications(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
			return nil
		}

		// Make sure that is the player turn
		// TODO: maybe it will not display the user summon rows after changing turns
		if gs.GetUserWallet() != gs.BoardStatus.CurrentPlayer {
			logger.LogInfo("[client] not your turn to summon")
			return nil
		}

		// The user can only summon 3 cards
		if totalSummons >= 3 {
			logger.LogInfo("[client] all the summons were used")
			gs.notificationMessages = append(gs.notificationMessages, "all the summon actions were already used")
			if err := gs.updateNotifications(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
			return nil
		}

		gs.CurrentAction = SummonAction
		yStart := int64(0)
		yEnd := int64(1)
		if gs.GetUserWallet() == gs.BoardStatus.PlayerTwo {
			yStart = int64(8)
			yEnd = int64(9)

		}

		for x := int64(0); x <= 9; x++ {
			for y := yStart; y <= yEnd; y++ {
				gs.setMovementPosition(x, y, gocui.ColorCyan)
			}
		}

		return nil
	})

	return nil
}

func delMsg(g *gocui.Gui, _ *gocui.View) error {
	return g.DeleteView("msg")
}

func (gs *GameState) boardMouseActionsHandler(_ *gocui.Gui, v *gocui.View) error {
	xy := strings.Replace(v.Name(), "board", "", 1)
	x, err := strconv.ParseInt(string(xy[0]), 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse x")
	}
	y, err := strconv.ParseInt(string(xy[1]), 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse y")
	}

	switch gs.CurrentAction {
	case SummonAction:
		if v.BgColor == gocui.ColorCyan {
			// Summon available
			msg := messages.PlaceCard{
				MsgType: "placecard",
				CardID:  gs.UnitSelected,
				X:       x,
				Y:       y,
			}
			err := gs.Ws.WriteJSON(msg)
			if err != nil {
				logger.LogError(fmt.Sprintf("[client] could not send place card message: %s", err))
			}
			gs.notificationMessages = append(gs.notificationMessages, "sending summon transaction")
			if err := gs.updateNotifications(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}

			gs.CurrentAction = EmptyAction
			gs.UnitSelected = ""

			if err := gs.updateBoard(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
			if err := gs.updatePlayerActions(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
			if err := gs.updateCardInfo(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
		} else {
			gs.notificationMessages = append(gs.notificationMessages, "summon canceled")
			if err := gs.updateNotifications(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}

			gs.CurrentAction = EmptyAction
			gs.UnitSelected = ""
			if err := gs.updateBoard(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
			if err := gs.updatePlayerActions(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
			if err := gs.updateCardInfo(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
		}
	case MoveAction:
		switch v.BgColor {
		case gocui.ColorYellow:
			// Move
			logger.LogInfo("[client] move card to pos")
			// TODO: this may fail if the user selects another card from the left table, make sure to clean the board background on unitselected changes
			msg := messages.MoveCard{
				MsgType: "movecard",
				CardID:  gs.UnitSelected,
				X:       x,
				Y:       y,
			}
			err := gs.Ws.WriteJSON(msg)
			if err != nil {
				logger.LogError(fmt.Sprintf("[client] could not send move card message: %s", err))
			}

			gs.notificationMessages = append(gs.notificationMessages, "sending move transaction")

			if err := gs.updateNotifications(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}

			gs.CurrentAction = EmptyAction
			gs.UnitSelected = ""
			if err := gs.updateBoard(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
			if err := gs.updatePlayerActions(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
			if err := gs.updateCardInfo(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
		case attackBackgroundColor:
			// Move
			logger.LogInfo("[client] attack to pos")
			// TODO: this may fail if the user selects another card from the left table, make sure to clean the board background on unitselected changes
			msg := messages.Attack{
				MsgType: "attack",
				CardID:  gs.UnitSelected,
				X:       x,
				Y:       y,
			}
			err := gs.Ws.WriteJSON(msg)
			if err != nil {
				logger.LogError(fmt.Sprintf("[client] could not send attack message: %s", err))
			}
			gs.notificationMessages = append(gs.notificationMessages, "sending attack transaction")
			if err := gs.updateNotifications(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}

			gs.CurrentAction = EmptyAction
			gs.UnitSelected = ""
			if err := gs.updateBoard(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
			if err := gs.updatePlayerActions(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
			if err := gs.updateCardInfo(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
		default:
			gs.notificationMessages = append(gs.notificationMessages, "stopping current action")
			if err := gs.updateNotifications(); err != nil {
				// TODO: move all the updates to a game cycle
				return err
			}
			gs.CurrentAction = EmptyAction
			gs.selectCard(x, y)
		}

	default:
		gs.selectCard(x, y)
	}

	return nil
}

func (gs *GameState) selectCard(x int64, y int64) {
	if err := gs.updateBoard(); err != nil {
		// TODO: move all the updates to a game cycle
		return
	}

	gs.UnitSelected = ""
	for _, value := range gs.BoardStatus.Cards {
		if value.Position.X == x && value.Position.Y == y {
			gs.UnitSelected = value.ID
			gs.notificationMessages = append(gs.notificationMessages, fmt.Sprintf("card selected, pos x:%d, y:%d", x, y))

			if err := gs.updateNotifications(); err != nil {
				// TODO: move all the updates to a game cycle
				return
			}

			// Update tables
			if err := gs.updateCardInfo(); err != nil {
				// TODO: move all the updates to a game cycle
				return
			}
			if err := gs.updatePlayerActions(); err != nil {
				// TODO: move all the updates to a game cycle
				return
			}

			// If the unit owner is the user, display the movement places
			for _, card := range gs.GetUserCards() {
				if card.ID == value.ID {
					if gs.BoardStatus.CurrentMana <= 1 {
						gs.CurrentAction = EmptyAction
						gs.notificationMessages = append(gs.notificationMessages, "not enough mana to execute any action, please end turn")
						if err := gs.updateNotifications(); err != nil {
							// TODO: move all the updates to a game cycle
							return
						}
						return
					}
					if !card.ActionReady {
						gs.CurrentAction = EmptyAction
						logger.LogDebug("[client] the selected card already used its action")
						gs.notificationMessages = append(gs.notificationMessages, "selected card already executed this turn attack action")
						if err := gs.updateNotifications(); err != nil {
							// TODO: move all the updates to a game cycle
							return
						}
						return
					}
					gs.CurrentAction = MoveAction
					gs.drawMovementPlaces(x, y, card.MovementSpeed)
					gs.drawAttackPlaces(x, y)
					return
				}
			}
		}
	}
	// Update tables
	if err := gs.updateCardInfo(); err != nil {
		// TODO: move all the updates to a game cycle
		return
	}
	if err := gs.updatePlayerActions(); err != nil {
		// TODO: move all the updates to a game cycle
		return
	}
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}
