package json_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/genya0407/confession-server/usecase"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func mustNewUUID() uuid.UUID {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err.Error())
	}
	return id
}

func generateMockUsecaseCreateChat(expectAccountID uuid.UUID, chatID uuid.UUID, expectedToken string) usecase.CreateChat {
	return func(accountID uuid.UUID, text string) (usecase.CreateChatResultDTO, *usecase.CreateChatError) {
		if accountID != expectAccountID {
			return usecase.CreateChatResultDTO{}, &usecase.AccountNotFound
		}

		return usecase.CreateChatResultDTO{
			AnonymousLoginInfo: usecase.AnonymousLoginInfoDTO{
				SessionToken: expectedToken,
			},
			Chat: usecase.ChatDTO{
				ChatID: chatID,
				Messages: []usecase.MessageDTO{
					usecase.MessageDTO{
						MessageID:   mustNewUUID(),
						Text:        text,
						ByAnonymous: true,
						SentAt:      time.Now(),
					},
				},
			},
		}, nil
	}
}

func generateMockCreateChatRouter(accountID uuid.UUID, chatID uuid.UUID, sessionToken string) *httprouter.Router {
	uc := generateMockUsecaseCreateChat(accountID, chatID, sessionToken)
	handler := GenerateCreateChat(uc)
	router := httprouter.New()
	router.POST("/anonymous/account/:account_id/chats", handler)
	return router
}

func TestCreateChat(t *testing.T) {
	accountID := mustNewUUID()
	chatID := mustNewUUID()
	expectedSessionToken := "aaaaa"
	router := generateMockCreateChatRouter(accountID, chatID, expectedSessionToken)

	reqBody, err := json.Marshal(CreateChatRequestJSON{
		BeginningMessageText: "This is first text.",
	})
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(
		"POST",
		fmt.Sprintf("http://confession.com/anonymous/account/%s/chats", accountID.String()),
		bytes.NewReader(reqBody),
	)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("Invalid status")
		return
	}

	response := &CreateChatResponseJSON{}
	err = json.NewDecoder(w.Body).Decode(response)
	if err != nil {
		panic(err)
	}

	if response.AnonymousSessionToken == "" {
		t.Error("Empty token")
		return
	}

	if response.ChatID != chatID {
		t.Error("Invalid chatID")
		return
	}

	if response.AnonymousSessionToken != expectedSessionToken {
		t.Error("Invalid session token")
		return
	}
}

// func TestGetAccountInfoNotFound(t *testing.T) {
// 	mockname := "Mock name"
// 	accountID, err := uuid.NewUUID()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	router := generateMockGetAccountInfoRouter(mockname, accountID)

// 	notExistAccountID, err := uuid.NewUUID()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	req := httptest.NewRequest("GET", "http://confession.com/account/"+notExistAccountID.String(), nil)
// 	w := httptest.NewRecorder()

// 	router.ServeHTTP(w, req)
// 	if w.Code != http.StatusNotFound {
// 		t.Error("Invalid status")
// 	}
// }
