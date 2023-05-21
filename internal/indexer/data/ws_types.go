package data

type PlacedCards struct {
	P1cards int64 `json:"p1cards"`
	P2cards int64 `json:"p2cards"`
}

type Position struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type MatchData struct {
	MatchID           string `json:"matchid"`
	PlayerOne         string `json:"playerone"`
	PlayerTwo         string `json:"playertwo"`
	PlayerOneUsermane string `json:"playeroneusername"`
	PlayerTwoUsermane string `json:"playertwousername"`
	CurrentTurn       int64  `json:"currenturn"`
	CurrentPlayer     string `json:"currentplayer"`
	CurrentMana       int64  `json:"currentmana"`
	// PlacedCards   PlacedCards `json:"placedcards"`
	Cards []Card `json:"cards"`
	// P1Cards       []Card      `json:"p1cards"`
	// P2Cards       []Card      `json:"p2cards"`
	// P1Base        Base        `json:"p1base"`
	// P2Base        Base        `json:"p2base"`
}

type Card struct {
	ID            string   `json:"id"`
	Type          int64    `json:"type"`
	AttackDamage  int64    `json:"attackdamage"`
	MaxHp         int64    `json:"maxhp"`
	CurrentHp     int64    `json:"currenthp"`
	MovementSpeed int64    `json:"movementspeed"`
	Position      Position `json:"position"`
	Owner         string   `json:"owner"`
	ActionReady   bool     `json:"actionready"`
	Placed        bool     `json:"placed"`
}

type Base struct {
	ID        string `json:"id"`
	MaxHp     int64  `json:"maxhp"`
	CurrentHp int64  `json:"currenthp"`
}

// var DummyData MatchData = MatchData{
// 	MatchID:       "0x0000000000000000000000000000000000000000000000000000000000000044",
// 	PlayerOne:     "0x0000000000000000000000001fc55fec81842096b61404aec5403d8aceb5d590",
// 	PlayerTwo:     "0x0000000000000000000000001fc55fec81842096b61404aec5403d8aceb5d590",
// 	CurrentTurn:   0,
// 	CurrentPlayer: "0x0000000000000000000000001fc55fec81842096b61404aec5403d8aceb5d590",
// 	CurrentMana:   2,
// 	PlacedCards:   PlacedCards{P1cards: 1, P2cards: 0},
// 	P1Cards: []Card{
// 		{
// 			ID:            "0xefbfbd0e2eefbfbdefbfbdefbfbdefbfbd52efbfbd166b3333efbfbdefbfbdefbfbd5235efbfbdefbfbdefbfbd18efbfbd10efbfbd16efbfbdefbfbd17efbfbdefbfbd12",
// 			Type:          0,
// 			AttackDamage:  1,
// 			MaxHp:         10,
// 			CurrentHp:     10,
// 			MovementSpeed: 1,
// 			Position:      Position{X: 0, Y: 0},
// 		},
// 	},
// 	P2Cards: []Card{
// 		{
// 			ID:            "0x4aefbfbd6017efbfbd6b4266efbfbd78456529efbfbd736a4defbfbd75efbfbdefbfbdefbfbd4eefbfbdefbfbd723a4fefbfbd4132efbfbd",
// 			Type:          1,
// 			AttackDamage:  1,
// 			MaxHp:         10,
// 			CurrentHp:     10,
// 			MovementSpeed: 1,
// 			Position:      Position{X: 1, Y: 8},
// 		},
// 	},
// 	P1Base: Base{ID: "0x1234", MaxHp: 10, CurrentHp: 10},
// 	P2Base: Base{ID: "0x4321", MaxHp: 10, CurrentHp: 10},
// }
