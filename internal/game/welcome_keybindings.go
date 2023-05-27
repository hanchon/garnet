package game

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hanchon/garnet/internal/backend/messages"
	"github.com/hanchon/garnet/internal/gui"
	"github.com/hanchon/garnet/internal/logger"
	"github.com/jroimartin/gocui"
)

func (gs *GameState) WelcomeScreenKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(createGameViewName, gocui.MouseLeft, gocui.ModNone, gs.createMatch); err != nil {
		return err
	}

	if err := g.SetKeybinding(welcomeTablesViewName, gocui.MouseLeft, gocui.ModNone, gs.clickOnTable); err != nil {
		return err
	}

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

	return g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, gs.pgDnPressed)
}

func (gs *GameState) homePressed(_ *gocui.Gui, _ *gocui.View) error {
	gs.keyPressed = "HOME"
	return nil
}

func (gs *GameState) endPressed(_ *gocui.Gui, _ *gocui.View) error {
	gs.keyPressed = "END"
	return nil
}

func (gs *GameState) downPressed(_ *gocui.Gui, _ *gocui.View) error {
	gs.keyPressed = "DOWN"
	return nil
}

func (gs *GameState) upPressed(_ *gocui.Gui, _ *gocui.View) error {
	gs.keyPressed = "UP"
	return nil
}

func (gs *GameState) pgUpPressed(_ *gocui.Gui, _ *gocui.View) error {
	gs.keyPressed = "PGUP"
	return nil
}

func (gs *GameState) pgDnPressed(_ *gocui.Gui, _ *gocui.View) error {
	gs.keyPressed = "PGDN"
	return nil
}

func (gs *GameState) createMatch(g *gocui.Gui, _ *gocui.View) error {
	if gs.Connected {
		msg := `{"msgtype":"creatematch"}`
		if err := gs.Ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			return err
		}
		// TODO: disable keybindings, getting the new game will change the view to the board view
		if v, err := g.SetView("msg", leftOffset, welcomeLogoHeight+4, boardWidth, welcomeLogoHeight+9); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			v.Frame = false

			fmt.Fprintf(
				v,
				"%s%s%s\n",
				gui.ColorMagenta("\u2554"),
				gui.ColorMagenta(strings.Repeat("\u2550", 102)),
				gui.ColorMagenta("\u2557"),
			)
			fmt.Fprintf(
				v,
				"%s%s%s%s%s\n",
				gui.ColorMagenta("\u2551"),
				strings.Repeat(" ", 43),
				"CREATING GAME...",
				strings.Repeat(" ", 43),
				gui.ColorMagenta("\u2551"),
			)
			fmt.Fprintf(
				v,
				"%s%s%s\n",
				gui.ColorMagenta("\u255A"),
				gui.ColorMagenta(strings.Repeat("\u2550", 102)),
				gui.ColorMagenta("\u255D"),
			)

		}
		return nil
	}
	return nil
}

func (gs *GameState) clickOnTable(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	cx, cy := v.Cursor()
	gameID := ""
	msg := ""
	if cx > 71 && cx < 86 {
		msg = fmt.Sprintf(
			"%s%s%s%s%s",
			gui.ColorMagenta("\u2551"),
			strings.Repeat(" ", 47),
			gui.ColorYellow("JOINING "),
			strings.Repeat(" ", 47),
			gui.ColorMagenta("\u2551"),
		)
	} else if cx > 88 && cx < 103 {
		return nil

		// TODO: spectate
		// msg = fmt.Sprintf(
		// 	"%s%s%s%s%s",
		// 	gui.ColorMagenta("\u2551"),
		// 	strings.Repeat(" ", 46),
		// 	gui.ColorLightCyan("SPECTATING"),
		// 	strings.Repeat(" ", 46),
		// 	gui.ColorMagenta("\u2551"),
		// )
	}

	if cy > 2 && cy < 13 {
		if l, err := v.Line(cy); err == nil {
			gameID = l[4:70]
			// If the line is the end of the table, just ignore it
			if strings.Contains(gameID, "\u2550") {
				gameID = ""
			}
		}
	}

	if msg != "" && gameID != "" {
		if v, err := g.SetView("msg", leftOffset, welcomeLogoHeight+4, boardWidth, welcomeLogoHeight+9); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			v.Frame = false

			fmt.Fprintf(
				v,
				"%s%s%s\n",
				gui.ColorMagenta("\u2554"),
				gui.ColorMagenta(strings.Repeat("\u2550", 102)),
				gui.ColorMagenta("\u2557"),
			)
			fmt.Fprintf(v, "%s\n", msg)
			fmt.Fprintf(
				v,
				"%s%s%s%s%s\n",
				gui.ColorMagenta("\u2551"),
				strings.Repeat(" ", 18),
				gameID,
				strings.Repeat(" ", 18),
				gui.ColorMagenta("\u2551"),
			)
			fmt.Fprintf(
				v,
				"%s%s%s\n",
				gui.ColorMagenta("\u255A"),
				gui.ColorMagenta(strings.Repeat("\u2550", 102)),
				gui.ColorMagenta("\u255D"),
			)

			var msg messages.JoinMatch
			msg.MsgType = "joinmatch"
			msg.MatchID = gameID
			err = gs.Ws.WriteJSON(msg)
			logger.LogInfo(fmt.Sprintf("[client] sending join match request: %s", gameID))
			if err != nil {
				logger.LogError(fmt.Sprintf("[client] error sending join match request: %s. %s", gameID, err))
			}

		}
	}
	return nil
}

func (gs *GameState) UpdateMatches() {
	for {
		select {
		case <-gs.done:
			return
		case <-time.After(50 * time.Millisecond):
			rerender := false
			lastUpdate := gs.lastDBUpdate
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
					gs.yOffset++
				}

				if gs.keyPressed == "UP" {
					gs.yOffset--
				}

				if gs.keyPressed == "PGUP" {
					gs.yOffset -= maxLines
				}

				if gs.keyPressed == "PGDN" {
					gs.yOffset += maxLines
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

			if rerender {
				if len(gs.listOfAvailableGamesToRender) == 0 {
					continue
				}
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
