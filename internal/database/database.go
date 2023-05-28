package database

import (
	"fmt"
	"sync"

	"github.com/hanchon/garnet/internal/txbuilder"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password []byte
	Address  string
	Index    int
}

type InMemoryDatabase struct {
	Users *map[string]User
	mu    *sync.Mutex
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		Users: &map[string]User{},
		mu:    &sync.Mutex{},
	}
}

func (db *InMemoryDatabase) RegisterUser(username string, password string) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := (*db.Users)[username]; ok {
		return 0, fmt.Errorf("user already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	index := len(*db.Users) + 1

	_, account, err := txbuilder.GetWallet(index)
	if err != nil {
		return 0, err
	}

	(*db.Users)[username] = User{
		Username: username,
		Password: hashedPassword,
		Index:    index,
		Address:  account.Address.Hex(),
	}

	return index, nil
}

func (db *InMemoryDatabase) Login(username string, password string) (User, error) {
	v, exists := (*db.Users)[username]
	if !exists {
		return User{}, fmt.Errorf("user is not in the database")
	}
	// TODO: add a timer here so the Login is not spammed
	if err := bcrypt.CompareHashAndPassword(v.Password, []byte(password)); err != nil {
		return User{}, fmt.Errorf("incorrect password")
	}
	return v, nil
}
