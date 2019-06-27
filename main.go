package main

import (
	"fmt"
	"github.com/genya0407/confession-server/domain"
	"github.com/genya0407/confession-server/jsonapi"
	"github.com/genya0407/confession-server/repository"
	"github.com/genya0407/confession-server/usecase"
	"github.com/julienschmidt/httprouter"
	"github.com/k0kubun/pp"
	"log"
	"net/http"
	"net/url"
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

	joinChatAnonymousHandler := jsonapi.AuthorizeAnonymous(
		jsonapi.GenerateJoinChatAnonymous(
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
	joinChatAccountHandler := jsonapi.AuthorizeAccount(
		jsonapi.GenerateJoinChatAccount(
			joinChatAccount,
			sendMessageAccountToAnonymous,
		),
	)

	return joinChatAccountHandler
}

func main() {
	repo := repository.NewOnMemoryRepository()
	account := domain.NewAccount("Yusuke Sangenya", "yusuke.sangenya", "http://hogehoge.com/img.png")
	repo.AccountStorage[account.Token()] = account
	pp.Println(account)
	pp.Println(account.AccountID().String())
	chat := domain.NewChat(account, "Lets start some chat!")
	repo.StoreChat(chat)
	pp.Println(chat)
	pp.Println(chat.ChatID().String())

	accountWSURL := fmt.Sprintf(`ws://localhost:8080/connect/chat/%s?access_token=%s`, chat.ChatID(), url.QueryEscape(account.Token()))
	anonymousWSURL := fmt.Sprintf(`ws://localhost:8080/anonymous/account/%s/chat/%s?access_token=%s`, account.AccountID(), chat.ChatID(), url.QueryEscape(chat.Anonymous().Token()))

	fmt.Printf(`
acc = new WebSocket("%s");
acc.onmessage = function(msg) { console.log(msg) };

anon = new WebSocket("%s");
anon.onmessage = function(msg) { console.log(msg) };

`, accountWSURL, anonymousWSURL)

	router := httprouter.New()
	router.GET(`/anonymous/account/:account_id/chat/:chat_id`, joinChatAnonymousHandler(repo))
	router.GET(`/connect/chat/:chat_id`, joinChatAccountHandler(repo))

	fmt.Println("Start at: http://localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
