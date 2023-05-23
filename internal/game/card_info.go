package game

import (
	"fmt"

	"github.com/hanchon/garnet/internal/gui"
	"github.com/jroimartin/gocui"
)

const cardInfoViewName = "cardInfo"

type CardInfo struct {
	name      string
	symbol    string
	maxHealth string
	hp        string
	attack    string
	movement  string
}

var cardInfoEmpty = CardInfo{
	name:      "No card",
	symbol:    "?",
	maxHealth: "",
	hp:        "",
	attack:    "",
	movement:  "",
}

func cardInfo(pos ViewPosition, g *gocui.Gui) error {
	if v, err := g.SetView(cardInfoViewName, pos.startX, pos.startY, pos.endX, pos.endY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		drawCardInfo(v, cardInfoEmpty)
	}
	return nil
}

func (gs *GameState) updateCardInfo() error {
	v, err := gs.ui.View(cardInfoViewName)
	if err != nil {
		return err
	}

	v.Clear()
	card, err := gs.GetSelectedCard()
	if err == nil {
		temp := CardInfo{
			name:      TypeOfCards[card.Type].Name,
			symbol:    TypeOfCards[card.Type].Symbol,
			maxHealth: fmt.Sprintf("%d", card.MaxHp),
			hp:        fmt.Sprintf("%d", card.CurrentHp),
			attack:    fmt.Sprintf("%d", card.AttackDamage),
			movement:  fmt.Sprintf("%d", card.MovementSpeed),
		}
		drawCardInfo(v, temp)
	} else {
		drawCardInfo(v, cardInfoEmpty)
	}

	return nil
}

func drawCardInfo(v *gocui.View, cardInfo CardInfo) {
	fmt.Fprintln(v, "      Card Info      ")
	fmt.Fprintln(v, "─────────────────────")
	fmt.Fprintf(v, " %s %s\n", gui.ColorMagenta(fmt.Sprintf("(%s)", cardInfo.symbol)), gui.ColorMagenta(cardInfo.name))
	fmt.Fprintf(v, "  ◎ Health  : %s/%s %s\n", cardInfo.maxHealth, cardInfo.hp, drawHeart())
	fmt.Fprintf(v, "  ◎ Attack  : %s (2%s)\n", cardInfo.attack, drawMana())
	fmt.Fprintf(v, "  ◎ Movement: %s (2%s)\n", cardInfo.movement, drawMana())
	fmt.Fprintln(v, " ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ")
	fmt.Fprintln(v, " Ability:")
	fmt.Fprintf(v, "  ◎ -----------: (4%s)\n", drawMana())
}
