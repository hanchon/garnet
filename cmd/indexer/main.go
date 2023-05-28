package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hanchon/garnet/internal/backend"
	"github.com/hanchon/garnet/internal/backend/api"
	"github.com/hanchon/garnet/internal/indexer"
	"github.com/hanchon/garnet/internal/indexer/data"
)

const (
	port = 6666
)

func pingServer() error {
	client := &http.Client{
		Timeout: time.Second,
	}

	r, err := client.Get(fmt.Sprintf("http://localhost:%d/ping", port))
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("incorrect response: %d", r.StatusCode)
	}

	return nil
}

func createUser(username string, password string) error {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	body := api.RegistationParams{
		Username: username,
		Password: password,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/signup", port), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/json")

	res, err := client.Do(r)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("incorrect response: %d", r.Response.StatusCode)
	}

	return nil
}

func main() {
	// Set the log output to a file (stdin, stdout, stderror used by GUI)
	fileName := "logFile.log"
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)

	// Index the database
	quit := false
	database := data.NewDatabase()
	go indexer.Process(database, &quit)

	// Set up the GUI
	ui := NewDebugUI()
	defer ui.ui.Close()

	go ui.ProcessIncomingData(database)
	go ui.ProcessBlockchainInfo(database)
	go ui.ProcessLatestEvents(database)

	// Start the backend server
	go func() {
		_ = backend.StartGorillaServer(port, database)
	}()

	err = fmt.Errorf("")
	for err != nil {
		// This will ping the server with 1 second timeout until it's alive
		err = pingServer()
	}
	_ = createUser("user1", "password1")
	_ = createUser("user2", "password2")

	// Display the GUI
	ui.Run()

	// Exit program
	quit = true
}
