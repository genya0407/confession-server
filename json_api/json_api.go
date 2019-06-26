package json_api

import (
	"fmt"
	"github.com/genya0407/confession-server/usecase"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	// "log"
	"encoding/json"
	"net/http"
	"time"
)

type AccountJSON struct {
	AccountID  uuid.UUID `json:"account_id"`
	Name       string    `json:"name"`
	ScreenName string    `json:"screen_name"`
	ImageURL   string    `json:"image_url"`
}

type ChatJSON struct {
	ChatID     uuid.UUID     `json:"chat_id"`
	Account    AccountJSON   `json:"account"`
	Messages   []MessageJSON `json:"messages"`
	StartedAt  time.Time     `json:"started_at"` // Go uses ISO8601 format by default
	FinishedAt *time.Time    `json:"finished_at,omitempty"`
}

type MessageJSON struct {
	MessageID   uuid.UUID `json:"message_id"`
	Text        string    `json:"text"`
	ByAnonymous bool      `json:"by_anonymous"`
	SentAt      time.Time `json:"sent_at"`
}

type ChatAbstractJSON struct {
	ChatID               uuid.UUID  `json:"chat_id"`
	BeginningMessageText string     `json:"beginning_message_text"`
	StartedAt            time.Time  `json:"started_at"`
	FinishedAt           *time.Time `json:"finished_at,omitempty"`
}

type CreateChatRequestJSON struct {
	BeginningMessageText string `json:"beginning_message_text"`
}

type CreateChatResponseJSON struct {
	ChatID                uuid.UUID `json:"chat_id"`
	AnonymousSessionToken string    `json:"session_token"`
}

type EverybodyEndpoint = func(http.ResponseWriter, *http.Request, httprouter.Params)
type AccountAuthorizedEndpoint = func(http.ResponseWriter, *http.Request, httprouter.Params, usecase.AccountLoginInfoDTO)
type AnonymousAuthorizedEndpoint = func(http.ResponseWriter, *http.Request, httprouter.Params, usecase.AnonymousLoginInfoDTO)

func GenerateGetAccountInfo(getAccountInfo usecase.GetAccountInfo) EverybodyEndpoint {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		accountIDString := ps.ByName("account_id")
		accountID, err := uuid.Parse(accountIDString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid AccountID")
			return
		}

		account, ok := getAccountInfo(accountID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Account not found: %s", accountID.String())
			return
		}

		accountJSON := AccountJSON{
			AccountID:  account.AccountID,
			Name:       account.Name,
			ScreenName: account.ScreeName,
			ImageURL:   account.ImageURL,
		}
		json.NewEncoder(w).Encode(accountJSON)
	}
}

func GenerateCreateChat(createChat usecase.CreateChat) EverybodyEndpoint {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		accountIDString := ps.ByName("account_id")
		accountID, err := uuid.Parse(accountIDString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid AccountID")
			return
		}

		createChatJSON := &CreateChatRequestJSON{}
		err = json.NewDecoder(r.Body).Decode(createChatJSON)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid request")
			return
		}

		createChatResultDTO, createChatErr := createChat(accountID, createChatJSON.BeginningMessageText)
		if createChatErr != nil {
			switch *createChatErr {
			case usecase.AccountNotFound:
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "Account not found")
			default:
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "Internal server error")
			}
			return
		}

		response := &CreateChatResponseJSON{
			ChatID:                createChatResultDTO.Chat.ChatID,
			AnonymousSessionToken: createChatResultDTO.AnonymousLoginInfo.SessionToken,
		}
		json.NewEncoder(w).Encode(response)
	}
}
