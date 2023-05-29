package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hanchon/garnet/internal/backend/cors"
	"github.com/hanchon/garnet/internal/database"
	"github.com/hanchon/garnet/internal/logger"
	"github.com/hanchon/garnet/internal/txbuilder"
)

func RestRoutes(router *mux.Router, db *database.InMemoryDatabase) {
	router.HandleFunc(
		"/signup",
		func(response http.ResponseWriter, request *http.Request) {
			RegisterEndpoint(response, request, db)
		},
	).Methods("POST", "OPTIONS")

	router.HandleFunc("/ping", RegisterPing).Methods("GET", "POST", "OPTIONS")
}

func SendInternalErrorResponse(msg string, w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(*w, msg)
}

func SendBadRequestResponse(msg string, w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusBadRequest)
	fmt.Fprint(*w, msg)
}

func SendJSONResponse(message interface{}, w *http.ResponseWriter) error {
	v, err := json.Marshal(message)
	if err != nil {
		SendInternalErrorResponse("invalid encoding for response", w)
		return err
	}
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
	_, err = (*w).Write(v)
	return err
}

// RegistationParams is struct to read the request body
type RegistationParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistationResponse struct {
	Username string `json:"username"`
	Index    int    `json:"index"`
}

func RegisterPing(response http.ResponseWriter, request *http.Request) {
	logger.LogDebug("[api] getting a request to register an user")
	if cors.SetHandlerCorsForOptions(request, &response) {
		return
	}
	response.WriteHeader(http.StatusOK)
	_, err := response.Write([]byte("pong"))
	if err != nil {
		logger.LogDebug(fmt.Sprintf("[api] error sending pong: %s", err.Error()))
	}
}

func RegisterEndpoint(response http.ResponseWriter, request *http.Request, db *database.InMemoryDatabase) {
	logger.LogDebug("[api] getting a request to register an user")
	if cors.SetHandlerCorsForOptions(request, &response) {
		return
	}

	var registationRequest RegistationParams

	decoder := json.NewDecoder(request.Body)
	decoderErr := decoder.Decode(&registationRequest)
	defer request.Body.Close()

	if decoderErr != nil {
		SendInternalErrorResponse("decode error", &response)
		return
	}

	if registationRequest.Username == "" || registationRequest.Password == "" {
		SendBadRequestResponse("invalid params", &response)
		return
	}

	logger.LogDebug(fmt.Sprintf("[api] registering user %s", registationRequest.Username))

	index, err := db.RegisterUser(registationRequest.Username, registationRequest.Password)
	if err != nil {
		SendBadRequestResponse(err.Error(), &response)
		return

	}

	logger.LogDebug(fmt.Sprintf("[api] registered user %s with id %d", registationRequest.Username, index))

	registrationResponse := RegistationResponse{
		Username: registationRequest.Username,
		Index:    index,
	}

	if _, account, err := txbuilder.GetWallet(index); err == nil {
		hash, errFaucet := txbuilder.Faucet(account.Address.Hex())
		if errFaucet == nil {
			logger.LogError(fmt.Sprintf("[backend] coins sent to wallet %s", account.Address.Hex()))
			successful := false
			retry := 0

			for {
				if retry > 3 {
					break
				}
				successful, err = txbuilder.WasTransactionSuccessful(hash)
				if err == nil {
					break
				}
				time.Sleep(time.Second)
				retry++
			}

			if successful {
				logger.LogError(fmt.Sprintf("[backend] registering the wallet %s in the chain %s", account.Address.Hex(), registationRequest.Username))
				var nameBytes [32]byte
				copy(nameBytes[:], registationRequest.Username)
				err = txbuilder.SendTransaction(index, "register", nameBytes)
				if err != nil {
					logger.LogError(fmt.Sprintf("[backend] error registering wallet %s, %s", account.Address.Hex(), err.Error()))
				} else {
					logger.LogError(fmt.Sprintf("[backend] wallet registered correctly %s", account.Address.Hex()))
				}
			}
		} else {
			logger.LogError(fmt.Sprintf("[backend] error sending coins to wallet %s, %s", account.Address.Hex(), errFaucet.Error()))
		}
	}

	err = SendJSONResponse(registrationResponse, &response)
	if err != nil {
		logger.LogError(fmt.Sprintf("[backend] error sending response to client %s", err.Error()))
	}
}
