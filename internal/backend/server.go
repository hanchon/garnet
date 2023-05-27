package backend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hanchon/garnet/internal/backend/cors"
	"github.com/hanchon/garnet/internal/backend/messages"
	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/logger"
)

func StartGorillaServer(port int, database *data.Database) error {
	logger.LogInfo(fmt.Sprintf("[backend] starting server at port: %d\n", port))
	router := mux.NewRouter()
	g := messages.NewGlobalState(database)
	router.HandleFunc("/ws", g.WebSocketConnectionHandler).Methods("GET", "OPTIONS")
	go g.BroadcastUpdates()

	cors.ServerEnableCORS(router)

	server := &http.Server{
		Addr:              fmt.Sprint(":", port),
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}
	return server.ListenAndServe()
}
