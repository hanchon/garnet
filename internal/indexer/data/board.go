package data

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hanchon/garnet/internal/logger"
)

func (db *Database) GetUserMatchID(w *World, userWallet string) string {
	// Check if the user is the player one
	playerOne := w.GetTableByName("PlayerOne")
	if playerOne == nil {
		logger.LogError("[backend] table player one does not exist")
		return ""
	}

	for k, v := range *playerOne.Rows {
		// Player one always have 1 field
		logger.LogDebug(fmt.Sprintf("[backend] comparing %s and %s", v[0].Data.String(), userWallet[2:]))
		if strings.Contains(strings.ToLower(v[0].Data.String()), strings.ToLower(userWallet[2:])) {
			logger.LogDebug(fmt.Sprintf("[backend] player found in match %s", k))
			return k
		}
	}

	// Check if the user is the player two
	playerTwo := w.GetTableByName("PlayerTwo")
	if playerTwo == nil {
		logger.LogError("[backend] table player two does not exist")
		return ""
	}

	for k, v := range *playerTwo.Rows {
		// Player one always have 1 field
		logger.LogDebug(fmt.Sprintf("[backend] comparing %s and %s", v[0].Data.String(), userWallet[2:]))
		if strings.Contains(strings.ToLower(v[0].Data.String()), strings.ToLower(userWallet[2:])) {
			logger.LogDebug(fmt.Sprintf("[backend] player found in match %s", k))
			return k
		}
	}
	return ""
}

const baseType = 6

func (db *Database) GetBoard(w *World, matchID string) *MatchData {
	logger.LogInfo("[backend] getting board status...")
	ret := MatchData{MatchID: matchID}

	playerOne := w.GetTableByName("PlayerOne")
	if playerOne == nil {
		logger.LogError("[backend] table player one does not exist")
		return nil
	}

	playerOneValue, ok := (*playerOne.Rows)[matchID]
	if !ok {
		logger.LogError("[backend] match does not have a player one")
		return nil
	}
	// Player one has only one field
	ret.PlayerOne = playerOneValue[0].Data.String()

	playerTwo := w.GetTableByName("PlayerTwo")
	if playerTwo == nil {
		logger.LogError("[backend] table player two does not exist")
		return nil
	}
	playerTwoValue, ok := (*playerTwo.Rows)[matchID]
	if !ok {
		// Player two did not join, just return the current values
		return &ret
	}
	ret.PlayerTwo = playerTwoValue[0].Data.String()

	// CurrentTurn
	currentTurn := w.GetTableByName("CurrentTurn")
	if currentTurn == nil {
		logger.LogError("[backend] table current turn does not exist")
		return nil
	}

	currentTurnValue, ok := (*currentTurn.Rows)[matchID]
	if !ok {
		logger.LogError(fmt.Sprintf("[backend] current turn does not exist for match %s", matchID))
		return nil
	}

	currentTurnValueParsed, err := strconv.ParseInt(currentTurnValue[0].Data.String(), 10, 32)
	if err != nil {
		logger.LogError(fmt.Sprintf("[backend] could not parse curren turn value %s", currentTurnValue[0].Data.String()))
	}
	ret.CurrentTurn = currentTurnValueParsed

	// CurrentPlayer
	currentPlayer := w.GetTableByName("CurrentPlayer")
	if currentPlayer == nil {
		logger.LogError("[backend] table current player does not exist")
		return nil
	}

	currentPlayerValue, ok := (*currentPlayer.Rows)[matchID]
	if !ok {
		logger.LogError(fmt.Sprintf("[backend] current player does not exist for match %s", matchID))
		return nil
	}

	ret.CurrentPlayer = currentPlayerValue[0].Data.String()

	// CurrentMana
	currentMana := w.GetTableByName("CurrentMana")
	if currentMana == nil {
		logger.LogError("[backend] table current mana does not exist")
		return nil
	}

	currentManaValue, ok := (*currentMana.Rows)[matchID]
	if !ok {
		logger.LogError(fmt.Sprintf("[backend] current mana does not exist for match %s", matchID))
		return nil
	}

	currentManaValueParsed, err := strconv.ParseInt(currentManaValue[0].Data.String(), 10, 32)
	if err != nil {
		logger.LogError(fmt.Sprintf("[backend] could not parse curren mana value %s", currentTurnValue[0].Data.String()))
	}
	ret.CurrentMana = currentManaValueParsed

	matchCards, err := db.GetCardsFromMatch("UsedIn", matchID, w)
	if err != nil {
		logger.LogError(fmt.Sprintf("[backend] error getting match cards %s", err))
		return nil
	}

	cards := []Card{}

	for _, card := range matchCards {
		c := Card{ID: card}
		// IsBase
		isBase, err := db.GetBoolFromTable("IsBase", card, w)
		if err != nil {
			return nil
		}

		if isBase {
			logger.LogInfo(fmt.Sprintf("[backend] is base %s", card))
			continue
		}

		// Owner
		owner, err := db.GetBytes32FromTable("OwnedBy", card, w)
		if err != nil {
			logger.LogInfo(fmt.Sprintf("[backend] ownedby  error %s", err))
			return nil
		}
		logger.LogInfo(fmt.Sprintf("[backend] ownedby  %s", card))

		c.Owner = owner

		// MaxHp
		value, err := db.GetInt64FromTable("MaxHp", card, w)
		if err != nil {
			return nil
		}
		c.MaxHp = value
		logger.LogInfo(fmt.Sprintf("[backend] hp %s", card))

		// CurrentHp
		value, err = db.GetInt64FromTable("CurrentHp", card, w)
		if err != nil {
			return nil
		}
		c.CurrentHp = value
		logger.LogInfo(fmt.Sprintf("[backend] hp2  %s", card))

		// UnitType
		value, err = db.GetInt64FromTable("UnitType", card, w)
		if err != nil {
			return nil
		}
		c.Type = value
		logger.LogInfo(fmt.Sprintf("[backend] type  %s", card))

		if value == baseType {
			c.Placed = false
			c.Position = Position{X: -2, Y: -2}
			cards = append(cards, c)
			continue
		}

		// AttackDamage
		value, err = db.GetInt64FromTable("AttackDamage", card, w)
		if err != nil {
			return nil
		}
		c.AttackDamage = value
		logger.LogInfo(fmt.Sprintf("[backend] attk  %s", card))

		// MovementSpeed
		value, err = db.GetInt64FromTable("MovementSpeed", card, w)
		if err != nil {
			return nil
		}
		c.MovementSpeed = value
		logger.LogInfo(fmt.Sprintf("[backend] speed  %s", card))

		// ActionReady
		valueBool, err := db.GetBoolFromTable("ActionReady", card, w)
		if err != nil {
			return nil
		}
		c.ActionReady = valueBool
		logger.LogInfo(fmt.Sprintf("[backend] ActionReady  %s", card))

		// Position
		p, err := db.GetPosition(card, w)
		if err != nil {
			return nil
		}

		logger.LogInfo(fmt.Sprintf("[backend] position1  %s", card))
		if p.X == -2 && p.Y == -2 {
			c.Placed = false
		} else {
			c.Placed = true
		}

		c.Position = p

		cards = append(cards, c)
	}
	ret.Cards = cards

	return &ret
}

func (db *Database) GetPosition(rowID string, w *World) (Position, error) {
	table := w.GetTableByName("Position")
	if table == nil {
		errorMsg := fmt.Sprintf("[backend] table %s does not exist", "Position")
		logger.LogError(errorMsg)
		return Position{X: -2, Y: -2}, fmt.Errorf(errorMsg)
	}

	value, ok := (*table.Rows)[rowID]
	if !ok {
		return Position{X: -2, Y: -2}, nil
	}

	x, err := strconv.ParseInt(value[2].Data.String(), 10, 32)
	if err != nil {
		errorMsg := fmt.Sprintf("[backend] could not parse X from %s value %s", "Position", value[2].Data.String())
		logger.LogError(errorMsg)
		return Position{X: -2, Y: -2}, fmt.Errorf(errorMsg)
	}

	y, err := strconv.ParseInt(value[3].Data.String(), 10, 32)
	if err != nil {
		errorMsg := (fmt.Sprintf("[backend] could not parse X from %s value %s", "Position", value[3].Data.String()))
		logger.LogError(errorMsg)
		return Position{X: -2, Y: -2}, fmt.Errorf(errorMsg)
	}

	return Position{X: x, Y: y}, nil
}

func (db *Database) GetInt64FromTable(tableName string, rowID string, w *World) (int64, error) {
	table := w.GetTableByName(tableName)
	if table == nil {
		errorMsg := fmt.Sprintf("[backend] table %s does not exist", tableName)
		logger.LogError(errorMsg)
		return 0, fmt.Errorf(errorMsg)
	}

	value, ok := (*table.Rows)[rowID]
	if !ok {
		errorMsg := fmt.Sprintf("[backend] row in (%s) does not exist for id %s", tableName, rowID)
		logger.LogError(errorMsg)
		return 0, fmt.Errorf(errorMsg)
	}

	valueParsed, err := strconv.ParseInt(value[0].Data.String(), 10, 32)
	if err != nil {
		logger.LogError(fmt.Sprintf("[backend] could not parse %s value %s", tableName, value[0].Data.String()))
	}
	return valueParsed, nil
}

func (db *Database) GetBoolFromTable(tableName string, rowID string, w *World) (bool, error) {
	table := w.GetTableByName(tableName)
	if table == nil {
		errorMsg := fmt.Sprintf("[backend] table %s does not exist", tableName)
		logger.LogError(errorMsg)
		return false, fmt.Errorf(errorMsg)
	}

	value, ok := (*table.Rows)[rowID]
	if !ok {
		// If the row does not exists is false by default
		return false, nil
	}

	return value[0].Data.String() == "true", nil
}

func (db *Database) GetBytes32FromTable(tableName string, rowID string, w *World) (string, error) {
	// UsedIn
	usedIn := w.GetTableByName(tableName)
	if usedIn == nil {
		errorMsg := fmt.Sprintf("[backend] table %s does not exist", tableName)
		logger.LogError(errorMsg)
		return "", fmt.Errorf(errorMsg)
	}

	for k, v := range *usedIn.Rows {
		// The string method returns "0x..." so we use contains to compare
		logger.LogDebug(fmt.Sprintf("[backend] bytes32 element %s", v[0].Data.String()))
		if k == rowID {
			return v[0].Data.String(), nil
		}
	}

	return "", nil
}

func (db *Database) GetCardsFromMatch(tableName string, rowID string, w *World) ([]string, error) {
	// UsedIn
	usedIn := w.GetTableByName(tableName)
	if usedIn == nil {
		errorMsg := fmt.Sprintf("[backend] table %s does not exist", tableName)
		logger.LogError(errorMsg)
		return []string{}, fmt.Errorf(errorMsg)
	}

	IDs := []string{}
	for k, v := range *usedIn.Rows {
		// The string method returns "0x..." so we use contains to compare
		if strings.Contains(v[0].Data.String(), rowID) {
			IDs = append(IDs, k)
		}
	}

	return IDs, nil
}

func (db *Database) GetBoardStatus(worldId string, userWallet string) *MatchData {
	w, ok := db.Worlds[worldId]
	if !ok {
		logger.LogError("[backend] world not found")
		return nil
	}
	matchID := db.GetUserMatchID(w, userWallet)
	if matchID == "" {
		return nil
	}

	data := db.GetBoard(w, matchID)

	return data
}
