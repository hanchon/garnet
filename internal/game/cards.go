package game

// This can be read from the blockchain but it is static

type Card struct {
	Name      string
	MaxHealth int
	Attack    int
	Movement  int
	Symbol    string
}

var (
	AvailableCards = []Card{
		{
			Name:      "Vaan Strife",
			MaxHealth: 6,
			Attack:    4,
			Movement:  2,
			Symbol:    "(♚)",
		},
		{
			Name:      "Felguard",
			MaxHealth: 6,
			Attack:    4,
			Movement:  2,
			Symbol:    "(♛)",
		},
		{
			Name:      "Sakura",
			MaxHealth: 5,
			Attack:    3,
			Movement:  3,
			Symbol:    "(♜)",
		},
		{
			Name:      "Freya",
			MaxHealth: 5,
			Attack:    3,
			Movement:  3,
			Symbol:    "(♝)",
		},
		{
			Name:      "Lyra",
			MaxHealth: 5,
			Attack:    1,
			Movement:  1,
			Symbol:    "(♞)",
		},
		{
			Name:      "Madmartigan",
			MaxHealth: 10,
			Attack:    2,
			Movement:  1,
			Symbol:    "(♟)",
		},
	}
)
