package game

import (
	"time"

	"github.com/jroimartin/gocui"
)

func (gs *GameState) WelcomeScreenKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyHome, gocui.ModNone, gs.homePressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEnd, gocui.ModNone, gs.endPressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEnd, gocui.ModNone, gs.endPressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, gs.downPressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, gs.upPressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, gs.pgUpPressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, gs.pgDnPressed); err != nil {
		return err
	}

	return nil
}

func (gs *GameState) homePressed(g *gocui.Gui, v *gocui.View) error {
	gs.keyPressed = "HOME"
	return nil
}
func (gs *GameState) endPressed(g *gocui.Gui, v *gocui.View) error {
	gs.keyPressed = "END"
	return nil
}

func (gs *GameState) downPressed(g *gocui.Gui, v *gocui.View) error {
	gs.keyPressed = "DOWN"
	return nil
}

func (gs *GameState) upPressed(g *gocui.Gui, v *gocui.View) error {
	gs.keyPressed = "UP"
	return nil
}

func (gs *GameState) pgUpPressed(g *gocui.Gui, v *gocui.View) error {
	gs.keyPressed = "PGUP"
	return nil
}

func (gs *GameState) pgDnPressed(g *gocui.Gui, v *gocui.View) error {
	gs.keyPressed = "PGDN"
	return nil
}

func (gs *GameState) UpdateMatches() {
	for {
		select {
		case <-gs.done:
			return
		case <-time.After(50 * time.Millisecond):
			rerender := false
			lastUpdate := gs.lastDbUpdate
			if gs.lastRenderUpdate != lastUpdate {
				gs.lastRenderUpdate = lastUpdate
				gs.listOfAvailableGamesToRender = gs.ListOfAvailableGames
				rerender = true
			}

			if gs.keyPressed != "" {
				if gs.keyPressed == "HOME" {
					gs.yOffset = 0
				}

				if gs.keyPressed == "DOWN" {
					gs.yOffset = gs.yOffset + 1
				}

				if gs.keyPressed == "UP" {
					gs.yOffset = gs.yOffset - 1
				}

				if gs.keyPressed == "PGUP" {
					gs.yOffset = gs.yOffset - maxLines
				}

				if gs.keyPressed == "PGDN" {
					gs.yOffset = gs.yOffset + maxLines
				}

				if gs.keyPressed == "END" {
					gs.yOffset = len(gs.listOfAvailableGamesToRender) - maxLines
				}

				if gs.yOffset+maxLines > len(gs.listOfAvailableGamesToRender) {
					gs.yOffset = len(gs.listOfAvailableGamesToRender) - maxLines
				}

				if gs.yOffset < 0 {
					gs.yOffset = 0
				}

				gs.keyPressed = ""
				rerender = true
			}

			if rerender == true {
				gs.ui.Update(func(g *gocui.Gui) error {
					v, err := g.View(welcomeTablesViewName)
					if err != nil {
						return err
					}
					v.Clear()
					if gs.yOffset > len(gs.listOfAvailableGamesToRender) {
						if len(gs.listOfAvailableGamesToRender) < maxLines {
							gs.yOffset = 0
						} else {
							gs.yOffset = len(gs.listOfAvailableGamesToRender) - maxLines
						}
					}

					end := 0

					if gs.yOffset+maxLines > len(gs.listOfAvailableGamesToRender) {
						end = len(gs.listOfAvailableGamesToRender)
					} else {
						end = gs.yOffset + maxLines
					}

					// fmt.Fprintln(v, gui.ColorMagenta("Tables Info:"))
					temp := []string{}
					for i := gs.yOffset; i < end; i++ {
						temp = append(temp, gs.listOfAvailableGamesToRender[i])
						// fmt.Fprintln(v, gs.listOfAvailableGamesToRender[i])
					}
					RenderWelcomeTable(
						temp,
						v,
						gs.yOffset != 0,
						temp[len(temp)-1] != gs.listOfAvailableGamesToRender[len(gs.listOfAvailableGamesToRender)-1],
					)
					return nil
				})

			}
		}
	}
}
