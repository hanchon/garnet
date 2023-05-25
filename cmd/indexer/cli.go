package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hanchon/garnet/internal/gui"
	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/logger"
	"github.com/jroimartin/gocui"
)

const maxLinesToDisplay = 37

type DebugUI struct {
	done             chan (struct{})
	ui               *gocui.Gui
	xOffset          int
	yOffset          int
	searchTerm       string
	data             []string
	dataLastUpdate   time.Time
	searchTotalIndex int
	searchIndex      int
	keyPressed       string
}

func NewDebugUI() *DebugUI {
	return &DebugUI{
		done:             make(chan struct{}),
		ui:               ui(),
		xOffset:          0,
		yOffset:          0,
		searchTerm:       "",
		data:             []string{""},
		dataLastUpdate:   time.Unix(0, 0),
		keyPressed:       "",
		searchTotalIndex: 0,
		searchIndex:      0,
	}
}

func findWord(input string, values []string) []int {
	a := []int{}
	for k, v := range values {
		if strings.Contains(strings.ToLower(v), strings.ToLower(input)) {
			a = append(a, k)
		}
	}
	return a
}

func (ui *DebugUI) Run() {
	if err := ui.keybindings(ui.ui); err != nil {
		log.Panicln(err)
	}

	if err := ui.ui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	close(ui.done)
}

func (ui *DebugUI) ProcessLatestEvents(database *data.Database) {
	for {
		select {
		case <-ui.done:
			return
		case <-time.After(500 * time.Millisecond):
			ui.ui.Update(func(g *gocui.Gui) error {
				v, err := g.View("latestevents")
				if err != nil {
					return err
				}
				v.Clear()
				fmt.Fprintln(v, gui.ColorMagenta("Latest Events:"))
				fmt.Fprintln(v, strings.Repeat("─", logoWidth-logoOffsetX))

				end := 0
				start := len(database.Events) - 1
				if start < 0 {
					start = 0
				}
				if len(database.Events) >= 7 {
					end = len(database.Events) - 7
				}

				for i := start; i > end; i-- {
					fmt.Fprintln(v, "----EVENT----")
					fmt.Fprint(v, "Table:")
					fmt.Fprintln(v, database.Events[i].Table)
					// fmt.Fprint(v, "Row:")
					// fmt.Fprintln(v, database.Events[i].Row)
					fmt.Fprint(v, "Value:")
					fmt.Fprintln(v, database.Events[i].Value)
				}

				return nil
			})
		}

	}
}
func (ui *DebugUI) ProcessBlockchainInfo(database *data.Database) {
	for {
		select {
		case <-ui.done:
			return
		case <-time.After(500 * time.Millisecond):
			ui.ui.Update(func(g *gocui.Gui) error {
				v, err := g.View("blockchaininfo")
				if err != nil {
					return err
				}
				v.Clear()
				fmt.Fprintln(v, gui.ColorMagenta("Blockchain Info:"))
				fmt.Fprintln(v, strings.Repeat("─", logoWidth))
				fmt.Fprintf(v, " \u26d3 ChainID: %s\n", database.ChainID)
				fmt.Fprintf(v, " \u279a Height : %d\n", database.LastHeight)
				return nil
			})
		}

	}
}

func (ui *DebugUI) ProcessIncomingData(database *data.Database) {
	for {
		select {
		case <-ui.done:
			return
		case <-time.After(50 * time.Millisecond):
			// TODO: move the search status updates to another function
			// Update search status
			ui.ui.Update(func(g *gocui.Gui) error {
				if v, err := g.View("searchboxinfo"); err == nil {
					v.Clear()
					fmt.Fprintf(v, "Type to search. Control+n and Control+p to move arround. Res:%d/%d", ui.searchIndex+1, ui.searchTotalIndex)
				}
				return nil
			})

			rerender := false
			lastUpdate := database.LastUpdate
			if ui.dataLastUpdate != lastUpdate {
				ui.data = database.ToStringList(debugWindowWidth - debugWindowOffset)
				ui.dataLastUpdate = lastUpdate
				rerender = true
			}

			if ui.keyPressed != "" {
				if ui.keyPressed == "HOME" {
					ui.yOffset = 0
				}

				if ui.keyPressed == "DOWN" {
					ui.yOffset = ui.yOffset + 1
				}

				if ui.keyPressed == "UP" {
					ui.yOffset = ui.yOffset - 1
				}

				if ui.keyPressed == "PGUP" {
					ui.yOffset = ui.yOffset - maxLinesToDisplay
				}

				if ui.keyPressed == "PGDN" {
					ui.yOffset = ui.yOffset + maxLinesToDisplay
				}

				if ui.keyPressed == "END" {
					ui.yOffset = len(ui.data) - maxLinesToDisplay
				}

				if ui.keyPressed == "P" {
					if v, err := ui.ui.View("searchboxcontent"); err == nil {
						content := v.ViewBufferLines()
						if len(content) > 0 {
							// In case the user pressed enter
							v.Clear()
							fmt.Fprintf(v, content[len(content)-1])

							// Find the word in the data
							logger.LogDebug(fmt.Sprintf("[garnet] searching for: %s", content[0]))
							pos := findWord(content[len(content)-1], ui.data)
							ui.searchTotalIndex = len(pos)
							if len(pos) != 0 {
								logger.LogDebug(fmt.Sprintf("[garnet] %s in positions: %v", content[0], pos))
								if ui.searchIndex-1 < 0 || ui.searchIndex-1 >= len(pos) {
									ui.searchIndex = len(pos) - 1
								} else {
									ui.searchIndex = ui.searchIndex - 1
								}
								ui.yOffset = pos[ui.searchIndex]
							}
						}
					}
				}

				if ui.keyPressed == "N" {
					if v, err := ui.ui.View("searchboxcontent"); err == nil {
						content := v.ViewBufferLines()
						if len(content) > 0 {
							// In case the user pressed enter
							v.Clear()
							fmt.Fprintf(v, content[len(content)-1])

							// Find the word in the data
							logger.LogDebug(fmt.Sprintf("[garnet] searching for: %s", content[0]))
							pos := findWord(content[len(content)-1], ui.data)
							ui.searchTotalIndex = len(pos)
							if len(pos) != 0 {
								logger.LogDebug(fmt.Sprintf("[garnet] %s in positions: %v", content[0], pos))
								if ui.searchIndex+1 > len(pos)-1 {
									ui.searchIndex = 0
								} else {
									ui.searchIndex = ui.searchIndex + 1
								}
								ui.yOffset = pos[ui.searchIndex]
							}
						}
					}
				}

				if ui.yOffset+maxLinesToDisplay > len(ui.data) {
					ui.yOffset = len(ui.data) - maxLinesToDisplay
				}

				if ui.yOffset < 0 {
					ui.yOffset = 0
				}

				ui.keyPressed = ""
				rerender = true
			}

			if rerender == true {
				ui.ui.Update(func(g *gocui.Gui) error {
					v, err := g.View("debugui")
					if err != nil {
						return err
					}
					v.Clear()
					if ui.yOffset > len(ui.data) {
						if len(ui.data) < maxLinesToDisplay {
							ui.yOffset = 0
						} else {
							ui.yOffset = len(ui.data) - maxLinesToDisplay
						}
					}

					end := 0

					if ui.yOffset+maxLinesToDisplay > len(ui.data) {
						end = len(ui.data)
					} else {
						end = ui.yOffset + maxLinesToDisplay
					}

					fmt.Fprintln(v, gui.ColorMagenta("Tables Info:"))
					fmt.Fprintln(v, strings.Repeat("─", debugWindowWidth-debugWindowOffset-1))
					for i := ui.yOffset; i < end; i++ {
						fmt.Fprintln(v, ui.data[i])
					}
					return nil
				})

			}
		}
	}
}

func ui() *gocui.Gui {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	g.SetManagerFunc(layout)

	return g
}

const (
	logoOffsetX          = 4
	logoHeight           = 8
	logoWidth            = 40
	blockchainInfoHeight = logoHeight + 5
	blockchainInfoOffset = 2
	debugWindowHeight    = 40
	debugWindowWidth     = 120
	debugWindowOffset    = 4
)

func layout(g *gocui.Gui) error {
	logo := gui.ColorRed(`
 _____                       _
|  __ \                     | |
| |  \/ __ _ _ __ _ __   ___| |_
| | __ / _' | '__| '_ \ / _ \ __|
| |_\ \ (_| | |  | | | |  __/ |_
 \____/\__,_|_|  |_| |_|\___|\__|
 `)
	if v, err := g.SetView("logo", logoOffsetX, 0, logoWidth, logoHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintln(v, logo)
	}

	if v, err := g.SetView("blockchaininfo", blockchainInfoOffset, logoHeight, logoWidth, blockchainInfoHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, gui.ColorMagenta("Blockchain Info:"))
		fmt.Fprintln(v, strings.Repeat("─", logoWidth))
		fmt.Fprintln(v, "ChainID: ")
		fmt.Fprintln(v, "Height: ")
	}

	if v, err := g.SetView("latestevents", blockchainInfoOffset, blockchainInfoHeight+2, logoWidth, debugWindowHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Frame = true

		fmt.Fprintln(v, gui.ColorMagenta("Latest Events:"))
		fmt.Fprintln(v, strings.Repeat("─", logoWidth-logoOffsetX))
	}

	if _, err := g.SetView("debugui", logoWidth+debugWindowOffset, 0, logoWidth+debugWindowWidth, debugWindowHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	searchBoxOffset := 8

	if v, err := g.SetView("searchbox", blockchainInfoOffset, debugWindowHeight+1, blockchainInfoOffset+searchBoxOffset, debugWindowHeight+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintf(v, "Search:")
	}

	if v, err := g.SetView("searchboxcontent", blockchainInfoOffset+searchBoxOffset+1, debugWindowHeight+1, logoWidth+debugWindowWidth, debugWindowHeight+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Editable = true
		g.SetCurrentView("searchboxcontent")
		fmt.Fprintf(v, "")
	}

	if v, err := g.SetView("searchboxinfo", blockchainInfoOffset, debugWindowHeight+3, logoWidth+debugWindowWidth, debugWindowHeight+5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintf(v, "Type to search. Control+n and Control+p to move arround.")
	}

	return nil
}

func (ui *DebugUI) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyHome, gocui.ModNone, ui.homePressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEnd, gocui.ModNone, ui.endPressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEnd, gocui.ModNone, ui.endPressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, ui.downPressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, ui.upPressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, ui.pgUpPressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, ui.pgDnPressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone, ui.controlNPressed); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, ui.controlPPressed); err != nil {
		return err
	}

	return nil
}

func (ui *DebugUI) controlNPressed(g *gocui.Gui, v *gocui.View) error {
	ui.keyPressed = "N"
	return nil
}

func (ui *DebugUI) controlPPressed(g *gocui.Gui, v *gocui.View) error {
	ui.keyPressed = "P"
	return nil
}

func (ui *DebugUI) homePressed(g *gocui.Gui, v *gocui.View) error {
	ui.keyPressed = "HOME"
	return nil
}
func (ui *DebugUI) endPressed(g *gocui.Gui, v *gocui.View) error {
	ui.keyPressed = "END"
	return nil
}

func (ui *DebugUI) downPressed(g *gocui.Gui, v *gocui.View) error {
	ui.keyPressed = "DOWN"
	return nil
}

func (ui *DebugUI) upPressed(g *gocui.Gui, v *gocui.View) error {
	ui.keyPressed = "UP"
	return nil
}

func (ui *DebugUI) pgUpPressed(g *gocui.Gui, v *gocui.View) error {
	ui.keyPressed = "PGUP"
	return nil
}

func (ui *DebugUI) pgDnPressed(g *gocui.Gui, v *gocui.View) error {
	ui.keyPressed = "PGDN"
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
