package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hanchon/garnet/internal/backend/messages"
	"github.com/hanchon/garnet/internal/logger"
	"github.com/jroimartin/gocui"
)

func (gs *GameState) GameKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	// Board cells keybindings
	for i := 0; i <= 9; i = i + 1 {
		for j := 0; j <= 9; j = j + 1 {
			key := fmt.Sprintf("%s%d%d", boardViewName, i, j)
			if err := g.SetKeybinding(key, gocui.MouseLeft, gocui.ModNone, gs.showMovementPlaces); err != nil {
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

	if err := g.SetKeybinding("msg", gocui.MouseLeft, gocui.ModNone, delMsg); err != nil {
		return err
	}
	return nil
}

func (gs *GameState) selectCardFromPlayerActions(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	if cy < 3 || cy > 8 {
		return nil
	}
	cardType := cy - 3
	userCards := gs.GetUserCards()
	for _, card := range userCards {
		if card.Type == int64(cardType) {
			gs.UnitSelected = card.ID
		}
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
		// TODO: make sure that the card was not already summoned
		// Make sure we have at least 3 mana
		gs.CurrentAction = SummonAction
		for x := int64(0); x <= 9; x++ {
			for y := int64(0); y <= 1; y++ {
				if x != 4 && x != 5 {
					setBackgroundBoardPosition(x, y, gocui.ColorCyan, gs.ui)
				}
			}
		}

		return nil
	})
	// maxX, maxY := g.Size()
	// if v, err := g.SetView("msg", maxX/2-10, maxY/2, maxX/2+10, maxY/2+2); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	fmt.Fprintln(v, )
	// }
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	return nil
}

func (gs *GameState) clickOnGameActions(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	if l, err := v.Line(cy); err == nil {
		if strings.Contains(l, "CREATE") {
			// CREATE GAME
		}

		if strings.Contains(l, "QUIT") {

			g.SetManagerFunc(GameLayout)
			if err := gs.GameKeybindings(g); err != nil {
				panic(err)
			}

			// if err := g.DeleteView(gameActionsViewName); err != nil {
			// 	return err
			// }
			// return nil

		}
	}
	return nil
}

func (gs *GameState) showMovementPlaces(g *gocui.Gui, v *gocui.View) error {
	xy := strings.Replace(v.Name(), "board", "", 1)
	x, err := strconv.ParseInt(string(xy[0]), 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse x")
	}
	y, err := strconv.ParseInt(string(xy[1]), 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse y")
	}

	if gs.CurrentAction == SummonAction {
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
			// TODO: unselect all squares and unselec card
		}
	} else {
		drawMovementPlaces(x, y, 3, g)
	}

	// maxX, maxY := g.Size()
	// if v2, err := g.SetView("msg", maxX/2-10, maxY/2, maxX/2+10, maxY/2+2); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	fmt.Fprintf(v2, "%s", v.Name())
	// }
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
