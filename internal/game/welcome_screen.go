package game

import (
	"fmt"
	"strings"

	"github.com/hanchon/garnet/internal/gui"
	"github.com/jroimartin/gocui"
)

const (
	welcomeLogoViewName   = "welcomelogo"
	welcomeLogoLeftOffset = leftOffset + 20
	welcomeLogoHeight     = 15
	// Create Button
	createGameViewName   = "creategameview"
	createGameLeftOffset = leftOffset + 45
	// Tables
	welcomeTablesViewName = "tableswelcome"
)

const (
	maxLines = 10
)

func WelcomeScreenLayout(g *gocui.Gui) error {
	if v, err := g.SetView(welcomeLogoViewName, welcomeLogoLeftOffset, topOffset, boardWidth, welcomeLogoHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false

		fmt.Fprintf(v, gui.ColorRed(`
___________ __                             .__
\_   _____//  |_  ___________  ____ _____  |  |
 |    __)_\   __\/ __ \_  __ \/    \\__  \ |  |
 |        \|  | \  ___/|  | \/   |  \/ __ \|  |__
/_______  /|__|  \___  >__|  |___|  (____  /____/
        \/           \/           \/     \/
             .____                                    .___
             |    |    ____   ____   ____   ____    __| _/______
             |    |  _/ __ \ / ___\_/ __ \ /    \  / __ |/  ___/
             |    |__\  ___// /_/  >  ___/|   |  \/ /_/ |\___ \
             |_______ \___  >___  / \___  >___|  /\____ /____  >
                     \/   \/_____/      \/     \/      \/    \/
        `))

	}

	if v, err := g.SetView(createGameViewName, createGameLeftOffset, welcomeLogoHeight, createGameLeftOffset+16, welcomeLogoHeight+4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false

		fmt.Fprintf(v, "%s%s%s\n", gui.ColorGreen("\u2554"), gui.ColorGreen(strings.Repeat("\u2550", 13)), gui.ColorGreen("\u2557"))
		fmt.Fprintf(v, "%s %s %s\n", gui.ColorGreen("\u2551"), gui.ColorLightCyan("CREATE GAME"), gui.ColorGreen("\u2551"))
		fmt.Fprintf(v, "%s%s%s\n", gui.ColorGreen("\u255A"), gui.ColorGreen(strings.Repeat("\u2550", 13)), gui.ColorGreen("\u255D"))
	}

	if v, err := g.SetView(welcomeTablesViewName, leftOffset, welcomeLogoHeight+4, boardWidth, boardHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
		RenderWelcomeTable([]string{}, v, false, false)
	}
	return nil
}

func RenderWelcomeTable(data []string, v *gocui.View, arrowTop bool, arrowButton bool) {
	fmt.Fprintf(
		v,
		"%s%s%s%s%s%s%s\n",
		gui.ColorGreen("\u2554"),
		gui.ColorGreen(strings.Repeat("\u2550", 69)),
		gui.ColorGreen("\u2566"),
		gui.ColorGreen(strings.Repeat("\u2550", 16)),
		gui.ColorGreen("\u2566"),
		gui.ColorGreen(strings.Repeat("\u2550", 16)),
		gui.ColorGreen("\u2557"),
	)
	fmt.Fprintf(
		v,
		"%s%s%s%s%s%s%s\n",
		gui.ColorGreen("\u2551"),
		gui.ColorLightCyan("                                MATCH                                "),
		gui.ColorGreen("\u2551"),
		gui.ColorLightCyan("      JOIN      "),
		gui.ColorGreen("\u2551"),
		gui.ColorLightCyan("    SPECTATE    "),
		gui.ColorGreen("\u2551"),
	)

	fmt.Fprintf(
		v,
		"%s%s%s%s%s%s%s\n",
		gui.ColorGreen("\u2560"),
		gui.ColorGreen(strings.Repeat("\u2550", 69)),
		gui.ColorGreen("\u256C"),
		gui.ColorGreen(strings.Repeat("\u2550", 16)),
		gui.ColorGreen("\u256C"),
		gui.ColorGreen(strings.Repeat("\u2550", 16)),
		gui.ColorGreen("\u2563"),
	)

	for k, value := range data {
		extra := " "
		if k == 0 && arrowTop {
			extra = "\u2191"
		}
		if k == len(data)-1 && arrowButton {
			extra = "\u2193"
		}
		fmt.Fprintf(
			v,
			"%s %s %s%s%s%s%s%s\n",
			gui.ColorGreen("\u2551"),
			value,
			extra,
			gui.ColorGreen("\u2551"),
			gui.ColorYellow("      JOIN      "),
			gui.ColorGreen("\u2551"),
			gui.ColorBlue("    SPECTATE    "),
			gui.ColorGreen("\u2551"),
		)
	}

	fmt.Fprintf(
		v,
		"%s%s%s%s%s%s%s\n",
		gui.ColorGreen("\u255A"),
		gui.ColorGreen(strings.Repeat("\u2550", 69)),
		gui.ColorGreen("\u2569"),
		gui.ColorGreen(strings.Repeat("\u2550", 16)),
		gui.ColorGreen("\u2569"),
		gui.ColorGreen(strings.Repeat("\u2550", 16)),
		gui.ColorGreen("\u255D"),
	)

}
