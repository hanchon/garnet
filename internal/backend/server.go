package backend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hanchon/garnet/internal/backend/api"
	"github.com/hanchon/garnet/internal/backend/cors"
	"github.com/hanchon/garnet/internal/backend/messages"
	"github.com/hanchon/garnet/internal/database"
	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/logger"
)

func StartGorillaServer(port int, mudDatabase *data.Database) error {
	logger.LogInfo(fmt.Sprintf("[backend] starting server at port: %d\n", port))
	router := mux.NewRouter()
	usersDatabase := database.NewInMemoryDatabase()
	g := messages.NewGlobalState(mudDatabase, usersDatabase)
	router.HandleFunc("/ws", g.WebSocketConnectionHandler).Methods("GET", "OPTIONS")
	api.RestRoutes(router, usersDatabase)
	go g.BroadcastUpdates()

	cors.ServerEnableCORS(router)

	server := &http.Server{
		Addr:              fmt.Sprint(":", port),
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}
	return server.ListenAndServe()
}
