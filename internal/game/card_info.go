package game

import (
	"fmt"

	"github.com/hanchon/garnet/internal/gui"
	"github.com/jroimartin/gocui"
)

const cardInfoViewName = "cardInfo"

func cardInfo(pos ViewPosition, g *gocui.Gui) error {
	if v, err := g.SetView(cardInfoViewName, pos.startX, pos.startY, pos.endX, pos.endY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "      Card Info      ")
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintln(v, " Name -> Warrior")
		fmt.Fprintf(v, "  ◎ Health  : 6 %s\n", drawHeart())
		fmt.Fprintf(v, "  ◎ Attack  : 4 (2%s)\n", drawMana())
		fmt.Fprintf(v, "  ◎ Movement: 2 (2%s)\n", drawMana())
		fmt.Fprintln(v, " ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ")
		fmt.Fprintln(v, " Ability:")
		fmt.Fprintf(v, "  ◎ Drain Sword: (4%s)\n", drawMana())
	}
	return nil
}

func (gs *GameState) updateCardInfo() error {
	v, err := gs.ui.View(cardInfoViewName)
	if err != nil {
		return err
	}
	name := ""
	symbol := ""
	maxHealth := ""
	hp := ""
	attack := ""
	movement := ""

	card, err := gs.GetSelectedCard()
	if err == nil {
		name = TypeOfCards[card.Type].Name
		symbol = TypeOfCards[card.Type].Symbol
		maxHealth = fmt.Sprintf("%d", card.MaxHp)
		hp = fmt.Sprintf("%d", card.CurrentHp)
		attack = fmt.Sprintf("%d", card.AttackDamage)
		movement = fmt.Sprintf("%d", card.MovementSpeed)
	}

	v.Clear()
	fmt.Fprintln(v, "      Card Info      ")
	fmt.Fprintln(v, "─────────────────────")
	fmt.Fprintf(v, " %s %s\n", gui.ColorMagenta(fmt.Sprintf("(%s)", symbol)), gui.ColorMagenta(name))
	fmt.Fprintf(v, "  ◎ Health  : %s/%s %s\n", maxHealth, hp, drawHeart())
	fmt.Fprintf(v, "  ◎ Attack  : %s (2%s)\n", attack, drawMana())
	fmt.Fprintf(v, "  ◎ Movement: %s (2%s)\n", movement, drawMana())
	fmt.Fprintln(v, " ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ")
	fmt.Fprintln(v, " Ability:")
	fmt.Fprintf(v, "  ◎ -----------: (4%s)\n", drawMana())
	return nil
}
