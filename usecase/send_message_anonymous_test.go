package usecase

import (
	"github.com/genya0407/confession-server/domain"
	"github.com/genya0407/confession-server/repository"
	"log"
	"testing"
)

type mockSocket struct {
	c chan domain.Message
}

func (ms *mockSocket) SendText(msg domain.Message) {
	ms.c <- msg
}

func (ms *mockSocket) Close() {
	close(ms.c)
}

type MockSocket struct {
	c chan domain.IMessage
}

func (ms *MockSocket) SendText(msg domain.IMessage) {
	ms.c <- msg
}

func (ms *MockSocket) Close() {
	close(ms.c)
}

func TestSendMessageAnonymousToAccount(t *testing.T) {
	repo := repository.NewOnMemoryRepository()
	account := domain.NewAccount("Yusuke Sangenya", "yusuke.sangenya", "http://hogehoge.com/img.png")
	repo.AccountStorage[account.Token()] = account
	chat := domain.NewChat(account, "Lets start some chat!")
	repo.StoreChat(chat)

	sendAnonymousMessageToAccountService := domain.GenerateSendAnonymousMessageToAccountService(
		repo.StoreChat,
	)
	sendMessageAnonymousToAccountUsecase := GenerateSendMessageAnonymousToAccount(
		sendAnonymousMessageToAccountService,
		repo.FindAnonymousByToken,
		repo.FindChatByID,
	)

	msgText := "abcdef"

	accChan := make(chan domain.IMessage)
	accSocket := &MockSocket{c: accChan}
	joinChatAccountService := domain.GenerateJoinChatAccountService(repo.FindChatByID, repo.StoreChat)
	joinChatAccountService(chat, accSocket)

	anonChan := make(chan domain.IMessage)
	anonSocket := &MockSocket{c: anonChan}
	joinChatAnonymousService := domain.GenerateJoinChatAnonymousService(repo.FindChatByID, repo.StoreChat)
	joinChatAnonymousService(chat, anonSocket)

	go func() {
		err := sendMessageAnonymousToAccountUsecase(
			AnonymousLoginInfoDTO{SessionToken: chat.Anonymous().Token()},
			chat.ChatID(),
			msgText,
		)
		if err != nil {
			log.Print(err.Error())
		}
	}()
	sentMessage := <-accChan
	if sentMessage.Text() != msgText {
		t.Error("Invalid message")
	}
}
