package entity

import (
	"github.com/google/uuid"
	"time"
)

// utils

func mustNewUUID() uuid.UUID {
	u, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	return u
}

// data structures

type TwitterAccount struct {
	TwitterAccountID string
	AccountID        uuid.UUID
}

type Account struct {
	AccountID  uuid.UUID
	Name       string
	ScreenName string
	ImageUrl   string
	Token      string
}

type Anonymous struct {
	AnonymousID uuid.UUID
	Token       string
}

type Message struct {
	MessageID   uuid.UUID
	Text        TextMessage
	ByAnonymous bool
	SentAt      time.Time
}

type TextMessage = string

type Socket interface {
	SendText(msg TextMessage)
	Close()
}

type ChatID = uuid.UUID

type Chat struct {
	ChatID          ChatID
	Account         Account
	Anonymous       Anonymous
	Messages        []Message
	StartedAt       time.Time
	FinishedAt      *time.Time
	AccountSocket   Socket
	AnonymousSocket Socket
	StoreMessage    StoreMessage
	FinishChat      FinishChat
}

func (c *Chat) Close() {
	c.FinishChat(c.ChatID, time.Now())

	if c.AccountSocket != nil {
		c.AccountSocket.Close()
	}

	if c.AnonymousSocket != nil {
		c.AnonymousSocket.Close()
	}
}

func (c *Chat) SendAnonymousMessageToAccount(text TextMessage) {
	msg := Message{
		MessageID:   mustNewUUID(),
		Text:        text,
		ByAnonymous: true,
		SentAt:      time.Now(),
	}
	c.StoreMessage(c.ChatID, msg)
	c.Messages = append(c.Messages, msg)

	if c.AccountSocket != nil {
		c.AccountSocket.SendText(msg.Text)
	}
}

func (c *Chat) SendAccountMessageToAnonymous(text TextMessage) {
	msg := Message{
		MessageID:   mustNewUUID(),
		Text:        text,
		ByAnonymous: false,
		SentAt:      time.Now(),
	}
	c.StoreMessage(c.ChatID, msg)
	c.Messages = append(c.Messages, msg)

	if c.AnonymousSocket != nil {
		c.AnonymousSocket.SendText(msg.Text)
	}
}

// repositories

type FindAccountByID = func(uuid.UUID) Account
type FindAccountByToken = func(string) Account
type FindAnonymousByToken = func(string) Anonymous
type FindChatByID = func(uuid.UUID) Chat
type StoreMessage = func(uuid.UUID, Message)
type FinishChat = func(uuid.UUID, time.Time)
