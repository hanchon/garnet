package game

import (
	"fmt"
	"strings"

	"github.com/hanchon/garnet/internal/gui"
	"github.com/jroimartin/gocui"
)

const playerActionsViewName = "playerActions"

func playerActions(pos ViewPosition, g *gocui.Gui) error {
	if v, err := g.SetView(playerActionsViewName, pos.startX, pos.startY, pos.endX, pos.endY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, "\n\n─────────────────────\n\n\n\n\n\n\n\n\n")
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "─────────────────────")
	}
	return nil
}

func currentMana(mana int64) string {
	manaString := fmt.Sprintf("%d", mana)
	if mana < 10 {
		manaString = fmt.Sprintf(" %d", mana)
	}

	return fmt.Sprintf(" Current Mana: %s%s\n", manaString, drawMana())
}

type UiCard struct {
	Name   string
	ID     int64
	Symbol string
}

var TypeOfCards = []UiCard{
	{
		Name:   "Vaan Strife",
		ID:     0,
		Symbol: "\u2694",
	},

	{
		Name:   "Felguard",
		ID:     1,
		Symbol: "\u2692",
	},

	{
		Name:   "Sakura",
		ID:     2,
		Symbol: "\u2698",
	},

	{
		Name:   "Freya",
		ID:     3,
		Symbol: "\u27b9",
	},

	{
		Name:   "Lyra",
		ID:     4,
		Symbol: "\u273a",
	},

	{
		Name:   "Madmartigan",
		ID:     5,
		Symbol: "\u26e8",
	},
}

func renderUnit(summoned bool, name string, symbol string, selected bool) string {
	symbolSummon := "◉"
	if summoned {
		symbolSummon = "◯"
	}

	symbolString := fmt.Sprintf("(%s)", symbol)
	if selected {
		symbolSummon = gui.ColorMagenta(symbolSummon)
		symbolString = gui.ColorMagenta(symbolString)
		name = gui.ColorMagenta(name)
	} else {
		symbolString = gui.ColorGreen(symbolString)
	}

	return fmt.Sprintf("  %s  %s %s\n", symbolSummon, symbolString, name)
}

func renderUnits(v *gocui.View, summoned map[int64]bool, optionWithBackground int) {
	for k, card := range TypeOfCards {
		val, ok := summoned[card.ID]
		if !ok {
			val = false
		}
		text := renderUnit(val, card.Name, card.Symbol, k == optionWithBackground)
		fmt.Fprintf(v, text)
	}
}

func (gs *GameState) updatePlayerActions() error {
	endTurn := "    ✔ END TURN ✔     "
	if gs.BoardStatus.CurrentPlayer == gs.GetUserWallet() {
		// TODO: check if already clicked
		if gs.CurrentAction == EndTurn {
			endTurn = gui.ColorLightCyan(endTurn)
		} else {
			endTurn = gui.ColorGreen(endTurn)
		}
	} else {
		endTurn = gui.ColorRed(endTurn)
	}

	userCards := gs.GetUserCards()
	s := map[int64]bool{}
	optionWithBackground := -1

	for _, v := range userCards {
		s[v.Type] = v.Position.X != -2
		if v.ID == gs.UnitSelected {
			optionWithBackground = int(v.Type)
		}
	}

	v, err := gs.ui.View(playerActionsViewName)
	if err != nil {
		return err
	}

	summoned := 0
	for _, temp := range s {
		if temp {
			summoned++
		}
	}

	v.Clear()
	fmt.Fprintf(v, currentMana(gs.BoardStatus.CurrentMana))
	fmt.Fprintln(v, "─────────────────────", optionWithBackground)
	fmt.Fprintf(v, " Summon (3%s): %s%s\n", drawMana(), strings.Repeat("◯", summoned), strings.Repeat("◉", 3-summoned))
	renderUnits(v, s, optionWithBackground)
	fmt.Fprintln(v, "─────────────────────")
	fmt.Fprintf(v, "%s\n", endTurn)
	fmt.Fprintln(v, "─────────────────────")

	return nil
}
