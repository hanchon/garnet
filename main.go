package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func colorRed(value string) string {
	return fmt.Sprintf("\033[31;1m%s\033[0m", value)
}

func colorGreen(value string) string {
	return fmt.Sprintf("\033[1;32m%s\033[0m", value)
}

func drawMana() string {
	return fmt.Sprint("\033[1;34m◆\033[0m")
}

func drawHeart() string {
	return colorRed("♥")
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

const (
	topOffset      = 2
	leftOffset     = 2
	cardInfoWidth  = 24
	cardInfoHeight = 13
)

func cardInfo(g *gocui.Gui) error {
	if v, err := g.SetView("cardInfo", leftOffset, topOffset, cardInfoWidth, cardInfoHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "      Card Info      ")
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintln(v, " Type -> Warrior")
		fmt.Fprintf(v, "  ◎ Health  : 6 %s\n", drawHeart())
		fmt.Fprintf(v, "  ◎ Attack  : 4 (2%s)\n", drawMana())
		fmt.Fprintf(v, "  ◎ Movement: 2 (2%s)\n", drawMana())
		fmt.Fprintln(v, " ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ")
		fmt.Fprintln(v, " Ability:")
		fmt.Fprintf(v, "  ◎ Drain Sword: (4%s)\n", drawMana())
	}
	return nil
}

const (
	playerActionsTopOffset = cardInfoHeight + 1
	playerActionsWidth     = cardInfoWidth
	playerActionsHeight    = cardInfoHeight + 13
)

func playerActions(g *gocui.Gui) error {
	if v, err := g.SetView("playerActions", leftOffset, playerActionsTopOffset, playerActionsWidth, playerActionsHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintf(v, " Current Mana: 10%s\n", drawMana())
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintf(v, " Summon (3%s): ◯◉◉\n", drawMana())
		fmt.Fprintf(v, "  ◯  %s Vaan Strife\n", colorGreen("(♚)"))
		fmt.Fprintf(v, "  ◉  %s Felguard\n", colorGreen("(♛)"))
		fmt.Fprintf(v, "  ◉  %s Makimachi\n", colorGreen("(♜)"))
		fmt.Fprintf(v, "  ◉  %s Freya\n", colorGreen("(♝)"))
		fmt.Fprintf(v, "  ◉  %s Madmartigan\n", colorGreen("(♞)"))
		fmt.Fprintf(v, "  ◉  %s Jaina\n", colorGreen("(♟)"))
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintln(v, "    ✔ END TURN ✔     ")
		fmt.Fprintln(v, "─────────────────────")
	}

	return nil
}

const (
	gameActionsTopOffset = playerActionsHeight + 1
	gameActionsWidth     = cardInfoWidth
	gameActionsHeight    = playerActionsHeight + 8
)

func gameActions(g *gocui.Gui) error {
	if v, err := g.SetView("gameActions", leftOffset, gameActionsTopOffset, gameActionsWidth, gameActionsHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "     GAME OPTIONS    ")
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintln(v, "   ✔ CREATE GAME ✔   ")
		fmt.Fprintln(v, "    ✔ JOIN GAME ✔    ")
		fmt.Fprintln(v, "─────────────────────")
		fmt.Fprintln(v, "      ✔ QUIT ✔       ")
	}

	return nil
}

const (
	boardLeftOffset = leftOffset + cardInfoWidth
	boardTopOffset  = topOffset
	boardWidth      = boardLeftOffset + 82
	boardHeight     = topOffset + 32
)

const (
	mulX = 8
	mulY = 3
)

func board(g *gocui.Gui) error {

	if _, err := g.SetView("board", boardLeftOffset, boardTopOffset, boardWidth, boardHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		offsetX := boardLeftOffset + 1
		offsetY := boardTopOffset + 1
		endX := offsetX + mulX
		endY := offsetY + mulY
		for i := 0; i <= 9; i = i + 1 {
			for j := 0; j <= 9; j = j + 1 {
				if v, err := g.SetView(fmt.Sprintf("board%d%d", i, j), offsetX, offsetY, endX, endY); err != nil {
					if err != gocui.ErrUnknownView {
						return err
					}
					// fmt.Fprintf(v, "%d%d", i, j)
					if j == 0 && i == 0 {
						fmt.Fprintln(v, "10\u26A1")
						fmt.Fprintln(v, "     ♖")
					}
					if j == 0 && i == 1 {
						fmt.Fprintf(v, "%d %s\n", 10, colorRed("♥"))
						fmt.Fprintln(v, "     P")
					}
				}
				offsetX = endX
				endX = offsetX + mulX
			}
			offsetX = boardLeftOffset + 1
			endX = offsetX + mulX
			offsetY = endY
			endY = offsetY + mulY
		}
	}
	return nil
}

func layout(g *gocui.Gui) error {
	if err := cardInfo(g); err != nil {
		return err
	}
	if err := playerActions(g); err != nil {
		return err
	}
	if err := gameActions(g); err != nil {
		return err
	}

	if err := board(g); err != nil {
		return err
	}

	// if v, err := g.SetView("but1", 2, 2, 22, 7); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	v.Highlight = true
	// 	v.SelBgColor = gocui.ColorGreen
	// 	v.SelFgColor = gocui.ColorBlack
	// 	fmt.Fprintln(v, "Button 1 - line 1")
	// 	fmt.Fprintln(v, "Button 1 - line 2")
	// 	fmt.Fprintln(v, "Button 1 - line 3")
	// 	fmt.Fprintln(v, "Button 1 - line 4")
	// }
	// if v, err := g.SetView("but2", 24, 2, 44, 4); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	v.Highlight = true
	// 	v.SelBgColor = gocui.ColorGreen
	// 	v.SelFgColor = gocui.ColorBlack
	// 	fmt.Fprintln(v, "Button 2 - line 1")
	// }
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	for _, n := range []string{"but1", "but2"} {
		if err := g.SetKeybinding(n, gocui.MouseLeft, gocui.ModNone, showMsg); err != nil {
			return err
		}
	}
	if err := g.SetKeybinding("msg", gocui.MouseLeft, gocui.ModNone, delMsg); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func showMsg(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-10, maxY/2, maxX/2+10, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, l)
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	return nil
}
