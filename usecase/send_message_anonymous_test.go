package usecase

import (
	"github.com/genya0407/confession-server/entity"
	"github.com/google/uuid"
	"testing"
	"time"
)

func mustNewUUID() uuid.UUID {
	u, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return u
}

type mockSocket struct {
	c chan entity.Message
}

func (ms *mockSocket) SendText(msg entity.Message) {
	ms.c <- msg
}

func (ms *mockSocket) Close() {
	close(ms.c)
}

func TestSendMessageAnonymousToAccount(t *testing.T) {
	anon := AnonymousLoginInfoDTO{SessionToken: "aaaaa"}
	expectedChatID := mustNewUUID()
	c := make(chan entity.Message)
	findChatAnonymous := func(cID entity.ChatID, anon entity.Anonymous) (entity.Chat, bool) {
		if cID == expectedChatID {
			return entity.Chat{
				ChatID:        expectedChatID,
				AccountSocket: &mockSocket{c: c},
				Anonymous: entity.Anonymous{
					Token: anon.Token,
				},
			}, true
		}

		return entity.Chat{}, false
	}

	msgText := "abcdef"
	msg := MessageDTO{
		MessageID:   mustNewUUID(),
		Text:        msgText,
		SentAt:      time.Now(),
		ByAnonymous: true,
	}

	sendMessageAnonymousToAccount := GenerateSendMessageAnonymousToAccount(findChatAnonymous)
	go func() {
		sendMessageAnonymousToAccount(anon, expectedChatID, msg)
	}()
	sentMessage := <-c
	if sentMessage.Text != msgText {
		t.Error("Invalid message")
	}
}
