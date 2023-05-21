package data

import (
	"fmt"
	"strings"

	"github.com/hanchon/garnet/internal/gui"
)

var SystemTables = []string{
	"schema",
	"StoreMetadata",
	"Hooks",
	"NamespaceOwner",
	"InstalledModules",
	"ResourceAccess",
	"Systems",
	"FunctionSelectors",
	"SystemRegistry",
	"ResourceType",
	"KeysWithValue",
}

func isSystemTable(key string) bool {
	for _, v := range SystemTables {
		if key == v {
			return true
		}
	}
	return false
}

func processTable(ret *[]string, vT *Table) {
	*ret = append(*ret, fmt.Sprintf("\u2727 Table %s", vT.Metadata.TableName))
	*ret = append(*ret, "  \u274a Rows:")
	for kR, vR := range *vT.Rows {
		// key := hexutil.Encode([]byte(kR))
		*ret = append(*ret, fmt.Sprintf("    \u2609 ID    : %s", kR))
		*ret = append(*ret, fmt.Sprintf("      Values:"))
		for _, b := range vR {
			*ret = append(*ret, fmt.Sprintf("          \u26ad  %s", b.String()))
		}
		*ret = append(*ret, fmt.Sprintf(""))

	}
}

func separatorOffset(maxLenght int, wordLength int) int {
	offset := 0
	if (maxLenght-wordLength)%2 == 0 {
		offset = (maxLenght - wordLength) / 2
	}
	offset = (maxLenght - wordLength - 1) / 2
	return offset
}

func (db Database) ToStringList(maxLenght int) []string {
	// For each world create a new array
	ret := make([]string, 0)
	tempSysTables := make([]string, 0)
	for _, vW := range db.Worlds {

		// World title

		worldSeparator := strings.Repeat("+", (separatorOffset(maxLenght, 48)))

		ret = append(ret, gui.ColorYellow(strings.Repeat("=", maxLenght)))
		ret = append(ret, fmt.Sprintf("%s %s %s", gui.ColorGreen(worldSeparator), gui.ColorBlue(fmt.Sprintf("World %s", vW.Address)), gui.ColorGreen(worldSeparator)))
		ret = append(ret, gui.ColorYellow(strings.Repeat("=", maxLenght)))
		ret = append(ret, "")

		// Game tables
		titleGameTables := "Game tables"
		gameTablesSeparator := strings.Repeat("\u2632", (separatorOffset(maxLenght, len(titleGameTables)) - 1))
		ret = append(ret, fmt.Sprintf("%s %s %s", gui.ColorCyan(gameTablesSeparator), gui.ColorBlue(titleGameTables), gui.ColorCyan(gameTablesSeparator)))
		ret = append(ret, "")
		for _, vT := range vW.Tables {
			if isSystemTable(vT.Metadata.TableName) == false {
				processTable(&ret, vT)
			} else {
				processTable(&tempSysTables, vT)
			}
		}
		ret = append(ret, "")

		// System tables
		titleSystemTables := "System tables"
		systemTablesSeparator := strings.Repeat("\u2632", (separatorOffset(maxLenght, len(titleSystemTables)) - 1))
		ret = append(ret, fmt.Sprintf("%s %s %s", gui.ColorCyan(systemTablesSeparator), gui.ColorBlue(titleSystemTables), gui.ColorCyan(systemTablesSeparator)))
		ret = append(ret, "")
		for _, v := range tempSysTables {
			ret = append(ret, v)
		}
		ret = append(ret, "")
	}
	return ret
}
