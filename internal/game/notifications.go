package game

import (
	"fmt"
	"strings"

	"github.com/hanchon/garnet/internal/gui"
	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/jroimartin/gocui"
)

const (
	notificationsViewName = "notifications"
)

func (gs *GameState) updateNotifications() error {
	v, err := gs.ui.View(notificationsViewName)
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintln(v, generateHeaderNotifications(gs.BoardStatus.MatchID))
	// TODO: just add the condition to the for loop
	j := 0
	for i := len(gs.notificationMessages) - 1; i >= 0; i-- {
		fmt.Fprintf(v, " %s %s\n", gui.ColorYellow("\u27a4"), gs.notificationMessages[i])
		j++
		if j == 6 {
			break
		}
	}

	return nil
}

func generateHeaderNotifications(matchid string) string {
	maxLengthNotifications := boardWidth - leftOffset

	titleGameTables := fmt.Sprintf("MESSAGES: %s", matchid)
	gameTablesSeparator := strings.Repeat("\u2632", (data.SeparatorOffset(maxLengthNotifications, len(titleGameTables)) - 1))
	header := fmt.Sprintf("%s %s %s", gui.ColorCyan(gameTablesSeparator), gui.ColorBlue(titleGameTables), gui.ColorCyan(gameTablesSeparator))
	return header
}

func notifications(pos ViewPosition, g *gocui.Gui) error {
	if v, err := g.SetView(notificationsViewName, pos.startX, pos.startY, pos.endX, pos.endY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, generateHeaderNotifications(""))
	}
	return nil
}
