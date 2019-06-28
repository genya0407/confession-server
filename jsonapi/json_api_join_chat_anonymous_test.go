package jsonapi

import (
	"fmt"
	"github.com/genya0407/confession-server/domain"
	"github.com/genya0407/confession-server/repository"
	"github.com/genya0407/confession-server/usecase"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	// "github.com/k0kubun/pp"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func joinChatAnonymousHandler(repo *repository.OnMemoryRepository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	joinChatAnonymousService := domain.GenerateJoinChatAnonymousService(
		repo.FindChatByID,
		repo.StoreChat,
	)
	joinChatAnonymous := usecase.GenerateJoinChatAnonymous(
		joinChatAnonymousService,
		repo.FindAnonymousByToken,
		repo.FindChatByID,
	)

	sendAnonymousMessageToAccountService := domain.GenerateSendAnonymousMessageToAccountService(
		repo.StoreChat,
	)
	sendMessageAnonymousToAccount := usecase.GenerateSendMessageAnonymousToAccount(
		sendAnonymousMessageToAccountService,
		repo.FindAnonymousByToken,
		repo.FindChatByID,
	)

	joinChatAnonymousHandler := AuthorizeAnonymous(
		GenerateJoinChatAnonymous(
			joinChatAnonymous,
			sendMessageAnonymousToAccount,
		),
	)

	return joinChatAnonymousHandler
}

func joinChatAccountHandler(repo *repository.OnMemoryRepository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	joinChatAccountService := domain.GenerateJoinChatAccountService(repo.FindChatByID, repo.StoreChat)
	joinChatAccount := usecase.GenerateJoinChatAccount(
		joinChatAccountService,
		repo.FindAccountByToken,
		repo.FindChatByID,
	)
	sendAccountMessageToAnonymousService := domain.GenerateSendAccountMessageToAnonymousService(
		repo.StoreChat,
	)
	sendMessageAccountToAnonymous := usecase.GenerateSendMessageAccountToAnonymous(
		sendAccountMessageToAnonymousService,
		repo.FindAccountByToken,
		repo.FindChatByID,
	)
	joinChatAccountHandler := AuthorizeAccount(
		GenerateJoinChatAccount(
			joinChatAccount,
			sendMessageAccountToAnonymous,
		),
	)

	return joinChatAccountHandler
}

func TestChat(t *testing.T) {
	repo := repository.NewOnMemoryRepository()
	account := domain.NewAccount("Yusuke Sangenya", "yusuke.sangenya", "http://hogehoge.com/img.png")
	repo.AccountStorage[account.Token()] = account
	// pp.Println(account)
	// pp.Println(account.AccountID().String())
	chat := domain.NewChat(account, "Lets start some chat!")
	repo.StoreChat(chat)
	// pp.Println(chat)
	// pp.Println(chat.ChatID().String())

	router := httprouter.New()
	router.GET(`/anonymous/account/:account_id/chat/:chat_id`, joinChatAnonymousHandler(repo))
	router.GET(`/connect/chat/:chat_id`, joinChatAccountHandler(repo))

	s := httptest.NewServer(router)
	defer s.Close()

	accountWSURL := fmt.Sprintf(
		`ws://%s/connect/chat/%s?access_token=%s`,
		strings.TrimPrefix(s.URL, "http://"),
		chat.ChatID(),
		url.QueryEscape(account.Token()),
	)
	anonymousWSURL := fmt.Sprintf(
		`ws://%s/anonymous/account/%s/chat/%s?access_token=%s`,
		strings.TrimPrefix(s.URL, "http://"),
		account.AccountID(),
		chat.ChatID(),
		url.QueryEscape(chat.Anonymous().Token()),
	)

	accountWS, _, err := websocket.DefaultDialer.Dial(accountWSURL, nil)
	if err != nil {
		panic(err.Error())
	}
	defer accountWS.Close()
	anonymousWS, _, err := websocket.DefaultDialer.Dial(anonymousWSURL, nil)
	if err != nil {
		panic(err.Error())
	}
	defer anonymousWS.Close()

	text := "some test message"
	err = accountWS.WriteJSON(NewMessageJSON{Text: text})
	if err != nil {
		panic(err.Error())
	}
	msg := &MessageJSON{}
	err = anonymousWS.ReadJSON(msg)
	if err != nil {
		panic(err.Error())
	}
	if msg.Text != text {
		t.Error("Invalid message received")
		return
	}

	err = accountWS.ReadJSON(msg)
	if err != nil {
		panic(err.Error())
	}
	if msg.Text != text {
		t.Error("Invalid sent-back message received")
		return
	}
}
